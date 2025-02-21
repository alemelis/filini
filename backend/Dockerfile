# Use official Golang image
FROM golang:1.23 AS build

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
# Copy the rest of the source code
COPY . .
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o filini ./cmd/main.go

# Use a lightweight image for production
FROM debian:bookworm-slim
WORKDIR /app

# Install FFmpeg (for GIF processing)
RUN apt update && apt install -y ffmpeg

# Copy the built binary from the previous stage
COPY --from=build /app/filini /usr/local/bin/filini
COPY --from=build /app/.env .
RUN chmod +x /usr/local/bin/filini

# Expose the port your app runs on
EXPOSE 8080

# Run Filini
CMD ["/usr/local/bin/filini"]
