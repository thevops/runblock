package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"runblock/pkg/logger"
)

func Exec(interpreter string, code string) {
	// Save the code block content to a temporary file
	fileNamePattern := fmt.Sprintf("runblock-*.%s", interpreter)
	tmpFile, err := os.CreateTemp("", fileNamePattern)
	if err != nil {
		logger.Log.Fatalf("Failed to create temporary file: %v", err)
	}
	logger.Log.Debug("Created temporary file: ", tmpFile.Name())
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(code); err != nil {
		logger.Log.Fatalf("Failed to write to temporary file: %v", err)
	}

	cmdExec := exec.Command(interpreter, tmpFile.Name())
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr
	cmdExec.Stdin = os.Stdin
	if err := cmdExec.Run(); err != nil {
		logger.Log.Fatalf("Failed to execute code block: %v", err)
	}
}
