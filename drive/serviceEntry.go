package drive

import (
	"context"
	"log"
	"uuid"

	"github.com/Kartikkala/drive_svc/shared"
)

func NewDriveServiceWithHooks(inner shared.IDrive) *DriveServiceWithHooks {
	return &DriveServiceWithHooks{
		IDrive: inner,
	}
}

func (h *DriveServiceWithHooks) CreateDrive(ctx context.Context, UserID uint64) (uuid.UUID, error) {
	rootNodeID, err := h.IDrive.CreateDrive(ctx, UserID)
	if err != nil {
		return rootNodeID, err
	}
	for _, hook := range h.createDriveHooksAfter {
		if stop, err := hook(ctx, UserID, rootNodeID); err != nil {
			log.Println("CreateDrive after-hook error:", err)
			if stop {
				break
			}
		}
	}
	return rootNodeID, nil
}

func (h *DriveServiceWithHooks) RegisterCreateDriveAfterHook(f CreateDriveAfterHook) {
	h.createDriveHooksAfter = append(h.createDriveHooksAfter, f)
}
