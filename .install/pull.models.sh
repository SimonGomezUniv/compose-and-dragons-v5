#!/bin/bash
docker model pull ai/qwen2.5:3B-F16
docker model pull huggingface.co/menlo/jan-nano-gguf:q4_k_m
docker model pull huggingface.co/qwen/qwen2.5-0.5b-instruct-gguf:Q4_K_M
docker model pull huggingface.co/unsloth/qwen3.5-0.8b-gguf:Q4_K_M

docker model pull ai/embeddinggemma:latest
docker model pull ai/qwen2.5:0.5B-F16
docker model pull ai/qwen2.5:1.5B-F16
docker model pull huggingface.co/unsloth/nvidia-nemotron-3-nano-4b-gguf:Q4_K_M