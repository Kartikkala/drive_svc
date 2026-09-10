package drive

import (
	"context"
	"errors"
	"uuid"

	"github.com/Kartikkala/drive_svc/shared"
	"gorm.io/gorm"
)

func NewDriveService(repository *DriveRepository) *DriveService {
	return &DriveService{
		repository: repository,
	}
}

func (svc *DriveService) CreateDrive(ctx context.Context, UserID uint64) (uuid.UUID, error) {
	rootNodeID := uuid.New()
	err := svc.repository.Create(ctx, UserID, rootNodeID)
	if err != nil {
		return uuid.Nil(), err
	}
	return rootNodeID, nil
}

func (svc *DriveService) Update(ctx context.Context, UserID uint64, DTO UpdateDriveDTO) error {
	err := svc.repository.Update(ctx, UserID, DTO)
	if err != nil {
		return err
	}
	return nil
}

func (svc *DriveService) UpdateStorageProvisioningStatusToTrue(ctx context.Context, RootNodeID string) error {
	var update UpdateDriveDTO
	storageProvisioned := true
	update.StorageProvisioned = &storageProvisioned
	return svc.repository.UpdateByNodeID(ctx, RootNodeID, update)
}

func (svc *DriveService) UpdateAuthorizationStatusToTrue(ctx context.Context, UserID uint64) error {
	var update UpdateDriveDTO
	authorized := true
	update.Authorized = &authorized
	return svc.Update(ctx, UserID, update)
}

func (svc *DriveService) GetRootNodeId(ctx context.Context, UserID uint64) (string, error) {
	drive, err := svc.repository.FindFirstDriveByUserId(ctx, UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", shared.ErrDriveNotFound
		}
		return "", err
	}

	return drive.RootNodeID, err
}
