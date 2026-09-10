package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/Kartikkala/drive_svc/shared"
	"github.com/nats-io/nats.go"
)

func NewDriveEventHandler(svc shared.IDrive) *SubscribeEventHandlers {
	return &SubscribeEventHandlers{
		svc: svc,
	}
}

func (h *SubscribeEventHandlers) HandleUserRegistration(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user User
	if err := json.Unmarshal(msg.Data, &user); err != nil {
		log.Printf("failed to parse authentication service payload: %v\n", err)
		return
	}

	rootNodeId, err := h.svc.CreateDrive(ctx, user.UserID)
	if err != nil {
		log.Printf("failed to create drive for user %d: %v\n", user.UserID, err)
		return
	}
	err = msg.Ack()
	if err != nil {
		log.Println("error while sending ack to NATS", err)
	}
	log.Printf("successfully initialized drive (root: %s) for user %d", rootNodeId, user.UserID)
}

func (h *SubscribeEventHandlers) HandleAuthorizationStatusUpdate(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var regStatus AuthorizationStatus
	if err := json.Unmarshal(msg.Data, &regStatus); err != nil {
		log.Printf("failed to parse registration status JSON: %v\n", err)
		return
	}

	if regStatus.Success {
		if err := h.svc.UpdateAuthorizationStatusToTrue(ctx, regStatus.UserID); err != nil {
			log.Printf("failed to change authorization status of user %v: %v\n", regStatus.UserID, err)
			return
		}
	}

	if err := msg.Ack(); err != nil {
		log.Printf("error while sending ack to NATS for user %d: %v\n", regStatus.UserID, err)
	}
}

func (h *SubscribeEventHandlers) HandleStorageProvisioningStatusUpdate(msg *nats.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var storageStatus StorageProvisionStatus
	if err := json.Unmarshal(msg.Data, &storageStatus); err != nil {
		log.Printf("failed to parse storage provision status JSON: %v\n", err)
		return
	}

	if storageStatus.Success {
		if err := h.svc.UpdateStorageProvisioningStatusToTrue(ctx, storageStatus.RootNodeID); err != nil {
			log.Printf("failed to change storage provisioning status of root nodeID: %s\n, %v\n", storageStatus.RootNodeID, err)
			return
		}
	}

	if err := msg.Ack(); err != nil {
		log.Printf("error while sending ack to NATS for root node ID %s: %v\n", storageStatus.RootNodeID, err)
	}
}

func (h *SubscribeEventHandlers) HandleUserInfoQuery(msg *nats.Msg) {
	var user User

	if err := json.Unmarshal(msg.Data, &user); err != nil {
		log.Printf("failed to parse user JSON: %v\n", err)
		return
	}

	root, err := h.svc.GetRootNodeId(context.TODO(), user.UserID)
	var userInfo *UserInfo = &UserInfo{
		RootNodeID: root,
	}
	if err != nil {
		if errors.Is(err, shared.ErrDriveNotFound) {
			userInfo.Found = false
		}
		userInfo.Error = err.Error()
		log.Printf("error while fetching user's root node ID: %v\n", err)
	}

	res, err := json.Marshal(userInfo)
	if err != nil {
		log.Printf("error while marshaling user's info %v\n", err)
		return
	}
	err = msg.Respond(res)
	if err != nil {
		log.Printf("error while responding to NATS with user info: %v\n", err)
	}
}
