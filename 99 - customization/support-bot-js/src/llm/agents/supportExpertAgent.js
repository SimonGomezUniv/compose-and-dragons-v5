const { OllamaClient } = require('../ollamaClient');

class SupportExpertAgent {
  constructor() {
    this.ollama = new OllamaClient();
  }

  /**
   * Analyser un ticket de support
   * @param {Object} currentTicket - Ticket à analyser
   * @param {string} expertContext - Contexte MD d'expertise
   * @param {Array} ragResults - Résultats RAG (tickets similaires)
   * @returns {Promise<Object>} AnalysisResult
   */
  async analyzeTicket(currentTicket, expertContext, ragResults = []) {
    // Préparer le contexte RAG
    const ragContext =
      ragResults.length > 0
        ? 'Similar past tickets:\n' +
          ragResults
            .map(
              (r) =>
                `- [${r.ticketId}] ${r.title} (${(r.similarity * 100).toFixed(0)}% similar)\n  Resolution: ${r.resolution || 'N/A'}`,
            )
            .join('\n')
        : 'No similar past tickets found.';

    const systemPrompt = `You are an expert support ticket analyzer. 
You have deep knowledge of support categories and resolution patterns:

${expertContext}

Analyze tickets objectively and provide:
1. Category suggestion
2. Confidence level (0-1)
3. Specific resolution suggestions
4. List of info needed to fully resolve
5. Brief reasoning

Respond ONLY with valid JSON.`;

    const userPrompt = `Analyze this support ticket:

TICKET ID: ${currentTicket.id}
TITLE: ${currentTicket.title}
DESCRIPTION: ${currentTicket.description}
PRIORITY: ${currentTicket.priority || 'not specified'}
TAGS: ${(currentTicket.tags || []).join(', ') || 'none'}

${ragContext}

Respond with JSON:
{
  "ticketId": "${currentTicket.id}",
  "suggestedCategory": "category name",
  "confidence": 0.85,
  "suggestions": ["suggestion 1", "suggestion 2"],
  "additionalInfoNeeded": ["info 1", "info 2"],
  "reasoning": "brief explanation"
}`;

    try {
      const schema = `{
        ticketId: string,
        suggestedCategory: string,
        confidence: number,
        suggestions: string[],
        additionalInfoNeeded: string[],
        reasoning: string
      }`;

      const result = await this.ollama.generateStructuredCompletion(
        systemPrompt,
        userPrompt,
        schema,
      );

      return {
        ticketId: result.ticketId,
        suggestedCategory: result.suggestedCategory,
        confidence: Math.min(1, Math.max(0, result.confidence)),
        suggestions: result.suggestions || [],
        additionalInfoNeeded: result.additionalInfoNeeded || [],
        reasoning: result.reasoning,
        relatedPastTickets: ragResults.map((r) => ({
          id: r.ticketId,
          title: r.title,
          similarity: r.similarity,
          resolution: r.resolution,
        })),
      };
    } catch (error) {
      console.error('Agent analysis error:', error.message);
      throw error;
    }
  }
}

module.exports = { SupportExpertAgent };
