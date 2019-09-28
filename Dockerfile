FROM resin/rpi-raspbian:latest

RUN apt-get -q update \
  && apt-get upgrade -qy

EXPOSE 8080

ADD pi-relay /

CMD ["/pi-relay"]
