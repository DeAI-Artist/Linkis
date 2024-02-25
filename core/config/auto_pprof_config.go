package config

import (
	"github.com/DeAI-Artist/MintAI/core/utils"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type AutoPprof interface {
	BlockProfileRate() int
	CPUProfileRate() int
	Enabled() bool
	GatherDuration() commonconfig.Duration
	GatherTraceDuration() commonconfig.Duration
	GoroutineThreshold() int
	MaxProfileSize() utils.FileSize
	MemProfileRate() int
	MemThreshold() utils.FileSize
	MutexProfileFraction() int
	PollInterval() commonconfig.Duration
	ProfileRoot() string
}
