FROM golang:1.14

WORKDIR /go/src/github.com/igresc/cv-blog

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN pwd
RUN ls -la .

EXPOSE 80

CMD ["cv-blog"]
