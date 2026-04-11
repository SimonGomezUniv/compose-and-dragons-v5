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
# Chat with a NPC

## Chat instructions

`dnd.chat.system.instructions.md`

## Conversational memory

- We keep all the messages exchanged during the conversation in memory to provide context for future interactions.
- This helps the NPC respond in a manner consistent with previous exchanges.
- ✋ But be **cautious**, as too much context can lead to <span class="indianred">**performance issues**</span> or <span class="indianred">**exceed token limits**</span>.