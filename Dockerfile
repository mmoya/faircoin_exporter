FROM        frolvlad/alpine-glibc
MAINTAINER  Maykel Moya <mmoya@mmoya.org>

COPY        faircoin_exporter /bin/faircoin_exporter

EXPOSE      9132
ENTRYPOINT  [ "/bin/faircoin_exporter" ]
