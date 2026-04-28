# 🚀 Getting Started - Étapes Concrètes

**Objectif:** Plan d'action détaillé pour démarrer l'implémentation du bot de support Node.js  
**Date:** 2026-04-28

---

## 📋 Vue d'ensemble des Documents Créés

Vous disposez maintenant de 3 documents de référence:

| Document | Contenu | Utilité |
|----------|---------|---------|
| **PLAN-NODEJS-SUPPORT-BOT.md** | Plan complet 11 phases | Roadmap d'implémentation |
| **GO-vs-NODEJS-COMPARISON.md** | Mapping Go → Node.js | Comprendre la transposition |
| **DATA-STRUCTURE-GUIDE.md** | Structure YAML + RAG | Préparer les données |
| **GETTING-STARTED.md** | Ce document | Démarrer maintenant |

---

## 🎯 Phase 0: Setup Immédiat (Jour 1-2)

### Étape 0.1: Créer la Structure de Base

```bash
# À la racine du projet
mkdir -p support-bot
cd support-bot

# Structure initiale
mkdir -p src/{api,rag,llm,models,storage,utils}
mkdir -p data/{tickets,context}
mkdir -p config

# Fichiers essentiels
touch package.json
touch tsconfig.json
touch .env
```

### Étape 0.2: Initialiser Git et Dependencies

**Fichier:** `support-bot/package.json`

```json
{
  "name": "support-bot",
  "version": "0.1.0",
  "description": "Support ticket analysis bot with RAG",
  "main": "dist/main.js",
  "scripts": {
    "dev": "ts-node src/main.ts",
    "build": "tsc",
    "start": "node dist/main.js",
    "test": "jest",
    "index": "ts-node src/scripts/indexTickets.ts"
  },
  "dependencies": {
    "express": "^4.18.2",
    "axios": "^1.6.0",
    "yaml": "^2.3.1",
    "dotenv": "^16.3.1",
    "cors": "^2.8.5"
  },
  "devDependencies": {
    "typescript": "^5.1.6",
    "@types/express": "^4.17.17",
    "@types/node": "^20.3.1",
    "ts-node": "^10.9.1"
  }
}
```

**Fichier:** `support-bot/tsconfig.json`

```json
{
  "compilerOptions": {
    "target": "ES2020",
    "module": "commonjs",
    "lib": ["ES2020"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "moduleResolution": "node"
  },
  "include": ["src"],
  "exclude": ["node_modules"]
}
```

**Fichier:** `support-bot/.env`

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

### Étape 0.3: Installer les Dépendances

```bash
cd support-bot
npm install
```

**Résultat attendu:**
- ✅ `node_modules/` créé
- ✅ `package-lock.json` généré
- ✅ TypeScript et dependencies installées

---

## 📂 Phase 1: Préparer les Données (Jour 2-3)

### Étape 1.1: Créer les Données YAML de Test

Créer 5 fichiers dans `support-bot/data/tickets/`:

**Fichier:** `data/tickets/ticket-001.yml`
```yaml
id: "TICKET-001"
title: "Cannot login to account"
description: |
  User reports being unable to login.
  Tried resetting password but getting 401 error.
  Last successful login was 3 days ago.
priority: high
category: account
tags: [login, account, 401-error]
created_at: "2026-04-15T10:30:00Z"
resolution: "Password reset token was expired. Issued new reset link."
```

*(Voir DATA-STRUCTURE-GUIDE.md pour les 4 autres fichiers)*

### Étape 1.2: Créer le Contexte d'Expertise

**Fichier:** `data/context/support-expert.md`

```markdown
# Support Expert Knowledge Base

## Catégories
- **Account**: Login, password, permissions
- **Billing**: Charges, refunds, subscriptions
- **Technical**: API errors, performance, outages
- **Feature**: Requests, enhancements

## Common Patterns

### Password Issues
1. Verify account exists
2. Check if locked (> 5 attempts)
3. Send reset link
4. Document in CRM

### Billing
1. Pull transaction history
2. Identify duplicates
3. Process refund if < 30 days
4. Investigate root cause

...
```

### Étape 1.3: Valider les Données

```bash
# Vérifier YAML valide
node -e "const yaml = require('yaml'); 
         const fs = require('fs');
         const file = fs.readFileSync('./data/tickets/ticket-001.yml', 'utf8');
         console.log(yaml.parse(file));"

# Résultat attendu: JSON valide sans erreur
```

---

## 🧠 Phase 2: Modèles TypeScript (Jour 3-4)

### Étape 2.1: Créer les Interfaces

**Fichier:** `src/models/ticket.ts`

```typescript
export interface SupportTicket {
  id: string;
  title: string;
  description: string;
  priority?: 'low' | 'medium' | 'high' | 'critical';
  category?: string;
  tags?: string[];
  created_at?: string;
  resolution?: string;
}

export interface AnalysisResult {
  ticketId: string;
  suggestedCategory: string;
  confidence: number;
  suggestions: string[];
  additionalInfoNeeded: string[];
  reasoning: string;
  relatedPastTickets: RelatedTicket[];
}

export interface RelatedTicket {
  id: string;
  title: string;
  similarity: number;
  resolution?: string;
}
```

