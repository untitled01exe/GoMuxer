# Start the Go app build
FROM golang:latest AS build

# Copy source
WORKDIR /go/src/my-golang-source-code
COPY . .

# Get required modules (assumes packages have been added to ./vendor)
RUN go get -d -v ./...

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# New build phase -- create binary-only image
FROM alpine:latest

# Add support for HTTPS
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR src

# Copy files from previous build container
COPY --from=build /go/src/my-golang-source-code/main ./

# Add environment variables
ENV LOGGLY_TOKEN **
ENV AWS_ACCESS_KEY_ID **
ENV AWS_SECRET_ACCESS_KEY **

# Check results
RUN env && pwd && find .

# Start the application
CMD ["./main"]
