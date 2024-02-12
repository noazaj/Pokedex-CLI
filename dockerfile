# Create the base image
FROM golang:1.21

# Create and change to current working directory
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Build the docker image for linux
RUN CGO_ENABLED=0 GOOS=linux go build -o /pokedex-cli

# Run the docker image CMD
CMD [ "/pokedex-cli" ]