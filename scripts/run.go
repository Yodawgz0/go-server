// scripts/run.go

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Run your custom tasks
	startServer()
}

func startServer() {
	fmt.Println("Starting the server...")

	// Example: Run the main.go file
	cmd := exec.Command("go", "run", "main.go", "services.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
