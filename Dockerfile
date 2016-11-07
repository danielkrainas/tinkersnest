FROM frolvlad/alpine-glibc
MAINTAINER Danny Krainas <me@danielkrainas.com>

ENV TINKERS_CONFIG_PATH /etc/tinkersnest.default.yml

COPY ./dist /bin/tinkersnest
COPY ./tinkerctl/dist /bin/tinkerctl
COPY ./config.default.yml /etc/tinkersnest.default.yml

ENTRYPOINT ["/bin/tinkersnest"]
