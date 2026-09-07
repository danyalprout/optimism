package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompletedSetupRepairsNestedMountpoint(t *testing.T) {
	root := t.TempDir()
	l1 := filepath.Join(root, "l1")
	l2 := filepath.Join(root, "l2")
	t.Setenv("OUTPUT_DIR", l1)
	t.Setenv("L2_OUTPUT_DIR", l2)
	for _, path := range []string{filepath.Join(l1, ".setup-complete"), filepath.Join(l1, "el/genesis.json"), filepath.Join(l1, "cl/genesis.ssz"), filepath.Join(l2, ".setup-complete")} {
		if err := write(path, []byte("preserve")); err != nil {
			t.Fatal(err)
		}
	}
	args := os.Args
	os.Args = []string{"op-deployer"}
	defer func() { os.Args = args }()
	if err := run(); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(filepath.Join(l1, "l2"))
	if err != nil || !info.IsDir() {
		t.Fatalf("nested read-only mountpoint missing: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(l1, "el/genesis.json"))
	if err != nil || string(data) != "preserve" {
		t.Fatal("resume changed existing genesis")
	}
}
