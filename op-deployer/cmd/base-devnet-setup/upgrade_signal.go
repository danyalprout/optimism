package main

import (
	"encoding/json"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Fixed devnet-only address, independent of deployer nonces or chain timestamps.
const upgradeSignalAddress = "0x000000000000000000000000000000000000ba5e"

type contractArtifact struct {
	Code          hexutil.Bytes `json:"code"`
	UpgradeIDs    []string      `json:"upgradeIds"`
	StorageLayout struct {
		Storage []struct {
			Label  string `json:"label"`
			Slot   string `json:"slot"`
			Offset int    `json:"offset"`
			Type   string `json:"type"`
		} `json:"storage"`
		Types map[string]struct {
			Encoding      string `json:"encoding"`
			NumberOfBytes string `json:"numberOfBytes"`
			Base          string `json:"base"`
			Label         string `json:"label"`
		} `json:"types"`
	} `json:"storageLayout"`
}

func installUpgradeSignal(genesis *core.Genesis, rollup map[string]json.RawMessage, artifactJSON []byte) error {
	var artifact contractArtifact
	if err := json.Unmarshal(artifactJSON, &artifact); err != nil {
		return err
	}
	var base map[string]*uint64
	if raw, ok := rollup["base"]; ok {
		if err := json.Unmarshal(raw, &base); err != nil {
			return err
		}
	}
	var schedule []uint64
	for _, id := range artifact.UpgradeIDs {
		value := base[id]
		if raw, ok := rollup[id+"_time"]; ok {
			if err := json.Unmarshal(raw, &value); err != nil {
				return err
			}
		}
		timestamp := uint64(0)
		if value != nil {
			timestamp = *value
			if timestamp == 0 {
				timestamp = genesis.Timestamp
			}
		}
		schedule = append(schedule, timestamp)
	}
	minimumString := os.Getenv("UPGRADE_SIGNAL_MIN_PROTOCOL_VERSION")
	if minimumString == "" {
		minimumString = "4294967296"
	}
	minimum, ok := new(big.Int).SetString(minimumString, 10)
	if !ok || minimum.Sign() < 0 || minimum.BitLen() > 256 {
		return fmt.Errorf("invalid minimum protocol version: %s", minimumString)
	}
	storage := make(map[common.Hash]common.Hash)
	fields := 0
	for _, field := range artifact.StorageLayout.Storage {
		slot, ok := new(big.Int).SetString(field.Slot, 10)
		if !ok || field.Offset != 0 {
			return fmt.Errorf("unsupported storage layout for %s", field.Label)
		}
		typ := artifact.StorageLayout.Types[field.Type]
		switch field.Label {
		case "minimumVersion":
			if typ.Label != "uint256" {
				return fmt.Errorf("minimumVersion must be uint256")
			}
			storage[common.BigToHash(slot)] = common.BigToHash(minimum)
		case "schedule":
			if typ.Encoding != "dynamic_array" || artifact.StorageLayout.Types[typ.Base].Label != "uint64" {
				return fmt.Errorf("schedule must be uint64[]")
			}
			storage[common.BigToHash(slot)] = common.BigToHash(big.NewInt(int64(len(schedule))))
			start := new(big.Int).SetBytes(crypto.Keccak256(common.BigToHash(slot).Bytes()))
			for i, timestamp := range schedule {
				key := common.BigToHash(new(big.Int).Add(start, big.NewInt(int64(i/4))))
				word := storage[key].Big()
				word.Or(word, new(big.Int).Lsh(new(big.Int).SetUint64(timestamp), uint(i%4)*64))
				storage[key] = common.BigToHash(word)
			}
		default:
			return fmt.Errorf("unsupported contract storage field: %s", field.Label)
		}
		fields++
	}
	if fields != 2 || len(artifact.Code) == 0 || len(schedule) == 0 {
		return fmt.Errorf("incomplete upgrade-signal artifact")
	}
	address := common.HexToAddress(upgradeSignalAddress)
	if _, exists := genesis.Alloc[address]; exists {
		return fmt.Errorf("upgrade-signal address already allocated: %s", address)
	}
	genesis.Alloc[address] = types.Account{Code: artifact.Code, Storage: storage, Balance: new(big.Int), Nonce: 1}
	return nil
}
