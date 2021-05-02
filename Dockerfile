# -----------------------------------------------------------------------------
# Build backend
# -----------------------------------------------------------------------------
FROM golang:alpine AS backend

RUN apk update && \
    apk add --no-cache git gcc musl-dev

WORKDIR /src

RUN go get -u github.com/swaggo/swag/cmd/swag@v1.6.7
COPY go.mod /src/go.mod
COPY go.sum /src/go.sum
RUN go mod download

COPY . /src
RUN mkdir -p /out/doc
RUN swag init --output /out/doc/ --generalInfo swagger.go --dir ./pkg/api/
RUN go build -o /out/backupmonitor
COPY ./doc /out/

# -----------------------------------------------------------------------------
# Build frontend
# -----------------------------------------------------------------------------
FROM node:latest AS frontend

WORKDIR /app/src
COPY ./client/package.json /app/src
COPY ./client/package-lock.json /app/src
RUN npm install

COPY ./client/ /app/src
RUN npm install
RUN npm run build

# -----------------------------------------------------------------------------
# Build runtime image
# -----------------------------------------------------------------------------
FROM alpine:latest
WORKDIR /app
COPY --from=backend /out/backupmonitor /app/backupmonitor
COPY --from=backend /out/doc/ /app/doc/
COPY --from=frontend /app/www/ /app/www/

ENTRYPOINT [ "/app/backupmonitor" ]
