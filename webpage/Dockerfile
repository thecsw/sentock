FROM abiosoft/caddy

RUN apk update 
RUN apk upgrade 
RUN apk add npm 
RUN npm i -g jsdoc
COPY ./Caddyfile /etc
COPY ./front /srv
WORKDIR /srv
RUN /usr/bin/jsdoc index.js
RUN ls
RUN jsdoc index.js 
