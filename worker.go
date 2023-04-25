package topspin

type (
	Worker interface {
		Name() string
		Setup() error
		Start() error
		Teardown() error
		Stop() error
		Log() Logger
	}
)

type (
	SimpleWorker struct {
		name     string
		didSetup bool
		didStart bool
		log      Logger
	}
)

func NewWorker(name string, log Logger) *SimpleWorker {
	name = GenName(name, "worker")

	return &SimpleWorker{
		name: name,
		log:  log,
	}
}

func (sw SimpleWorker) Name() string {
	return sw.name
}

func (sw SimpleWorker) SetName(name string) {
	sw.name = name
}

func (sw SimpleWorker) Setup() error {
	sw.Log().Info("Setup")
	return nil
}

func (sw SimpleWorker) Start() error {
	sw.Log().Info("Start")
	return nil
}

func (sw SimpleWorker) Teardown() error {
	sw.Log().Info("Teardown")
	return nil
}

func (sw SimpleWorker) Stop() error {
	sw.Log().Info("Stop")
	return nil
}
func (sw SimpleWorker) Log() Logger {
	return sw.log
}
