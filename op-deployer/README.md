# Base op-deployer

This branch provides the focused offline Base devnet generator:

```sh
just op-deployer/build
op-deployer/bin/op-deployer --help
```

See [command documentation](cmd/op-deployer/README.md) for supported inputs,
required runtime tools, output behavior, artifact provenance, and validation.
The upstream init/apply/inspect/bootstrap/upgrade/verify CLI is not supported by
this fork's executable. Shared deployment libraries remain available internally.
