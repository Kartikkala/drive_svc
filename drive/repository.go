package drive

import (
	"context"
	"uuid"

	"gorm.io/gorm"
)

func NewDriveRepository(DB *gorm.DB) *DriveRepository {
	DB.AutoMigrate(&Drive{})
	return &DriveRepository{
		db: DB,
	}
}

func (repository *DriveRepository) Create(ctx context.Context, UserID uint64, rootNodeID uuid.UUID) error {
	err := repository.db.WithContext(ctx).Create(&Drive{
		UserID:     UserID,
		RootNodeID: rootNodeID.String(),
	}).Error

	if err != nil {
		return err
	}
	return nil
}

func (repository *DriveRepository) Update(
	ctx context.Context,
	UserID uint64,
	dto UpdateDriveDTO,
) error {
	updates := make(map[string]any)

	if dto.Authorized != nil {
		updates["authorized"] = *dto.Authorized
	}

	if dto.StorageProvisioned != nil {
		updates["storage_provisioned"] = *dto.StorageProvisioned
	}

	if len(updates) == 0 {
		return nil
	}

	result := repository.db.WithContext(ctx).
		Model(&Drive{}).
		Where("user_id = ?", UserID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repository *DriveRepository) UpdateByNodeID(
	ctx context.Context,
	RootNodeID string,
	dto UpdateDriveDTO,
) error {
	updates := make(map[string]any)

	if dto.Authorized != nil {
		updates["authorized"] = *dto.Authorized
	}

	if dto.StorageProvisioned != nil {
		updates["storage_provisioned"] = *dto.StorageProvisioned
	}

	if len(updates) == 0 {
		return nil
	}

	result := repository.db.WithContext(ctx).
		Model(&Drive{}).
		Where("root_node_id = ?", RootNodeID).
		Updates(updates)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (repository *DriveRepository) FindFirstDriveByUserId(
	ctx context.Context,
	userID uint64,
) (*Drive, error) {
	var status Drive
	err := repository.db.WithContext(ctx).Where("user_id = ?", userID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (repository *DriveRepository) FindFirstDriveByDriveId(
	ctx context.Context,
	nodeID string,
) (*Drive, error) {
	var status Drive
	err := repository.db.WithContext(ctx).Where("root_node_id = ?", nodeID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}
