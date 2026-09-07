// extract-artifacts is a build-time utility; it is not included in the runtime image.
package main

import (
	"fmt"
	"github.com/ethereum-optimism/optimism/op-deployer/pkg/deployer/artifacts"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: extract-artifacts DIRECTORY")
		os.Exit(1)
	}
	if _, err := artifacts.ExtractEmbedded(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
