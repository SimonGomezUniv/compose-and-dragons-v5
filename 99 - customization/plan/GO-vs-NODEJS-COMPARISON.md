# 🔄 Comparaison: Architecture Go vs Node.js

**Objectif:** Montrer comment les concepts de D&D (Go) se transposent à Node.js pour le bot de support

---

## 1️⃣ Architecture RAG: Convergences et Divergences

### Go Implementation (D&D - Repo 05)

```
┌─────────────────────┐
│  Markdown Sections  │
└──────────┬──────────┘
           │
           ▼
┌──────────────────────────────┐
│  Metadata Extractor Agent    │  ◄── Structured Agent
│  (Keywords, Topic, Category) │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Chunk Enrichment            │
│  [METADATA] + Content        │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Embedding Generation        │  ◄── embeddinggemma
│  (on enriched content)       │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  RAG Vector Store (JSON)     │
│  - ID, Embedding, Metadata   │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Interactive Chat Loop       │
│  - Context Compression       │
│  - Similarity Search (top-7) │
└──────────────────────────────┘
```

### Node.js Implementation (Support Bot)

```
┌─────────────────────┐
│   YAML Tickets      │
└──────────┬──────────┘
           │
           ▼
┌──────────────────────────────┐
│  Ticket Loader               │
│  (Parse YAML → SupportTicket)│
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Metadata Extractor Agent    │  ◄── Structured Agent
│  (Keywords, Category, etc.)  │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Enrichment Service          │
│  Ticket + Metadata           │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Embedding Service           │  ◄── embeddinggemma
│  (via Ollama client)         │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  Vector Store (JSON)         │
│  - ticketId, embedding, meta │
└──────────┬───────────────────┘
           │
           ▼
┌──────────────────────────────┐
│  API REST Endpoints          │
│  - POST /api/analyze         │
│  - Search + Analysis         │
└──────────────────────────────┘
```

### Similarités ✅

| Aspect | Go | Node.js | Notes |
|--------|----|---------| -----|
| **Métadonnées** | Keywords + Topic + Category | Keywords + Category + Importance | Même approche structurée |
| **Enrichissement** | Injection avant embedding | Injection avant embedding | Même stratégie RAG |
| **Storage** | JSON file-based | JSON file-based | Persistance légère |
| **Similarité** | Cosine Similarity | Cosine Similarity | Même algorithme |
| **Modèles** | qwen2 + embeddinggemma | qwen2 + embeddinggemma | Identique |

### Différences 🔀

| Aspect | Go (D&D) | Node.js (Support) | Raison |
|--------|----------|-------------------|--------|
| **Entrée** | Markdown Sections | YAML Tickets | Format source différent |
| **Agents** | Character generation + Roleplay | Ticket analysis + Support | Cas d'usage distinct |
| **Flow** | Chat interactif TUI | API REST synchrone | Interface différente |
| **Context Compression** | Automatique en chat | Non utilisé (API) | Pas besoin en mode API |
| **Output** | Streaming responses | JSON responses | API REST standard |

---

## 2️⃣ Métadonnées: Comment les Adapter

### Pattern Go (KeywordMetadata)

```go
type KeywordMetadata struct {
    Keywords  []string `json:"keywords"`
    MainTopic string   `json:"main_topic"`
    Category  string   `json:"category"`
}

// Usage en Go:
extractionPrompt := fmt.Sprintf(`Analyze content and extract:
- Keywords: 4 important terms
- Main topic: primary subject
- Category: type of content`)
```

### Pattern Node.js (TicketMetadata)

```typescript
interface TicketMetadata {
  keywords: string[];           // 4 mots-clés
  mainCategory: string;         // account|billing|technical|...
  subcategory: string;          // Catégorie spécifique
  importance: 'low'|'medium'|'high'; // Priorité estimée
}

// Usage en Node.js:
const extractionPrompt = `Analyze ticket and extract:
- Keywords: 4 important terms from title and description
- Main category: one of [account, billing, technical, feature, docs]
- Subcategory: specific type
- Importance: estimated priority level`;
```

### Processus d'Extraction (Identique)

```
Go:
┌─ Go Structured Agent ─┐
│ (generalize + parse)  │
└─ KeywordMetadata JSON ┘

Node.js:
┌─ Node.js Structured Agent ─┐
│ (same pattern)              │
└─ TicketMetadata JSON ──────┘

Les deux:
1. Envoyer texte à agent
2. Parser output JSON
3. Enrichir document
4. Générer embedding
```

---

## 3️⃣ RAG Search: Implémentation

### Similarity Search - Go

