const express = require('express');
const router = express.Router();
const { SupportExpertAgent } = require('../llm/agents/supportExpertAgent');
const { EmbeddingService } = require('../rag/embeddingService');
const { VectorStore } = require('../models/embedding');
const { searchSimilar } = require('../rag/similaritySearch');
const { loadTicketsFromYAML, loadExpertContext } = require('../storage/ticketLoader');
const config = require('../config');

// Initialiser services
const agent = new SupportExpertAgent();
const embeddingService = new EmbeddingService();
const vectorStore = new VectorStore();

let expertContext = '';

/**
 * Initialiser les données
 */
async function initializeData() {
  try {
    // Charger expertise
    expertContext = await loadExpertContext(config.data.contextFile);

    // Charger tickets existants
    const tickets = await loadTicketsFromYAML(config.data.ticketsDir);
    console.log(`📚 Loaded ${tickets.length} tickets`);

    // Indexer les tickets
    for (const ticket of tickets) {
      const content = `${ticket.title}\n${ticket.description}\n${(ticket.tags || []).join(' ')}`;
      const embedding = await embeddingService.generateEmbeddingForTicket(
        ticket.title,
        ticket.description,
        ticket.tags,
      );

      vectorStore.addRecord({
        ticketId: ticket.id,
        content: content,
        embedding: embedding,
        metadata: {
          title: ticket.title,
          category: ticket.category,
          priority: ticket.priority,
          resolution: ticket.resolution,
          tags: ticket.tags || [],
        },
      });
    }

    // Sauvegarder vector store
    await vectorStore.saveToFile(config.rag.vectorStorePath);
    console.log(`🔍 Indexed ${vectorStore.size()} tickets`);
  } catch (error) {
    console.error('Initialization error:', error);
  }
}

/**
 * POST /api/analyze - Analyser un ticket
 */
router.post('/analyze', async (req, res) => {
  try {
    const { id, title, description, tags = [], priority = 'medium' } = req.body;

    if (!id || !title || !description) {
      return res.status(400).json({
        error: 'Missing required fields: id, title, description',
      });
    }

    // Générer embedding pour la requête
    const queryEmbedding = await embeddingService.generateEmbeddingForTicket(
      title,
      description,
      tags,
    );

    // Rechercher tickets similaires
    const allRecords = vectorStore.getAllRecords();
    const similarTickets = searchSimilar(
      allRecords,
      queryEmbedding,
      config.rag.similarityThreshold,
      config.rag.topN,
    );

    // Analyser avec agent
    const analysis = await agent.analyzeTicket(
      { id, title, description, tags, priority },
      expertContext,
      similarTickets,
    );

    return res.json({
      success: true,
      ticketId: id,
      analysis,
      similarTickets: similarTickets.map((t) => ({
        id: t.ticketId,
        title: t.metadata.title,
        similarity: t.similarity,
      })),
    });
  } catch (error) {
    console.error('Analysis error:', error);
    res.status(500).json({
      error: 'Analysis failed',
      message: error.message,
    });
  }
});

/**
 * GET /api/health - Health check
 */
router.get('/health', (req, res) => {
  res.json({
    status: 'ok',
    vectorStoreSize: vectorStore.size(),
    embeddingCacheSize: embeddingService.cacheSize(),
  });
});

/**
 * GET /api/stats - Statistiques
 */
router.get('/stats', (req, res) => {
  res.json({
    vectorRecords: vectorStore.size(),
    embeddingCacheSize: embeddingService.cacheSize(),
    expertContextLoaded: expertContext.length > 0,
  });
});

// Initialiser data au démarrage du routeur
initializeData().catch((err) => console.error('Failed to initialize:', err));

module.exports = router;
