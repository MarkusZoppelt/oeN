package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/MarkusZoppelt/oen/pkg/agent"
	"github.com/MarkusZoppelt/oen/pkg/tools"
)

func main() {
	client := anthropic.NewClient()

	scanner := bufio.NewScanner(os.Stdin)
	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	toolsList := []agent.ToolDefinition{
		tools.ReadFileDefinition,
		tools.ListFilesDefinition,
		tools.EditFileDefinition,
		tools.MakeDirectoryDefinition,
		tools.RemoveDirectoryDefinition,
		tools.RenameDirectoryDefinition,
	}

	ag := agent.NewAgent(client, getUserMessage, toolsList)
	if err := ag.Run(context.TODO()); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}
