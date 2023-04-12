FROM mcr.microsoft.com/vscode/devcontainers/go:1.20-bullseye AS build
ARG version="0.0.0"
WORKDIR /build
COPY ./ ./

RUN --mount=type=cache,target=/home/vscode/.cache/go-build go mod download
RUN go build -o ./frontdoor-ingress -ldflags="-X 'pkg.Version=${version}'" ./cmd/

FROM golang:1.20-bullseye

RUN adduser -u 1001 appuser
USER 1001
WORKDIR /app
COPY --from=build /build/frontdoor-ingress ./

ENTRYPOINT [ "/app/frontdoor-ingress" ]