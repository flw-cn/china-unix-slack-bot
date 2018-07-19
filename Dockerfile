FROM ubuntu
LABEL name="slack-bot"
LABEL maintainer="flw <flw@cpan.org>"

# Layer 1: Base OS
RUN sed -i 's/archive.ubuntu.com/mirrors.ustc.edu.cn/g' /etc/apt/sources.list
RUN apt update
RUN apt install -y locales locales-all
RUN apt install -y docker.io golang
RUN apt install -y fortune fortune-zh
RUN apt clean

ENV LANG en_US.UTF8
ENV TERM xterm-256color

# Layer 2: Base Golang development environment
ENV GOPATH "/go"
ENV PATH "$PATH:/go/bin"
RUN go get -v github.com/sirupsen/logrus
RUN go get -v github.com/flw-cn/go-smartConfig

# Layer 3: Slack bot development environment
RUN go get -v github.com/flw-cn/slack
RUN go get -v github.com/flw-cn/go-slackbot
RUN go get -v github.com/flw-cn/slack-bot

WORKDIR /slack-bot
ADD config.yaml.default /slack-bot/config.yaml

VOLUME ["/slack-bot/playground/"]

ENTRYPOINT ["slack-bot"]
