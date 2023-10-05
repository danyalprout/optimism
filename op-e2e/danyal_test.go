package op_e2e

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMinesBlocksPreShanghaiUsingV2(t *testing.T) {
	cfg := DefaultSystemConfig(t)

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	for i := 0; i < 5; i++ {
		l1Block, err := l1Client.BlockByNumber(context.Background(), nil)
		require.Nil(t, err)

		l2Block, err := l2Client.BlockByNumber(context.Background(), nil)
		require.Nil(t, err)

		require.Nil(t, err)
		fmt.Println("DANYAL L1", l1Block.Number().Uint64())
		fmt.Println("DANYAL L2", l2Block.Number().Uint64())

		time.Sleep(3 * time.Second)
	}
}

func TestMinesFromGenesis(t *testing.T) {
	cfg := DefaultSystemConfig(t)
	s := hexutil.Uint64(0)
	cfg.DeployConfig.L2GenesisRegolithTimeOffset = &s
	cfg.DeployConfig.L2GenesisCanyonTimeOffset = &s

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()

	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	for i := 0; i < 5; i++ {
		l1Block, err := l1Client.BlockByNumber(context.Background(), nil)
		require.Nil(t, err)

		l2Block, err := l2Client.BlockByNumber(context.Background(), nil)
		require.Nil(t, err)

		require.Nil(t, err)
		fmt.Println("DANYAL L1", l1Block.Number().Uint64())
		fmt.Println("DANYAL L2", l2Block.Number().Uint64())
		fmt.Println("DANYAL WI", l2Block.Withdrawals() == nil)

		time.Sleep(3 * time.Second)
	}
}

func TestTransitionsToCanyon(t *testing.T) {
	cfg := DefaultSystemConfig(t)
	s := hexutil.Uint64(0)
	c := hexutil.Uint64(20)

	cfg.DeployConfig.L2GenesisRegolithTimeOffset = &s
	cfg.DeployConfig.L2GenesisCanyonTimeOffset = &c

	start := time.Now()

	sys, err := cfg.Start(t)
	require.Nil(t, err, "Error starting up system")
	defer sys.Close()
	//
	l1Client := sys.Clients["l1"]
	l2Client := sys.Clients["sequencer"]

	for i := 0; i < 20; i++ {
		now := time.Now()
		fmt.Println("DANYAL SECS", now.Sub(start).Seconds())

		l1Block, err := l1Client.BlockByNumber(context.Background(), nil)
		require.Nil(t, err)

		l2Block, err := l2Client.BlockByNumber(context.Background(), nil)
		require.Nil(t, err)

		require.Nil(t, err)
		fmt.Println("DANYAL L1", l1Block.Number().Uint64())
		fmt.Println("DANYAL L2", l2Block.Number().Uint64())
		fmt.Println("DANYAL WI", l2Block.Withdrawals() == nil)

		time.Sleep(2 * time.Second)
	}
}
