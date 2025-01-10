package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pkg "runblock/pkg"
	logger "runblock/pkg/logger"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func RunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Interactively select and run a code block from a Markdown file",
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

			var codeBlockOptions []huh.Option[string]
			// Create a list of options for the select form
			// fullDesc will be displayed, but name will be used as the value
			for name, content := range codeBlocks {
				fullDesc := fmt.Sprintf("%s - %s", name, content.Attributes.Description)
				codeBlockOptions = append(codeBlockOptions, huh.NewOption(fullDesc, name))
			}

			// Handle Ctrl+C to gracefully exit the loop
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				<-sigs
				os.Exit(0)
			}()

			for {
				// Use huh to create an interactive select form
				var selectedName string
				form := huh.NewForm(
					huh.NewGroup(
						huh.NewSelect[string]().
							Title("Select a code block to run:").
							Options(codeBlockOptions...).
							Value(&selectedName),
					),
				)

				if err := form.Run(); err != nil {
					logger.Log.Fatalf("Failed to run form: %v", err)
				}

				blockContent := codeBlocks[selectedName].Content
				blockLanguage := codeBlocks[selectedName].Language

				if blockContent == "" {
					logger.Log.Fatalf("Code block for '%s' is empty", selectedName)
				}

				// Execute the code block
				pkg.Exec(blockLanguage, blockContent)

				fmt.Printf("--- Code block finished ---\n")
				fmt.Printf("Press Ctrl+C to exit or any key to run another code block\n")
				fmt.Scanln()
			}
		},
	}

	return cmd
}
