FROM golang

MAINTAINER Alice Di Nunno

ADD . ./

RUN ls

RUN make

RUN go build -o main ./cmd

CMD ["./main"]
