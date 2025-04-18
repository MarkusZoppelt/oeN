package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MarkusZoppelt/oen/pkg/agent"
)

// ReadFileDefinition allows reading file contents
var ReadFileDefinition = agent.ToolDefinition{
	Name:        "read_file",
	Description: "Read the contents of a given relative file path. Use this when you want to see what's inside a file. Do not use this with directory names.",
	InputSchema: GenerateSchema[ReadFileInput](),
	Function:    ReadFile,
}

// ReadFileInput holds input for read_file tool
type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

// ReadFileInputSchema holds the schema for read_file input
var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

// ReadFile reads the contents of a file
func ReadFile(input json.RawMessage) (string, error) {
	var in ReadFileInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}
	content, err := os.ReadFile(in.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ListFilesDefinition allows listing files in a directory
var ListFilesDefinition = agent.ToolDefinition{
	Name:        "list_files",
	Description: "List files and directories at a given path. If no path is provided, lists files in the current directory.",
	InputSchema: GenerateSchema[ListFilesInput](),
	Function:    ListFiles,
}

// ListFilesInput holds input for list_files tool
type ListFilesInput struct {
	Path string `json:"path,omitempty" jsonschema_description:"Optional relative path to list files from. Defaults to current directory if not provided."`
}

// ListFilesInputSchema holds the schema for list_files input
var ListFilesInputSchema = GenerateSchema[ListFilesInput]()

// ListFiles lists files and directories
func ListFiles(input json.RawMessage) (string, error) {
	var in ListFilesInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}

	dir := "."
	if in.Path != "" {
		dir = in.Path
	}

	var files []string
	err := filepath.Walk(dir, func(pathStr string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dir, pathStr)
		if err != nil {
			return err
		}

		if relPath != "." {
			if info.IsDir() {
				files = append(files, relPath+"/")
			} else {
				files = append(files, relPath)
			}
		}
		return nil
	})
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// EditFileDefinition allows editing file contents
var EditFileDefinition = agent.ToolDefinition{
	Name:        "edit_file",
	Description: `Make edits to a text file.

Replaces 'old_str' with 'new_str' in the given file. 'old_str' and 'new_str' MUST be different from each other.

If the file specified with path doesn't exist, it will be created.
`,
	InputSchema: GenerateSchema[EditFileInput](),
	Function:    EditFile,
}

// EditFileInput holds input for edit_file tool
type EditFileInput struct {
	Path   string `json:"path" jsonschema_description:"The path to the file"`
	OldStr string `json:"old_str" jsonschema_description:"Text to search for - must match exactly and must only have one match exactly"`
	NewStr string `json:"new_str" jsonschema_description:"Text to replace old_str with"`
}

// EditFileInputSchema holds the schema for edit_file input
var EditFileInputSchema = GenerateSchema[EditFileInput]()

// EditFile edits or creates a file
func EditFile(input json.RawMessage) (string, error) {
	var in EditFileInput
	if err := json.Unmarshal(input, &in); err != nil {
		return "", err
	}

	if in.Path == "" || in.OldStr == in.NewStr {
		return "", fmt.Errorf("invalid input parameters")
	}

	content, err := os.ReadFile(in.Path)
	if err != nil {
		if os.IsNotExist(err) && in.OldStr == "" {
			return createNewFile(in.Path, in.NewStr)
		}
		return "", err
	}

	oldContent := string(content)
	newContent := strings.Replace(oldContent, in.OldStr, in.NewStr, -1)

	if oldContent == newContent && in.OldStr != "" {
		return "", fmt.Errorf("old_str not found in file")
	}

	if err := os.WriteFile(in.Path, []byte(newContent), 0644); err != nil {
		return "", err
	}
	return "OK", nil
}

// createNewFile creates a new file with content
func createNewFile(filePath, content string) (string, error) {
	dir := path.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", fmt.Errorf("failed to create directory: %w", err)
		}
	}
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	return fmt.Sprintf("Successfully created file %s", filePath), nil
}
