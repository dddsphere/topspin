// Package service provides application resources for managing todo lists.
package service

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/config"
	"github.com/dddsphere/topspin/examples/todo/core"
	"github.com/dddsphere/topspin/examples/todo/repo"
)

type (
	Todo struct {
		*topspin.SimpleWorker
		config      *config.Config
		cqrs        *topspin.CQRSManager
		repoRead    repo.ListRead
		repoWrite   repo.ListWrite
		listService *core.List
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
	return t.repoWrite.Create(ctx, core.List{
		ID:          uid,
		UserID:      uuid.New(),
		Slug:        slug,
		TenantID:    uuid.New(),
		OrgID:       uuid.New(),
		OwnerID:     uuid.New(),
		Name:        "list name",
		Description: "list description",
		Items:       []core.Item{},
	})
}
