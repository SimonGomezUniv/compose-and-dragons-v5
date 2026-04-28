# 🤖 Plan d'Implémentation - Bot de Support avec RAG (Node.js/JavaScript)

**Date:** 2026-04-28  
**Objectif:** Réimplémenter un bot intelligent de catégorisation de tickets de support en Node.js (JavaScript pur), en utilisant les concepts RAG (Retrieval Augmented Generation) éprouvés du projet D&D.  
**Stack:** Node.js + JavaScript (ES6+) + Ollama (embeddinggemma + qwen2) + Express.js  
**Localisation:** `99 - customization/support-bot-js/`  

---

## 📋 Vue d'ensemble

Le bot doit :
1. **Indexer** les tickets de support antérieurs (YAML) avec des embeddings
2. **Stocker** les index dans un fichier JSON persistent
3. **Fournir une API** qui reçoit les détails d'un ticket courant
4. **Enrichir le contexte** avec les expériences passées via RAG
5. **Générer une analyse** avec catégorisation, suggestions et demandes d'infos

---

## 🎯 Phase 1 : Architecture et Structure du Projet

### 1.1 Structure des Répertoires

```
99 - customization/
├── support-bot-js/              # Application Node.js/JavaScript
│   ├── src/
│   │   ├── api/
│   │   │   ├── routes.js        # Express routes
│   │   │   ├── controllers.js   # Logique métier
│   │   │   └── middleware.js    # Auth, validation
│   │   ├── rag/
│   │   │   ├── embeddingService.js  # Gestion des embeddings
│   │   │   ├── vectorStore.js       # Stockage vectoriel JSON
│   │   │   └── similaritySearch.js  # Recherche de similarité
│   │   ├── llm/
│   │   │   ├── ollamaClient.js      # Client Ollama
│   │   │   ├── agents/
│   │   │   │   ├── supportExpertAgent.js
│   │   │   │   ├── metadataExtractorAgent.js
│   │   │   │   └── ticketAnalyzerAgent.js
│   │   ├── models/
│   │   │   ├── ticket.js        # Interfaces (commentaires JSDoc)
│   │   │   ├── embedding.js
│   │   │   └── analysisResult.js
│   │   ├── storage/
│   │   │   ├── ticketLoader.js  # Charger tickets YAML
│   │   │   └── initialization.js # Initialisation RAG
│   │   └── utils/
│   │       ├── logger.js
│   │       └── fileHelpers.js
│   ├── data/
│   │   ├── tickets/             # Tickets YAML antérieurs
│   │   │   ├── ticket-001.yml
│   │   │   └── ...
│   │   ├── context/
│   │   │   └── support-expert.md    # Description de l'expertise
│   │   └── vectorStore.json         # Store persistant (généré)
│   ├── .env
│   ├── .env.example
│   ├── .gitignore
│   ├── package.json
│   └── main.js
├── docker-compose.yml           # Ollama (existant)
└── INDEX.md                     # Ce plan
```

### 1.2 Stack Technologique

| Composant | Technologie | Rôle |
|-----------|-------------|------|
| **Serveur** | Express.js 4.18+ | API REST |
| **Langage** | JavaScript (ES6+) | Pur JS, pas de build |
| **LLM** | Ollama (qwen2:0.5b) | Génération de texte |
| **Embeddings** | Ollama (embeddinggemma) | Vecteurs pour RAG |
| **Stockage** | JSON + YAML | Persistance |
| **Runtime** | Node.js 18+ | Exécution directe `.js` |
| **Similarité** | Cosine Similarity | Recherche vectorielle |

### 1.3 Dépendances NPM

```json
{
  "dependencies": {
    "express": "^4.18.2",
    "axios": "^1.6.0",
    "yaml": "^2.3.1",
    "dotenv": "^16.3.1",
    "cors": "^2.8.5"
  },
  "devDependencies": {
    "@types/node": "^20.3.1",
    "nodemon": "^3.0.1"
  }
}
```

**Notes:**
- Zéro TypeScript → Pas de `tsc` ou compilation
- Nodemon pour dev avec hot-reload
- Execution directe: `node src/main.js`

---

## 🔧 Phase 2 : Configuration Docker et Services

### 2.1 Docker-Compose (Ollama)

**Utiliser le fichier existant** `docker-compose.yml` (déjà en place dans `99 - customization/`) :
- ✅ Service Ollama (qwen2:0.5b) sur port 11434
- ✅ Service Ollama Embedding (embeddinggemma) sur port 11435
- ✅ Volumes persistants pour les modèles

