# GO BUILD
FROM golang:1.22 AS go-builder

WORKDIR /app
COPY go/ .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o darwin2


# VUE BUILD
FROM node:20 AS vue-builder

WORKDIR /app
COPY vue/ .

RUN npm install
RUN npm install bootstrap bootstrap-icons selenium-webdriver
RUN npm run build


# FINAL IMAGE
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y \
    wget gnupg ca-certificates git curl perl perl-modules \
    nmap unzip tree net-tools vim \
    libnet-ssleay-perl libio-socket-ssl-perl libjson-perl libxml-writer-perl libwhisker2-perl \
    libatk1.0-0 libatk-bridge2.0-0 libcups2 libxcomposite1 libxdamage1 libxrandr2 \
    libgbm1 libnss3 libnspr4 libcairo2 libpango-1.0-0 libgdk-pixbuf2.0-0 libgtk-3-0 \
    libasound2t64 \
    --no-install-recommends \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go

# Chrome install
RUN mkdir selenium

RUN wget https://storage.googleapis.com/chrome-for-testing-public/122.0.6261.128/linux64/chrome-linux64.zip \
    && unzip chrome-linux64.zip -d selenium \
    && rm chrome-linux64.zip

RUN wget https://storage.googleapis.com/chrome-for-testing-public/122.0.6261.128/linux64/chromedriver-linux64.zip \
    && unzip chromedriver-linux64.zip -d selenium \
    && rm chromedriver-linux64.zip

ENV CHROME_BIN=/go/selenium/chrome-linux64/chrome
ENV PATH=$PATH:/go/selenium/chromedriver-linux64


# Nikto
RUN git clone https://github.com/sullo/nikto.git nikto


# copy binaries
COPY --from=go-builder /app/darwin2 .
COPY --from=vue-builder /app/dist /vue/dist


EXPOSE 8080

CMD ["./darwin2"]