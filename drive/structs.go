package drive

import "gorm.io/gorm"

type DriveRepository struct {
	db *gorm.DB
}

type DriveService struct {
	repository *DriveRepository
}

type UpdateDriveDTO struct {
	Authorized         *bool
	StorageProvisioned *bool
}