package gotasker

type MsgeType string

const (
	MsgType_UpdateTaskStatus MsgeType = "update_task_status"
)

type MsgUpdateTaskStatus struct {
	ID     string
	Status Status
}

type Message struct {
	Type MsgeType
	Data interface{}
}

type MessageHandler func(msg Message)
