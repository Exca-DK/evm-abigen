package launcher

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/Exca-DK/evm-abigen/internal/binder"
	"github.com/urfave/cli/v2"
)

func Launch(_ []string) error {
	app := &cli.App{
		Name:  "ABIBinder",
		Usage: "Generate Go bindings from multiple ABIs using Docker",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Usage:   "Path to the YAML config file",
				Value:   "config.yaml",
				Aliases: []string{"c"},
			},
			&cli.StringFlag{
				Name:    "version",
				Usage:   "Toolchain version for abigen",
				Value:   "latest",
				Aliases: []string{"v"},
			},
		},
		Action: generateBindings,
	}

	return app.Run(os.Args)
}

func generateBindings(ctx *cli.Context) error {
	configFile := ctx.String("config")
	cfg, err := Load(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	toolchainVersion := ctx.String("version")
	if toolchainVersion == "" {
		toolchainVersion = "latest"
	}

	for _, abi := range cfg.ABIs {
		data, err := binder.LoadABI(abi.ABI, abi.Type)
		if err != nil {
			return fmt.Errorf("failed to load abi for %s: %w", abi.ABI, err)
		}

		if err := binder.RunAbigen(data, toolchainVersion, abi.Output, abi.Package); err != nil {
			return fmt.Errorf("failed to generate bindings for %s: %w", abi.ABI, err)
		}

		if abi.DeployedBytecode {
			if abi.Type != binder.Foundry {
				return errors.New("deployed bytecode is allowed only for foundry type")
			}
			data, err := binder.LoadDeployedBytecodeFromFoundryJSON(abi.ABI)
			if err != nil {
				return fmt.Errorf("failed to load foundry %s: %w", abi.ABI, err)
			}
			err = binder.GenerateGoFileWithBytecode(abi.Package, path.Join(path.Dir(abi.Output), "vars.go"), "DeployedBytecode", data)
			if err != nil {
				return fmt.Errorf("failed to generate bytecode: %w", err)
			}
		}
	}

	log.Println("✅ All bindings generated successfully!")
	return nil
}
