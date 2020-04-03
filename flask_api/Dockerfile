FROM python:alpine
WORKDIR /code
COPY requirements.txt .
RUN pip install -r requirements.txt
WORKDIR /code/src
COPY src .
ENTRYPOINT [ "python" ]
CMD [ "app.py" ]

#FROM python:alpine
#WORKDIR /code
#COPY requirements.txt .
#RUN pip install -r requirements.txt
#WORKDIR /code/src
#COPY src .
#EXPOSE 5000
#CMD ["python", "app.py"]

# FROM ubuntu:16.04
# RUN apt-get update -y && \
# apt-get install -y python-pip python-dev
# # We copy just the requirements.txt first to leverage Docker cache