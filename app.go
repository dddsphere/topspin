package topspin

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"
	"time"
)

type (
	// App description
	App struct {
		name    string
		version string

		log Logger

		cancel context.CancelFunc
	}
)

func NewApp(name, version string, log Logger) *App {
	name = GenName(name, "app")

	return &App{
		name:    name,
		version: version,
		log:     log,
	}
}

func (app *App) Name() string {
	return app.name
}

func (app *App) Version() string {
	return app.version
}

func (app *App) Log() Logger {
	return app.log
}

func GenName(name, defName string) string {
	if strings.Trim(name, " ") == "" {
		return fmt.Sprintf("%s-%s", defName, nameSufix())
	}
	return name
}

func nameSufix() string {
	digest := hash(time.Now().String())
	return digest[len(digest)-8:]
}

func hash(s string) string {
	h := fnv.New32a()
	h.Write([]byte(s))
	return fmt.Sprintf("%d", h.Sum32())
}
