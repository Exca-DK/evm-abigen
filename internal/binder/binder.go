package binder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// RunAbigen generates Go bindings for a single ABI configuration using Docker.
func RunAbigen(abiData []byte, toolchainVersion, outPath, outputPackageName string) error {
	if len(abiData) == 0 {
		return errors.New("empty abi data")
	}

	// Write ABI data to a temporary file. This is needed because we mount this file.
	abiFilePath, err := writeABIDataToFile(abiData)
	if err != nil {
		return fmt.Errorf("failed to write ABI data to file: %w", err)
	}
	defer os.Remove(abiFilePath)

	outPath, err = filepath.Abs(outPath)
	if err != nil {
		return fmt.Errorf("failed to resolve output path: %w", err)
	}

	if err := os.MkdirAll(path.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	containerReq, err := createContainerRequest(toolchainVersion, abiFilePath, outPath, outputPackageName)
	if err != nil {
		return fmt.Errorf("failed to create container request: %w", err)
	}

	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: containerReq,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start container: %w", err)
	}
	defer container.Terminate(context.Background())

	log.Printf("✅ Successfully generated Go bindings in %s", outPath)
	return nil
}

// writeABIDataToFile writes the provided ABI data to a temporary file.
func writeABIDataToFile(abiData []byte) (string, error) {
	tmpFile, err := os.CreateTemp("", "*.abi")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary ABI file: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(abiData); err != nil {
		return "", fmt.Errorf("failed to write ABI data to file: %w", err)
	}

	return tmpFile.Name(), nil
}

// createMounts creates the necessary mount points for the Docker container.
func createMounts(abiFilePath, outputPath string) []testcontainers.ContainerMount {
	// Mount ABI directory as read-only.
	abiMount := testcontainers.BindMount(filepath.Dir(abiFilePath), testcontainers.ContainerMountTarget("/contracts"))
	abiMount.ReadOnly = true

	// Mount output directory.
	outputMount := testcontainers.BindMount(filepath.Dir(outputPath), testcontainers.ContainerMountTarget("/output"))

	return []testcontainers.ContainerMount{abiMount, outputMount}
}

// createContainerRequest generates the container request for running the abigen command.
func createContainerRequest(toolchainVersion, abiFile, outputPath, outputPackageName string) (testcontainers.ContainerRequest, error) {
	outputFile := filepath.Base(outputPath)
	return testcontainers.ContainerRequest{
		Image:      fmt.Sprintf("ethereum/client-go:alltools-%s", toolchainVersion),
		Cmd:        []string{"abigen", "--abi=" + path.Join("/contracts", filepath.Base(abiFile)), "--pkg=" + outputPackageName, "--out=" + path.Join("/output", outputFile)},
		AutoRemove: true,
		WaitingFor: wait.ForExit().WithPollInterval(time.Second),
		Mounts:     testcontainers.Mounts(createMounts(abiFile, outputPath)...),
		User:       fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
	}, nil

}
