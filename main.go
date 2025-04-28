package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func main() {
	// Define the Falco script path
	falcoScript := "./falco.sh" // Replace with the actual path to your script

	// Ensure the script exists and is executable
	if _, err := os.Stat(falcoScript); os.IsNotExist(err) {
		log.Fatalf("Falco script not found at %s. Please ensure the script exists.", falcoScript)
	}
	if err := os.Chmod(falcoScript, 0755); err != nil {
		log.Fatalf("Failed to make the script executable: %v", err)
	}

	// Start the Falco script
	log.Println("Starting the Falco script...")
	cmd := exec.Command(falcoScript)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the command and check for errors
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start the Falco script: %v", err)
	}

	// Capture the process ID
	falcoPID := cmd.Process.Pid
	log.Printf("Falco script started with PID: %d", falcoPID)

	// Handle signals to stop the Falco script gracefully
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Stopping the Falco script...")
		if err := cmd.Process.Kill(); err != nil {
			log.Printf("Failed to stop the Falco script: %v", err)
		} else {
			log.Println("Falco script stopped successfully.")
		}
		os.Exit(0)
	}()

	// Wait for the script to finish
	if err := cmd.Wait(); err != nil {
		log.Printf("Falco script exited with error: %v", err)
	} else {
		log.Println("Falco script exited successfully.")
	}
}
