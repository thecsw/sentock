FROM nginx
RUN \
	apt update -y && \
	apt upgrade -y && \
	apt install -y npm && \
	npm install -g jsdoc
WORKDIR /var/www
COPY front .
RUN \
	jsdoc *.js && \
	mv out docs

# WORKDIR /var/www/docs-frontend
# COPY --from=dev /code/docs-frontend .
# COPY docs-backend /var/www/docs-backend
COPY nginx.conf /etc/nginx/
EXPOSE 80