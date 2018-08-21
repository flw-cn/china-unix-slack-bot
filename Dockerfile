FROM ubuntu
LABEL name="slack-bot"
LABEL maintainer="flw <flw@cpan.org>"

ENV DEBIAN_FRONTEND noninteractive

# Layer 1: Base OS
RUN    apt-get update \
    && apt-get install -y --no-install-recommends apt-utils \
    && apt-get install -y locales locales-all git golang fortune fortune-zh \
    && apt-get clean

ENV LANG=en_US.UTF8 \
    TERM=xterm-256color \
    DEBIAN_FRONTEND=teletype

# Layer 2: Base Golang development environment
ENV GOPATH="/go" \
    PATH="$PATH:/go/bin"

RUN    go get github.com/sirupsen/logrus \
    && go get github.com/flw-cn/go-smartConfig

# Layer 3: Slack bot development environment
RUN    go get github.com/flw-cn/slack \
    && go get github.com/flw-cn/go-slackbot \
    && go get github.com/flw-cn/slack-bot

WORKDIR /slack-bot
RUN cp $GOPATH/src/github.com/flw-cn/slack-bot/config.yaml.sample /slack-bot/config.yaml

VOLUME ["/slack-bot/playground/"]

ENTRYPOINT ["slack-bot"]
