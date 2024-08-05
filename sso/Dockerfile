################################################################################
ARG GO_VERSION=1.22.0
FROM --platform=$BUILDPLATFORM amd64/golang:${GO_VERSION}-alpine AS build-sso
WORKDIR /src

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOARCH=$TARGETARCH GOOS=linux go build -o /bin/server ./cmd/server
################################################################################
FROM amd64/alpine:latest AS final

RUN --mount=type=cache,target=/var/cache/apk \
    apk --update add \
        ca-certificates \
        tzdata \
        && \
        update-ca-certificates

ARG UID=10001
USER root

RUN mkdir /sso
COPY --from=build-sso /bin/server /sso
COPY config/config.yaml /sso

EXPOSE 9999

ENTRYPOINT ["/sso/server"]