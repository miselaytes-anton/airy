# Compile stage
FROM golang:alpine AS build-env
ENV CGO_ENABLED 0

COPY . /app_src
WORKDIR /app_src
RUN go build -gcflags "all=-N -l" -o /app ./cmd/processor

# Final stage
FROM alpine:latest

COPY --from=build-env /app /

# Run
CMD ["/app"]