package drive

import (
	"context"
	"uuid"

	"github.com/Kartikkala/drive_svc/shared"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type DriveRepository struct {
	db *gorm.DB
}

type DriveService struct {
	repository *DriveRepository
}

type CreateDriveAfterHook func(ctx context.Context, UserID uint64, RootNodeID uuid.UUID) (bool, error)

type DriveServiceWithHooks struct {
	shared.IDrive
	createDriveHooksAfter []CreateDriveAfterHook
}

type UpdateDriveDTO struct {
	Authorized         *bool
	StorageProvisioned *bool
}

type HooksForDrive struct {
	nc *nats.Conn
}
