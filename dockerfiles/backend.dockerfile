FROM public.ecr.aws/docker/library/golang:1.23-bullseye AS build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY /cmd/backend /app/cmd/backend
COPY /internal /app/internal
COPY /pkg /app/pkg

ARG LDFLAGS

RUN go build -o /bin/backend -ldflags="$LDFLAGS" /app/cmd/backend/backend.go

FROM public.ecr.aws/ubuntu/ubuntu:20.04

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=build /bin/backend /bin/backend

COPY /data /data

CMD ["/bin/backend"]
