package cmd

import (
	"os"

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

			// Get block by name
			namedCodeBlock, err := pkg.GetNamedCodeBlock(codeBlocks, name)
			if err != nil {
				logger.Log.Fatalf("Failed to get named code block: %v", err)
			}

			pkg.Exec(namedCodeBlock.Language, namedCodeBlock.Content)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the code block to execute")
	cmd.MarkFlagRequired("name")

	return cmd
}
