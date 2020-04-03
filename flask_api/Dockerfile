FROM python:alpine
WORKDIR /code
COPY requirements.txt .
RUN pip install -r requirements.txt
WORKDIR /code/src
COPY src .
EXPOSE 5000
CMD ["python", "app.py"]