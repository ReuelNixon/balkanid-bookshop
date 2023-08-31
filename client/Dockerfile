FROM node:alpine AS build
WORKDIR /app
COPY ["package.json", "yarn.lock", "/app/"]
RUN yarn