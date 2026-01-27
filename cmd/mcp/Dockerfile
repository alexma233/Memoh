FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETARCH
ARG COMMIT_HASH=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags "-s -w -X main.commitHash=${COMMIT_HASH}" -o /out/mcp ./cmd/mcp

FROM busybox:latest
COPY --from=build /out/mcp /mcp
ENTRYPOINT ["/mcp"]
