package repo

import (
	"context"

	"github.com/google/uuid"

	"github.com/dddsphere/topspin/examples/todo/internal/core"
)

type (
	// ListRepo interface
	ListWrite interface {
		Create(ctx context.Context, list core.List) error
		Update(ctx context.Context, list core.List) (err error)
		Delete(ctx context.Context, uid uuid.UUID) error
		DeleteBySlug(ctx context.Context, slug string) error
		DeleteByName(ctx context.Context, listname string) error
	}

	ListRead interface {
		GetAll(ctx context.Context) (lists []core.List, err error)
		Get(ctx context.Context, uid uuid.UUID) (list core.List, err error)
		GetBySlug(ctx context.Context, slug string) (list core.List, err error)
		GetByName(ctx context.Context, listname string) (core.List, error)
	}
)
