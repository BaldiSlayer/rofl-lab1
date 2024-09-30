#!/bin/bash

# Start Ollama in the background.
/bin/ollama serve &
# Record Process ID.
pid=$!

# Pause for Ollama to start.
sleep 5

echo "🔴 Retrieve LLAMA3 model..."
ollama pull llama3
echo "🟢 Done!"

echo "🔴 Retrieve LLAMA3 model for embeddings..."
ollama pull nomic-embed-text
echo "🟢 Done!"

echo "🟢 Ready for usage!"

# Wait for Ollama process to finish.
wait $pid
