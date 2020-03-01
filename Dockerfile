FROM ubuntu:latest

LABEL maintainer "flyrainning <flyrainning8@163.com>"

RUN apt-get update -y \
  && apt-get install -y \
    wget \
    gnupg2 


RUN echo "deb [arch=amd64] https://dl.google.com/linux/chrome/deb/ stable main" > /etc/apt/sources.list.d/google-chrome.list
RUN bash -c "wget -q -O - https://dl.google.com/linux/linux_signing_key.pub  | apt-key add -"


RUN apt-get update -y \
  && apt-get install -y \
    google-chrome-stable \
    fonts-wqy-microhei \
    ttf-wqy-zenhei \
  && apt-get autoclean \
  && apt-get autoremove \
  && rm -rf /var/lib/apt/lists/*

ADD web2png /app/web2png
RUN chmod a+x /app/web2png
WORKDIR /app

ENV VERSION 1
ENV PATH "/app:$PATH"

EXPOSE 80 9222
ENTRYPOINT ["/app/web2png"]
