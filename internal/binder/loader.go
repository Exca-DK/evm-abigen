package binder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type TypeABI string

const (
	Foundry TypeABI = "foundry"
	File    TypeABI = "file"
)

// LoadABIFromFile loads ABI data from a .abi file
func LoadABI(path string, t TypeABI) ([]byte, error) {
	switch t {
	case Foundry:
		return loadABIFromFoundryJSON(path)
	case File:
		return loadABIFromFile(path)
	}
	return nil, errors.New("unknown type")
}

// loadABIFromFile loads ABI data from a .abi file
func loadABIFromFile(abiFilePath string) ([]byte, error) {
	if _, err := os.Stat(abiFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	abiData, err := os.ReadFile(abiFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %w", err)
	}

	return abiData, nil
}

// loadABIFromFoundryJSON loads ABI data from a Foundry JSON output file
func loadABIFromFoundryJSON(foundryJSONFilePath string) ([]byte, error) {
	if _, err := os.Stat(foundryJSONFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	fileContent, err := os.ReadFile(foundryJSONFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Foundry JSON file: %w", err)
	}

	var jsonData map[string]any
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse Foundry JSON: %w", err)
	}

	abiData, ok := jsonData["abi"].([]any)
	if !ok {
		return nil, fmt.Errorf("foundry JSON ABI field is missing or incorrectly formatted")
	}

	abiJSON, err := json.Marshal(abiData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal ABI data from Foundry JSON: %w", err)
	}

	return abiJSON, nil
}
