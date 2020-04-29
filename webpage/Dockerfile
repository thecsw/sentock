FROM abiosoft/caddy

RUN apk update 
RUN apk upgrade 
RUN apk add npm 
RUN npm i -g jsdoc
COPY ./front /srv
COPY ./Caddyfile /etc
WORKDIR /srv
RUN sudo jsdoc index.js 
RUN ls
#RUN mv /srv/out /srv/docs/front_docs
