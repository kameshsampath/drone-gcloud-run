#syntax=docker/dockerfile:1.3-labs

FROM golang:1.18-alpine AS builder
ENV CGO_ENABLED=0
RUN apk add --update make git bash\
    && mkdir /build
WORKDIR /build
COPY go.* .
RUN go install github.com/goreleaser/goreleaser@latest 
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make bin-all

FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM alpine:3.6
ENV GODEBUG netdns=go
ARG TARGETARCH
ARG DOCKER_CREDENTIAL_GCR_VERSION=2.1.6

RUN apk add --update --no-cache jq curl

COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build /build

RUN <<EOT
    jq -r --arg target_arch $TARGETARCH  '.[] | select(.goarch==$target_arch) | .path' /build/dist/artifacts.json > /build/pluginfile.txt
    cp $(cat /build/pluginfile.txt) /bin/plugin
    chmod +x /bin/plugin
EOT

## docker gcr credential helper.
## https://cloud.google.com/artifact-registry/docs/docker/authentication#standalone-helper
RUN curl -fsSL "https://github.com/GoogleCloudPlatform/docker-credential-gcr/releases/download/v${DOCKER_CREDENTIAL_GCR_VERSION}/docker-credential-gcr_linux_${TARGETARCH}-${DOCKER_CREDENTIAL_GCR_VERSION}.tar.gz" \
| tar xz docker-credential-gcr \
&& chmod +x docker-credential-gcr && mv docker-credential-gcr /usr/local/bin

RUN rm -rf /build

CMD ["/bin/plugin"]