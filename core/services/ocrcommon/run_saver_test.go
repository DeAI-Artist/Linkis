package ocrcommon

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/DeAI-Artist/MintAI/core/logger"
	"github.com/DeAI-Artist/MintAI/core/services/pipeline"
	"github.com/DeAI-Artist/MintAI/core/services/pipeline/mocks"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"
)

func TestRunSaver(t *testing.T) {
	pipelineRunner := mocks.NewRunner(t)
	rs := NewResultRunSaver(
		pipelineRunner,
		logger.TestLogger(t),
		1000,
		100,
	)
	servicetest.Run(t, rs)
	for i := 0; i < 100; i++ {
		d := i
		pipelineRunner.On("InsertFinishedRun", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil).
			Run(func(args mock.Arguments) {
				args.Get(0).(*pipeline.Run).ID = int64(d)
			}).
			Once()
		rs.Save(&pipeline.Run{ID: int64(i)})
	}
}
