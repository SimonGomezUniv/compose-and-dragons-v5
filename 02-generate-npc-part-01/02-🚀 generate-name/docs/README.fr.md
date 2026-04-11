# Générateur de Personnages NPC pour D&D

## Description

Ce programme génère automatiquement des personnages non-joueurs (NPC) pour Donjons & Dragons en utilisant un agent IA structuré basé sur le framework Nova SDK.

## Fonctionnement

Le programme crée un agent IA capable de générer des personnages D&D avec des informations structurées :
- **Prénom** et **Nom de famille** respectant les conventions de la race
- **Race** (Nain, Elfe, Humain)
- **Classe** (Guerrier, Mage, Rôdeur, Clerc, Voleur, Paladin, etc.)
- **Genre** (masculin/féminin)

### Architecture

```mermaid
graph TB
    A[Démarrage du Programme] --> B[Chargement des Règles de Nommage]
    B --> C[Chargement des Instructions Système]
    C --> D[Création de l'Agent Structuré]

    D --> E[Requête 1: Nain]
    D --> F[Requête 2: Elfe Ranger Féminine]
    D --> G[Requête 3: Paladin Humain Masculin]

    E --> H[Génération Structurée]
    F --> H
    G --> H

    H --> I[Affichage des Résultats]

    style A fill:#4A90E2,stroke:#2E5C8A,color:#fff
    style B fill:#7B68EE,stroke:#5A4BC2,color:#fff
    style C fill:#7B68EE,stroke:#5A4BC2,color:#fff
    style D fill:#50C878,stroke:#3A9B5C,color:#fff
    style E fill:#FF9500,stroke:#CC7700,color:#fff
    style F fill:#FF9500,stroke:#CC7700,color:#fff
    style G fill:#FF9500,stroke:#CC7700,color:#fff
    style H fill:#E74C3C,stroke:#C0392B,color:#fff
    style I fill:#9B59B6,stroke:#7D3C98,color:#fff
```

## Composants Principaux

### 1. Structure de Données
```go
type NPCCharacter struct {
    FirstName  string  // Prénom
    FamilyName string  // Nom de famille
    Race       string  // Race (Dwarf/Elf/Human)
    Class      string  // Classe D&D
    Gender     string  // Genre (male/female)
}
```

### 2. Base de Connaissances
- **Règles de nommage** (`dnd.naming.rules.md`) : Conventions de noms par race
- **Instructions système** (`dnd.system.instructions.md`) : Directives pour l'IA

### 3. Agent IA
- Agent Nova de type **structured** : `structured.NewAgent`
- Utilise le type `NPCCharacter` pour la génération structurée
- Utilise le modèle `nvidia_nemotron-mini-4b-instruct`
- Configuration créative (`temperature: 0.7`, `topP: 0.9`, `topK: 40`)
- Génère des sorties structurées au format JSON

## Flux d'Exécution

1. **Initialisation**
   - Lecture des règles de nommage D&D
   - Injection des règles dans les instructions système

2. **Création de l'Agent**
   - Configuration du modèle LLM
   - Définition du schéma de sortie structuré

3. **Génération des Personnages**
   - Traite 3 cas de test différents
   - Chaque requête génère un personnage complet
   - Affiche les résultats formatés

## Exemple de Sortie

```
🎲 Request 1: Generate a dwarf character
🔄 Generating NPC...

🧙 Generated NPC Summary:
Name       : Thorin Ironforge
Race       : Dwarf
Class      : Warrior
Gender     : male
```

## Technologies Utilisées

- **Langage** : Go
- **Framework** : Nova SDK
- **Modèle IA** : Nemotron Mini 4B (quantisé Q4_K_M)
- **Moteur** : Docker Model Runner avec endpoint llama.cpp (`http://localhost:12434/engines/llama.cpp/v1`)

## Exécution

```bash
go run main.go
```

Le programme génère automatiquement 3 personnages de test et affiche leurs caractéristiques.
