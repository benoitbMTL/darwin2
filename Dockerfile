# Syntaxe Dockerfile v1.4

# Image de base Node
FROM node:18-alpine as node

WORKDIR /app

# Installation de Vue CLI et des dépendances
COPY package.json .
RUN npm install --global @vue/cli
RUN npm install

# Construction du frontend VueJS
RUN npm run build

# Image de base Go
FROM golang:1.19-alpine as go

WORKDIR /go/src/darwin2

# Copie du code Go
COPY go .

# Téléchargement des modules Go
RUN go mod download

# Compilation du backend Go
RUN go build -o main

# Image de base Nginx
FROM nginx:1.23-alpine as nginx

WORKDIR /usr/share/nginx/html

# Copie du frontend
COPY ../vue/dist .

# Configuration Nginx
RUN echo "server_names_hash_bucket_size 64;" >> /etc/nginx/conf.d/default.conf

# Installation d'Echo et des dépendances
RUN apk add --no-cache ca-certificates curl
RUN go install github.com/labstack/echo/v4@latest

# Copie du binaire Echo
COPY --from=go /go/bin/main /usr/bin/

CMD ["/usr/sbin/nginx", "-g", "daemon off;"]