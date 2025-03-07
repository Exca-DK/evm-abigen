# EVM ABIGEN

`ABIGEN` is a tool to generate Go bindings from multiple ABIs using Docker. It can automate the process of compiling and binding multiple Ethereum contract ABIs into Go code, simplifying integration with smart contracts.


```sh
go run cmd/abibinder/main.go
```

## COMMANDS:
- **help, h**  Shows a list of commands or help for one command

## GLOBAL OPTIONS:
- `--config value, -c value`  Path to the YAML config file (default: `"config.yaml"`)
- `--version value, -v value` Toolchain version for abigen (default: `"latest"`)
- `--help, -h`                Show help


## Example Configuration

This example demonstrates how to define ABI files and their associated configuration for generating Go code.

- simple.yaml
``` yaml
abis:
- abi: "foo/quoter.abi"
    package: "quoter"
    output: "foo_out/quoter.go"
    type: file
```

- foundry.yaml
``` yaml
abis:
- abi: "foo/quoter.abi"
    package: "quoter"
    output: "foo_out/quoter.go"
    type: foundry
    deployed_bytecode: true # additionally creates vars.go with embedded deployed_bytecode
```