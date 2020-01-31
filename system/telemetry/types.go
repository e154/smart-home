package telemetry

type ITelemetry interface {
	Broadcast(interface{})
	BroadcastOne(interface{})
}

type Device struct {
	Id          int64
	ElementName string
}

type Node struct{}

type WorkflowScenario struct {
	WorkflowId int64
	ScenarioId int64
}
