################################################################################
ARG GO_VERSION=1.22.0
FROM --platform=$BUILDPLATFORM amd64/golang:${GO_VERSION}-alpine AS build-backend
WORKDIR /src

RUN apk --no-cache update && \
apk --no-cache add git gcc libc-dev

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH
ARG TARGETOS

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=1 GOARCH=$TARGETARCH  GOOS=$TARGETOS go build -tags musl -ldflags "-s -w -X main.Version=$GIT_TAG_NAME"  -o /bin/server ./cmd/server
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

RUN mkdir /msg
COPY --from=build-backend /bin/server /msg
COPY config/config.yaml /msg

RUN apk --no-cache update && \
apk --no-cache add git gcc libc-dev


EXPOSE 8888
RUN chmod a+x /msg/server

ENTRYPOINT ["/msg/server"]