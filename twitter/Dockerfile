FROM python:alpine
WORKDIR /code
#api for twitter
COPY requirements.txt .
RUN pip install -r requirements.txt
WORKDIR /code/src
COPY src .
#RUN python twitter.py