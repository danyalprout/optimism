package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/ethereum-optimism/optimism/op-chain-ops/addresses"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/artifacts"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/state"
	"github.com/ethereum/go-ethereum/common"
)

// baseIntent fixes protocol/deployment settings to the Base devnet contract.
// Only chain IDs and role addresses are caller inputs; no intent file or generic
// deployment overrides are accepted.
func baseIntent(l1ID string, l2ID uint64) (*state.Intent, error) {
	l1, err := strconv.ParseUint(l1ID, 10, 64)
	if err != nil || l1 == 0 || l2ID == 0 {
		return nil, fmt.Errorf("chain IDs must be positive integers")
	}
	roles := make(map[string]common.Address)
	for _, key := range []string{"DEPLOYER_ADDR", "SEQUENCER_ADDR", "BATCHER_ADDR", "PROPOSER_ADDR", "CHALLENGER_ADDR"} {
		value := os.Getenv(key)
		if !common.IsHexAddress(value) {
			return nil, fmt.Errorf("invalid %s", key)
		}
		roles[key] = common.HexToAddress(value)
	}
	locator, err := artifacts.NewLocatorFromURL("embedded")
	if err != nil {
		return nil, err
	}
	owner := roles["DEPLOYER_ADDR"]
	return &state.Intent{
		ConfigType: state.IntentTypeCustom, L1ChainID: l1, FundDevAccounts: true,
		L1ContractsLocator: locator, L2ContractsLocator: locator,
		SuperchainRoles: &addresses.SuperchainRoles{SuperchainProxyAdminOwner: owner, SuperchainGuardian: owner, ProtocolVersionsOwner: owner, Challenger: roles["CHALLENGER_ADDR"]},
		Chains: []*state.ChainIntent{{
			ID:                    common.BigToHash(new(big.Int).SetUint64(l2ID)),
			BaseFeeVaultRecipient: owner, L1FeeVaultRecipient: owner, SequencerFeeVaultRecipient: owner, OperatorFeeVaultRecipient: owner,
			Eip1559DenominatorCanyon: 250, Eip1559Denominator: 50, Eip1559Elasticity: 6, GasLimit: 60000000,
			ChainFeesRecipient: owner, MinBaseFee: 1000000000,
			Roles: state.ChainRoles{L1ProxyAdminOwner: owner, L2ProxyAdminOwner: owner, SystemConfigOwner: owner,
				UnsafeBlockSigner: roles["SEQUENCER_ADDR"], Batcher: roles["BATCHER_ADDR"], Proposer: roles["PROPOSER_ADDR"], Challenger: roles["CHALLENGER_ADDR"]},
		}},
	}, nil
}
