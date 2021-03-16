FROM ubuntu:20.04

COPY . /source

RUN apt-get update
RUN apt-get -y install git zip

# https://github.com/jkehler/awslambda-psycopg2
RUN git clone https://github.com/jkehler/awslambda-psycopg2.git
RUN cp -R awslambda-psycopg2/psycopg2-3.8 /source/psycopg2
RUN cd /source && zip -r lambda_function.zip .
