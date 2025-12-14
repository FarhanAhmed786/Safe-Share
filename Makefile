# Makefile for Safe Share Server
# Handles SSL certificate generation and application execution.

# Default target executed when running 'make' without arguments
all: run

# ==============================================================================
# TARGET: run
# Description: Checks for SSL certificates and starts the Go application.
# Dependencies: server.crt, server.key
# ==============================================================================
run: server.crt server.key
	@echo "Starting Safe Share Server..."
	go run .

# ==============================================================================
# TARGET: server.crt / server.key
# Description: Generates a self-signed SSL certificate and private key.
# Triggered automatically if either file is missing.
# ==============================================================================
server.crt server.key:
	@echo "SSL certificates not found. Generating new self-signed certificates..."
	@openssl req -x509 -newkey rsa:4096 \
		-keyout server.key \
		-out server.crt \
		-days 365 \
		-nodes \
		-subj "/CN=localhost" 2>/dev/null
	@echo "Certificates generated successfully."

# ==============================================================================
# TARGET: clean
# Description: Removes the generated SSL certificates.
# Use this to reset the environment or force certificate regeneration.
# ==============================================================================
clean:
	@echo "Removing SSL certificates..."
	rm -f server.crt server.key
	@echo "Cleanup complete."

# .PHONY ensures that 'run' and 'clean' are treated as commands, 
# not files, preventing conflicts if files with these names exist.
.PHONY: all run clean
