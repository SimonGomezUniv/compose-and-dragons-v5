---
marp: true
theme: default
paginate: true
---
<style>
.dodgerblue {
  color: dodgerblue;
}
.indianred {
  color: indianred;
}
</style>
# Create NPC Characters

## 2 Kinds of completions

### Structured (name, race, class, gender)

- Creative configuration (`temperature: 0.7`, `topP: 0.9`, `topK: 40`)
- Generates JSON format output

### Chat (narrative character sheet == **story**)

- Highly creative configuration (`temperature: 0.8`, `topP: 0.95`)
- Produces complete narrative character sheet