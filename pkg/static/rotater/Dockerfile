FROM ubuntu:20.04
COPY . /source
RUN apt-get update
RUN apt-get -y install python3 python3-pip zip python3-dev libpq-dev
RUN pip3 install -r /source/requirements.txt --target /source
RUN cd /source && zip -r lambda_function.zip .
