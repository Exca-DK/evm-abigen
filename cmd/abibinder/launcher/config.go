package launcher

import (
	"os"

	"github.com/Exca-DK/evm-abigen/internal/binder"
	"github.com/go-yaml/yaml"
)

type ABIConfig struct {
	ABI              string         `yaml:"abi"`
	Package          string         `yaml:"package"`
	Output           string         `yaml:"output"`
	Type             binder.TypeABI `yaml:"type"`
	DeployedBytecode bool           `yaml:"deployed_bytecode"`
}

type Config struct {
	ABIs []ABIConfig `yaml:"abis"`
}

func Load(filePath string) (*Config, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
