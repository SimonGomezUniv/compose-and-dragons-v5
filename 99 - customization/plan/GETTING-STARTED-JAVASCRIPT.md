# 🚀 Getting Started - Étapes Concrètes (JavaScript)

**Objectif:** Plan d'action détaillé pour démarrer l'implémentation du bot de support en Node.js/JavaScript  
**Localisation:** `99 - customization/support-bot-js/`  
**Langage:** JavaScript ES6+ (pas de TypeScript, pas de build)  
**Date:** 2026-04-28

---

## 🎯 Phase 0: Setup Immédiat (Jour 1-2)

### Étape 0.1: Créer la Structure de Base

```bash
# Dans 99 - customization/
cd "99 - customization"

# Créer structure
mkdir -p support-bot-js/src/{api,rag,llm,models,storage,utils}
mkdir -p support-bot-js/data/{tickets,context}
cd support-bot-js

# Fichiers essentiels
touch package.json
touch .env
touch .gitignore
touch src/main.js
```

**Note:** Pas de `tsconfig.json` (JavaScript pur, execution directe)

---

### Étape 0.2: Initialiser Dependencies

**Fichier:** `support-bot-js/package.json`

```json
{
  "name": "support-bot-js",
  "version": "0.1.0",
  "description": "Support ticket analysis bot with RAG (JavaScript)",
  "main": "src/main.js",
  "type": "commonjs",
  "scripts": {
    "dev": "nodemon src/main.js",
    "start": "node src/main.js",
    "index": "node src/scripts/indexTickets.js"
  },
  "dependencies": {
    "express": "^4.18.2",
    "axios": "^1.6.0",
    "yaml": "^2.3.1",
    "dotenv": "^16.3.1",
    "cors": "^2.8.5",
    "uuid": "^9.0.0"
  },
  "devDependencies": {
    "nodemon": "^3.0.1"
  }
}
```

**Avantages du setup JavaScript:**
- ✅ Pas de compilation (execution directe)
- ✅ Hot-reload avec nodemon en dev
- ✅ Fichiers `.js` directement
- ✅ Pas de `dist/` folder
- ✅ Démarrage immédiat: `npm run dev`

---

### Étape 0.3: Configuration .env

**Fichier:** `support-bot-js/.env`

```env
# Ollama Configuration
OLLAMA_URL=http://localhost:11434/v1
OLLAMA_EMBEDDING_URL=http://localhost:11435/v1
OLLAMA_LLM_MODEL=qwen2:0.5b
OLLAMA_EMBEDDING_MODEL=embeddinggemma:latest

# API Configuration
PORT=3000
NODE_ENV=development

# Data Paths
DATA_TICKETS_DIR=./data/tickets
DATA_CONTEXT_FILE=./data/context/support-expert.md
DATA_VECTOR_STORE=./data/vectorStore.json

# RAG Configuration
RAG_SIMILARITY_THRESHOLD=0.4
RAG_TOP_N=5
```

---

### Étape 0.4: Installer les Dépendances

```bash
cd support-bot-js
npm install
```

**Résultat attendu:**
- ✅ `node_modules/` créé
- ✅ `package-lock.json` généré
- ✅ Dependencies installées
- ✅ Prêt à lancer: `npm run dev`

---

## 📂 Phase 1: Préparer les Données (Jour 2-3)

### Étape 1.1: Créer les Données YAML de Test

Voir `DATA-STRUCTURE-GUIDE.md` pour les 5 fichiers tickets YAML.

Créer dans `support-bot-js/data/tickets/`:
- `ticket-001.yml`
- `ticket-002.yml`
- ...

### Étape 1.2: Créer le Contexte d'Expertise

Voir `DATA-STRUCTURE-GUIDE.md` pour le template complet.

**Fichier:** `support-bot-js/data/context/support-expert.md`

---

## 🧠 Phase 2: Modèles de Données (Jour 3-4)

### Étape 2.1: Structures (JSDoc)

**Fichier:** `src/models/ticket.js`

