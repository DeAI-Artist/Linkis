package resolver

import (
	"testing"

	"github.com/DeAI-Artist/MintAI/core/internal/testutils/configtest"
	"github.com/DeAI-Artist/MintAI/core/services/chainlink"
)

func Test_ToFeatures(t *testing.T) {
	query := `
	{
		features {
			... on Features {
				csa
				feedsManager
			}	
		}
	}`

	testCases := []GQLTestCase{
		unauthorizedTestCase(GQLTestCase{query: query}, "features"),
		{
			name:          "success",
			authenticated: true,
			before: func(f *gqlTestFramework) {
				f.App.On("GetConfig").Return(configtest.NewGeneralConfig(t, func(c *chainlink.Config, s *chainlink.Secrets) {
					t, f := true, false
					c.Feature.UICSAKeys = &f
					c.Feature.FeedsManager = &t
				}))
			},
			query: query,
			result: `
			{
				"features": {
					"csa": false,
					"feedsManager": true
				}
			}`,
		},
	}

	RunGQLTests(t, testCases)
}
