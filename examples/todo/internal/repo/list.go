package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/dddsphere/topspin/examples/todo/internal/domain"
)

type (
	// ListRepo interface
	ListWrite interface {
		Create(ctx context.Context, list domain.List) error
		Update(ctx context.Context, list domain.List) (err error)
		Delete(ctx context.Context, uid uuid.UUID) error
		DeleteBySlug(ctx context.Context, slug string) error
		DeleteByName(ctx context.Context, listname string) error
	}

	ListRead interface {
		GetAll(ctx context.Context) (lists []domain.List, err error)
		Get(ctx context.Context, uid uuid.UUID) (list domain.List, err error)
		GetBySlug(ctx context.Context, slug string) (list domain.List, err error)
		GetByName(ctx context.Context, listname string) (domain.List, error)
	}
)
