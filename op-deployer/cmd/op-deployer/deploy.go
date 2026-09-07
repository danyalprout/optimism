package main

import (
	"fmt"
	"os"

	"github.com/ethereum-optimism/optimism/op-chain-ops/foundry"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/artifacts"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/broadcaster"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/opcm"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/pipeline"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"
	scriptenv "github.com/ethereum-optimism/optimism/op-deployer/pkg/env"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

// deploy executes only the stages required by the fixed Base intent. There is
// no live deployment, resume state, multi-chain loop, broadcasting, or prestate.
func deploy(intent *state.Intent, st *state.State, logger log.Logger, cache string) error {
	if err := intent.Check(); err != nil {
		return err
	}
	if len(intent.Chains) != 1 {
		return fmt.Errorf("expected exactly one Base chain")
	}
	var fs foundry.StatDirFs
	if path := os.Getenv("BASE_DEVNET_ARTIFACTS"); path != "" {
		fs = os.DirFS(path).(foundry.StatDirFs)
	} else {
		var err error
		fs, err = artifacts.ExtractEmbedded(cache)
		if err != nil {
			return err
		}
	}
	bcaster := broadcaster.NoopBroadcaster()
	deployer := common.Address{0x01}
	host, err := scriptenv.DefaultScriptHost(bcaster, logger, deployer, fs)
	if err != nil {
		return err
	}
	scripts, err := opcm.NewScripts(host)
	if err != nil {
		return err
	}
	pEnv := &pipeline.Env{L1ScriptHost: host, Deployer: deployer, Logger: logger, Broadcaster: bcaster, Scripts: scripts, StateWriter: pipeline.NoopStateWriter()}
	if err := pipeline.InitGenesisStrategy(pEnv, intent, st); err != nil {
		return err
	}
	if err := pipeline.DeploySuperchain(pEnv, intent, st); err != nil {
		return err
	}
	if err := pipeline.DeployImplementations(pEnv, intent, st); err != nil {
		return err
	}
	chainID := intent.Chains[0].ID
	if err := pipeline.DeployOPChain(pEnv, intent, st, chainID); err != nil {
		return err
	}
	if err := pipeline.GenerateL2Genesis(pEnv, intent, pipeline.ArtifactsBundle{L1: fs, L2: fs}, st, chainID); err != nil {
		return err
	}
	if err := pipeline.PreinstallL1DevGenesis(pEnv, intent, st); err != nil {
		return err
	}
	dump, err := host.StateDump()
	if err != nil {
		return err
	}
	st.L1StateDump = &state.GzipData[foundry.ForgeAllocs]{Data: dump}
	st.AppliedIntent = intent
	return nil
}
