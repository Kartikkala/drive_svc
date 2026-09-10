package transport

type StorageRequest struct {
	RootNodeID string `json:"root_node_id"`
}

type AuthorizationGrantRequest struct {
	RootNodeID string `json:"root_node_id"`
	UserID     uint64 `json:"user_id"`
	Previlege  string `json:"previlege"`
}

type UserInfo struct {
	RootNodeID string `json:"root_node_id"`
	Found      bool   `json:"found"`
	Error      string `json:"error"`
}

type StorageProvisionStatus struct {
	Success    bool   `json:"success"`
	RootNodeID string `json:"root_node_id"`
}

type AuthorizationStatus struct {
	UserID  uint64 `json:"user_id"`
	Success bool   `json:"success"`
}

type DeadLetterQueueRequest struct {
	UserID     uint64 `json:"user_id"`
	Status     bool   `json:"status"`
	RootNodeID string `json:"root_node_id"`
}

type User struct {
	UserID uint64 `json:"user_id"`
}
