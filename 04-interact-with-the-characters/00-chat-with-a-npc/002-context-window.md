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
.seagreen {
  color: seagreen;
}
</style>

# <span class="dodgerblue">Context Size</span> / Context Window

## What is it?

The **total number of tokens** the model can process at once
- Think of it as the model's **short-term memory**
- Includes **everything**: system prompt, user messages, history, documents, AND the generated response

## Example

If a model has a **32k context window**:
```
input tokens + history tokens + output tokens ≤ 32,000 tokens
```

