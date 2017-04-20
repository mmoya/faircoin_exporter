FROM        frolvlad/alpine-glibc
MAINTAINER  Maykel Moya <mmoya@mmoya.org>

COPY        faircoin2_exporter /bin/faircoin2_exporter

EXPOSE      9132
ENTRYPOINT  [ "/bin/faircoin2_exporter" ]
