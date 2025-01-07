#!/bin/bash

# Define the log file location
LOG_FILE=""

# Check if Falco is installed
if ! command -v falco &> /dev/null; then
    echo "Falco is not installed. Please install it first."
    exit 1
fi

# Ensure the script has proper permissions to write to the log file
if ! touch "$LOG_FILE" &> /dev/null; then
    echo "Error: Cannot write to $LOG_FILE. Check file permissions."
    exit 1
fi

# Start Falco and redirect its output to the log file in real-time
echo "Starting Falco and logging alerts to $LOG_FILE..."
sudo stdbuf -oL falco -o output_format=json | stdbuf -oL tee "$LOG_FILE" &
falco_pid=$!

# Function to clean up and stop Falco when the script exits
cleanup() {
    echo "Stopping Falco..."
    kill "$falco_pid" &> /dev/null
    exit 0
}

# Trap signals (SIGINT, SIGTERM) to trigger the cleanup function
trap cleanup SIGINT SIGTERM

# Continuously monitor the log file in real-time
echo "Monitoring Falco alerts. Press Ctrl+C to stop."
tail -f "$LOG_FILE" &

# Wait for the Falco process to terminate (or until the script is stopped)
wait "$falco_pid"
