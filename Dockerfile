FROM golang:1.10.3-alpine3.7 as builder

COPY . /go/src/github.com/go-park-mail-ru/2018_2_codeloft
ENV gopath /go
RUN cd /go/src/github.com/go-park-mail-ru/2018_2_codeloft && go build -o goapp


FROM ubuntu:18.04

#COPY --from=builder /go/src/2018_2_codeloft/.env /app/
#CMD source /app/.env

# Обвновление списка пакетов
ARG USERNAME
ARG PASSWORD
RUN apt-get -y update
ENV USERNAME $USERNAME
ENV PASSWORD $PASSWORD
#
# Установка postgresql
#
ENV PGVER 10
RUN apt-get install -y postgresql-$PGVER

# Run the rest of the commands as the ``postgres`` user created by the ``postgres-$PGVER`` package when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER $USERNAME WITH SUPERUSER PASSWORD '$PASSWORD';" &&\
    createdb -O codeloft codeloft &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
#VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

#FROM alpine
WORKDIR /app
COPY --from=builder /go/src/github.com/go-park-mail-ru/2018_2_codeloft/goapp /app/
#COPY .env .
#RUN source ./.env
EXPOSE 8080
CMD service postgresql start && ./goapp

#CMD ["2018_2_codeloft"]

#FROM scratch
#COPY --from=builder /app ./
#EXPOSE 8080
#ENTRYPOINT ["./app"]


#CMD ["2018_2_codeloft"]

