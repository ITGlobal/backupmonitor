FROM golang:alpine AS build

RUN apk update && \
    apk add --no-cache git gcc musl-dev

WORKDIR /go/src/github.com/itglobal/backupmonitor

RUN go get -u github.com/swaggo/swag/cmd/swag

COPY . .
RUN mkdir -p /out/doc
RUN go get
RUN swag init --output /out/doc/ --generalInfo swagger.go --dir ./pkg/api/
RUN go build -o /out/backupmonitor
COPY ./doc /out/

FROM alpine:latest
WORKDIR /app
COPY --from=build /out/backupmonitor /app/backupmonitor
COPY --from=build /out/doc /app/

ENTRYPOINT [ "/app/backupmonitor" ]
