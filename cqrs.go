package topspin

import (
	"context"
)

type (
	CQRSManager struct {
		Commands CommandSet
		Queries  QuerySet
	}

	CommandSet map[string]Command
	QuerySet   map[string]Query

	Command interface {
		Name() string
		HandleFunc() func(ctc context.Context, data interface{}) (err error)
	}

	Query interface {
		Name() string
	}

	BaseCommand struct {
		name string
	}

	BaseQuery struct {
		name string
	}
)

var (
	SampleCommand = BaseCommand{name: "sample-command"}
	SampleQuery   = BaseQuery{name: "sample-query"}
)

func NewCQRSManager() *CQRSManager {
	return &CQRSManager{
		Commands: CommandSet{},
		Queries:  QuerySet{},
	}
}

func NewBaseCommand(name string) *BaseCommand {
	return &BaseCommand{
		name: name,
	}
}

func (cqrs CQRSManager) AddCommand(cmd Command) {
	cqrs.Commands[cmd.Name()] = cmd
}

func (cqrs CQRSManager) FindCommand(name string) (cmd Command, ok bool) {
	cmd, ok = cqrs.Commands[name]
	return cmd, ok
}

func (cqrs CQRSManager) AddQuery(qry Query) {
	cqrs.Queries[qry.Name()] = qry
}

func (cqrs CQRSManager) FindQuery(name string) (qry Query, ok bool) {
	qry, ok = cqrs.Queries[name]
	return qry, ok
}

func (bc *BaseCommand) Name() string {
	return bc.name
}

func (bc *BaseCommand) HandleFunc() (f func(ctx context.Context, data interface{}) (err error)) {
	return f
}