### 2.2 Configuration de l'Application

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

## 📚 Phase 3 : Modèles de Données

### 3.1 Interfaces/Structures (JSDoc)

**Fichier:** `src/models/ticket.js`

```javascript
/**
 * @typedef {Object} SupportTicket
 * @property {string} id - Identifiant unique
 * @property {string} title - Titre du ticket
 * @property {string} description - Description détaillée
 * @property {string} [category] - Catégorie (optionnelle)
 * @property {'low'|'medium'|'high'|'critical'} [priority] - Priorité
 * @property {Date} [submittedAt] - Date de soumission
 */

const SupportTicket = {
  id: string;
  title: string;
  description: string;
  category?: string;
  priority?: 'low' | 'medium' | 'high' | 'critical';
  submittedAt: Date;
  metadata?: Record<string, any>;
}

// Demande d'analyse
export interface AnalysisRequest {
  currentTicket: SupportTicket;
  expertContext: string;
}

// Résultat d'analyse
export interface AnalysisResult {
  ticketId: string;
  suggestedCategory: string;
  confidence: number;
  suggestions: string[];
  additionalInfoNeeded: string[];
  reasoning: string;
  relatedPastTickets: RelatedTicket[];
}

// Ticket passé lié
export interface RelatedTicket {
  id: string;
  title: string;
  similarity: number;
  resolution?: string;
}
```

### 3.2 Embeddings et Vecteurs

**Fichier:** `src/models/embedding.ts`

```typescript
export interface VectorRecord {
  id: string;
  ticketId: string;
  content: string;
  embedding: number[];
  metadata: {
    keywords: string[];
    category: string;
    importance: 'low' | 'medium' | 'high';
  };
  timestamp: Date;
}

export interface VectorStore {
  records: Map<string, VectorRecord>;
  version: string;
  lastUpdated: Date;
}
```

---

## 🧠 Phase 4 : Implémentation du Système RAG

### 4.1 Service d'Embeddings

**Fichier:** `src/rag/embeddingService.ts`

**Fonctionnalités:**
- Appeler Ollama pour générer embeddings (embeddinggemma)
- Cacher les résultats pour éviter requêtes répétées
- Supporter du texte long via chunking (si nécessaire)

**Méthodologie inspirée de Go:**

```typescript
async generateEmbedding(text: string): Promise<number[]> {
  // 1. Appeler endpoint Ollama
  // 2. Retourner vecteur 384-dim pour embeddinggemma
  // 3. Gérer erreurs et fallbacks
}

async generateEmbeddingForTicket(ticket: SupportTicket): Promise<number[]> {
  // Combiner title + description
  const text = `${ticket.title}\n${ticket.description}`;
  return this.generateEmbedding(text);
}
```

### 4.2 Vector Store Persistant

**Fichier:** `src/rag/vectorStore.ts`

**Fonctionnalités:**
- Charger/sauvegarder `vectorStore.json`
- Gérer la collection d'embeddings
- Supporter add/update/delete de records

```typescript
class VectorStore {
  private records: Map<string, VectorRecord>;

  loadFromFile(path: string): void {
    // Lire vectorStore.json
    // Parser et repeupler records
  }

  saveToFile(path: string): void {
    // Sérialiser records
    // Écrire en JSON pretty-printed
  }

  addRecord(record: VectorRecord): void {
    // Ajouter ou mettre à jour
  }

  getAllRecords(): VectorRecord[] {
    // Retourner tous les records
  }
}
```

### 4.3 Recherche de Similarité

**Fichier:** `src/rag/similaritySearch.ts`

**Algorithme:** Cosine Similarity (comme en Go)

```typescript
function cosineSimilarity(a: number[], b: number[]): number {
  // Calcul du cosinus
  // Return valeur entre -1 et 1 (généralement 0 à 1)
}

async function searchSimilar(
  query: string,
  threshold: number = 0.4,
  topN: number = 5
): Promise<VectorRecord[]> {
  // 1. Générer embedding pour la requête
  // 2. Chercher tous les records avec similarité >= threshold
  // 3. Trier par similarité descendante
  // 4. Retourner top N
}
```

### 4.4 Enrichissement avec Métadonnées

**Inspiré de la démarche Go:**

Avant d'indexer un ticket, extraire et enrichir:

```typescript
interface TicketMetadata {
  keywords: string[];
  mainCategory: string;
  subcategory: string;
  importance: 'low' | 'medium' | 'high';
}

async function extractTicketMetadata(
  ticket: SupportTicket,
  metadataAgent: StructuredAgent
): Promise<TicketMetadata> {
  // Utiliser agent structuré pour extraire
  // Retourner métadonnées typées
}
```

