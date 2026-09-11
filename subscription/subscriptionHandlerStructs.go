package subscription

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

type User struct {
	UserID uint64 `json:"user_id"`
}
