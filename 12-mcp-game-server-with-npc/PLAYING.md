# Comment lancer et jouer

## Prérequis

- **Docker Desktop** avec le **model runner** activé
- Les modèles suivants téléchargés dans Docker Desktop :
  - `hf.co/menlo/jan-nano-gguf:q4_k_m` (Dungeon Master / outils)
  - `hf.co/menlo/lucy-gguf:q4_k_m` (narrateur Sylph)
  - `ai/qwen2.5:1.5B-F16` (NPCs)
  - `ai/embeddinggemma:latest` (RAG embeddings)
  - `ai/qwen2.5:0.5B-F16` (compresseur de contexte)
  - `huggingface.co/menlo/jan-nano-gguf:q4_k_m` (extraction de métadonnées)

---

## Structure du projet

```
12-mcp-game-server-with-npc/
├── compose.yml                  ← orchestration principale
├── dungeon-master.yaml          ← configuration docker-agent (Dungeon Master)
├── dungeon-mcp-server/          ← serveur MCP de la logique de jeu
├── non-player-characters/       ← services NPC (Nain, Elfe, Humain, Sphinx)
│   └── sheets/                  ← fiches de personnage (.md + .json)
└── docker-agent/                ← sources docker-agent (image dungeon-master)
```

---

## Démarrage

### 1. Lancer tous les services

```bash
docker compose up --build -d
```

Cela démarre dans l'ordre :
1. `dungeon-mcp-server` (logique de jeu, port 6060)
2. `mcp-gateway` (sécurisation MCP, port 9011)
3. `dwarf-warrior`, `elf-mage`, `human-rogue`, `sphinx-end-of-level-boss` (NPCs, ports 9091-9094)
4. `dungeon-master` (docker-agent, se connecte aux services ci-dessus)
5. `mcp-inspector` (interface de debug MCP, port 6274)

> Le démarrage complet peut prendre quelques minutes le temps que les modèles
> soient chargés par le model runner.

### 2. Vérifier que les services sont prêts

```bash
# Santé du serveur MCP
curl http://localhost:6060/health

# Santé d'un NPC (le Nain)
curl http://localhost:9091/health

# Spec OpenAPI d'un NPC (générée dynamiquement depuis la fiche de perso)
curl http://localhost:9091/openapi.json | jq .
```

---

## Jouer

### Ouvrir une session interactive avec le Dungeon Master

```bash
docker compose attach dungeon-master
```

ou, si vous préférez relancer une nouvelle session :

```bash
docker compose run --rm dungeon-master run --yolo /work/dungeon-master.yaml
```

Vous êtes maintenant en conversation avec **Zephyr**, le Dungeon Master IA.

---

## Commandes de jeu (à taper dans le chat)

Zephyr comprend le langage naturel. Voici des exemples de ce que vous pouvez lui demander.

### Explorer le donjon

```
Montre-moi la carte du donjon.
Décris la salle dans laquelle je suis.
Quel est l'état actuel de la partie ?
```

### Se déplacer

```
Déplace-toi vers le nord.
Va vers l'est, puis regarde autour de toi.
Explore toutes les salles accessibles.
```

### Combattre

```
Commence le combat contre le monstre.
Attaque le gobelin.
Bois une potion pour récupérer de la vie.
```

### Collecter des objets

```
Collecte les objets dans cette salle.
Montre mon inventaire.
```

### Parler aux NPCs

Lorsque vous êtes dans une salle avec un NPC, Zephyr appelle automatiquement
le service NPC via l'outil OpenAPI généré depuis sa fiche de personnage.

```
Parle au NPC présent dans cette salle.
Demande au Nain s'il connaît un passage secret.
```

Une fois que Zephyr a initié le contact, vous pouvez converser directement
avec le NPC en lui posant vos questions. Le NPC répond en restant dans son
personnage (RAG sur sa fiche + historique de conversation).

### Le Sphinx — fin de niveau

