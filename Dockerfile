FROM golang:1.9.3-stretch

WORKDIR /app

ADD . /app

RUN GOOS=linux go build wiki.go
RUN chmod +x wiki

EXPOSE 80

CMD ["./wiki"]
