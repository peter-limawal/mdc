---
name: qmd
description: Search MDC's local Markdown roadmap, decisions, and runbooks with QMD. Use when a question may be answered by indexed project notes before searching the web.
compatibility: Requires the qmd CLI and the local mdc collection.
allowed-tools: Bash(qmd:*)
---

# QMD

Search the local MDC knowledge base before web research when the answer may be in project Markdown.

```bash
qmd search "exact terms" -c mdc -n 5
qmd query $'intent: Find the MDC decision about durable job state.\nlex: SQLite job store persistence\nvec: control-plane persistence decision' -c mdc -n 5
```

Retrieve sources after searching; do not answer from snippets alone:

```bash
qmd get "<path-or-docid>"
qmd multi-get "<docid-1>,<docid-2>" --format md
```

Use `qmd search` for exact terms. Use structured `qmd query` for conceptual recall after embeddings are available.

## Maintenance

Only update the local index when asked:

```bash
qmd update
qmd embed
```
