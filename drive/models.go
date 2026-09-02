package drive

type Drive struct {
	UserID                uint64 `gorm:"column:user_id;primaryKey"`
	RootNodeID            string `gorm:"column:root_node_id;type:varchar(255);uniqueIndex;not null"`
	Authorized            bool   `gorm:"column:authorized;default:false;not null"`
	StorageProvisioned    bool   `gorm:"column:storage_provisioned;default:false;not null"`
	AuthorizationTries    uint8  `gorm:"column:authorization_tries;default:0;not null"`
	StorageProvisionTries uint8  `gorm:"column:storage_provision_tries;default:0;not null"`
}
