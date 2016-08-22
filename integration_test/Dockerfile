FROM us.gcr.io/replay-gaming/go-eventsource

RUN apk add -U --no-cache bash netcat-openbsd

ADD ./scripts/retry /retry

EXPOSE 80

CMD /retry
