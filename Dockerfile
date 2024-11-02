# Use the official Go base image
FROM golang:1.22.2 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download 

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# # Use a smaller base image for the final image
# FROM alpine:latest

# # Set the working directory
# WORKDIR /app

# # Copy the built binary from the builder stage
# COPY --from=builder /app/main .

# RUN ls -la

# # Copy the Firebase credentials file (if needed)
# # COPY path/to/your/keyfile.json . 

# # Expose the port your application listens on
# EXPOSE 8080

# # Set the entrypoint to run your application
CMD ["./main"]