package gotasker

import (
	"context"
	"testing"

	"github.com/google/uuid"
)

func TestWorkflow(t *testing.T) {

	cases := []workflow{
		{
			works: make([]Work, 5),
			task: &Task{
				ID: uuid.NewString(),
				Status: Status{
					Step:   0,
					State:  State_Scheduling,
					Action: Action_Unknown,
				},
			},
		},
		{
			works: make([]Work, 5),
			task: &Task{
				ID: uuid.NewString(),
				Status: Status{
					Step:   4,
					State:  State_Scheduling,
					Action: Action_Rollback,
				},
			},
		},
	}

	for _, wf := range cases {
		t.Run(wf.ID(), func(t *testing.T) {
			wf.init()

			wf.run(context.TODO())
		})
	}
}
