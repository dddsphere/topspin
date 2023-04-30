package core

import (
	"github.com/google/uuid"
)

type (
	List struct {
		ID          uuid.UUID
		UserID      uuid.UUID
		Slug        string
		TenantID    uuid.UUID
		OrgID       uuid.UUID
		OwnerID     uuid.UUID
		Name        string
		Description string
		Items       []Item
	}
)
