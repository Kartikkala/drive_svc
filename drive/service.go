package drive

import (
	"uuid"

	"gorm.io/gorm"
)

func NewDriveService(DB *gorm.DB) *DriveService {
	DB.AutoMigrate(&Drive{})
	return &DriveService{
		db: DB,
	}
}

func (svc *DriveService) NewDrive(UserID uint64) (uuid.UUID, error) {
	rootNodeID := uuid.New()
	err := svc.db.Create(&Drive{
		UserID:     UserID,
		RootNodeID: rootNodeID.String(),
	}).Error

	if err != nil {
		return uuid.Nil(), err
	}
	return rootNodeID, nil
}

func (svc *DriveService) UpdateStatus(
	userId uint64,
	authorized bool,
	storageProvisioned bool,
	authorizationTries uint8,
	storageProvisionTries uint8,
) error {
	result := svc.db.Model(&Drive{}).
		Where("user_id = ?", userId).
		Select(
			"Authorized",
			"StorageProvisioned",
			"AuthorizationTries",
			"StorageProvisionTries",
		).
		Updates(Drive{
			Authorized:            authorized,
			StorageProvisioned:    storageProvisioned,
			AuthorizationTries:    authorizationTries,
			StorageProvisionTries: storageProvisionTries,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (svc *DriveService) GetStatus(
	userID uint64,
) (*Drive, error) {
	var status Drive
	err := svc.db.Where("user_id = ?", userID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (svc *DriveService) GetStatusUsingNodeID(
	nodeID string,
) (*Drive, error) {
	var status Drive
	err := svc.db.Where("root_node_id = ?", nodeID).First(&status).Error
	if err != nil {
		return nil, err
	}
	return &status, nil
}
