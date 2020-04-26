FROM python:3.8.2-buster
WORKDIR /code
COPY requirements.txt .
RUN \
 apt-get install -y gcc && \
 python3 -m pip install -r requirements.txt --no-cache-dir && \
 pip3 install pycco
WORKDIR /code/src
COPY src .
CMD python twippy.py
