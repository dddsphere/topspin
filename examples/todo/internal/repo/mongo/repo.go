// package mongo provides a Mongo based implementation of UserRepo interface
package mongo

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/dddsphere/topspin"
	db "github.com/dddsphere/topspin/db/mongo"
	"github.com/dddsphere/topspin/examples/todo/internal/config"
)

type (
	Repo struct {
		*topspin.SimpleWorker
		config     *config.Config
		conn       *db.Client
		collection string
	}
)

func NewRepo(name string, conn *db.Client, collection string, cfg *config.Config, log topspin.Logger) *Repo {
	return &Repo{
		SimpleWorker: topspin.NewWorker(name, log),
		config:       cfg,
		conn:         conn,
		collection:   collection,
	}
}

func (r *Repo) Conn() (c *db.Client) {
	return r.conn
}

func (r *Repo) Client() (c *mongo.Client, err error) {
	if r.conn.Client == nil {
		return c, errors.New("no MongoDB client")
	}

	return r.conn.Client, nil
}

func (r *Repo) Session() (s mongo.Session, err error) {
	c, err := r.Client()
	if err != nil {
		return s, err
	}

	return c.StartSession()
}

func (r *Repo) Collection() (coll *mongo.Collection, err error) {
	c, err := r.Client()
	if err != nil {
		return coll, err
	}

	return c.Database(r.conn.Db()).Collection(r.collection), nil
}
