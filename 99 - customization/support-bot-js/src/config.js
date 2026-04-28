require('dotenv').config();

const config = {
  ollama: {
    baseUrl: process.env.OLLAMA_URL || 'http://localhost:11434/v1',
    embeddingUrl: process.env.OLLAMA_EMBEDDING_URL || 'http://localhost:11435/v1',
    llmModel: process.env.OLLAMA_LLM_MODEL || 'qwen2:0.5b',
    embeddingModel: process.env.OLLAMA_EMBEDDING_MODEL || 'embeddinggemma:latest',
  },
  api: {
    port: parseInt(process.env.PORT || '3000', 10),
    host: '0.0.0.0',
  },
  rag: {
    vectorStorePath: process.env.DATA_VECTOR_STORE || './data/vectorStore.json',
    similarityThreshold: parseFloat(process.env.RAG_SIMILARITY_THRESHOLD || '0.4'),
    topN: parseInt(process.env.RAG_TOP_N || '5', 10),
  },
  data: {
    ticketsDir: process.env.DATA_TICKETS_DIR || './data/tickets',
    contextFile: process.env.DATA_CONTEXT_FILE || './data/context/support-expert.md',
  },
};

module.exports = config;
