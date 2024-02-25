package workflows

import (
	"github.com/DeAI-Artist/MintAI/core/capabilities/targets"
	"github.com/DeAI-Artist/MintAI/core/chains/legacyevm"
	"github.com/DeAI-Artist/MintAI/core/logger"
	"github.com/DeAI-Artist/MintAI/core/services/job"
	"github.com/DeAI-Artist/MintAI/core/services/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
)

type Delegate struct {
	registry types.CapabilitiesRegistry
	logger   logger.Logger
}

var _ job.Delegate = (*Delegate)(nil)

func (d *Delegate) JobType() job.Type {
	return job.Workflow
}

func (d *Delegate) BeforeJobCreated(spec job.Job) {}

func (d *Delegate) AfterJobCreated(jb job.Job) {}

func (d *Delegate) BeforeJobDeleted(spec job.Job) {}

func (d *Delegate) OnDeleteJob(jb job.Job, q pg.Queryer) error { return nil }

// ServicesForSpec satisfies the job.Delegate interface.
func (d *Delegate) ServicesForSpec(spec job.Job) ([]job.ServiceCtx, error) {
	engine, err := NewEngine(d.logger, d.registry)
	if err != nil {
		return nil, err
	}
	return []job.ServiceCtx{engine}, nil
}

func NewDelegate(logger logger.Logger, registry types.CapabilitiesRegistry, legacyEVMChains legacyevm.LegacyChainContainer) *Delegate {
	// NOTE: we temporarily do registration inside NewDelegate, this will be moved out of job specs in the future
	_ = targets.InitializeWrite(registry, legacyEVMChains)

	return &Delegate{logger: logger, registry: registry}
}
