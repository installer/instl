FROM golang:latest as builder
WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main


FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add curl
COPY --from=builder /build/main ./main
ADD ./scripts /scripts
EXPOSE 80
HEALTHCHECK CMD curl --fail http://localhost:80 || exit 1
ENTRYPOINT ["./main"]
