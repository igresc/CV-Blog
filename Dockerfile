FROM golang:1.11
EXPOSE 80
COPY ./src/ /asdf/
RUN ls -la /asdf/*
COPY ./bin/cv /usr/local/bin/
CMD ["cv"]

# FROM alpine:latest
# RUN apk --no-cache add ca-certificates
# WORKDIR /root/
# COPY --from=0 /go/src/app .
# CMD ["./app"]
