FROM golang:latest as builder
WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main


FROM alpine:latest
WORKDIR /app
COPY --from=builder /build/main ./main
ADD ./scripts /scripts
EXPOSE 80
ENTRYPOINT ["./main"]
