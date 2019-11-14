#!/bin/sh

# Run this after docker container is running

curl -X POST -F file=@wonder.wav localhost:8080/rec