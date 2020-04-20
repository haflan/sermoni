### Builder stage
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git npm
COPY . $GOPATH/src/sermoni/

# Build website files and move html.go
WORKDIR $GOPATH/src/sermoni/ui/
RUN npm install; \
    npm run build; \
    $GOPATH/src/sermoni/ui/generate.sh ; \
    mv $GOPATH/src/sermoni/ui/dist/html.go $GOPATH/src/sermoni/internal/http/
# Build the sermoni binary
WORKDIR $GOPATH/src/sermoni/
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go get -d ./... ; \
    go build \
        -ldflags="-w -s" \
        -o /go/bin/sermoni \
        -tags PRODUCTION \
        ./cmd/sermoni/
# Empty directory for database
RUN mkdir -p /data
### Production image
FROM scratch
COPY --from=builder /data /data
COPY --from=builder /go/bin/sermoni /sermoni
ENTRYPOINT ["/sermoni", "-d", "/data/sermoni.db"]
