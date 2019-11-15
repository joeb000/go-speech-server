export DEEPSPEECH_CONFIG=$DEEPSPEECH/deepspeech-0.5.1-models

./go-speech-server \
-model $DEEPSPEECH_CONFIG/output_graph.pbmm \
-alphabet $DEEPSPEECH_CONFIG/alphabet.txt \
-lm $DEEPSPEECH_CONFIG/lm.binary \
-trie $DEEPSPEECH_CONFIG/trie