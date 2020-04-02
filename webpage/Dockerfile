FROM nginx
WORKDIR /var/www
COPY front .
RUN ls
# WORKDIR /var/www/docs-frontend
# COPY --from=dev /code/docs-frontend .
# COPY docs-backend /var/www/docs-backend
COPY nginx.conf /etc/nginx/
EXPOSE 10000