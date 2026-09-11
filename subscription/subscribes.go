package subscription

import (
	"github.com/nats-io/nats.go"
)

func AttachEvents(nc *nats.Conn, evHndlrs *SubscribeEventHandlers) {
	nc.Subscribe("drive.new.user", evHndlrs.HandleUserRegistration)
	nc.Subscribe("drive.update.authorization.status", evHndlrs.HandleAuthorizationStatusUpdate)
	nc.Subscribe("drive.update.storage.status", evHndlrs.HandleStorageProvisioningStatusUpdate)
	nc.Subscribe("drive.user.info", evHndlrs.HandleUserInfoQuery)
}
