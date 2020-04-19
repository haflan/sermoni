### Builder stage
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git npm
COPY . $GOPATH/src/sermoni/

# Build website files and move html.go
WORKDIR $GOPATH/src/sermoni/ui/
RUN npm install; \
    npm run build; \
    mv $GOPATH/src/sermoni/ui/dist/html.go $GOPATH/src/sermoni/internal/http/
# Build the sermoni binary
WORKDIR $GOPATH/src/sermoni/
RUN go get -d ./... ; \
    GOOS=linux GOOARCH=amd64 go build \
        -ldflags="-w -s" \
        -o /go/bin/sermoni \
        -tags PRODUCUTION \
        ./cmd/sermoni/

### Production image
FROM scratch
COPY --from=builder /go/bin/sermoni /sermoni
RUN mkdir -p /data
ENTRYPOINT ["/sermoni", "-d", "/data/sermoni.db"]
