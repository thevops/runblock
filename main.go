package main

import (
	"os"

	"github.com/spf13/cobra"

	cmd "runblock/internal/cmd"
	logger "runblock/pkg/logger"
)

var filePath string

func main() {
	rootCmd := &cobra.Command{
		Use: "runblock",
	}

	// Define the global file flag
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "Path to the Markdown file")
	rootCmd.MarkPersistentFlagRequired("file")

	// Add the print command
	rootCmd.AddCommand(cmd.ListCmd())
	rootCmd.AddCommand(cmd.PrintCmd())
	rootCmd.AddCommand(cmd.ExecCmd())
	rootCmd.AddCommand(cmd.RunCmd())

	if err := rootCmd.Execute(); err != nil {
		logger.Log.Fatalf("Error executing command: %v", err)
		os.Exit(1)
	}
}
