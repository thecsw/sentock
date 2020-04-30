FROM caddy:alpine

RUN apk update 
RUN apk upgrade 
RUN apk add npm 
RUN npm i -g jsdoc
<<<<<<< HEAD
WORKDIR /srv
COPY ./front .
COPY ./Caddyfile /etc/caddy
RUN jsdoc *.js 
RUN ls out
RUN mv out front_docs
=======
COPY ./Caddyfile /etc
COPY ./front /srv
WORKDIR /srv
RUN /usr/bin/jsdoc index.js
RUN ls
RUN jsdoc index.js 
>>>>>>> parent of 04f8eaf... finished getting documentation for the backend stuff, now just need to figure out the frontend docs
