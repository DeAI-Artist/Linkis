package chainlink

import "github.com/DeAI-Artist/MintAI/core/config/toml"

type thresholdConfig struct {
	s toml.ThresholdKeyShareSecrets
}

func (t *thresholdConfig) ThresholdKeyShare() string {
	if t.s.ThresholdKeyShare == nil {
		return ""
	}
	return string(*t.s.ThresholdKeyShare)
}
