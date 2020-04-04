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
RUN cat app.html
RUN mv app.html ../templates/.
RUN mv pycco.css ../static/css/.
RUN cp ../static/css/pycco.css .

WORKDIR /code/src
RUN ls
RUN ls -R static
ENTRYPOINT [ "python" ]
CMD [ "app.py" ]