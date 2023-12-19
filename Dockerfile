# Compile stage
FROM golang:alpine AS build-env
ENV CGO_ENABLED 0

COPY . /app_src
WORKDIR /app_src
RUN go build -gcflags "all=-N -l" -o /processor ./cmd/processor
RUN go build -gcflags "all=-N -l" -o /server ./cmd/server

# Final stage
FROM alpine:latest

COPY --from=build-env /server /
COPY --from=build-env /processor /
