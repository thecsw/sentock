FROM python:alpine
WORKDIR /code
COPY requirements.txt .
RUN pip install -r requirements.txt
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