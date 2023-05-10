package command

import (
	"context"
	"errors"
	"fmt"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/internal/service"
)

type (
	CreateListCommandData struct {
		UserID      string
		Name        string
		Description string
	}

	CreateListCommand struct {
		*topspin.SimpleWorker
		*topspin.BaseCommand
		todoService *service.Todo
	}
)

func NewCreateListCommand(todoService *service.Todo, log topspin.Logger) *CreateListCommand {
	if todoService == nil {
		panic("nil Todo service")
	}

	return &CreateListCommand{
		SimpleWorker: topspin.NewWorker("create-list-command", log),
		BaseCommand:  topspin.NewBaseCommand("create-list"),
		todoService:  todoService,
	}
}

func (cmd CreateListCommand) Name() string {
	return cmd.BaseCommand.Name()
}

func (c *CreateListCommand) HandleFunc() (f func(ctx context.Context, data interface{}) error) {
	return c.handle
}

func (c *CreateListCommand) handle(ctx context.Context, data interface{}) (err error) {
	switch d := data.(type) {
	case CreateListCommandData:
		c.Log().Debugf("Processing %s with %+v", c.Name(), d)

		err = c.todoService.CreateList(ctx, d.Name, d.Description)
		if err != nil {
			return fmt.Errorf("%s handle error: %w", c.Name(), err)
		}

	default:
		return errors.New("create list wrong command data")
	}

	return nil
}
