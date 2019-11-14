#!/bin/sh
echo "Downloading data model files to config directory..."
curl -o models.tar -L https://github.com/mozilla/DeepSpeech/releases/download/v0.5.1/deepspeech-0.5.1-models.tar.gz
tar xvfz models.tar

mv deepspeech-0.5.1-models config
rm -rf deepspeech-0.5.1-model*
rm -rf models.tar

# build docker image - could take a couple min
echo "Building docker image, this may take a couple minutes..."
docker build -t test-server .

echo "Running server in detached docker container on port 8080..."
docker run -d -p 8080:8080 test-server

echo "DONE"
