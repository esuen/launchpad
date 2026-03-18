FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /bin/launchpad ./cmd/server

FROM alpine:3.21

RUN addgroup -S app && adduser -S app -G app

COPY --from=build /bin/launchpad /bin/launchpad

USER app

EXPOSE 8080

ENTRYPOINT ["/bin/launchpad"]
