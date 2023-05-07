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

FROM gcr.io/distroless/static-debian11 as prod

# only copy files from bin over
COPY --from=dev go/bin/app /

COPY --from=dev app/docs/docgen.md /docs/docgen.md

# calling app will run the executable app
CMD ["/app"]