---

## 🤖 Phase 5 : Agents LLM

### 5.1 Support Expert Agent

**Fichier:** `src/llm/agents/supportExpertAgent.ts`

**Rôle:** Agent principal qui analyse les tickets

**System Instructions:**
```
Tu es un expert en support client senior avec 10+ ans d'expérience.

Ton rôle:
1. Analyser les tickets entrants
2. Proposer des catégories appropriées
3. Suggérer des solutions basées sur l'historique
4. Identifier les informations manquantes

Context disponible:
{expertContext}
{ragContext}

Ticket courant:
{currentTicket}

Fournis une analyse structurée avec:
- Catégorie sugérée (avec confiance 0-100%)
- 2-3 suggestions de résolution
- 1-2 informations supplémentaires à demander
- Explication brève (2-3 lignes)
```

### 5.2 Metadata Extractor Agent

**Fichier:** `src/llm/agents/metadataExtractorAgent.ts`

**Rôle:** Agent structuré pour extraction déterministe

**Output JSON:**
```typescript
{
  keywords: ["keyword1", "keyword2", "keyword3", "keyword4"],
  mainCategory: "billing|technical|account|other",
  subcategory: "specific subcategory",
  importance: "low|medium|high"
}
```

### 5.3 Ticket Analyzer Agent

**Fichier:** `src/llm/agents/ticketAnalyzerAgent.ts`

**Rôle:** Agent pour analyse détaillée des tickets

---

## 🔗 Phase 6 : API REST

### 6.1 Endpoints

**POST /api/analyze**
```json
{
  "ticket": {
    "id": "TICKET-123",
    "title": "Cannot login to account",
    "description": "...",
    "priority": "high"
  }
}

// Response
{
  "ticketId": "TICKET-123",
  "suggestedCategory": "account",
  "confidence": 0.92,
  "suggestions": [
    "Reset password and check email",
    "Verify account status in system"
  ],
  "additionalInfoNeeded": [
    "Error message exact text",
    "Last successful login date"
  ],
  "reasoning": "...",
  "relatedPastTickets": [
    {
      "id": "TICKET-089",
      "title": "Login failed - password reset",
      "similarity": 0.87
    }
  ]
}
```

**POST /api/index-tickets**
```json
{
  "message": "Indexed 42 tickets from directory"
}
```

**GET /api/health**
```json
{
  "status": "healthy",
  "services": {
    "ollama": "connected",
    "vectorStore": "ready"
  }
}
```

### 6.2 Middleware et Validation

- ✅ Validation des requêtes (joi/zod)
- ✅ Gestion d'erreurs globale
- ✅ Logging structuré
- ✅ Rate limiting (optionnel)

---

## 💾 Phase 7 : Chargement Initial des Données

### 7.1 Loader de Tickets YAML

**Fichier:** `src/storage/ticketLoader.ts`

```typescript
async function loadTicketsFromDirectory(dir: string): Promise<SupportTicket[]> {
  // 1. Lire tous les fichiers .yml du répertoire
  // 2. Parser YAML
  // 3. Mapper vers SupportTicket[]
  // 4. Retourner tickets
}
```

### 7.2 Processus d'Initialisation

**Fichier:** `src/storage/initialization.ts`

```
1. Vérifier si vectorStore.json existe
2. Si NON:
   a. Charger tous les tickets depuis data/tickets/*.yml
   b. Pour chaque ticket:
      - Extraire métadonnées
      - Générer embedding
      - Enrichir avec métadonnées
      - Ajouter au store
   c. Sauvegarder vectorStore.json
3. Si OUI:
   - Charger depuis le fichier (fast path)
```

---

## 📋 Phase 8 : Contexte d'Expertise Support

### 8.1 Fichier `support-expert.md`

**Chemin:** `data/context/support-expert.md`

**Contenu (à adapter selon vos processus):**

```markdown
# Expertise Support Client

## Catégories de Tickets
- **Account**: Problèmes de compte, accès, permissions
- **Billing**: Facturation, paiements, souscriptions
- **Technical**: Bugs, erreurs, performances
- **Feature Request**: Demandes de nouvelles fonctionnalités
- **Documentation**: Questions sur la documentation

## Processus de Résolution Rapide
1. Vérifier les logs du système
2. Consulter les tickets similaires
3. Appliquer solution éprouvée si existe
4. Demander infos supplémentaires si besoin

## SLA par Catégorie
- Critical: 1h
- High: 4h
- Medium: 24h
- Low: 72h

## Partenaires et Escalades
...
```