```javascript
/**
 * @typedef {Object} SupportTicket
 * @property {string} id - Identifiant unique
 * @property {string} title - Titre du ticket
 * @property {string} description - Description détaillée
 * @property {string} [category] - Catégorie (optionnelle)
 * @property {'low'|'medium'|'high'|'critical'} [priority] - Priorité
 * @property {string[]} [tags] - Tags pour recherche
 * @property {Date} [created_at] - Date de soumission
 * @property {string} [resolution] - Résolution si trouvée
 */

/**
 * @typedef {Object} AnalysisResult
 * @property {string} ticketId
 * @property {string} suggestedCategory
 * @property {number} confidence - Entre 0 et 1
 * @property {string[]} suggestions - Suggestions de résolution
 * @property {string[]} additionalInfoNeeded - Infos manquantes
 * @property {string} reasoning - Explication brève
 * @property {RelatedTicket[]} relatedPastTickets - Tickets similaires
 */

/**
 * @typedef {Object} RelatedTicket
 * @property {string} id
 * @property {string} title
 * @property {number} similarity - Score de similarité
 * @property {string} [resolution]
 */

/**
 * Valider un ticket
 * @param {SupportTicket} ticket
 * @returns {boolean}
 */
function validateTicket(ticket) {
  if (!ticket.id || !ticket.title || !ticket.description) {
    throw new Error('Missing required fields: id, title, description');
  }
  return true;
}

module.exports = { validateTicket };
```

**Fichier:** `src/models/embedding.js`

```javascript
/**
 * @typedef {Object} VectorRecord
 * @property {string} id - UUID unique du record
 * @property {string} ticketId - ID du ticket
 * @property {string} content - Contenu enrichi avec metadata
 * @property {number[]} embedding - Vecteur d'embedding (1024 dimensions)
 * @property {Object} metadata - Métadonnées extraites
 * @property {string[]} metadata.keywords - 4 mots-clés
 * @property {string} metadata.mainCategory - Catégorie principale
 * @property {string} metadata.subcategory - Sous-catégorie
 * @property {'low'|'medium'|'high'} metadata.importance - Importance
 * @property {Date} timestamp - Timestamp création
 */

class VectorStore {
  constructor() {
    this.records = new Map();
  }

  /**
   * Charger le store depuis un fichier JSON
   * @param {string} path - Chemin du fichier
   */
  async loadFromFile(path) {
    const fs = require('fs').promises;
    try {
      if (require('fs').existsSync(path)) {
        const data = JSON.parse(await fs.readFile(path, 'utf-8'));
        this.records = new Map(Object.entries(data));
      }
    } catch (error) {
      console.error('Error loading vector store:', error);
    }
  }

  /**
   * Sauvegarder le store en JSON
   * @param {string} path
   */
  async saveToFile(path) {
    const fs = require('fs').promises;
    const data = Object.fromEntries(this.records);
    await fs.writeFile(path, JSON.stringify(data, null, 2));
  }

  /**
   * Ajouter un record
   * @param {Omit<VectorRecord, 'id'>} record
   * @returns {VectorRecord}
   */
  addRecord(record) {
    const { v4: uuidv4 } = require('uuid');
    const fullRecord = {
      id: uuidv4(),
      ...record,
    };
    this.records.set(fullRecord.id, fullRecord);
    return fullRecord;
  }

  getAllRecords() {
    return Array.from(this.records.values());
  }
}

module.exports = { VectorStore };
```

---

### Étape 2.2: Configuration

**Fichier:** `src/config.js`

```javascript
require('dotenv').config();

const config = {
  ollama: {
    baseUrl: process.env.OLLAMA_URL || 'http://localhost:11434/v1',
    embeddingUrl: process.env.OLLAMA_EMBEDDING_URL || 'http://localhost:11435/v1',
    llmModel: 'qwen2:0.5b',
    embeddingModel: 'embeddinggemma:latest',
  },
  api: {
    port: parseInt(process.env.PORT || '3000', 10),
    host: '0.0.0.0',
  },
  rag: {
    vectorStorePath: './data/vectorStore.json',
    similarityThreshold: parseFloat(process.env.RAG_SIMILARITY_THRESHOLD || '0.4'),
    topN: parseInt(process.env.RAG_TOP_N || '5', 10),
  },
  data: {
    ticketsDir: './data/tickets',
    contextFile: './data/context/support-expert.md',
  },
};

module.exports = config;
```

---

## 🔗 Phase 3: Services Core RAG (Jour 4-6)

### Étape 3.1: Embedding Service

**Fichier:** `src/rag/embeddingService.js`

```javascript
const axios = require('axios');
const config = require('../config');

class EmbeddingService {
  constructor() {
    this.cache = new Map();
  }

  /**
   * Générer embedding pour un texte
   * @param {string} text
   * @returns {Promise<number[]>}
   */
  async generateEmbedding(text) {
    // Vérifier cache
    if (this.cache.has(text)) {
      return this.cache.get(text);
    }

    try {
      const response = await axios.post(
        `${config.ollama.embeddingUrl}/embeddings`,
        {
          model: config.ollama.embeddingModel,
          prompt: text,
        }
      );

      const embedding = response.data.embedding;
      this.cache.set(text, embedding);
      return embedding;
    } catch (error) {
      console.error('Error generating embedding:', error);
      throw error;
    }
  }

  /**
   * Générer embedding pour un ticket
   * @param {string} title
   * @param {string} description
   * @returns {Promise<number[]>}
   */
  async generateEmbeddingForTicket(title, description) {
    const text = `${title}\n${description}`;
    return this.generateEmbedding(text);
  }
}

module.exports = { EmbeddingService };
```

