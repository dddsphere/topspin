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

// Dispatch is a WIP: This can be improved,
// Maybe delegating the action to the command itself right here
// but removing first the concrete actions from the OpenAPI specifications.
// so the handlers don't have to adjust to the interface (response/request interface)
// and code and checks duplication is avoided.
func (rm *RequestManager) Dispatch(w http.ResponseWriter, r *http.Request, commandName string) {
	reqID := genReqID(r)

	// WIP: Sending gob data for now
	// TODO: Command should validate payload
	// before enqueuing.
	payload, err := body(r)
	if err != nil {
		err := fmt.Errorf("send command error: %w", err)
		rm.Error(err, w)
	}

	switch commandName {
	case "create-list":
		rm.bus.SendCommand("create-list-command", payload, reqID)

	//case "update-list":
	//	rm.UpdateList(w, r)

	//case "delete-list":
	//	rm.DeleteList(w, r)

	//case "get-all-lists":
	//	rm.GetAllLists(w, r)

	default:
		err := fmt.Errorf("command '%s' not found", commandName)
		rm.Error(err, w)
	}
}

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
