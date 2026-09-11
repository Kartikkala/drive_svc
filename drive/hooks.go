package drive

import (
	"context"
	"encoding/json"
	"log"
	"uuid"

	"github.com/nats-io/nats.go"
)

func newHooksForDrive(nc *nats.Conn) *HooksForDrive {
	return &HooksForDrive{
		nc: nc,
	}
}

func (hook *HooksForDrive) provisionStorage(ctx context.Context, UserID uint64, RootNodeID uuid.UUID) (bool, error) {
	storageReq := &StorageRequest{
		RootNodeID: RootNodeID.String(),
	}
	msg, err := json.Marshal(storageReq)
	if err != nil {
		log.Println("error marshaling msg to json for storage svc")
		return false, err
	}
	err = hook.nc.Publish("storage.new.root", msg)
	if err != nil {
		log.Println("req to storage svc failed!")
		return false, err
	}
	return false, nil
}

func (hook *HooksForDrive) grantAuthorizationPrevileges(ctx context.Context, UserID uint64, RootNodeID uuid.UUID) (bool, error) {
	authReq := &AuthorizationGrantRequest{
		UserID:     UserID,
		RootNodeID: RootNodeID.String(),
		Previlege:  "owner",
	}
	msg, err := json.Marshal(authReq)
	if err != nil {
		log.Println("error marshaling msg to json for authorization svc")
		return false, err
	}
	err = hook.nc.Publish("authorization.new.user", msg)
	if err != nil {
		log.Println("req to authorization svc failed!")
		return false, err
	}
	return false, nil
}

func RegisterAllDriveHooks(nc *nats.Conn, svcWithHooks *DriveServiceWithHooks) {
	hooks := newHooksForDrive(nc)
	svcWithHooks.RegisterCreateDriveAfterHook(hooks.provisionStorage)
}
