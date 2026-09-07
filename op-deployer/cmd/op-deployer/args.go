package main

import (
	"flag"
	"fmt"
	"os"
)

// parseArgs exposes only inputs used by Base devnets. Environment variables
// remain supported so Compose does not have to duplicate keys in its command.
func parseArgs(args []string) error {
	flags := flag.NewFlagSet("op-deployer", flag.ContinueOnError)
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage: op-deployer [flags]\nGenerate a fresh offline Base L1/L2 devnet. Environment defaults are supported.")
		flags.PrintDefaults()
	}
	bindings := map[string]string{}
	for _, spec := range []struct{ name, key, help string }{
		{"output-dir", "OUTPUT_DIR", "L1 output directory (OUTPUT_DIR; default /output)"},
		{"l2-output-dir", "L2_OUTPUT_DIR", "L2 output directory (L2_OUTPUT_DIR; default /devnet/l2/configs)"},
		{"l1-chain-id", "CHAIN_ID", "L1 chain ID (CHAIN_ID or L1_CHAIN_ID; default 1337)"},
		{"l2-chain-id", "L2_CHAIN_ID", "L2 chain ID (L2_CHAIN_ID; default 84538453)"},
		{"slot-duration", "SLOT_DURATION", "L1 slot duration in seconds (SLOT_DURATION; default 12)"},
		{"activation-admin", "L2_ACTIVATION_ADMIN_ADDR", "Activation admin address (defaults to SEQUENCER_ADDR)"},
		{"isthmus-block", "L2_ISTHMUS_BLOCK", "Isthmus activation block; unset preserves upstream configuration"},
		{"azul-block", "L2_BASE_AZUL_BLOCK", "Azul activation block"},
		{"beryl-block", "L2_BASE_BERYL_BLOCK", "Beryl activation block"},
		{"cobalt-block", "L2_BASE_COBALT_BLOCK", "Cobalt activation block"},
		{"denim-block", "L2_BASE_DENIM_BLOCK", "Denim activation block"},
		{"zenith-block", "L2_BASE_ZENITH_BLOCK", "Zenith activation block"},
	} {
		flags.String(spec.name, "", spec.help)
		bindings[spec.name] = spec.key
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("unexpected arguments: %v; usage: op-deployer [flags]", flags.Args())
	}
	var err error
	flags.Visit(func(f *flag.Flag) {
		if err == nil {
			err = os.Setenv(bindings[f.Name], f.Value.String())
		}
	})
	return err
}
