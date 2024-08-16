ARG NODE_VERSION=22.5
ARG NGINX_VERSION=1.17.1
ARG NODE_BASE_IMAGE=alpine 

FROM node:${NODE_VERSION}-${NODE_BASE_IMAGE} as build

WORKDIR /usr/local/app

COPY ./frontend /usr/local/app

RUN npm install 

RUN npm run build

FROM nginx:${NGINX_VERSION}-${NODE_BASE_IMAGE}

RUN rm -rf /usr/share/nginx/html/*

COPY --from=build  /usr/local/app/dist/frontend/browser /usr/share/nginx/html

EXPOSE 80
