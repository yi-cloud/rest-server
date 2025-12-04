FROM ubuntu:22.04

COPY ./rest-server /usr/local/bin/
COPY ./etc/rest-server /etc/rest-server

ENV TZ=Asia/Shanghai
ENV LANG=zh_CN
ENV LC_CTYPE=zh_CN.UTF-8

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl locales tzdata \
	&& locale-gen zh_CN.UTF-8 \
	&& ln -snf /usr/share/zoneinfo/$TZ /etc/localtime \
	&& echo $TZ > /etc/timezone \
	&& chmod +x /usr/local/bin/rest-server \
    && mv /etc/rest-server/config.linux.yaml /etc/rest-server/config.yaml

WORKDIR /usr/local/bin/

CMD ["rest-server"]