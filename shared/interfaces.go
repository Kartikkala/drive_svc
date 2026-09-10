package shared

import (
	"context"
	"uuid"
)

type IDrive interface {
	CreateDrive(ctx context.Context, UserID uint64) (uuid.UUID, error)
	UpdateAuthorizationStatusToTrue(ctx context.Context, UserID uint64) error
	UpdateStorageProvisioningStatusToTrue(ctx context.Context, RootNodeID string) error
	GetRootNodeId(ctx context.Context, UserID uint64) (string, error)
}