```go
// File: rag.go
func (mvs *MemoryVectorStore) SearchTopNSimilarities(
    embeddingFromQuestion VectorRecord,
    limit float64,     // threshold (ex: 0.4)
    max int,           // topN (ex: 7)
) ([]VectorRecord, error) {
    records, err := mvs.SearchSimilarities(embeddingFromQuestion, limit)
    // Filter by threshold
    return getTopNVectorRecords(records, max)
    // Sort by similarity desc, return top max
}

func CosineSimilarity(a []float32, b []float32) float64 {
    // Standard cosine similarity formula
}
```

### Similarity Search - Node.js

```typescript
// File: src/rag/similaritySearch.ts
async function searchSimilar(
  queryEmbedding: number[],
  threshold: number = 0.4,
  topN: number = 5
): Promise<VectorRecord[]> {
  const results = [];
  
  for (const record of store.getAllRecords()) {
    const similarity = cosineSimilarity(queryEmbedding, record.embedding);
    
    if (similarity >= threshold) {
      results.push({ ...record, similarity });
    }
  }
  
  // Sort by similarity desc
  return results
    .sort((a, b) => b.similarity - a.similarity)
    .slice(0, topN);
}

function cosineSimilarity(a: number[], b: number[]): number {
  // Same formula as Go
  const dotProduct = a.reduce((sum, x, i) => sum + x * b[i], 0);
  const magnitudeA = Math.sqrt(a.reduce((sum, x) => sum + x * x, 0));
  const magnitudeB = Math.sqrt(b.reduce((sum, x) => sum + x * x, 0));
  return dotProduct / (magnitudeA * magnitudeB);
}
```

### Résultats Similaires ✅

```
Both return:
[
  { id, ticketId, embedding, metadata, similarity: 0.87 },
  { id, ticketId, embedding, metadata, similarity: 0.79 },
  { id, ticketId, embedding, metadata, similarity: 0.73 }
]

Ordre: Décroissant par similarité
Filtrage: >= threshold
Limite: TopN résultats
```

---

## 4️⃣ Agents LLM: Mapping

### Go - D&D Agents

1. **NPC Generation Agent** → Generates character sheets
2. **Metadata Extractor** → Extracts keywords + topic
3. **Compressor Agent** → Summarizes conversation
4. **Roleplay Agent** → Responds as NPC character

### Node.js - Support Bot Agents

1. **Metadata Extractor** → Extracts keywords + category ✅ (1:1)
2. **Support Expert Agent** → Analyzes tickets
3. **Ticket Analyzer Agent** → Deep analysis + suggestions
4. ❌ **No Compressor** (pas besoin en API REST)

### Détails Agents

#### Structured Agent (Identique concept)

```go
// Go
agent, err := structured.NewAgent[KeywordMetadata](...)
metadata, _, err := agent.GenerateStructuredData(messages)
```

```typescript
// Node.js
class StructuredAgent<T> {
  async generateStructuredData(messages: Message[]): Promise<T> {
    // Same pattern as Go
    // Parse JSON from LLM output
  }
}
```

#### Roleplay/Expertise Agent

Go:
```go
roleplaySystemInstructions := fmt.Sprintf(
  "You are %s, a %s %s %s...",
  npc.FirstName, npc.Gender, npc.Race, npc.Class,
)
npcAgent.GenerateStreamCompletion(messages, callback)
```

Node.js:
```typescript
const systemInstructions = `You are an expert support agent.
Use the expertise context and past tickets to analyze.
Provide structured analysis with categories and suggestions.`;

await supportAgent.generateCompletion(messages)
```

---

## 5️⃣ Stockage Persistant: Comparaison

### Go - Vector Store

```go
type VectorRecord struct {
    Id           string    `json:"id"`
    Prompt       string    `json:"prompt"`     // Content
    Embedding    []float32 `json:"embedding"` // Vector
    CosineSimilarity float64 // Runtime only
}

// Persistence:
mvs.SaveJSONToFile(filename)    // Save to JSON
mvs.LoadFromJSONFile(filename)  // Load from JSON
```

**Structure JSON Go:**
```json
{
  "uuid-123": {
    "id": "uuid-123",
    "prompt": "## Backstory\nBorn in Ironforge...",
    "embedding": [0.123, 0.456, ...],
    "cosineSimilarity": 0
  }
}
```

### Node.js - Vector Store

```typescript
interface VectorRecord {
  id: string;
  ticketId: string;
  content: string;
  embedding: number[];
  metadata: {
    keywords: string[];
    category: string;
    importance: string;
  };
  timestamp: Date;
}

// Persistence:
await vectorStore.saveToFile(path)  // Save to JSON
await vectorStore.loadFromFile(path) // Load from JSON
```

