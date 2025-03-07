package binder

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"text/template"
)

// LoadDeployedBytecodeFromFoundryJSON loads the deployed bytecode from a Foundry JSON output file
func LoadDeployedBytecodeFromFoundryJSON(foundryJSONFilePath string) ([]byte, error) {
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

	deployedBytecodeObject, ok := jsonData["deployedBytecode"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("deployedBytecode field is missing or incorrectly formatted")
	}

	deployedBytecode, ok := deployedBytecodeObject["object"].(string)
	if !ok {
		return nil, fmt.Errorf("deployedBytecode.object field is missing or incorrectly formatted")
	}

	bytecode, err := hexStringToByteSlice(deployedBytecode)
	if err != nil {
		return nil, fmt.Errorf("failed to convert deployed bytecode to []byte: %w", err)
	}

	return bytecode, nil
}

func hexStringToByteSlice(hexStr string) ([]byte, error) {
	if len(hexStr) > 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	return hex.DecodeString(hexStr)
}

// GenerateGoFileWithBytecode generates a Go file containing the deployed bytecode variable
func GenerateGoFileWithBytecode(packageName, outputFilePath, bytecodeVarName string, bytecode []byte) error {
	const goFileTemplate = `
// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package {{.PackageName}}

var {{.BytecodeVarName}} = []byte{
	{{.Bytecode}}
}
`

	var bytecodeFormatted string
	for i, b := range bytecode {
		if i%512 == 0 && i != 0 {
			bytecodeFormatted += "\n\t"
		}
		bytecodeFormatted += fmt.Sprintf("0x%02x,", b)
	}

	data := struct {
		PackageName     string
		BytecodeVarName string
		Bytecode        string
	}{
		PackageName:     packageName,
		BytecodeVarName: bytecodeVarName,
		Bytecode:        bytecodeFormatted,
	}

	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("failed to create Go file: %w", err)
	}
	defer outputFile.Close()

	tmpl, err := template.New("goFile").Parse(goFileTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	return nil
}
