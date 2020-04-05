FROM python:3.6-alpine
WORKDIR /code
#api for twitter
COPY requirements.txt .
RUN \
 apk add --no-cache postgresql-libs && \
 apk add --no-cache --virtual .build-deps gcc musl-dev postgresql-dev && \
 python3 -m pip install -r requirements.txt --no-cache-dir && \
 apk --purge del .build-deps
WORKDIR /code/src
COPY src .
CMD python twitter.py