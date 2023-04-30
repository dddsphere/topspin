package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo"
	"github.com/dddsphere/topspin/examples/todo/adapter/rest"

	"github.com/dddsphere/topspin/examples/todo/cqrs/bus/nats"
	"github.com/dddsphere/topspin/examples/todo/repo/mongo"
	"github.com/dddsphere/topspin/examples/todo/service"

	db "github.com/dddsphere/topspin/db/mongo"
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

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	initExitMonitor(ctx, cancel)

	// Database
	mgo := db.NewMongoClient("mongo-client", db.Config{
		Host:       cfg.Mongo.Host,
		Port:       cfg.Mongo.Port,
		User:       cfg.Mongo.User,
		Pass:       cfg.Mongo.Pass,
		Database:   cfg.Mongo.Database,
		MaxRetries: cfg.Mongo.MaxRetriesUInt64(),
	}, log)

	// Repo
	lrr := mongo.NewListRead("list-read-repo", mgo, &cfg, log)

	lwr := mongo.NewListWrite("list-write-repo", mgo, &cfg, log)

	// Service
	ts, err := service.NewTodo("todo-app-service", lrr, lwr, &cfg, log)
	if err != nil {
		exit(err)
	}

	a.TodoService = &ts

	// Server
	a.RESTServer = rest.NewServer("rest-server", &cfg, log)

	// Bus
	nc := nats.NewClient("nats-client", &cfg, log)

	a.Bus = nats.NewBusManager("nats-bus", &cfg, nc, log)

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
