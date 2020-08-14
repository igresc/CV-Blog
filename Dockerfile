FROM golang:latest
WORKDIR /go/src/
RUN go get encoding/json fmt html/template io/ioutil net/http
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 ./app .
CMD ["./app"]
