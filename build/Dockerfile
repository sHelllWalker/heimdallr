FROM --platform=$BUILDPLATFORM golang:1.24.0-alpine3.21 AS build

WORKDIR /app

COPY ./go.* .

RUN go mod download

COPY . .

RUN go build -ldflags="-s" -o /app/heimdallr ./cmd/heimdallr

FROM --platform=$BUILDPLATFORM alpine:3.21

WORKDIR /app
COPY --from=build /app/heimdallr /app

CMD ["/app/heimdallr"]
