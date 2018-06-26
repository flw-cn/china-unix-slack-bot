FROM golang
LABEL name="slack-bot"
LABEL maintainer="flw@cpan.org"

RUN apt update && apt install -y fortune fortune-zh locales locales-all
RUN go get github.com/flw-cn/slack-bot

ENV LANG en_US.UTF8

ENTRYPOINT ["slack-bot", "--debug", "--token"]
