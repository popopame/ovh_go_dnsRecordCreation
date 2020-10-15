FROM golang:alpine AS builder

#Set necessary env file, you can set the needed argument for the script here if you wish
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

# Copy and download dependency
COPY go.mod .
COPY go.sum .
RUN go mod download 
# copy the code into the container
COPY . . 

# Build the app
RUN go build -o script .

# Move the file output to another folder
WORKDIR /output
RUN cp /build/script .

#Build from scratch as we only need the binary
FROM alpine

COPY --from=builder /build/script /

# Command to run at contianer startup
ENTRYPOINT [ "/script" ]
