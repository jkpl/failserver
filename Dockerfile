FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN apk update && \
    apk add git && \
    go-wrapper download && \
    go-wrapper install

CMD ["go-wrapper", "run"]

EXPOSE 8080