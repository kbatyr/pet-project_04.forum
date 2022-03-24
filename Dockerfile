#syntax=docker/dockerfile:1

FROM golang:1.17.6

WORKDIR /forum

COPY . ./

RUN go build -o /docker-forum ./cmd/web

EXPOSE 8080

CMD [ "/docker-forum" ]