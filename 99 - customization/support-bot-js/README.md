# 🤖 Support Bot - Quick Start

Bot intelligent de catégorisation de tickets de support avec RAG (Retrieval Augmented Generation).

---

## 🚀 Démarrage Rapide

### 1️⃣ Installation

```bash
npm install
```

### 2️⃣ Lancer le serveur

```bash
npm run dev
```

**Résultat attendu:**
```
╔════════════════════════════════════════╗
║   🤖 Support Bot Running               ║
║   Server: http://0.0.0.0:3000      
║   ✅ Loaded 5 tickets
║   🔍 Indexed 5 tickets
╚════════════════════════════════════════╝
```

Le serveur indexe automatiquement les tickets à chaque démarrage.

---

## 📂 Où Mettre les Documents à Indexer

### Structure des données

```
support-bot-js/
├── data/
│   ├── tickets/              ← Tickets à analyser (FORMAT YAML)
│   │   ├── ticket-001.yml
│   │   ├── ticket-002.yml
│   │   └── ... autres tickets
│   ├── context/
│   │   └── support-expert.md ← Contexte d'expertise
│   └── vectorStore.json      ← Index auto-généré (NE PAS ÉDITER)
```

### Format YAML des Tickets

**Fichier:** `data/tickets/ticket-001.yml`

```yaml
id: TICKET-001
title: Impossible de réinitialiser mot de passe
description: |
  Utilisateur signale que le bouton "Réinitialiser mot de passe" ne fonctionne pas.
  Aucun email n'est envoyé.
  Problème sur deux comptes différents.
category: authentification
priority: high
tags:
  - password-reset
  - email
  - critical-path
resolution: Vérifier le service SMTP, redéployer API d'authentification.
```

### Ajouter Vos Propres Tickets

1. Créez `data/tickets/ticket-NNN.yml` (remplacez NNN par le numéro)
2. Remplissez les champs: `id`, `title`, `description` (obligatoires)
3. Redémarrez le serveur → Il indexe automatiquement
4. Fichier `data/vectorStore.json` est créé/mis à jour automatiquement

---

## 🔍 Utiliser l'API

### Tester la santé du serveur

```bash
curl http://localhost:3000/api/health
```

**Réponse:**
```json
{
  "status": "ok",
  "vectorStoreSize": 5,
  "embeddingCacheSize": 4
}
```

### Analyser un ticket

```bash
curl -X POST http://localhost:3000/api/analyze \
  -H "Content-Type: application/json" \
  -d '{
    "id": "NEW-001",
    "title": "Application très lente",
    "description": "L'\''application met 5+ secondes à charger",
    "tags": ["performance", "slow"],
    "priority": "high"
  }'
```

**Réponse:**
```json
{
  "success": true,
  "ticketId": "NEW-001",
  "analysis": {
    "ticketId": "NEW-001",
    "suggestedCategory": "performance",
    "confidence": 0.92,
    "suggestions": [
      "Vérifier bundle size",
      "Profiler avec DevTools"
    ],
    "relatedPastTickets": [
      {
        "id": "TICKET-003",
        "title": "Interface gelée après mise à jour",
        "similarity": 0.87
      }
    ]
  }
}
```

---

## 📝 Contexte d'Expertise

### Fichier: `data/context/support-expert.md`

Contient les catégories de support et patterns de résolution.

**Format:**
```markdown
# Support Expert Knowledge Base

## 1. Authentification & Sécurité
**Subcategories:**
- Password reset
- 2FA issues
- Login failures

**Common Patterns:**
- Email delivery failures → Check SMTP
- Token expiration → Verify JWT settings
```

Le contexte est chargé au démarrage et utilisé pour l'analyse RAG.

---

## 📊 Endpoints Disponibles

| Méthode | URL | Description |
|---------|-----|-------------|
| `GET` | `/` | Info service |
| `GET` | `/api/health` | Status + nombre de documents indexés |
| `POST` | `/api/analyze` | Analyser un ticket |
| `GET` | `/api/stats` | Statistiques du store |

---

## 🔧 Architecture

```
Ticket YAML
    ↓
Chargement + Enrichissement
    ↓
Embedding (embeddinggemma)
    ↓
Vector Store (JSON)
    ↓
RAG Search (cosine similarity)
    ↓
LLM Analysis (qwen2 + expert context)
    ↓
Catégorisation + Suggestions
```

---

## 🐛 Dépannage

**Erreur: "Cannot find module"**
→ Lancez `npm install` à nouveau

**Erreur: "EADDRINUSE port 3000"**
→ Tuez les processus Node: `Get-Process node | Stop-Process -Force`

**Erreur: "Embedding error"**
→ Vérifiez que Ollama tourne: `docker-compose ps`

---

## 📚 Documentation Complète

Voir dossier `../plan/`:
- `GETTING-STARTED-JAVASCRIPT.md` - Guide détaillé jour par jour
- `PLAN-NODEJS-SUPPORT-BOT.md` - Architecture complète (11 phases)
- `DATA-STRUCTURE-GUIDE.md` - Format des données
- `GO-vs-NODEJS-COMPARISON.md` - Comparaison des patterns

---

## ✅ Checklist pour Démarrer

- [ ] `npm install` → dépendances installées
- [ ] `docker-compose up -d` → Ollama en marche
- [ ] `npm run dev` → serveur démarré
- [ ] Ajouter tickets dans `data/tickets/*.yml`
- [ ] `curl http://localhost:3000/api/health` → voir les documents indexés
- [ ] Envoyer requête POST `/api/analyze` → obtenir analyse

**C'est prêt! 🚀**
