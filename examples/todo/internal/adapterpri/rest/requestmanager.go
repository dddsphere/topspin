package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/google/uuid"

	"github.com/dddsphere/topspin"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs/bus/nats"
	"github.com/dddsphere/topspin/examples/todo/internal/cqrs/command"
)

type (
	RequestManager struct {
		*topspin.SimpleWorker
		cqrs *topspin.CQRSManager
		bus  *nats.BusManager
	}
)

func NewRequestManager(cqrs *topspin.CQRSManager, bus *nats.BusManager, log topspin.Logger) (rm *RequestManager) {
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

	// WIP: Sending gob data for now
	// TODO: Command should do a basic validation of the payload
	// before enqueuing.
	payload, err := body(r)
	if err != nil {
		err := fmt.Errorf("send command error: %w", err)
		rm.Error(err, w)
	}

	rm.Log().Debug("1")
	switch commandName {
	case "create-list":
		rm.Log().Debug("2")
		err = rm.bus.SendCommand("create-list-command", payload, reqID)

	//case "update-list":
	//	rm.bus.SendCommand("update-list-command", payload, reqID)

	//case "delete-list":
	//	rm.bus.SendCommand("delete-list-command", payload, reqID)

	default:
		rm.Log().Debug("3")
		err := fmt.Errorf("command '%s' not found", commandName)
		rm.Error(err, w)
	}

	if err != nil {
		rm.Log().Debug("4")
		rm.Log().Errorf("Dispatch command error: %w", err)
	}

	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write([]byte("202 - Temporary OK message"))
	if err != nil {
		rm.Log().Debug("5")
		rm.Log().Errorf("Dispatch command error: %w", err)
	}
}

// CreateList is a WIP: Only for now the request manager is responsible for
// executing commands received from the bus.
func (rm *RequestManager) CreateList(w http.ResponseWriter, r *http.Request) {
	name := "create-list"

	cmd, ok := rm.cqrs.FindCommand(name)
	if !ok {
		err := fmt.Errorf("command '%s' not found", name)
		rm.Error(err, w)
		return
	}

	switch cmd := cmd.(type) {
	case *command.CreateListCommand:
		data, err := ToCreateListCommandData(r)
		if err != nil {
			err := fmt.Errorf("wrong '%s' data: %+v", cmd.Name(), data)
			rm.Error(err, w)
			return
		}

		err = cmd.HandleFunc()(r.Context(), data)
		if err != nil {
			err := fmt.Errorf("error: %s", err.Error())
			rm.Error(err, w)
			return
		}

	default:
		rm.Log().Errorf("wrong command: %+v", cmd)
	}
}

// ToCreateListCommandData command
func ToCreateListCommandData(r *http.Request) (cmdData command.CreateListCommandData, err error) {
	err = json.NewDecoder(r.Body).Decode(&cmdData)
	if err != nil {
		return cmdData, err
	}

	return cmdData, nil
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