**Fichier:** `src/models/embedding.ts`

```typescript
export interface VectorRecord {
  id: string;
  ticketId: string;
  content: string;
  embedding: number[];
  metadata: {
    keywords: string[];
    mainCategory: string;
    subcategory: string;
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

### Étape 2.2: Créer Config

**Fichier:** `src/config/config.ts`

```typescript
import dotenv from 'dotenv';

dotenv.config();

export const config = {
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
```

---

## 🔗 Phase 3: Services Core RAG (Jour 4-6)

### Étape 3.1: Embedding Service

**Fichier:** `src/rag/embeddingService.ts`

```typescript
import axios from 'axios';
import { config } from '../config/config';

export class EmbeddingService {
  private cache: Map<string, number[]> = new Map();

  async generateEmbedding(text: string): Promise<number[]> {
    // Check cache first
    if (this.cache.has(text)) {
      return this.cache.get(text)!;
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

  async generateEmbeddingForTicket(title: string, description: string): Promise<number[]> {
    const text = `${title}\n${description}`;
    return this.generateEmbedding(text);
  }
}
```

### Étape 3.2: Vector Store

**Fichier:** `src/rag/vectorStore.ts`

```typescript
import fs from 'fs';
import { VectorRecord } from '../models/embedding';
import { v4 as uuidv4 } from 'uuid';

export class VectorStore {
  private records: Map<string, VectorRecord> = new Map();

  async loadFromFile(path: string): Promise<void> {
    try {
      if (fs.existsSync(path)) {
        const data = JSON.parse(fs.readFileSync(path, 'utf-8'));
        this.records = new Map(Object.entries(data));
      }
    } catch (error) {
      console.error('Error loading vector store:', error);
    }
  }

  async saveToFile(path: string): Promise<void> {
    const data = Object.fromEntries(this.records);
    fs.writeFileSync(path, JSON.stringify(data, null, 2));
  }

  addRecord(record: Omit<VectorRecord, 'id'>): VectorRecord {
    const fullRecord: VectorRecord = {
      id: uuidv4(),
      ...record,
    };
    this.records.set(fullRecord.id, fullRecord);
    return fullRecord;
  }

  getAllRecords(): VectorRecord[] {
    return Array.from(this.records.values());
  }
}
```

### Étape 3.3: Similarity Search

**Fichier:** `src/rag/similaritySearch.ts`

```typescript
import { VectorRecord } from '../models/embedding';

export function cosineSimilarity(a: number[], b: number[]): number {
  const dotProduct = a.reduce((sum, x, i) => sum + x * b[i], 0);
  const magnitudeA = Math.sqrt(a.reduce((sum, x) => sum + x * x, 0));
  const magnitudeB = Math.sqrt(b.reduce((sum, x) => sum + x * x, 0));
  return magnitudeA > 0 && magnitudeB > 0 ? dotProduct / (magnitudeA * magnitudeB) : 0;
}

export function searchSimilar(
  records: VectorRecord[],
  queryEmbedding: number[],
  threshold: number = 0.4,
  topN: number = 5
): VectorRecord[] {
  const results = records
    .map(record => ({
      ...record,
      similarity: cosineSimilarity(queryEmbedding, record.embedding),
    }))
    .filter(r => r.similarity >= threshold)
    .sort((a, b) => b.similarity - a.similarity)
    .slice(0, topN);

  return results;
}
```

---

## 🤖 Phase 4: Agents LLM (Jour 6-7)

### Étape 4.1: Base LLM Client

**Fichier:** `src/llm/ollamaClient.ts`

```typescript
import axios from 'axios';
import { config } from '../config/config';

export class OllamaClient {
  private baseUrl = config.ollama.baseUrl;

  async generateCompletion(
    systemPrompt: string,
    userPrompt: string,
    model: string = config.ollama.llmModel
  ): Promise<string> {
    try {
      const response = await axios.post(`${this.baseUrl}/chat/completions`, {
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

  async generateStructuredCompletion(
    systemPrompt: string,
    userPrompt: string,
    jsonSchema: any
  ): Promise<any> {
    // Pour output JSON structured
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
```

### Étape 4.2: Support Expert Agent

**Fichier:** `src/llm/agents/supportExpertAgent.ts`

```typescript
import { OllamaClient } from '../ollamaClient';
import { SupportTicket, AnalysisResult } from '../../models/ticket';
import { VectorRecord } from '../../models/embedding';

export class SupportExpertAgent {
  private ollama: OllamaClient;

  constructor() {
    this.ollama = new OllamaClient();
  }

  async analyzeTicket(
    currentTicket: SupportTicket,
    expertContext: string,
    ragResults: VectorRecord[]
  ): Promise<AnalysisResult> {
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
      .map(r => `- Ticket ${r.ticketId}: ${r.content} (similarity: ${r.metadata})`)
      .join('\n');

    const userPrompt = `Analyze this ticket:
Title: ${currentTicket.title}
Description: ${currentTicket.description}
Priority: ${currentTicket.priority}

Related past tickets:
${ragContext}

Provide analysis and suggestions.`;

    const result = await this.ollama.generateStructuredCompletion(
      systemPrompt,
      userPrompt,
      {}
    );

    return {
      ticketId: currentTicket.id,
      suggestedCategory: result.suggestedCategory,
      confidence: result.confidence,
      suggestions: result.suggestions,
      additionalInfoNeeded: result.additionalInfoNeeded,
      reasoning: result.reasoning,
      relatedPastTickets: ragResults.map(r => ({
        id: r.ticketId,
        title: currentTicket.title,
        similarity: 0.85, // À calculer correctement
      })),
    };
  }
}
```

---

## 🔌 Phase 5: API REST (Jour 7-8)

### Étape 5.1: Express Setup

**Fichier:** `src/api/routes.ts`

```typescript
import express, { Request, Response } from 'express';
import { analyzeTicketController } from './controllers';

const router = express.Router();

router.post('/api/analyze', async (req: Request, res: Response) => {
  try {
    const { ticket } = req.body;
    const result = await analyzeTicketController(ticket);
    res.json(result);
  } catch (error) {
    res.status(500).json({ error: String(error) });
  }
});

router.get('/api/health', (req: Request, res: Response) => {
  res.json({ status: 'healthy' });
});

export default router;
```

### Étape 5.2: Main Application

**Fichier:** `src/main.ts`

```typescript
import express from 'express';
import cors from 'cors';
import { config } from './config/config';
import routes from './api/routes';

const app = express();

app.use(express.json());
app.use(cors());
app.use(routes);

app.listen(config.api.port, () => {
  console.log(`🚀 Server running on port ${config.api.port}`);
});
```

---

## ✅ Checklist de Démarrage

### Jour 1-2: Setup
- [ ] Créer structure de répertoires
- [ ] `npm install` succès
- [ ] `.env` configuré
- [ ] TypeScript compilation ✅

### Jour 2-3: Données
- [ ] 5 fichiers YAML créés
- [ ] support-expert.md prêt
- [ ] YAML valide (test parsing)

### Jour 3-4: Modèles
- [ ] Interfaces TypeScript définies
- [ ] Config chargeable
- [ ] Types compilables

### Jour 4-6: RAG
- [ ] Embedding service fonctionnel
- [ ] Vector store persistence OK
- [ ] Similarity search testée
- [ ] Cache embeddings ✅

### Jour 6-7: Agents
- [ ] Ollama client fonctionne
- [ ] Support Expert Agent répond
- [ ] Output JSON valide

### Jour 7-8: API
- [ ] Express routes définis
- [ ] `/api/analyze` endpoint
- [ ] `/api/health` endpoint
- [ ] CORS configuré

### Jour 8: Tests
- [ ] Tester `/api/health`
- [ ] Tester `/api/analyze` avec ticket
- [ ] Vérifier RAG results
- [ ] Validation output

---

## 🧪 Test Rapide (Jour 8)

```bash
# Terminal 1: Démarrer Ollama
docker-compose up

# Terminal 2: Démarrer app
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
{
  "ticketId": "TEST-001",
  "suggestedCategory": "account",
  "confidence": 0.92,
  "suggestions": [...],
  ...
}
```

---

## 🎯 Priorités

### ✅ MVP (Minimal Viable Product)
1. Load tickets YAML
2. Generate embeddings
3. Store vector
4. Search similar
5. Basic analysis API

### ⭐ Phase 2 Enhancements
1. Metadata extraction agent
2. Context compression
3. History tracking
4. Admin dashboard

### 🚀 Phase 3 Advanced
1. Fine-tuning on custom data
2. Multiple model support
3. Batch processing
4. Production deployment

---

## 📞 Support Pendant Développement

**Questions à se poser:**

1. ✅ **Embeddings générés?** → Vérifier Ollama URL et modèle
2. ✅ **Vector store persiste?** → Vérifier chemin fichier
3. ✅ **LLM répond?** → Vérifier API completion
4. ✅ **Results pertinents?** → Ajuster threshold RAG

---

## 🎓 Ressources Pratiques

Dans ce répertoire (99 - customization):
- ✅ PLAN-NODEJS-SUPPORT-BOT.md - Référence complète
- ✅ GO-vs-NODEJS-COMPARISON.md - Clarifications
- ✅ DATA-STRUCTURE-GUIDE.md - Format données
- ✅ GETTING-STARTED.md - Ce document

---

## 🚀 Prochaine Étape

**Maintenant:**
1. Lire PLAN-NODEJS-SUPPORT-BOT.md complètement
2. Discuter des points d'clarification
3. Valider approche avec l'équipe
4. Démarrer Phase 0 de ce document

**Puis:**
1. Créer structure repos
2. Initialiser Node.js
3. Commencer implémentation progressive

---

**Document créé:** 2026-04-28  
**Status:** 🟢 Prêt pour démarrage immédiat
