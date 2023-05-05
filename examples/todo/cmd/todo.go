package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dddsphere/topspin"
	db "github.com/dddsphere/topspin/db/mongo"
	"github.com/dddsphere/topspin/examples/todo"
	"github.com/dddsphere/topspin/examples/todo/internal/adapterpri/rest"
	nats2 "github.com/dddsphere/topspin/examples/todo/internal/cqrs/bus/nats"
	"github.com/dddsphere/topspin/examples/todo/internal/repo/mongo"
	"github.com/dddsphere/topspin/examples/todo/internal/service"
)

const (
	name    = "todo"
	version = "v0.0.1"
)

var (
	a   *todo.App
	log topspin.Logger
)

func main() {
	log = topspin.NewLogger("debug", true)

	// App
	a := todo.NewApp(name, version, log)
	cfg := a.LoadConfig()

	// WIP: Verifying new configuration mechanism
	config := topspin.NewConfig(name, log)
	_, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	config.List()

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	initExitMonitor(ctx, cancel)

	// Databases
	mgo := db.NewMongoClient("mongo-client", db.Config{
		Host:       cfg.Mongo.Host,
		Port:       cfg.Mongo.Port,
		User:       cfg.Mongo.User,
		Pass:       cfg.Mongo.Pass,
		Database:   cfg.Mongo.Database,
		MaxRetries: cfg.Mongo.MaxRetriesUInt64(),
	}, log)

	// Repos
	lrr := mongo.NewListRead("list-read-repo", mgo, &cfg, log)

	lwr := mongo.NewListWrite("list-write-repo", mgo, &cfg, log)

	// Services
	ts, err := service.NewTodo("todo-app-service", lrr, lwr, &cfg, log)
	if err != nil {
		exit(err)
	}

	a.TodoService = &ts

	// Server
	a.RESTServer = rest.NewServer("rest-server", &cfg, log)

	// Bus
	nc := nats2.NewClient("nats-client", &cfg, log)

	a.Bus = nats2.NewBusManager("nats-bus", &cfg, nc, log)

	// Init & Start
	err = a.InitAndStart()
	if err != nil {
		exit(err)
	}

	log.Errorf("%s stopped: %s (%s)", a.Name(), a.Version(), err)
}

func exit(err error) {
	log.Fatal(err)
}

func initExitMonitor(ctx context.Context, cancel context.CancelFunc) {
	go checkSigterm(cancel)
	go checkCancel(ctx)
}

func checkSigterm(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()
}

func checkCancel(ctx context.Context) {
	<-ctx.Done()
	a.Stop()
	os.Exit(1)
}
