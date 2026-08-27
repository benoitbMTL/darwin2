# GO BUILD
FROM golang:1.26.7-bookworm AS go-builder

WORKDIR /app

COPY go/go.mod go/go.sum ./
RUN go mod download

COPY go/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o darwin2 .


# VUE BUILD
FROM node:24.19.0-bookworm-slim AS vue-builder

WORKDIR /app

COPY vue/package.json vue/package-lock.json ./
RUN npm ci

COPY vue/ ./
RUN npm run build


# NIKTO
FROM ubuntu:24.04 AS nikto-builder

ARG NIKTO_COMMIT=6e38da24cd9ec20d9239685071999c79288e89df

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates git \
    && git clone --filter=blob:none https://github.com/sullo/nikto.git /nikto \
    && git -C /nikto checkout "$NIKTO_COMMIT" \
    && rm -rf /var/lib/apt/lists/* /nikto/.git


# FINAL IMAGE
FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl perl perl-modules unzip \
    nmap \
    libnet-ssleay-perl libio-socket-ssl-perl libjson-perl libxml-writer-perl libwhisker2-perl \
    libatk1.0-0 libatk-bridge2.0-0 libcups2 libxcomposite1 libxdamage1 libxrandr2 \
    libgbm1 libnss3 libnspr4 libcairo2 libpango-1.0-0 libgdk-pixbuf2.0-0 libgtk-3-0 \
    libasound2t64 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go

# Chrome install
ARG CHROME_VERSION=152.0.7977.64
ARG CHROME_SHA256=8b592f066af71f054aab2cc80fc26f73c775c6d44ebb99d16ade924b24756c2e
ARG CHROMEDRIVER_SHA256=2457e3d1e204ca712d650e1f13c2b524270682471e371b4750fdbe4f15c1f2dc

RUN mkdir selenium \
    && curl -fsSLo chrome.zip "https://storage.googleapis.com/chrome-for-testing-public/${CHROME_VERSION}/linux64/chrome-linux64.zip" \
    && echo "${CHROME_SHA256}  chrome.zip" | sha256sum -c - \
    && unzip -q chrome.zip -d selenium \
    && curl -fsSLo chromedriver.zip "https://storage.googleapis.com/chrome-for-testing-public/${CHROME_VERSION}/linux64/chromedriver-linux64.zip" \
    && echo "${CHROMEDRIVER_SHA256}  chromedriver.zip" | sha256sum -c - \
    && unzip -q chromedriver.zip -d selenium \
    && rm chrome.zip chromedriver.zip

ENV CHROME_BIN=/go/selenium/chrome-linux64/chrome
ENV PATH=$PATH:/go/selenium/chromedriver-linux64


# Nikto
COPY --from=nikto-builder /nikto ./nikto


# copy binaries
COPY --from=go-builder /app/darwin2 .
COPY --from=vue-builder /app/dist /vue/dist

RUN useradd --create-home --uid 10001 darwin2 \
    && chown -R darwin2:darwin2 /go

USER darwin2

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD curl -fsS http://127.0.0.1:8080/config >/dev/null || exit 1

CMD ["./darwin2"]
