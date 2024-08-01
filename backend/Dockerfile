################################################################################
ARG GO_VERSION=1.22.0
FROM --platform=$BUILDPLATFORM amd64/golang:${GO_VERSION}-alpine AS build-backend
WORKDIR /src

#RUN #apk add clang lld #apk add --no-cache libgit2 libgit2-dev git gcc g++ pkgconfig

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH
ARG TARGETOS

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
     GOARCH=$TARGETARCH  GOOS=$TARGETOS go build  -o /bin/server ./cmd/server
# \CGO_ENABLED=0 GO111MODULE=on \
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

EXPOSE 8888
RUN chmod a+x /messager/server

ENTRYPOINT ["/messager/server"]
#ENTRYPOINT ["top"]