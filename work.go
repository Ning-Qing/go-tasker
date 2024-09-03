package gotasker

// Work 
type Work interface {
	Execute(wf Workflow, in []byte) (out []byte, err error)
	Rollback(wf Workflow, in []byte) (out []byte, err error)
	Save() []byte
	Restore(data []byte)
}
