FROM golang:alpine3.13
WORKDIR /
COPY . .
RUN go build -o /hash .
ENTRYPOINT /hash
EXPOSE 11992