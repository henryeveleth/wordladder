FROM golang:latest
LABEL maintainer="Henry Eveleth <henryeveleth@gmail.com>"
WORKDIR $GOPATH/src/wordladder-app
COPY ./ $GOPATH/src/wordladder-app
RUN go get
RUN go build .
EXPOSE 8080
CMD ["./wordladder-app"]
