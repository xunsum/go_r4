FROM alpine
FROM golang:1.20.1

COPY . ./opt
WORKDIR ./opt

ENV GOPROXY=https://goproxy.cn

RUN go build -o main ./main/main.go

CMD ["/opt/main/main"]

EXPOSE 9060