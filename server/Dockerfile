FROM golang:1.15.0-alpine3.12
WORKDIR /app
COPY . /app
RUN go get github.com/cespare/reflex
EXPOSE 8080
CMD ["./hot-reload.sh"]