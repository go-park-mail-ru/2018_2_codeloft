FROM ubuntu:18.04

# Обвновление списка пакетов
RUN apt-get -y update && apt-get -y upgrade

ARG USERNAME
ARG PASSWORD
ENV USERNAME $USERNAME
ENV PASSWORD $PASSWORD

ENV GOVER 1.10
RUN apt-get install -y golang-$GOVER

# Выставляем переменную окружения для сборки проекта
ENV GOROOT /usr/lib/go-$GOVER
ENV GOPATH /opt/go
ENV PATH $GOROOT/bin:$GOPATH/bin:/usr/local/go/bin:$PATH

# Копируем исходный код в Docker-контейнер
WORKDIR $GOPATH/src/github.com/go-park-mail-ru/2018_2_codeloft/game
COPY . $GOPATH/src/github.com/go-park-mail-ru/2018_2_codeloft/game

RUN go install .
EXPOSE 8082

CMD game $USERNAME $PASSWORD