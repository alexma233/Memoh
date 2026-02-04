FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
ARG TARGETARCH
ARG COMMIT_HASH=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath -ldflags "-s -w -X github.com/memohai/memoh/internal/version.CommitHash=${COMMIT_HASH}" -o /out/mcp ./cmd/mcp

FROM alpine:latest
RUN apk add --no-cache grep
WORKDIR /app
COPY --from=build /out/mcp /app/mcp
ENTRYPOINT ["/app/mcp"]
