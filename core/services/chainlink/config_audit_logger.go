package chainlink

import (
	"github.com/DeAI-Artist/MintAI/core/build"
	"github.com/DeAI-Artist/MintAI/core/config/toml"
	"github.com/DeAI-Artist/MintAI/core/store/models"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
)

type auditLoggerConfig struct {
	c toml.AuditLogger
}

func (a auditLoggerConfig) Enabled() bool {
	return *a.c.Enabled
}

func (a auditLoggerConfig) ForwardToUrl() (commonconfig.URL, error) {
	return *a.c.ForwardToUrl, nil
}

func (a auditLoggerConfig) Environment() string {
	if !build.IsProd() {
		return "develop"
	}
	return "production"
}

func (a auditLoggerConfig) JsonWrapperKey() string {
	return *a.c.JsonWrapperKey
}

func (a auditLoggerConfig) Headers() (models.ServiceHeaders, error) {
	return *a.c.Headers, nil
}