### Étape 3.2: Similarity Search

**Fichier:** `src/rag/similaritySearch.js`

```javascript
/**
 * Calculer cosine similarity entre deux vecteurs
 * @param {number[]} a
 * @param {number[]} b
 * @returns {number} Valeur entre -1 et 1
 */
function cosineSimilarity(a, b) {
  const dotProduct = a.reduce((sum, x, i) => sum + x * b[i], 0);
  const magnitudeA = Math.sqrt(a.reduce((sum, x) => sum + x * x, 0));
  const magnitudeB = Math.sqrt(b.reduce((sum, x) => sum + x * x, 0));
  return magnitudeA > 0 && magnitudeB > 0 ? dotProduct / (magnitudeA * magnitudeB) : 0;
}

/**
 * Chercher records similaires
 * @param {VectorRecord[]} records
 * @param {number[]} queryEmbedding
 * @param {number} threshold - Seuil de similarité
 * @param {number} topN - Nombre de résultats max
 * @returns {VectorRecord[]}
 */
function searchSimilar(records, queryEmbedding, threshold = 0.4, topN = 5) {
  const results = records
    .map((record) => ({
      ...record,
      similarity: cosineSimilarity(queryEmbedding, record.embedding),
    }))
    .filter((r) => r.similarity >= threshold)
    .sort((a, b) => b.similarity - a.similarity)
    .slice(0, topN);

  return results;
}

module.exports = { cosineSimilarity, searchSimilar };
```

---

## 🤖 Phase 4: Agents LLM (Jour 6-7)

### Étape 4.1: Ollama Client

**Fichier:** `src/llm/ollamaClient.js`

```javascript
const axios = require('axios');
const config = require('../config');

class OllamaClient {
  /**
   * Générer une completion
   * @param {string} systemPrompt
   * @param {string} userPrompt
   * @param {string} model
   * @returns {Promise<string>}
   */
  async generateCompletion(systemPrompt, userPrompt, model = config.ollama.llmModel) {
    try {
      const response = await axios.post(`${config.ollama.baseUrl}/chat/completions`, {
        model,
        messages: [
          { role: 'system', content: systemPrompt },
          { role: 'user', content: userPrompt },
        ],
        temperature: 0.7,
        stream: false,
      });

      return response.data.choices[0].message.content;
    } catch (error) {
      console.error('Error generating completion:', error);
      throw error;
    }
  }

  /**
   * Générer completion structurée (JSON)
   * @param {string} systemPrompt
   * @param {string} userPrompt
   * @returns {Promise<Object>}
   */
  async generateStructuredCompletion(systemPrompt, userPrompt) {
    const completion = await this.generateCompletion(
      systemPrompt,
      userPrompt + '\nRespond ONLY with valid JSON.',
      config.ollama.llmModel
    );

    try {
      return JSON.parse(completion);
    } catch {
      throw new Error(`Invalid JSON from LLM: ${completion}`);
    }
  }
}

module.exports = { OllamaClient };
```

### Étape 4.2: Support Expert Agent

**Fichier:** `src/llm/agents/supportExpertAgent.js`

```javascript
const { OllamaClient } = require('../ollamaClient');

class SupportExpertAgent {
  constructor() {
    this.ollama = new OllamaClient();
  }

  /**
   * Analyser un ticket
   * @param {SupportTicket} currentTicket
   * @param {string} expertContext
   * @param {VectorRecord[]} ragResults
   * @returns {Promise<AnalysisResult>}
   */
  async analyzeTicket(currentTicket, expertContext, ragResults) {
    const systemPrompt = `You are an expert support analyst with 10+ years experience.
Your role is to analyze support tickets and provide actionable insights.

Expert Context:
${expertContext}

Respond with a valid JSON object with this structure:
{
  "suggestedCategory": "category",
  "confidence": 0.95,
  "suggestions": ["suggestion1", "suggestion2"],
  "additionalInfoNeeded": ["info1", "info2"],
  "reasoning": "brief explanation"
}`;

    const ragContext = ragResults
      .map((r) => `- Ticket ${r.ticketId}: ${r.content.substring(0, 100)}... (similarity: ${r.similarity})`)
      .join('\n');

    const userPrompt = `Analyze this ticket:
