FROM python:3.8.2-buster
WORKDIR /code
COPY requirements.txt .
RUN \
 apt-get install -y gcc && \
 python3 -m pip install -r requirements.txt --no-cache-dir
WORKDIR /code/src
COPY src .
RUN pycco twippy.py
WORKDIR /docs
RUN mv /code/src/docs/* .
WORKDIR /code/src
CMD python twippy.py
