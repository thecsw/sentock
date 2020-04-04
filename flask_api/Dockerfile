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
RUN mv app.html ../templates/.
RUN mv pycco.css ../static/css/.
RUN cp ../static/css/pycco.css .
#RUN cp pycco.css templates/.
WORKDIR /code/src
ENTRYPOINT [ "python" ]
CMD [ "app.py" ]