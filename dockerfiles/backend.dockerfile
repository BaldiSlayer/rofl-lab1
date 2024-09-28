FROM golang:1.23-bullseye AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY /cmd/backend /app/cmd/backend
COPY /internal/app/backend /app/internal/app/backend

RUN go build -o /bin/backend /app/cmd/backend/backend.go

FROM ubuntu:20.04

COPY --from=build /bin/backend /bin/backend

CMD ["/bin/backend"]