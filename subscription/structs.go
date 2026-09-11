package subscription

import (
	"github.com/Kartikkala/drive_svc/shared"
	"github.com/nats-io/nats.go"
)

type SubscribeEventHandlers struct {
	svc shared.IDrive
}

type PublishEventHandlers struct {
	nc *nats.Conn
}
