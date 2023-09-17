FROM golang:latest as builder
WORKDIR /build
COPY ./ ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main -v -ldflags="-X 'github.com/installer/instl/internal/pkg/global.Version=`git log -1 --pretty=format:'%h - %ci'`'"


FROM alpine
WORKDIR /app
COPY --from=builder /build/main ./main
ADD static /app/static
EXPOSE 80
ENTRYPOINT ["./main"]
