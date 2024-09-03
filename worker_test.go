package gotasker_test

import (
	"testing"

	tasker "github.com/Ning-Qing/go-tasker"
)

type NormalWork struct{}

func (w *NormalWork) Execute(wf tasker.Workflow, in []byte) (out []byte, err error) {
	return
}
func (w *NormalWork) Rollback(wf tasker.Workflow, in []byte) (out []byte, err error) {
	return
}
func (w *NormalWork) Save() []byte {
	return nil
}
func (w *NormalWork) Restore(data []byte) {
}

func NormalWorkBuilder() tasker.Work {
	return &NormalWork{}
}

func TestRegisterWork(t *testing.T){
	
}