FROM alpine:3.4

ENV DOCKER_MACHNE_VERSION=0.7.0 \
  DELMO_VERSION=0.6.0

RUN apk add --update curl nodejs py-pip \
      && rm -rf /var/cache/apk/*

RUN pip install awscli docker-compose

# Install machine-share
RUN npm install -g machine-share

# Install docker-machine
RUN curl -L https://github.com/docker/machine/releases/download/v${DOCKER_MACHNE_VERSION}/docker-machine-`uname -s`-`uname -m` > /usr/local/bin/docker-machine \
    && chmod +x /usr/local/bin/docker-machine \
    && mkdir -p /root/.docker/machine/machines \
    && mkdir -p /root/.docker/machine/certs

# Install delmo
RUN curl -L https://github.com/bodymindarts/delmo/releases/download/v${DELMO_VERSION}/delmo-linux-amd64 > /usr/local/bin/delmo \
      && chmod +x /usr/local/bin/delmo

ADD entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
CMD ["delmo", "-v"]
