package drive

import "gorm.io/gorm"

type NewUser struct {
	UserID uint64 `json:"user_id"`
}

type DriveService struct {
	db *gorm.DB
}

type StorageRequest struct {
	RootNodeID string `json:"root_node_id"`
}

type AuthorizationGrantRequest struct {
	RootNodeID string `json:"root_node_id"`
	UserID     uint64 `json:"user_id"`
	Previlege  string `json:"previlege"`
}

type StorageProvisionStatus struct {
	Success    bool   `json:"success"`
	RootNodeID string `json:"root_node_id"`
}

type AuthorizationStatus struct {
	UserID     uint64 `json:"user_id"`
	Success     bool   `json:"status"`
	RootNodeID string `json:"root_node_id"`
}

type DeadLetterQueueRequest struct {
	UserID     uint64 `json:"user_id"`
	Status     bool   `json:"status"`
	RootNodeID string `json:"root_node_id"`
}
