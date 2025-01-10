package cmd

import (
	"fmt"
	"os"

	pkg "runblock/pkg"
	logger "runblock/pkg/logger"

	"github.com/spf13/cobra"
)

func ListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all code blocks from a Markdown file",
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

			for blockName, blockContent := range codeBlocks {
				if len(blockContent.Attributes.Description) == 0 {
					fmt.Printf("%s\n", blockName)
				} else {
					fmt.Printf("%s - %s\n", blockName, blockContent.Attributes.Description)
				}
			}
		},
	}

	return cmd
}
