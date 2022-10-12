## Build
FROM golang:1.19-alpine as build

WORKDIR /usr/src/bot

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN CGO_ENABLED=0 go build -v -o /bot github.com/NEKETSKY/footg-bot/cmd/footg-bot

## Deploy
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /bot /
COPY config.yaml /

EXPOSE 8080

USER nonroot:nonroot

CMD ["/bot"]