Pour passer le Sphinx, vous devez trouver et fournir les trois mots secrets
des trois autres NPCs (Nain, Elfe, Humain). Chaque NPC connaît un mot secret.

```
Pose la question au Sphinx pour accéder au niveau suivant.
Donne au Sphinx les mots secrets : <mot1> <mot2> <mot3>
```

Si les mots sont corrects, le Sphinx vous accordera le passage. 🎉

### Sauvegarder la partie

```
Sauvegarde la partie.
```

### Quitter

Tapez `/bye` ou fermez le terminal (`Ctrl+C`).

---

## Rapports narratifs (Sylph)

Après chaque action significative, Zephyr délègue à **Sylph** (sub-agent `buddy`)
pour générer un rapport narratif illustré en markdown. Sylph raconte l'aventure
avec humour et emojis.

Pour demander un rapport explicitement :

```
Génère un rapport de l'aventure jusqu'ici.
Fais un résumé de ce qui s'est passé.
```

---

## Interface de debug MCP

L'**MCP Inspector** est disponible sur [http://localhost:6274](http://localhost:6274).

Il permet de tester manuellement les outils du `dungeon-mcp-server` sans passer
par le Dungeon Master.

**Configuration de connexion :**
- Transport Type : `Streamable HTTP`
- URL : `http://localhost:6060/mcp` (connexion directe au serveur MCP)
  ou `http://localhost:9011/mcp` (via la gateway)

---

## API directe des NPCs

Chaque service NPC expose trois endpoints utilisables indépendamment :

| Endpoint | Méthode | Description |
|----------|---------|-------------|
| `/health` | GET | Santé du service |
| `/models` | GET | Modèles utilisés |
| `/openapi.json` | GET | Spec OpenAPI (auto-générée depuis la fiche) |
| `/chat` | POST | Conversation non-streaming (JSON) |
| `/completion` | POST | Conversation streaming (SSE) |

**Exemple : conversation directe avec le Nain (port 9091)**

```bash
curl -X POST http://localhost:9091/chat \
  -H "Content-Type: application/json" \
  -d '{"message": "Bonjour Thorin, tu connais un passage secret ?"}'
```

Réponse :
```json
{"response": "Harrumph ! Un aventurier curieux... Oui, je connais bien des passages..."}
```

**Exemple : stream SSE (format natif nova-sdk)**

```bash
curl -N -X POST http://localhost:9091/completion \
  -H "Content-Type: application/json" \
  -H "Accept: text/event-stream" \
  -d '{"data": {"message": "Parle-moi de ton clan."}}'
```

---

## Arrêter les services

```bash
# Arrêter proprement
docker compose down

# Arrêter et supprimer les volumes (repart de zéro)
docker compose down -v
```

---

## Dépannage

### Le Dungeon Master ne démarre pas

Vérifiez que le model runner Docker Desktop est actif et que les modèles
`jan-nano-gguf` et `lucy-gguf` sont disponibles :

```bash
curl http://model-runner.docker.internal/engines/llama.cpp/v1/models
```

Si cette commande échoue depuis votre machine, vérifiez l'accès depuis l'intérieur
d'un container :

```bash
docker run --rm alpine/curl \
  curl http://model-runner.docker.internal/engines/llama.cpp/v1/models
```

### Un NPC ne répond pas / openapi.json introuvable

Le NPC doit avoir chargé sa fiche de personnage au démarrage. Vérifiez les logs :

```bash
docker compose logs dwarf-warrior
```

L'endpoint `/openapi.json` n'est disponible qu'une fois le serveur NPC
entièrement démarré (la goroutine d'enregistrement s'exécute quelques
millisecondes après le démarrage du serveur HTTP). Si vous appelez trop tôt,
attendez quelques secondes et réessayez.

### Réinitialiser l'historique de conversation d'un NPC

```bash
curl -X POST http://localhost:9091/memory/reset
```

### Voir la taille du contexte d'un NPC

```bash
curl http://localhost:9091/memory/messages/context-size
```
