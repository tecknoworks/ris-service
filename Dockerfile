# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go get -u github.com/go-http-utils/logger && go get -u github.com/PuerkitoBio/goquery

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main ./src

# Expose port 3003 to the outside world
EXPOSE 3006

# Command to run the executable
CMD ["./main"]