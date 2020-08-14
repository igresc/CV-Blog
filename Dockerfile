FROM golang:latest
WORKDIR /go/src/
RUN go get encoding/json fmt html/template io/ioutil net/http
COPY main.go .
RUN go build -o cv
CMD ["cv"]

# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=0 /go/src/app .
# CMD ["./app"]
