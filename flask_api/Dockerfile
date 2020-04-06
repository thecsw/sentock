FROM python:3.6-alpine
WORKDIR /code
COPY requirements.txt .
RUN \
 apk add --no-cache postgresql-libs && \
 apk add --no-cache --virtual .build-deps gcc musl-dev postgresql-dev && \
 python3 -m pip install -r requirements.txt --no-cache-dir && \
 apk --purge del .build-deps
WORKDIR /code/src
COPY src .
RUN pycco *.py
RUN mkdir templates
RUN mkdir -p static/css
WORKDIR /code/src/docs
RUN echo "<style>" >> app.html
RUN cat pycco.css >> app.html
RUN echo "</style>" >> app.html
RUN mv app.html ../templates/.
WORKDIR /code/src
ENTRYPOINT [ "python" ]
CMD [ "app.py" ]