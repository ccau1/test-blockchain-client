FROM golang:1.20.3-alpine as dev
# install air for hot-reload
RUN go install github.com/cosmtrek/air@latest

# copy files from local to vm
WORKDIR /app
COPY . /app/

# install mod dependencies
RUN go mod download

# build app into bin
RUN CGO_ENABLED=0 go build -o /go/bin/app

# ----------------------------------------------------------------------
# a small image to get wget
FROM busybox AS wget

ARG BUSYBOX_VERSION=1.31.0-i686-uclibc
ADD https://busybox.net/downloads/binaries/$BUSYBOX_VERSION/busybox_WGET /wget
RUN chmod a+x /wget

# ----------------------------------------------------------------------

FROM gcr.io/distroless/static-debian11 as prod

# install curl for health check
COPY --from=wget /wget /usr/bin/wget

# only copy files from bin over
COPY --from=dev go/bin/app /

COPY --from=dev app/docs/docgen.md /docs/docgen.md

# calling app will run the executable app
CMD ["/app"]