---

## 🧪 Phase 9 : Intégration et Tests

### 9.1 Tests Unitaires

- ✅ Tests des services d'embedding
- ✅ Tests de similarité
- ✅ Tests des agents (mocking Ollama)
- ✅ Tests des routes API

### 9.2 Tests d'Intégration

- ✅ Flow complet: indexation → recherche → analyse
- ✅ Persistance vectorStore
- ✅ Load & performance

### 9.3 Données de Test

**Fichier:** `data/tickets/sample-*.yml`
- 5-10 tickets d'exemple pour test initial
- Variété de catégories et priorités

---

## 🚀 Phase 10 : Déploiement et Exécution

### 10.1 Commandes de Build

```bash
# Build
npm run build

# Tests
npm test

# Start en dev
npm run dev

# Start en prod
npm start
```

### 10.2 Docker Compose Complet

**Fichier:** `docker-compose.yml` (à mettre à jour)

```yaml
version: '3.8'

services:
  # Ollama existing (keep as is)
  ollama:
    # ...
  
  ollama-embedding:
    # ...
  
  # Nouveau service Node.js
  support-bot:
    build: .
    ports:
      - "3000:3000"
    environment:
      OLLAMA_URL: "http://ollama:11434/v1"
      OLLAMA_EMBEDDING_URL: "http://ollama-embedding:11435/v1"
    depends_on:
      - ollama
      - ollama-embedding
    volumes:
      - ./data:/app/data
```

### 10.3 Scripts de Démarrage

**Fichier:** `package.json`

```json
{
  "scripts": {
    "dev": "ts-node src/main.ts",
    "build": "tsc",
    "start": "node dist/main.js",
    "test": "jest",
    "index": "ts-node src/utils/indexTickets.ts"
  }
}
```

---

## 📊 Phase 11 : Monitoring et Optimisations

### 11.1 Métriques à Tracker

- ✅ Temps de réponse API
- ✅ Qualité des recommandations (feedback)
- ✅ Similarité moyenne des résultats RAG
- ✅ Cache hit rate des embeddings

### 11.2 Optimisations Futures

1. **Caching d'embeddings**: Redis pour cache distribué
2. **Batch processing**: Indexer par batch plutôt qu'une par une
3. **Fine-tuning du modèle**: Sur vos données spécifiques
4. **Streaming responses**: Pour l'API
5. **Database**: PostgreSQL/MongoDB pour historique à grande échelle

---

## ✅ Checklist de Réalisation

### Phase 1-2: Setup de Base
- [ ] Structure de répertoires créée
- [ ] Package.json avec dépendances
- [ ] Docker-compose up & functional
- [ ] Config centralisée

### Phase 3-4: RAG Core
- [ ] Modèles TypeScript définis
- [ ] Service d'embeddings implémenté
- [ ] Vector store JSON persistant
- [ ] Recherche de similarité fonctionnelle

### Phase 5-6: Agents & API
- [ ] Agents LLM créés
- [ ] Routes API implémentées
- [ ] Validation des requêtes
- [ ] Gestion d'erreurs

### Phase 7-8: Données
- [ ] Tickets YAML chargés
- [ ] Contexte d'expertise prêt
- [ ] Processus d'initialisation automatique

### Phase 9-11: Polish
- [ ] Tests complets
- [ ] Docker compose fonctionnel
- [ ] Monitoring en place
- [ ] Documentation

---

## 📚 Ressources et Références

### Concepts clés dans ce projet (D&D):
- **05-solve-the-context-size-problem/03**: RAG + Metadata extraction
- **Metadata enrichment**: Améliore précision recherche
- **Vector store JSON**: Persistance légère vs DB complexe
- **Cosine similarity**: Métrique standard pour embeddings

### Documentation supplémentaire à créer:
- [ ] API.md - Documentation OpenAPI détaillée
- [ ] RAG-ARCHITECTURE.md - Détails techniques du RAG
- [ ] TESTING.md - Guide des tests
- [ ] DEPLOYMENT.md - Instructions de déploiement

---

## 🎯 Prochaines Étapes

1. **Validation du plan** avec l'équipe
2. **Démarrage Phase 1**: Mise en place structure et config
3. **Itérations** par phase avec tests et ajustements
4. **Intégration continue** au fur et à mesure

---

**Plan rédigé:** 2026-04-28  
**Statut:** 🟡 Draft - En attente de révision et approbation
