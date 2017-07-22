FROM        frolvlad/alpine-glibc
MAINTAINER  Maykel Moya <mmoya@mmoya.org>

RUN         apk --no-cache add libzmq
COPY        faircoin_exporter.docker /bin/faircoin_exporter

EXPOSE      9132
ENTRYPOINT  [ "/bin/faircoin_exporter" ]
