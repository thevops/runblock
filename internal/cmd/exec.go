package cmd

import (
	"fmt"
	"os"
	"os/exec"

	pkg "runblock/pkg"
	logger "runblock/pkg/logger"

	"github.com/spf13/cobra"
)

func ExecCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "exec",
		Short: "Execute a code block by name from a Markdown file",
		Run: func(cmd *cobra.Command, args []string) {
			// Access the global file flag
			filePath, err := cmd.Flags().GetString("file")
			if err != nil {
				logger.Log.Fatalf("Failed to get file flag: %v", err)
			}

			markdownContent, err := os.ReadFile(filePath)
			if err != nil {
				logger.Log.Fatalf("Failed to read file: %v", err)
			}

			codeBlocks, err := pkg.LoadNamedCodeBlocks(markdownContent)
			if err != nil {
				logger.Log.Fatalf("Failed to extract code blocks: %v", err)
			}

			blockContent := codeBlocks[name].Content
			blockLanguage := codeBlocks[name].Language

			if blockContent == "" {
				logger.Log.Fatalf("Code block for '%s' is empty", name)
			}

			// Save the code block content to a temporary file
			fileNamePattern := fmt.Sprintf("runblock-*.%s", blockLanguage)
			tmpFile, err := os.CreateTemp("", fileNamePattern)
			if err != nil {
				logger.Log.Fatalf("Failed to create temporary file: %v", err)
			}
			logger.Log.Debug("Created temporary file: ", tmpFile.Name())
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(blockContent); err != nil {
				logger.Log.Fatalf("Failed to write to temporary file: %v", err)
			}

			// Execute the code block content in the current shell
			cmdExec := exec.Command(blockLanguage, tmpFile.Name())
			cmdExec.Stdout = os.Stdout
			cmdExec.Stderr = os.Stderr
			cmdExec.Stdin = os.Stdin
			if err := cmdExec.Run(); err != nil {
				logger.Log.Fatalf("Failed to execute code block: %v", err)
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the code block to execute")
	cmd.MarkFlagRequired("name")

	return cmd
}
