package config

import "github.com/DeAI-Artist/MintAI/core/chains/evm/config/toml"

type balanceMonitorConfig struct {
	c toml.BalanceMonitor
}

func (b *balanceMonitorConfig) Enabled() bool {
	return *b.c.Enabled
}
