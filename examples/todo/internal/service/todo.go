// Package service provides application resources for managing todo lists.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/dddsphere/topspin"
	core2 "github.com/dddsphere/topspin/examples/todo/internal/domain"
	"github.com/dddsphere/topspin/examples/todo/internal/repo"
	"github.com/dddsphere/topspin/examples/todo/pkg/config"
)

type (
	Todo struct {
		*topspin.SimpleWorker
		config      *config.Config
		cqrs        *topspin.CQRSManager
		repoRead    repo.ListRead
		repoWrite   repo.ListWrite
		listService *core2.List
	}
)

func NewTodo(name string, rr repo.ListRead, rw repo.ListWrite, cfg *config.Config, log topspin.Logger) (Todo, error) {
	var svc Todo

	if rr == nil {
		return svc, errors.New("no read repo")
	}

	if rw == nil {
		return svc, errors.New("no write repo")
	}

	svc = Todo{
		SimpleWorker: topspin.NewWorker(name, log),
		config:       cfg,
		repoRead:     rr,
		repoWrite:    rw,
	}

	return svc, nil
}

func (t *Todo) CreateList(ctx context.Context, name, description string) error {
	t.Log().Infof("CreateList name: '%s', description: '%s'", name, description)

	uid := uuid.New()
	slug := strings.Split(uid.String(), "-")[4]

	// WIP: Filling empty fields with fake data
	return t.repoWrite.Create(ctx, core2.List{
		ID:          uid,
		UserID:      uuid.New(),
		Slug:        slug,
		TenantID:    uuid.New(),
		OrgID:       uuid.New(),
		OwnerID:     uuid.New(),
		Name:        "list name",
		Description: "list description",
		Items:       []core2.Item{},
	})
}
