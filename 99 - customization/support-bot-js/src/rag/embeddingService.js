const axios = require('axios');
const config = require('../config');

class EmbeddingService {
  constructor() {
    this.cache = new Map();
    this.client = axios.create({
      baseURL: config.ollama.embeddingUrl,
      timeout: 30000,
    });
  }

  /**
   * Générer un embedding pour un texte
   * @param {string} text
   * @returns {Promise<number[]>}
   */
  async generateEmbedding(text) {
    if (!text || text.trim().length === 0) {
      throw new Error('Empty text cannot be embedded');
    }

    // Vérifier le cache
    if (this.cache.has(text)) {
      return this.cache.get(text);
    }

    try {
      const response = await this.client.post('/embeddings', {
        model: config.ollama.embeddingModel,
        input: text,
      });

      if (!response.data || !response.data.data || !response.data.data[0]) {
        throw new Error('Invalid embedding response from Ollama');
      }

      const embedding = response.data.data[0].embedding;
      this.cache.set(text, embedding);
      return embedding;
    } catch (error) {
      console.error('Embedding error:', error.message);
      throw error;
    }
  }

  /**
   * Générer embedding pour un ticket avec métadonnées enrichies
   * @param {string} title
   * @param {string} description
   * @param {string[]} tags
   * @returns {Promise<number[]>}
   */
  async generateEmbeddingForTicket(title, description, tags = []) {
    // Enrichir le texte avec métadonnées
    const enrichedText = [
      `TITLE: ${title}`,
      `DESCRIPTION: ${description}`,
      tags.length > 0 ? `TAGS: ${tags.join(', ')}` : '',
    ]
      .filter(Boolean)
      .join('\n');

    return this.generateEmbedding(enrichedText);
  }

  /**
   * Vider le cache
   */
  clearCache() {
    this.cache.clear();
  }

  /**
   * Taille du cache
   * @returns {number}
   */
  cacheSize() {
    return this.cache.size;
  }
}

module.exports = { EmbeddingService };
