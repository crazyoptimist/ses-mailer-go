BINARY_NAME=mailer
BUILD_DIR=build
SRC_PATH=./cmd/mailer

.PHONY: run forlinux forwindows clean

# Run the application
run:
	go run $(SRC_PATH)/*.go

# Build for Linux (AMD64)
forlinux: clean
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_PATH)

# Build for Windows (AMD64)
forwindows: clean
	mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(SRC_PATH)

# Clean up build artifacts
clean:
	rm -rf $(BUILD_DIR)/*
