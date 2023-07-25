# Start with a lightweight base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go server source code into the container
COPY . .

# Build the Go server executable
RUN go build

# Expose the port your Go server listens on
EXPOSE 8080

# Run your Go server when the container starts
CMD ["./asset_service"]