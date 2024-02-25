//go:build integration

package ocr2_test

import (
	"testing"

	"github.com/DeAI-Artist/MintAI/core/config/env"
)

func TestIntegration_OCR2_plugins(t *testing.T) {
	t.Setenv(string(env.MedianPlugin.Cmd), "chainlink-feeds")
	testIntegration_OCR2(t)
}
