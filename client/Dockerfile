# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM ubuntu:latest

# Add Maintainer Info
LABEL maintainer="Joe B <joebernitt00@gmail.com>"

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    curl \
    wget \
    software-properties-common \
    git \
    gcc \
    sox \
    swig \
    cmake \
    pkg-config \
    libpng-dev \
    libsox-dev \
    g++



RUN apt-get install -y --no-install-recommends libasound-dev portaudio19-dev libportaudio2 libportaudiocpp0

#RUN mkdir /pa
#WORKDIR /pa
#RUN wget http://portaudio.com/archives/pa_stable_v190600_20161030.tgz
#RUN tar xvf pa_stable_v190600_20161030.tgz
#WORKDIR /pa/portaudio
#RUN ./configure && make
#RUN make install

RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get install -y --no-install-recommends golang-go
RUN mkdir /gopath
RUN mkdir /gopath/src
RUN mkdir /gopath/src/local

ENV GOPATH "/gopath"

# Set the Current Working Directory inside the container
WORKDIR /gopath/src/local

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN go get

#RUN go get
RUN go build -o client .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ./client -host 192.168.0.118:8080

