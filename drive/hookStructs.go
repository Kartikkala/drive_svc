package drive

type StorageRequest struct {
	RootNodeID string `json:"root_node_id"`
}

type AuthorizationGrantRequest struct {
	RootNodeID string `json:"root_node_id"`
	UserID     uint64 `json:"user_id"`
	Previlege  string `json:"previlege"`
}