package core

import (
	"github.com/DeAI-Artist/MintAI/core/chains/evm/types"
	"github.com/DeAI-Artist/MintAI/core/gethwrappers/generated/automation_utils_2_1"
	iregistry21 "github.com/DeAI-Artist/MintAI/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	"github.com/DeAI-Artist/MintAI/core/gethwrappers/generated/i_log_automation"
	"github.com/DeAI-Artist/MintAI/core/gethwrappers/generated/streams_lookup_compatible_interface"
)

var UtilsABI = types.MustGetABI(automation_utils_2_1.AutomationUtilsABI)
var RegistryABI = types.MustGetABI(iregistry21.IKeeperRegistryMasterABI)
var StreamsCompatibleABI = types.MustGetABI(streams_lookup_compatible_interface.StreamsLookupCompatibleInterfaceABI)
var ILogAutomationABI = types.MustGetABI(i_log_automation.ILogAutomationABI)
