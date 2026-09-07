package main

import (
	"flag"
	"os"
	"testing"
)

func TestOnlyBaseSetupArguments(t *testing.T) {
	for _, args := range [][]string{{"apply"}, {"inspect", "genesis"}, {"--deployment-target", "live"}, {"--l1-rpc-url", "http://localhost:8545"}, {"extract-artifacts", "/tmp/artifacts"}} {
		if err := parseArgs(args); err == nil {
			t.Fatalf("accepted unsupported arguments %v", args)
		}
	}
	if err := parseArgs([]string{"--help"}); err != flag.ErrHelp {
		t.Fatalf("help error=%v", err)
	}
}

func TestFlagsOverrideEnvironment(t *testing.T) {
	t.Setenv("L2_BASE_ZENITH_BLOCK", "100")
	t.Setenv("OUTPUT_DIR", "/old")
	if err := parseArgs([]string{"--zenith-block=", "--output-dir=/new"}); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("OUTPUT_DIR") != "/new" || os.Getenv("L2_BASE_ZENITH_BLOCK") != "" {
		t.Fatal("flags did not override environment")
	}
}
