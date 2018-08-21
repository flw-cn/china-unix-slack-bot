FROM ubuntu
LABEL name="slack-bot"
LABEL maintainer="flw <flw@cpan.org>"

ENV DEBIAN_FRONTEND=noninteractive

# Layer 1: Base OS
RUN apt-get update
RUN apt-get install -y apt-utils locales locales-all
RUN apt-get install -y git golang
RUN apt-get install -y fortune fortune-zh
RUN apt-get clean

ENV LANG en_US.UTF8
ENV TERM xterm-256color

# Layer 2: Base Golang development environment
ENV GOPATH "/go"
ENV PATH "$PATH:/go/bin"
RUN go get github.com/sirupsen/logrus
RUN go get github.com/flw-cn/go-smartConfig

# Layer 3: Slack bot development environment
RUN go get github.com/flw-cn/slack
RUN go get github.com/flw-cn/go-slackbot
RUN go get github.com/flw-cn/slack-bot

WORKDIR /slack-bot
RUN cp $GOPATH/src/github.com/flw-cn/slack-bot/config.yaml.sample /slack-bot/config.yaml

VOLUME ["/slack-bot/playground/"]

ENTRYPOINT ["slack-bot"]
