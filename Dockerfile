FROM frolvlad/alpine-glibc
MAINTAINER Danny Krainas <me@danielkrainas.com>

ENV CSENSE_CONFIG_PATH /etc/tinkersnest.default.yml

COPY ./dist /bin/tinkersnest
COPY ./config.default.yml /etc/tinkersnest.default.yml

ENTRYPOINT ["/bin/tinkersnest"]
