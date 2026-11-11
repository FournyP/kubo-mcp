FROM --platform=$BUILDPLATFORM golang:1.27-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev

WORKDIR /usr/local/go/src/github.com/FournyP/kubo-mcp

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./ ./

WORKDIR /usr/local/go/src/github.com/FournyP/kubo-mcp/src/kubo_mcp_host

# Cross-compiles instead of emulating the target platform.
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -trimpath \
    -ldflags "-s -w -X github.com/FournyP/kubo-mcp/src/kubo_mcp_host/configurations.Version=${VERSION}" \
    -o /app .

FROM scratch AS final

COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs /etc/ssl/certs

ENV HOST__PORT=3000

EXPOSE 3000

USER 65534:65534

ENTRYPOINT ["/app"]
