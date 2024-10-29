# Use the official Go image to build the Go code
FROM golang:1.20 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go source code (assuming it's in the same directory as the Dockerfile)
COPY peripheral.go .

# Initialize a new Go module within the Dockerfile
RUN go mod init example.com/my-peripheral

# Download the necessary dependencies
RUN go mod tidy

# Build the Go binary
RUN go build -o peripheral .

# Use a minimal base image for the final container
FROM debian:bookworm-slim

# Install necessary packages for Bluetooth and `tini`
RUN apt-get update && apt-get install -y \
    bluetooth \
    bluez \
    dbus \
    # tini \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary from the builder stage
COPY --from=builder /app/peripheral /app/peripheral

# Ensure `tini` is the entry point
# ENTRYPOINT ["/usr/bin/tini", "--"]

# Run the services and the peripheral code
CMD ["sh", "-c", "/bin/bash"]
#CMD ["sh", "-c", "service dbus start && bluetoothd & hciconfig hci0 down && hciconfig hci0 up && /bin/bash"]

