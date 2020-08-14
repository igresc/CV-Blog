FROM golang:1.11
EXPOSE 80
COPY /go/ /go/
RUN ls -la /go/*
COPY ./bin/cv /usr/local/bin/
CMD ["cv"]

# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=0 /go/src/app .
# CMD ["./app"]
