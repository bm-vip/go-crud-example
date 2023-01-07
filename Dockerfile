# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.19-alpine

# Create and change to the app directory.
WORKDIR /app/go-crud-example

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the Go app
RUN go build -o ./out/go-crud-example .

# This container exposes port 3000 to the outside world
EXPOSE 3000
# Run the binary program produced by `go install`
CMD ["./out/go-crud-example"]