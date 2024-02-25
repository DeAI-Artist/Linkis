package docs

import (
	"testing"

	"github.com/DeAI-Artist/MintAI/core/services/chainlink/cfgtest"
)

func TestCoreDefaults_notNil(t *testing.T) {
	cfgtest.AssertFieldsNotNil(t, CoreDefaults())
}