Title: ${currentTicket.title}
Description: ${currentTicket.description}
Priority: ${currentTicket.priority || 'not specified'}

Related past tickets:
${ragContext || 'None found'}

Provide analysis and suggestions.`;

    const result = await this.ollama.generateStructuredCompletion(systemPrompt, userPrompt);

    return {
      ticketId: currentTicket.id,
      suggestedCategory: result.suggestedCategory,
      confidence: result.confidence,
      suggestions: result.suggestions,
      additionalInfoNeeded: result.additionalInfoNeeded,
      reasoning: result.reasoning,
      relatedPastTickets: ragResults.map((r) => ({
        id: r.ticketId,
        title: currentTicket.title,
        similarity: r.similarity,
      })),
    };
  }
}

module.exports = { SupportExpertAgent };
```

---

## 🔌 Phase 5: API REST (Jour 7-8)

### Étape 5.1: Express Routes

**Fichier:** `src/api/routes.js`

```javascript
const express = require('express');
const { analyzeTicketController } = require('./controllers');

const router = express.Router();

/**
 * POST /api/analyze
 * Analyser un ticket
 */
router.post('/api/analyze', async (req, res) => {
  try {
    const { ticket } = req.body;
    const result = await analyzeTicketController(ticket);
    res.json(result);
  } catch (error) {
    res.status(500).json({ error: String(error) });
  }
});

/**
 * GET /api/health
 * Health check
 */
router.get('/api/health', (req, res) => {
  res.json({ status: 'healthy' });
});

module.exports = router;
```

### Étape 5.2: Main Application

**Fichier:** `src/main.js`

```javascript
const express = require('express');
const cors = require('cors');
const config = require('./config');
const routes = require('./api/routes');

const app = express();

// Middleware
app.use(express.json());
app.use(cors());

// Routes
app.use(routes);

// Démarrer serveur
app.listen(config.api.port, () => {
  console.log(`🚀 Server running on http://localhost:${config.api.port}`);
});
```

---

## ✅ Checklist de Démarrage

### Jour 1-2: Setup
- [ ] Structure `support-bot-js/` créée
- [ ] `npm install` succès
- [ ] `.env` configuré
- [ ] `npm run dev` démarre sans erreur

### Jour 2-3: Données
- [ ] 5 fichiers YAML créés dans `data/tickets/`
- [ ] `support-expert.md` prêt
- [ ] YAML valide (test parsing)

### Jour 3-4: Modèles
- [ ] `src/models/ticket.js` créé
- [ ] `src/models/embedding.js` créé
- [ ] `src/config.js` chargeable

### Jour 4-6: RAG
- [ ] `src/rag/embeddingService.js` fonctionne
- [ ] `src/rag/similaritySearch.js` testée
- [ ] VectorStore persistence OK
- [ ] Cache embeddings actif

### Jour 6-7: Agents
- [ ] `src/llm/ollamaClient.js` fonctionne
- [ ] `src/llm/agents/supportExpertAgent.js` répond
- [ ] Output JSON valide

### Jour 7-8: API
- [ ] `src/api/routes.js` définis
- [ ] `/api/analyze` endpoint
- [ ] `/api/health` endpoint
- [ ] CORS configuré

---

## 🧪 Test Rapide

```bash
# Terminal 1: Démarrer Ollama
cd ..  # Retour à 99 - customization/
docker-compose up

# Terminal 2: Démarrer app
cd support-bot-js
npm run dev

# Terminal 3: Tester API
curl -X POST http://localhost:3000/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "ticket": {
      "id": "TEST-001",
      "title": "Cannot login",
      "description": "Getting 401 error"
    }
  }'

# Résultat attendu:
# {
#   "ticketId": "TEST-001",
#   "suggestedCategory": "account",
#   "confidence": 0.92,
#   ...
# }
```

---

## 🎯 Priorités

### ✅ MVP (Jour 8-10)
1. Load tickets YAML ✅
2. Generate embeddings ✅
3. Store vector ✅
4. Search similar ✅
5. Basic analysis API ✅

### ⭐ Phase 2 (Semaine 2)
1. Metadata extraction agent
2. Persister vector store automatiquement
3. Error handling complet

### 🚀 Phase 3 (Semaine 3+)
1. History tracking
2. Batch processing
3. Admin dashboard

---

**Document créé:** 2026-04-28  
**Version:** 1.0 JavaScript  
**Statut:** 🟢 Prêt pour démarrage
