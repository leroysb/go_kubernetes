FROM golang:1.22

# Set the Current Working Directory inside the container
WORKDIR /usr/src/app

# Download Go modules
COPY go.mod  ./
COPY go.sum ./
RUN go mod download && go mod verify

# Copy the source code into the container
COPY . .

# RUN go mod tidy

# Navigate to the cmd directory
# WORKDIR /usr/src/app/cmd

# build the binary
# RUN go build -v -o /usr/local/bin/app ./...

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
# CMD ["./api"]
