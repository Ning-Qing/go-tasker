package gotasker

import "time"

type States string

const (
	State_Waiting    States = "waiting"
	State_Scheduling States = "scheduling"
	State_Running    States = "running"
	State_Finished   States = "finished"
	State_Failed     States = "failed"
	State_Stopped    States = "stopped"
)

type Actions string

const (
	Action_Unknown  Actions = "unknown"
	Action_Execute  Actions = "execute"
	Action_Rollback Actions = "rollback"
)

type Status struct {
	Step   uint8
	State  States
	Action Actions
}

type Metadata struct {
	Name     string
	Describe string
}

type Task struct {
	ID   string
	Kind string

	Status   Status
	Metadata Metadata

	Policies []Policy

	Params []byte

	CreatedAt time.Time
	UpdatedAt time.Time
}
