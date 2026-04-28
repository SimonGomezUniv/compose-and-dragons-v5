const axios = require('axios');
const config = require('../config');

class OllamaClient {
  constructor() {
    this.client = axios.create({
      baseURL: config.ollama.baseUrl,
      timeout: 60000,
    });
  }

  /**
   * Générer une complétion avec Ollama
   * @param {string} systemPrompt - Rôle / contexte système
   * @param {string} userPrompt - Question utilisateur
   * @returns {Promise<string>} Réponse générée
   */
  async generateCompletion(systemPrompt, userPrompt) {
    try {
      const messages = [
        { role: 'system', content: systemPrompt },
        { role: 'user', content: userPrompt },
      ];

      const response = await this.client.post('/chat/completions', {
        model: config.ollama.llmModel,
        messages,
        temperature: 0.3, // Low temp pour réponses déterministes
        top_p: 0.9,
        stream: false,
      });

      if (!response.data || !response.data.choices || !response.data.choices[0]) {
        throw new Error('Invalid response from Ollama');
      }

      return response.data.choices[0].message.content;
    } catch (error) {
      console.error('Ollama completion error:', error.message);
      throw error;
    }
  }

  /**
   * Générer une complétion structurée (JSON)
   * @param {string} systemPrompt
   * @param {string} userPrompt
   * @param {string} jsonSchema - Description du format JSON attendu
   * @returns {Promise<Object>} Objet JSON parsé
   */
  async generateStructuredCompletion(systemPrompt, userPrompt, jsonSchema) {
    const enhancedPrompt = `${userPrompt}\n\nRespond ONLY with valid JSON matching this schema:\n${jsonSchema}`;

    let attempts = 0;
    const maxAttempts = 3;

    while (attempts < maxAttempts) {
      try {
        const response = await this.generateCompletion(systemPrompt, enhancedPrompt);
        const parsed = JSON.parse(response);
        return parsed;
      } catch (error) {
        attempts++;
        if (attempts >= maxAttempts) {
          throw new Error(`Failed to generate valid JSON after ${maxAttempts} attempts: ${error.message}`);
        }
        console.warn(`Attempt ${attempts} failed, retrying...`);
      }
    }
  }
}

module.exports = { OllamaClient };
