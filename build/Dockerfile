FROM golang:1.13.4-alpine3.10 AS Builder

# The WORKDIR instruction sets the working directory for any RUN, CMD, ENTRYPOINT, COPY and ADD instructions that
# follow it in the Dockerfile
WORKDIR /artemis

# Although ADD and COPY are functionally similar, generally speaking, COPY is preferred. That’s because it’s more
# transparent than ADD. COPY only supports the basic copying of local files into the container, while ADD has some
# features (like local-only tar extraction and remote URL support) that are not immediately obvious.
COPY go.mod go.sum ./
RUN go mod download
COPY . .

### Build Go binary with disabling crossing compilation and only for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o artemis

FROM gcr.io/distroless/base-debian10
COPY --from=Builder /artemis/artemis /artemis
COPY --from=Builder /artemis/configs /configs
ENTRYPOINT ["./artemis"]
