package tools

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/MarkusZoppelt/oen/pkg/agent"
)

// MakeDirectoryDefinition allows creating directories recursively
var MakeDirectoryDefinition = agent.ToolDefinition{
	Name:        "make_directory",
	Description: "Create a new directory at the given relative path, creating parent directories as needed.",
	InputSchema: GenerateSchema[MakeDirectoryInput](),
	Function:    MakeDirectory,
}

// MakeDirectoryInput holds input for make_directory tool
type MakeDirectoryInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of the directory to create."`
}

// MakeDirectory creates a directory and any necessary parents
func MakeDirectory(input json.RawMessage) (string, error) {
	var in MakeDirectoryInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}
	if in.Path == "" {
		return "", fmt.Errorf("path must not be empty")
	}
	if err := os.MkdirAll(in.Path, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	return fmt.Sprintf("Successfully created directory %s", in.Path), nil
}

// RemoveDirectoryDefinition allows removing directories
var RemoveDirectoryDefinition = agent.ToolDefinition{
	Name:        "remove_directory",
	Description: "Remove a directory at the given relative path. If recursive is true, remove all contents recursively; otherwise, only if empty.",
	InputSchema: GenerateSchema[RemoveDirectoryInput](),
	Function:    RemoveDirectory,
}

// RemoveDirectoryInput holds input for remove_directory tool
type RemoveDirectoryInput struct {
	Path      string `json:"path" jsonschema_description:"The relative path of the directory to remove."`
	Recursive bool   `json:"recursive" jsonschema_description:"Whether to remove directory recursively along with its contents."`
}

// RemoveDirectory removes a directory based on input
func RemoveDirectory(input json.RawMessage) (string, error) {
	var in RemoveDirectoryInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}
	if in.Path == "" {
		return "", fmt.Errorf("path must not be empty")
	}
	if in.Recursive {
		if err := os.RemoveAll(in.Path); err != nil {
			return "", fmt.Errorf("failed to remove directory recursively: %w", err)
		}
	} else {
		if err := os.Remove(in.Path); err != nil {
			return "", fmt.Errorf("failed to remove directory: %w", err)
		}
	}
	return fmt.Sprintf("Successfully removed directory %s", in.Path), nil
}

// RenameDirectoryDefinition allows renaming or moving directories
var RenameDirectoryDefinition = agent.ToolDefinition{
	Name:        "rename_directory",
	Description: "Rename or move a directory from the old path to the new path.",
	InputSchema: GenerateSchema[RenameDirectoryInput](),
	Function:    RenameDirectory,
}

// RenameDirectoryInput holds input for rename_directory tool
type RenameDirectoryInput struct {
	OldPath string `json:"old_path" jsonschema_description:"The current relative path of the directory."`
	NewPath string `json:"new_path" jsonschema_description:"The new relative path for the directory."`
}

// RenameDirectory renames or moves a directory
func RenameDirectory(input json.RawMessage) (string, error) {
	var in RenameDirectoryInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}
	if in.OldPath == "" || in.NewPath == "" {
		return "", fmt.Errorf("old_path and new_path must not be empty")
	}
	if err := os.Rename(in.OldPath, in.NewPath); err != nil {
		return "", fmt.Errorf("failed to rename directory: %w", err)
	}
	return fmt.Sprintf("Successfully renamed directory from %s to %s", in.OldPath, in.NewPath), nil
}
