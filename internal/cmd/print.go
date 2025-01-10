package cmd

import (
	"fmt"
	"os"

	pkg "runblock/pkg"
	logger "runblock/pkg/logger"

	"github.com/spf13/cobra"
)

func PrintCmd() *cobra.Command {
	var (
		name    string
		details bool
	)

	cmd := cobra.Command{
		Use:   "print",
		Short: "Print a code block by name from a Markdown file",
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
				if name != "" && blockName == name {
					if details {
						fmt.Printf("Name: %s\n", blockName)
						fmt.Printf("Language: %s\n", blockContent.Language)
						if len(blockContent.Attributes.Description) == 0 {
							fmt.Printf("Description: none\n")
						} else {
							fmt.Printf("Description: %s\n", blockContent.Attributes.Description)
						}
						fmt.Printf("---\n")
					}
					fmt.Printf("%s", blockContent.Content)
				}
			}
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Name of the code block to execute")
	cmd.MarkFlagRequired("name")
	cmd.Flags().BoolVar(&details, "details", false, "Include details in the output")

	return &cmd
}
