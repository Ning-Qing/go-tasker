package gotasker

import (
	"context"
	"fmt"

	"github.com/looplab/fsm"
)

// Workflow is an interface for
type Workflow interface {
	ID() string
	Context() context.Context
}

const (
	FSMState_Start    = "start"
	FSMState_Execute  = "execute"
	FSMState_Rollback = "rollback"
	FSMState_Stopped  = "stopped"
	FSMState_Finished = "finished"
	FSMState_Decided  = "decided"
)

const (
	FSMEvent_Execute      = "execute"
	FSMEvent_Rollback     = "rollback"
	FSMEvent_Stop         = "stop"
	FSMEvent_Finish       = "finish"
	FSMEvent_MakeDecision = "make_decision"
)

const (
	FSMAction_BeforeExecute  = "before_execute"
	FSMAction_AfterExecute   = "after_execute"
	FSMAction_EnterExecute   = "enter_execute"
	FSMAction_BeforeRollback = "before_rollback"
	FSMAction_AfterRollback  = "after_rollback"
	FSMAction_EnterRollback  = "enter_rollback"
	FSMAction_EnterStopped   = "enter_stopped"
	FSMAction_EnterFinished  = "enter_finished"
	FSMAction_EnterDecided   = "enter_decided"
)

var workflowEvents = fsm.Events{
	{Name: FSMEvent_Execute, Src: []string{FSMState_Start, FSMState_Decided}, Dst: FSMState_Execute},
	{Name: FSMEvent_Rollback, Src: []string{FSMState_Start, FSMState_Decided}, Dst: FSMState_Rollback},
	{Name: FSMEvent_Stop, Src: []string{FSMState_Start, FSMState_Execute, FSMState_Rollback}, Dst: FSMState_Stopped},
	{Name: FSMEvent_Finish, Src: []string{FSMState_Execute, FSMState_Rollback}, Dst: FSMState_Finished},
	{Name: FSMEvent_MakeDecision, Src: []string{FSMState_Execute, FSMState_Rollback}, Dst: FSMState_Decided},
}

type workflow struct {
	ctx    context.Context
	fsm    *fsm.FSM
	works  []Work
	index  uint8
	params map[string]interface{}

	topic MessageHandler
	task  *Task
}

func (w *workflow) ID() string {
	return w.task.ID
}

func (w *workflow) Context() context.Context {
	return w.ctx
}

func (w *workflow) init() {
	callbacks := map[string]fsm.Callback{
		FSMAction_BeforeExecute:  w.beforeExecute(),
		FSMAction_AfterExecute:   w.afterExecute(),
		FSMAction_EnterExecute:   w.enterExecute(),
		FSMAction_BeforeRollback: w.beforeRollback(),
		FSMAction_AfterRollback:  w.afterRollback(),
		FSMAction_EnterRollback:  w.enterRollback(),
		FSMAction_EnterStopped:   w.enterStopped(),
		FSMAction_EnterFinished:  w.enterFinished(),
		FSMAction_EnterDecided:   w.enterDecided(),
	}
	w.fsm = fsm.NewFSM(FSMState_Start, workflowEvents, callbacks)
}

func (w *workflow) run(ctx context.Context) {
	w.ctx = ctx
	switch w.task.Status.Action {
	case Action_Unknown:
		w.index = w.task.Status.Step
		w.fsm.Event(ctx, FSMEvent_Execute)
	case Action_Rollback:
		w.index = w.task.Status.Step + 1
		w.fsm.Event(ctx, FSMEvent_Rollback)
	}
}

func (w *workflow) beforeExecute() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		w.index++
		// restore snaphost

	}
}

func (w *workflow) enterExecute() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		fmt.Println("execute")
	}
}

func (w *workflow) afterExecute() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		if int(w.index) >= len(w.works) {
			w.fsm.Event(ctx, FSMEvent_Finish)
			return
		}
		w.fsm.Event(ctx, FSMEvent_MakeDecision)
	}
}

func (w *workflow) beforeRollback() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		w.index--
		// restore snaphost

	}
}

func (w *workflow) enterRollback() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		fmt.Println("rollback")
	}
}

func (w *workflow) afterRollback() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		if w.index == 0 {
			w.fsm.Event(ctx, FSMEvent_Finish)
			return
		}
		w.fsm.Event(ctx, FSMEvent_MakeDecision)
	}
}

func (w *workflow) enterDecided() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		select {
		case <-w.ctx.Done():
			w.fsm.Event(ctx, FSMEvent_Stop)
			return
		default:
			switch e.Src {
			case FSMState_Execute:
				w.fsm.Event(ctx, FSMEvent_Execute)
			case FSMState_Rollback:
				w.fsm.Event(ctx, FSMEvent_Rollback)
			}
		}
	}
}

func (w *workflow) enterStopped() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		fmt.Println("stopped....")
		
		msg := Message{
			Type: MsgType_UpdateTaskStatus,
			Data: MsgUpdateTaskStatus{
				ID:     w.ID(),
				Status: w.task.Status,
			},
		}
		w.topic(msg)
	}
}

func (w *workflow) enterFinished() fsm.Callback {
	return func(ctx context.Context, e *fsm.Event) {
		fmt.Println("finished....")
	}
}
