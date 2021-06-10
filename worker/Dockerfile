FROM golang:alpine3.13
WORKDIR /
COPY . .
RUN go build -o /worker .
ENTRYPOINT /worker
EXPOSE 11993