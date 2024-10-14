FROM golang:1.23-bullseye AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY /cmd/kbase /app/cmd/kbase
COPY /internal/app /app/internal/app

RUN go build -o /bin/kbase /app/cmd/kbase/kbase.go

FROM ubuntu:20.04

COPY --from=build /bin/kbase /bin/kbase

CMD ["/bin/bash"]
