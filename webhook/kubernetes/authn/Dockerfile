# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Sebastien Cabarrus <scabarrus@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /go/src

# Copy the source from the current directory to the Working Directory inside the container
COPY ./k8s.scabarrus.com ./k8s.scabarrus.com

RUN cd k8s.scabarrus.com/webhook ; go mod init; go mod vendor

# Build the Go app
RUN go build -o main  k8s.scabarrus.com/webhook/cmd/authn


# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./main"]

CMD ["-protocol=http","-port=8080"]