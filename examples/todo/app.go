package todo

import (
	"fmt"
	"sync"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/internal/adapterpri/rest"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs/bus/nats"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs/command"
	"github.com/dddsphere/topspin/examples/todo/internal/ports/openapi"
	"github.com/dddsphere/topspin/examples/todo/internal/service"
)

type (
	// App description
	App struct {
		*topspin.App

		// Service
		TodoService *service.Todo

		// CQRS
		CQRS *cqrs.Manager

		// Bus
		// NATS
		NATS *nats.Client
		Bus  *nats.BusManager

		RESTServer *rest.Server
		//WebServer     *web.Server
		//GRPCServer    *grpc.Server
	}
)

type (
	Command struct {
		name string
	}
)

const (
	cfgLoggingLevel = "logging.level"
	cfgJSONAPIPort  = "config.server.json-api-port"
)

func NewApp(name, version string, log topspin.Logger) *App {
	return &App{
		App:  topspin.NewApp(name, version, log),
		CQRS: cqrs.NewManager(log),
	}
}

func (app *App) SetLogLevel(level string) {
	app.Log().SetLevel(app.Cfg().GetString(cfgLoggingLevel))
}

// Init app
func (app *App) Init() (err error) {
	// Commands
	app.initCommands()

	// Router
	if app.RESTServer != nil {
		rm := rest.NewRequestManager(app.CQRS, app.Bus, app.Log())
		h := openapi.Handler(rm)

		app.RESTServer.InitRESTRouter(h)
	}

	return nil
}

// Start app
func (app *App) Start() error {
	var errREST error
	var errBus error
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		port := app.Cfg().GetInt(cfgJSONAPIPort)
		errREST = app.RESTServer.Start(port)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		errBus = app.Bus.Start()
		wg.Done()
	}()

	wg.Wait()

	if errREST != nil {
		return fmt.Errorf("cannot start server: %w", errREST)
	}

	if errBus != nil {
		return fmt.Errorf("cannot start server: %w", errBus)
	}

	return fmt.Errorf("cannot start server:\n\t%s\n\t%s\n", errREST.Error(), errBus.Error())
}

func (app *App) InitAndStart() error {
	err := app.Init()
	if err != nil {
		return fmt.Errorf("%s init error: %w", app.Name(), err)
	}

	err = app.Start()
	if err != nil {
		return fmt.Errorf("%s start error: %w", app.Name(), err)
	}

	return nil
}

func (app *App) Stop() {
	// TODO: Gracefully stop the app
}

func (app *App) initCommands() {
	log := app.Log()
	//app.AddCommand(&topspin.SampleCommand) // TODO: Remove
	app.AddCommand(command.NewCreateListCommand(app.TodoService, log))
	//app.AddCommand(command.NewAddItemCommand(app.TodoService))
	//app.AddCommand(command.NewGetItemCommand(app.TodoService))
	//app.AddCommand(command.NewUpdateItemCommand(app.TodoService))
	//app.AddCommand(command.NewDeleteItemCommand(app.TodoService))
	//app.AddCommand(command.NewDeleteListCommand(app.TodoService))
}

func (app *App) AddCommand(command topspin.Command) {
	app.CQRS.AddCommand(command)
}

func (app *App) AddQuery(query topspin.Query) {
	app.CQRS.AddQuery(query)
}

func (app *App) LoadConfig() (cfg *topspin.Config, err error) {
	return app.Cfg().Load()
}

func (app *App) EnableBus() {
	app.NATS = nats.NewClient("nats-client", app.Cfg(), app.Log())
	app.CQRS = cqrs.NewManager(app.Log())
	app.Bus = nats.NewBusManager("nats-bus", app.Cfg(), app.NATS, app.CQRS, app.Log())
}
