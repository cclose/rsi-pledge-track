FROM golang:1.21.3-alpine3.18 as base

# Copy the go.mod
WORKDIR /go/src/
COPY go.mod go.sum ./
RUN go mod download

FROM base as build
COPY . .
RUN go build -tags musl -o rsiPullFunding cmd/rsiPullFunding/main.go
RUN go build -tags musl -o rsiAPI cmd/rsiAPI/main.go

FROM alpine:3.18 as app
## Set the PORT env variable and expose it
ARG PORT="8080"
ENV PORT=$PORT
EXPOSE $PORT

## Build the service
WORKDIR /app
COPY --from=build /go/src/rsiPullFunding /go/src/rsiAPI /go/src/cmd/rsiAPI/wait-for-db.sh /app/

ENTRYPOINT ["/app/rsiAPI"]