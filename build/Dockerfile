FROM ubuntu:xenial
RUN apt-get update
RUN apt-get install tzdata
RUN echo "Asia/Shanghai" > /etc/timezone
RUN rm -f /etc/localtime
RUN dpkg-reconfigure -f noninteractive tzdata
RUN apt-get install -y ca-certificates
#RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai/etc/localtime
#RUN ln -fs /usr/share/zoneinfo/US/Pacific-New /etc/localtime && dpkg-reconfigure -f noninteractive tzdata
WORKDIR /root/deploy/facebook-spider

RUN mkdir -p /root/deploy/facebook-spider
COPY config.toml /root/deploy/facebook-spider/config.toml
COPY facebook-spider /root/deploy/facebook-spider/
RUN chmod +x /root/deploy/facebook-spider/facebook-spider

ENTRYPOINT ["/root/deploy/facebook-spider/facebook-spider"]
