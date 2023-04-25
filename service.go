package topspin

type (
	Service interface {
		Worker
	}

	SimpleService struct {
		*SimpleWorker
	}
)

func NewSimpleService(name string, log Logger) *SimpleService {
	return &SimpleService{
		SimpleWorker: NewWorker(name, log),
	}
}
