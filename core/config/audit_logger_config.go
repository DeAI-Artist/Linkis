package config

import (
	"github.com/DeAI-Artist/MintAI/core/store/models"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type AuditLogger interface {
	Enabled() bool
	ForwardToUrl() (commonconfig.URL, error)
	Environment() string
	JsonWrapperKey() string
	Headers() (models.ServiceHeaders, error)
}
