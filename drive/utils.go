package drive

import (
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

func provisionStorage(nc *nats.Conn, RootNodeID string) error {
	storageReq := &StorageRequest{
		RootNodeID: RootNodeID,
	}
	msg, err := json.Marshal(storageReq)
	if err != nil {
		log.Println("error marshaling msg to json for storage svc")
		return err
	}
	err = nc.Publish("storage.new.root", msg)
	if err != nil {
		log.Println("req to storage svc failed!")
		return err
	}
	return nil
}

func grantAuthorizationPrevileges(nc *nats.Conn, UserID uint64, RootNodeID string, Previlege string) error {
	authReq := &AuthorizationGrantRequest{
		UserID:     UserID,
		RootNodeID: RootNodeID,
		Previlege:  Previlege,
	}
	msg, err := json.Marshal(authReq)
	if err != nil {
		log.Println("error marshaling msg to json for authorization svc")
		return err
	}
	err = nc.Publish("authorization.new.user", msg)
	if err != nil {
		log.Println("req to authorization svc failed!")
		return err
	}
	return nil
}

func deadLetterQueue(nc *nats.Conn, UserID uint64, RootNodeID string) {
	var dlqReq *DeadLetterQueueRequest = &DeadLetterQueueRequest{
		Status:     false,
		UserID:     UserID,
		RootNodeID: RootNodeID,
	}
	msg, err := json.Marshal(*dlqReq)
	if err != nil {
		log.Println("error marshaling registration status to json")
	}
	nc.Publish("dlq.user.status", msg)
}