**Structure JSON Node.js:**
```json
{
  "vuuid-456": {
    "id": "vuuid-456",
    "ticketId": "TICKET-001",
    "content": "Cannot login to account...",
    "embedding": [0.234, 0.567, ...],
    "metadata": {
      "keywords": ["login", "account", "access"],
      "category": "account",
      "importance": "high"
    },
    "timestamp": "2026-04-28T10:30:00Z"
  }
}
```

### Similarités ✅

- Même format JSON
- Même structure clé-valeur (Map/Object)
- Même persistence disk-based
- Même runtime similarity calc

### Différences 🔀

- **Go**: Similarity score au runtime (non sauvegardé)
- **Node.js**: Metadata stocké avec le record (enrichissement persistant)

---

## 6️⃣ Flow de Traitement: Vue Complète

### Go Flow - D&D NPC Chat

```
1. Load Character Sheet (markdown)
2. Split into sections
3. For each section:
   a. Extract metadata (structured agent)
   b. Enrich with metadata
   c. Generate embedding
   d. Store record
4. Save vector store
5. Start chat loop:
   a. Get user question
   b. Generate embedding from question
   c. Search similar records
   d. Inject RAG context
   e. Generate response
```

### Node.js Flow - Support Bot Analysis

```
1. Initialize:
   a. Load config
   b. Connect to Ollama
   c. Load vector store OR
   d. Index all tickets from data/tickets/:
      - Load YAML
      - Extract metadata
      - Generate embedding
      - Store record
      - Save vector store
2. API Server:
   a. POST /api/analyze receives ticket
   b. Generate embedding from ticket
   c. Search similar past tickets
   d. Call support expert agent with:
      - Expert context (md)
      - Current ticket
      - RAG results
   e. Parse analysis result
   f. Return JSON response
```

---

## 7️⃣ Modèles et Services: Mapping

| Responsabilité | Go (D&D) | Node.js (Bot) |
|---|---|---|
| **Markdown parsing** | chunks.SplitMarkdownBySections() | - (pas utilisé) |
| **YAML parsing** | - (pas utilisé) | ticketLoader.load() |
| **Metadata extraction** | metadataExtractorAgent | metadataExtractorAgent ✅ |
| **Embeddings** | ragAgent.SaveEmbedding() | embeddingService.generate() |
| **Vector store** | rag.Agent (internal) | VectorStore class |
| **Similarity search** | ragAgent.SearchTopN() | similaritySearch.search() |
| **Chat/Analysis** | chat.Agent | supportExpertAgent |
| **API/Interface** | prompt (TUI) | express (REST) |

---

## 8️⃣ Erreurs Communes à Éviter

### From Go Experience:

1. ❌ **Ne pas enrichir les chunks avant embedding**
   - ✅ Solution: Injecter metadata dans le texte avant embedding

2. ❌ **Threshold trop haut** (résultats vides)
   - ✅ Solution: Commencer à 0.4 comme en Go

3. ❌ **Oublier le top-N limit**
   - ✅ Solution: Toujours limiter à N résultats

4. ❌ **Ne pas persister le vector store**
   - ✅ Solution: Sauvegarder après indexation

5. ❌ **Context explosion**
   - ✅ Solution: Go utilise compressor; Node.js peut utiliser summarization

---

## 9️⃣ Avantages Node.js pour ce Use Case

| Aspect | Avantage |
|--------|----------|
| **Speed** | Démarrage plus rapide, idéal pour API |
| **Ecosystem** | NPM a excellentes libs pour YAML, REST, etc. |
| **Async/Await** | Gestion I/O plus naturelle pour API |
| **JSON** | Premier ordre pour parsing/manipulation |
| **Frontend** | Même stack si future web UI |
| **Deployment** | Plus léger que Go pour serveur API |

---

## 🔟 Checklist de Transposition

- [ ] Copier architecture RAG (metadata + enrichment)
- [ ] Adapter structures Go → TypeScript interfaces
- [ ] Réutiliser même logic similarité (cosine)
- [ ] Même seuil threshold (0.4) et topN (5)
- [ ] Même format JSON pour vector store
- [ ] Structured agents même pattern
- [ ] Tests avec données similaires

---

## 📚 Fichiers de Référence Importants

**Dans le repo (Go):**
- `05-solve-the-context-size-problem/03/README.md` - Architecture RAG + metadata
- `05/.../rag.agent.go` - Implementation RAG
- `05/.../interactive.chat.go` - Chat loop avec RAG

**À créer (Node.js):**
- `src/rag/similaritySearch.ts` - Même pattern que Go
- `src/rag/vectorStore.ts` - Même structure JSON
- `src/llm/agents/metadataExtractorAgent.ts` - Même logic extraction

---

**Conclusion:** La transposition de Go à Node.js est directe grâce aux patterns proven du projet D&D. Les concepts RAG, metadata enrichment, et structured agents restent identiques; seule l'implémentation technique change.
