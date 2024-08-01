################################################################################
ARG GO_VERSION=1.22.0
FROM --platform=$BUILDPLATFORM amd64/golang:${GO_VERSION}-alpine AS build-backend
WORKDIR /src

#RUN #apk add clang lld #apk add --no-cache libgit2 libgit2-dev git gcc g++ pkgconfig
#ARG GIT_TAG_NAME

RUN apk --no-cache update && \
apk --no-cache add git gcc libc-dev

#ENV CGO_ENABLED 1
#ENV GOFLAGS -mod=vendor
#ENV GOOS=linux
#ENV GOARCH=amd64
#
#RUN export GO111MODULE=on

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

ARG TARGETARCH
ARG TARGETOS

RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=1 GOARCH=$TARGETARCH  GOOS=$TARGETOS go build -tags musl -ldflags "-s -w -X main.Version=$GIT_TAG_NAME"  -o /bin/server ./cmd/server
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

RUN apk --no-cache update && \
apk --no-cache add git gcc libc-dev


EXPOSE 8888
RUN chmod a+x /msg/server

ENTRYPOINT ["/msg/server"]
#ENTRYPOINT ["top"]