package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/cenkalti/backoff"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	retry struct {
		Client *mongo.Client
		Error  error
	}
)

func (c *Client) retryConnection() (r chan retry) {
	r = make(chan retry)
	bo := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), c.config.MaxRetries)

	go func() {
		defer close(r)

		url := c.URL()

		for i := 0; i <= int(c.config.MaxRetries); i++ {
			c.Log().Infof("dialing to mongo at %s", url)

			opts := options.Client().ApplyURI(url)

			client, err := mongo.Connect(context.TODO(), opts)
			if err != nil {
				c.Log().Infof("mongo connection error: %v\n", err)
			}

			err = client.Ping(context.TODO(), nil)

			if err != nil {
				c.Log().Infof("mongo ping error: %v", err)

				// Backoff
				next := bo.NextBackOff()
				if next == backoff.Stop {
					err := errors.New("max number of Mongo connection attempts reached")

					r <- retry{nil, err}

					bo.Reset()
					return
				}

				c.Log().Infof("connection attempt to Mongo failed: %v", err)
				c.Log().Infof("retrying connection to Mongo in %s seconds", next.String())

				time.Sleep(next)
				continue
			}

			r <- retry{client, nil}
		}
	}()

	return r
}
