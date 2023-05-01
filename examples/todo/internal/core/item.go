package core

import (
	"time"

	"github.com/google/uuid"
)

type (
	Item struct {
		ID          uuid.UUID
		Name        string
		Description string
		Meta        []Meta
	}
)

type (
	Meta struct {
		state      string
		tags       []string
		ownwership []uuid.UUID
		assigned   []uuid.UUID
		mentioned  []uuid.UUID
		audit      []Record
	}

	Record struct {
		timestamp   time.Time
		action      string
		description string
		createdBy   uuid.UUID
		updatedBy   uuid.UUID
	}
)
