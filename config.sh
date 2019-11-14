#!/bin/sh

curl https://github.com/mozilla/DeepSpeech/releases/download/v0.5.1/deepspeech-0.5.1-models.tar.gz

tar xvfz deepspeech-0.5.1-models.tar.gz

mv deepspeech-0.5.1-models/* config/.
rm -rf deepspeech-0.5.1-model*
