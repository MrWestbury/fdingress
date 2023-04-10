FROM mcr.microsoft.com/vscode/devcontainers/go:1.20-bullseye AS build

WORKDIR /build
COPY ./ ./

RUN go mod download
RUN go build -o ./frontdoor-ingress ./cmd/

FROM golang:1.20-bullseye

WORKDIR /app
COPY --from=build /build/frontdoor-ingress ./

RUN adduser -u 1001 appuser
USER 1001

ENTRYPOINT [ "/app/frontdoor-ingress" ]