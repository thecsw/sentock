FROM nginx
RUN apt update -y && apt upgrade -y
RUN apt install curl -y
WORKDIR /var/www
COPY front .
# WORKDIR /var/www/docs-frontend
# COPY --from=dev /code/docs-frontend .
# COPY docs-backend /var/www/docs-backend
COPY nginx.conf /etc/nginx/
EXPOSE 10000
