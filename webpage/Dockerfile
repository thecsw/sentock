FROM caddy:alpine

RUN apk update 
RUN apk upgrade 
RUN apk add npm 
RUN npm i -g jsdoc
WORKDIR /srv
COPY ./front .
COPY ./Caddyfile /etc/caddy
RUN jsdoc *.js 
RUN ls out
RUN mv out front_docs
