FROM alpine:3.10

EXPOSE 8080

ADD . /

CMD ["/py-relay"]
