FROM abiosoft/caddy

RUN \
	apt update -y && \
	apt upgrade -y && \
	apt install -y npm && \
	apt install -g jsdoc
COPY ./Caddyfile /etc
COPY ./front /srv
RUN \
	jsdoc *.js && \
	mv out docs

