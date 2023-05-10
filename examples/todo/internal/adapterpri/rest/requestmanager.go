package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs/bus/nats"
)

type (
	RequestManager struct {
		*topspin.SimpleWorker
		cqrs *cqrs.Manager
		bus  *nats.BusManager
	}
)

func NewRequestManager(cqrs *cqrs.Manager, bus *nats.BusManager, log topspin.Logger) (rm *RequestManager) {
	return &RequestManager{
		SimpleWorker: topspin.NewWorker("request-manager", log),
		cqrs:         cqrs,
		bus:          bus,
	}

}

// Dispatch is a WIP:
// Eventually the commands will return an HTTP 202
// including headers with relevant metadata (i.e.: request ID).
func (rm *RequestManager) Dispatch(w http.ResponseWriter, r *http.Request, commandName string) {
	reqID := genReqID(r)

	// TODO: Command should do a basic validation of the payload
	// before enqueuing.
	payload, err := body(r)
	if err != nil {
		err := fmt.Errorf("send command error: %w", err)
		rm.Error(err, w)
	}

	switch commandName {
	case "create-list":
		err = rm.bus.SendCommand("create-list", payload, reqID)

	//case "update-list":
	//	rm.bus.SendCommand("update-list-command", payload, reqID)

	//case "delete-list":
	//	rm.bus.SendCommand("delete-list-command", payload, reqID)

	default:
		err := fmt.Errorf("command '%s' not found", commandName)
		rm.Error(err, w)
	}

	if err != nil {
		rm.Log().Errorf("Dispatch command error: %w", err)
	}

	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte("202 - Temporary OK message"))
	if err != nil {
		rm.Log().Errorf("Dispatch command error: %w", err)
	}
}

func (rm *RequestManager) Error(err error, w http.ResponseWriter) {
	rm.Log().Error(err.Error())
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

// Helpers
func body(r *http.Request) (body []byte, err error) {
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return body, err
	}

	return body, nil
}

func genReqID(r *http.Request) (id string) {
	id = r.Header.Get("X-Request-ID")
	if id == "" {
		return uuid.New().String()
	}

	return id
}
