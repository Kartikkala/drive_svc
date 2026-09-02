package drive

import (
	"context"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

func AttachEvents(ctx context.Context, nc *nats.Conn, svc *DriveService) {
	nc.Subscribe("authentication.user.registered", handleNewUserRegistration(ctx, nc, svc))
	nc.Subscribe("authorization.user.status", updateAuthorizationStatus(ctx, nc, svc))
	nc.Subscribe("storage.new.rootnode", updateStorageProvisioningStatus(ctx, nc, svc))
}

func handleNewUserRegistration(ctx context.Context, nc *nats.Conn, svc *DriveService) func(*nats.Msg) {
	return func(msg *nats.Msg) {
		var user NewUser
		err := json.Unmarshal(msg.Data, &user)

		if err != nil {
			log.Println("failed to parse json")
			// No notification to NATS on parse failure
			// as user ID is not known
			return
		}

		rootNodeId, err := svc.NewDrive(user.UserID)
		if err != nil {
			log.Printf("failed to create new drive for user with ID: %d\n", user.UserID)
			deadLetterQueue(nc, user.UserID, rootNodeId.String())
			return
		}
		// Actual provisioning starts with provisioning storage
		// and authorizing the user
		status, err := svc.GetStatus(user.UserID)
		if err != nil {
			log.Printf("error while fetching status of the user with ID: %v\n", user.UserID)
			return
		}
		numRetries := 100
		if !status.StorageProvisioned && status.StorageProvisionTries < uint8(numRetries) {
			if err = provisionStorage(nc, rootNodeId.String()); err != nil {
				err = svc.UpdateStatus(
					user.UserID,
					status.Authorized,
					status.StorageProvisioned,
					status.AuthorizationTries,
					status.StorageProvisionTries+1,
				)
				if err != nil {
					log.Printf("error while updating provisioning status of user with ID: %v\n", user.UserID)
				}
				log.Printf("error while provisioning storage for user with ID: %v\n", user.UserID)
				return
			}
		} else {
			log.Printf("storage provision limit reached for user: %v\n", user.UserID)
			deadLetterQueue(nc, user.UserID, rootNodeId.String())
			return
		}
		log.Printf("request to storage svc sent for user ID %v\n", user.UserID)

		if !status.Authorized && status.AuthorizationTries < uint8(numRetries) {
			if err = grantAuthorizationPrevileges(nc, user.UserID, rootNodeId.String(), "owner"); err != nil {
				err = svc.UpdateStatus(
					user.UserID,
					status.Authorized,
					status.StorageProvisioned,
					status.AuthorizationTries+1,
					status.StorageProvisionTries,
				)
				if err != nil {
					log.Printf("error while updating authorization status of user with ID: %v\n", user.UserID)
				}
				log.Printf("error while authorizing user with ID: %v\n", user.UserID)
				return
			}
		} else {
			log.Printf("authorization limit reached for user: %v\n", user.UserID)
			deadLetterQueue(nc, user.UserID, rootNodeId.String())
			return
		}
		log.Printf("request to authorization svc sent for user ID %v\n", user.UserID)
	}
}

func updateAuthorizationStatus(ctx context.Context, nc *nats.Conn, svc *DriveService) func(*nats.Msg) {
	return func(msg *nats.Msg) {
		var regStatus AuthorizationStatus
		if err := json.Unmarshal(msg.Data, &regStatus); err != nil {
			log.Println("failed to parse registration status JSON")
			return
		}

		status, err := svc.GetStatus(regStatus.UserID)
		if err != nil {
			log.Printf("error fetching drive status for user ID %d: %v\n", regStatus.UserID, err)
			return
		}

		// Authorization succeeded
		if regStatus.Success {
			err = svc.UpdateStatus(
				regStatus.UserID,
				true,
				status.StorageProvisioned,
				status.AuthorizationTries,
				status.StorageProvisionTries,
			)
			if err != nil {
				log.Printf("failed to mark user ID %d as authorized: %v\n", regStatus.UserID, err)
				return
			}
			log.Printf("successfully updated authorization status to true for user ID %d\n", regStatus.UserID)
			return
		}

		// Authorization failed: increment attempts or drop to dead-letter queue
		const maxRetries uint8 = 100
		newTries := status.AuthorizationTries + 1

		if newTries >= maxRetries {
			log.Printf("authorization retry limit reached for user ID: %d\n", regStatus.UserID)
			deadLetterQueue(nc, regStatus.UserID, status.RootNodeID)
			return
		}

		err = svc.UpdateStatus(
			regStatus.UserID,
			false,
			status.StorageProvisioned,
			newTries,
			status.StorageProvisionTries,
		)
		if err != nil {
			log.Printf("error updating retry count for user ID %d: %v\n", regStatus.UserID, err)
			return
		}
	}
}

func updateStorageProvisioningStatus(ctx context.Context, nc *nats.Conn, svc *DriveService) func(*nats.Msg) {
	return func(msg *nats.Msg) {
		var provStatus StorageProvisionStatus
		if err := json.Unmarshal(msg.Data, &provStatus); err != nil {
			log.Println("failed to parse storage provision status JSON")
			return
		}

		status, err := svc.GetStatusUsingNodeID(provStatus.RootNodeID)
		if err != nil {
			log.Printf("error fetching drive status for node ID %s: %v\n", provStatus.RootNodeID, err)
			return
		}

		// Storage provisioning succeeded
		if provStatus.Success {
			err = svc.UpdateStatus(
				status.UserID,
				status.Authorized,
				true,
				status.AuthorizationTries,
				status.StorageProvisionTries,
			)
			if err != nil {
				log.Printf("failed to mark storage provisioned for user ID %d: %v\n", status.UserID, err)
				return
			}
			log.Printf("successfully updated storage provision status to true for user ID %d\n", status.UserID)
			return
		}

		// Storage provisioning failed: increment attempts or route to dead-letter queue
		const maxRetries uint8 = 100
		newTries := status.StorageProvisionTries + 1

		if newTries >= maxRetries {
			log.Printf("storage provisioning retry limit reached for user ID: %d (node ID: %s)\n", status.UserID, status.RootNodeID)
			deadLetterQueue(nc, status.UserID, status.RootNodeID)
			return
		}

		err = svc.UpdateStatus(
			status.UserID,
			status.Authorized,
			false,
			status.AuthorizationTries,
			newTries,
		)
		if err != nil {
			log.Printf("error updating storage retry count for user ID %d: %v\n", status.UserID, err)
			return
		}
	}
}
