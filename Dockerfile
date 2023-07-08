# Start from the specific golang base image
FROM golang:1.21rc2-bullseye AS builder

# Add Maintainer Info
LABEL maintainer="Nikolay Stankov <babati@duck.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first for layers caching
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Now copy the rest of the source code. If any changes are made in the source code, 
# Docker will invalidate the cache for this step and all subsequent steps.
# This is the "cache busting" step.
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./

# Start a new stage from scratch
FROM gcr.io/distroless/base-debian11

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Command to run the executable
ENTRYPOINT ["./main"]

# As a best practice, we don't EXPOSE the port in the Dockerfile, as it can lead to 
# misunderstandings about the function of the EXPOSE directive. The port binding should
# be specified at runtime using the -p flag.
