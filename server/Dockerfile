######################################################################
#### Stage 1: Compile the binary in a containerized Golang environment
######################################################################

# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.18.3-alpine3.16 as builder

# Create and change to our app directory.
WORKDIR /chess

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary!
# RUN CGO_ENABLED=0 GOOS=linux go build -o chess-server
RUN go build -v -o chess-server

######################################################################
#### Stage 2: Build the Chess-Server image proper
######################################################################

# Use a "scratch" image, which contains no distribution files
FROM scratch as image

# Copy the binary from the build container
COPY --from=builder /chess/chess-server /chess/chess-server

# Tell Docker we'll be using port 8080
EXPOSE 8080

# Tell Docker to execute this command on a "docker run"
CMD ["/chess/chess-server"]