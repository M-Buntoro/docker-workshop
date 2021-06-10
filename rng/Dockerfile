FROM golang:alpine3.13
WORKDIR /
COPY . .
RUN go build -o /rng .
ENTRYPOINT /rng
EXPOSE 11991