# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
#FROM golang:latest
FROM nvidia/cuda:10.0-cudnn7-devel-ubuntu18.04

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
    libsox-fmt-mp3 \
    htop \
    nano \
    swig \
    cmake \
    libboost-all-dev \
    zlib1g-dev \
    libbz2-dev \
    liblzma-dev \
    locales \
    pkg-config \
    libpng-dev \
    libsox-dev \
    libmagic-dev \
    libgsm1-dev \
    libltdl-dev \
    openjdk-8-jdk \
    bash-completion \
    g++ \
    unzip



RUN mkdir /tmp/lib 
WORKDIR /tmp/lib
RUN wget https://github.com/mozilla/DeepSpeech/releases/download/v0.5.1/native_client.amd64.cpu.linux.tar.xz

RUN tar xvf native_client.amd64.cpu.linux.tar.xz

RUN mkdir /tmp/include
WORKDIR /tmp/include
RUN wget https://github.com/mozilla/DeepSpeech/raw/v0.5.1/native_client/deepspeech.h




ENV CGO_LDFLAGS "-L/tmp/lib/"
ENV CGO_CXXFLAGS "-I/tmp/include/"
ENV LD_LIBRARY_PATH "/tmp/lib/"


RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get install -y --no-install-recommends golang-go

# Set the Current Working Directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


WORKDIR /app

# Copy go mod and sum files
#COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
#RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
RUN go mod download


WORKDIR /root/go/pkg/mod/github.com/asticode/go-astideepspeech@v0.0.0-20191027095326-e831da11d013
RUN wget https://github.com/mozilla/DeepSpeech/raw/v0.5.1/native_client/deepspeech.h
# Build the Go app
WORKDIR /app

RUN go build -o main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ./main -model /app/config/output_graph.pbmm -alphabet /app/config/alphabet.txt -lm /app/config/lm.binary -trie /app/config/trie

