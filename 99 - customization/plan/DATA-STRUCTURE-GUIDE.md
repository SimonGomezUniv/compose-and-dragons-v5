# 📋 Guide des Données de Test - Bot de Support

**Objectif:** Définir la structure des tickets YAML et des données contextuelles pour l'indexation RAG

---

## 1️⃣ Structure des Tickets YAML

### 1.1 Format Standard

**Chemin:** `data/tickets/ticket-XXX.yml`

```yaml
# Ticket de support - Format standard
id: "TICKET-001"
title: "Cannot login to account"
description: |
  User reports being unable to login to their account.
  Tried resetting password but still getting 401 error.
  Last successful login was 3 days ago.
  
priority: high
category: account  # Can be auto-filled by RAG
status: open
created_at: "2026-04-15T10:30:00Z"
customer:
  name: "John Doe"
  email: "john.doe@example.com"
  account_id: "ACC-12345"

# Tags pour recherche rapide
tags:
  - login
  - account
  - 401-error
  - password

# Contexte additionnel
environment:
  browser: "Chrome"
  os: "Windows 10"
  ip_location: "US"

# Résolution si existante
resolution: null  # ou texte de résolution
resolved_at: null
```

### 1.2 Variations Structurelles

Vous pouvez adapter selon vos besoins, mais garder ce noyau:

```yaml
# Minimal (si données limitées)
id: "TICKET-002"
title: "Payment failed"
description: "Stripe payment rejected during checkout"
priority: high
tags: [billing, payment, error]

---

# Enrichi (data supplémentaire)
id: "TICKET-003"
title: "Feature request: Dark mode"
description: |
  Users requesting dark mode in the application.
  Multiple tickets mentioning eye strain.
priority: low
category: feature-request
related_tickets: ["TICKET-089", "TICKET-102"]
customer_feedback_count: 23

---

# Technique (détails système)
id: "TICKET-004"
title: "Database connection timeout"
description: "Production database returning 504 errors"
priority: critical
severity: production-outage
affected_systems:
  - api-server
  - worker-queue
  - reporting-module
error_code: "DB_TIMEOUT_5000"
```

---

## 2️⃣ Ensemble de Données de Test

### 2.1 Exemple Complet - 5 Tickets Initiaux

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
created_at: "2026-04-15T10:30:00Z"
customer:
  name: "John Doe"
  email: "john.doe@example.com"
tags: [login, account, 401-error]
resolution: "Password reset token was expired. Issued new reset link via email."
resolved_at: "2026-04-15T14:22:00Z"
```

**Fichier:** `data/tickets/ticket-002.yml`

```yaml
id: "TICKET-002"
title: "Billing issue: Duplicate charges"
description: |
  Customer charged twice for same subscription.
  Charges occurred on April 10 and April 12.
  Both transactions completed successfully.
  Customer requests refund for duplicate charge.
  
priority: high
category: billing
created_at: "2026-04-14T08:15:00Z"
customer:
  name: "Sarah Johnson"
  email: "sarah.j@company.com"
tags: [billing, duplicate-charge, refund]
resolution: "Issued refund of $49.99 for duplicate charge. Investigated webhook retry issue."
resolved_at: "2026-04-14T16:45:00Z"
```

**Fichier:** `data/tickets/ticket-003.yml`

```yaml
id: "TICKET-003"
title: "API returning 500 error on /users endpoint"
description: |
  The /api/v1/users endpoint is returning 500 errors.
  This is affecting production.
  Error started ~2 hours ago.
  Database connection seems to be working fine.
  
priority: critical
category: technical
created_at: "2026-04-16T11:00:00Z"
tags: [api, 500-error, endpoint, production]
resolution: "Memory leak found in user service. Redeployed with memory limit fix. Added monitoring."
resolved_at: "2026-04-16T12:30:00Z"
environment:
  affected_endpoints: ["/api/v1/users", "/api/v1/users/{id}"]
  error_rate: "87%"
```

**Fichier:** `data/tickets/ticket-004.yml`

```yaml
id: "TICKET-004"
title: "Feature request: Export data to CSV"
description: |
  Multiple customers are asking for ability to export their data.
  Currently only JSON export is available.
  CSV format would improve integration with Excel and other tools.
  
priority: low
category: feature-request
created_at: "2026-04-10T13:20:00Z"
tags: [export, csv, feature-request]
related_customer_requests: 12
status: backlog
```

**Fichier:** `data/tickets/ticket-005.yml`

```yaml
id: "TICKET-005"
title: "Password reset email not arriving"
description: |
  User requested password reset but email didn't arrive.
  Checked spam folder - email not there.
  Tried resetting again - same issue.
  Other users report same problem.
  
priority: high
category: account
created_at: "2026-04-16T14:45:00Z"
customer:
  name: "Michael Chen"
  email: "m.chen@startup.io"
tags: [password, email, account]
resolution: "Email service provider had delivery issues. Escalated to email team. Manually reset user password."
resolved_at: "2026-04-16T15:20:00Z"
```

---

## 3️⃣ Contexte d'Expertise Support

### 3.1 Fichier Principal: `data/context/support-expert.md`

```markdown
# Support Expert Knowledge Base

## 1. Catégories de Tickets et Escalades

### Account Issues
- **Login failures**: Usually password reset or account lock
- **Permission issues**: Check user role and access level
- **Account deletion requests**: Follow GDPR compliance
- **SLA**: 4 hours

### Billing Issues
- **Duplicate charges**: Check webhook logs and payment history
- **Failed payments**: Verify card expiration and fraud blocks
- **Refund requests**: Process if within 30-day window
- **SLA**: 2 hours (critical for business)

### Technical Issues
- **API errors (5xx)**: Check server logs and database
- **Connection timeouts**: Verify network and resource limits
- **Data corruption**: Assess impact, consider rollback
- **SLA**: 1 hour (production outages)

### Feature Requests
- **Enhancement proposals**: Forward to product team
- **Integration requests**: Assess feasibility
- **Custom development**: Direct to sales
- **SLA**: 24 hours for acknowledgment

## 2. Common Resolution Patterns

### Pattern: Password Issues
1. Verify account exists and is active
2. Check if account is locked (> 5 failed attempts)
3. Send password reset link
4. If email fails, use alternative verification
5. Document in CRM

### Pattern: Billing Reconciliation
1. Pull transaction history from payment processor
2. Identify duplicate/failed charges
3. Map to subscription records
4. Process refund if needed (< 30 days)
5. Investigate root cause (webhook retry, race condition)

### Pattern: Production Outage
1. Confirm issue is affecting users
2. Check monitoring alerts and logs
3. Identify affected systems
4. Apply hotfix if available
5. Escalate to engineering if needed
6. Set customer expectations on timeline

## 3. Tools and Resources

- **CRM**: Available at crm.internal.com
- **Payment Dashboard**: stripe.internal.com
- **Server Logs**: logs.internal.com
- **Monitoring**: datadog.internal.com
- **Escalation**: #support-escalations on Slack

## 4. First Response Time Target

- **Critical**: 15 minutes
- **High**: 1 hour
- **Medium**: 4 hours
- **Low**: 24 hours

## 5. Customer Communication Templates

### Problem Confirmation
"Thank you for reporting this issue. I've confirmed [issue description].
I'm investigating and will update you within [timeframe]."

### Resolution Provided
"I've resolved your issue by [action taken].
Please test and confirm it's working. Let me know if you need anything else."

### Escalation Needed
"This requires engineering review. I've escalated to our technical team.
You can expect an update within [timeframe]."

## 6. Useful Queries and Checks

### Account Status Check
```
SELECT * FROM users WHERE email = ?
SELECT * FROM user_roles WHERE user_id = ?
SELECT * FROM login_attempts WHERE user_id = ? ORDER BY created_at DESC
```

### Billing Check
```
SELECT * FROM transactions WHERE customer_id = ?
SELECT * FROM subscriptions WHERE customer_id = ?
```

## 7. Common Error Codes

| Code | Meaning | Solution |
|------|---------|----------|
| 401 | Unauthorized | Password reset or token refresh |
| 403 | Forbidden | Check permissions, update role |
| 404 | Not found | Verify resource ID, check soft deletes |
| 429 | Rate limited | Wait 60s, check API quota |
| 500 | Server error | Check logs, restart service if needed |
| 504 | Gateway timeout | Database/service latency issue |

## 8. Escalation Criteria

Escalate to Engineering if:
- Production outage affecting > 10% users
- Data corruption or loss
- Security vulnerability
- Custom development needed

Escalate to Product if:
- Feature request consensus from > 5 customers
- Integration with third-party service
- Major UX issue

Escalate to Legal if:
- Data deletion request (GDPR)
- Terms of service violation
- Potential lawsuit mentioned
```

---

## 4️⃣ Structure du Répertoire Data

```
data/
├── tickets/
│   ├── ticket-001.yml          # Account login
│   ├── ticket-002.yml          # Billing duplicate charge
│   ├── ticket-003.yml          # Technical API error
│   ├── ticket-004.yml          # Feature request
│   ├── ticket-005.yml          # Email delivery
│   ├── ticket-006.yml          # (à ajouter)
│   └── ...
├── context/
│   └── support-expert.md       # Knowledge base (créé ci-dessus)
└── vectorStore.json            # Auto-généré lors indexation

# Format du vectorStore.json après indexation:
{
  "vuuid-1": {
    "id": "vuuid-1",
    "ticketId": "TICKET-001",
    "content": "Cannot login to account. User reports...",
    "embedding": [0.123, 0.456, ...],  # 1024 dimensions pour embeddinggemma
    "metadata": {
      "keywords": ["login", "account", "password", "401"],
      "mainCategory": "account",
      "subcategory": "login-failure",
      "importance": "high"
    },
    "timestamp": "2026-04-28T10:00:00Z"
  },
  "vuuid-2": { ... },
  ...
}
```

---

## 5️⃣ Guide d'Expansion des Données

### Phase 1: Minimum Viable
- 3-5 tickets d'exemple (déjà fourni ci-dessus)
- Couvre les 3 principales catégories (account, billing, technical)
- Simule réalisme avec résolutions

### Phase 2: Enrichissement (10-20 tickets)
- Ajouter variations par catégorie
- Inclure edge cases
- Augmenter diversité de langage

### Phase 3: Production (100+ tickets)
- Importer historique réel (anonymisé)
- Vérifier qualité données
- Ajouter métadonnées supplémentaires

### Données à Ajouter Progressivement

```yaml
# Exemple: Ticket avec contexte additionnel
id: "TICKET-006"
title: "Slow report generation"
description: |
  Custom analytics report taking > 5 minutes to generate.
  Previously took 30 seconds.
  Performance degradation started yesterday.
  
priority: medium
category: technical
sub_category: performance
performance_metrics:
  before_baseline: "30 seconds"
  current_time: "320 seconds"
  degradation: "10.6x slower"
  
root_cause: "Added JOIN to accounts table without index"
resolution: "Created composite index on (report_id, account_id)"
```

---

## 6️⃣ Exemple de Requête d'Analyse (Inférence)

### Input API

```json
{
  "ticket": {
    "id": "TICKET-NEW-001",
    "title": "Forgot password but email not arriving",
    "description": "I clicked forgot password, received no email",
    "priority": "medium",
    "customer": {
      "email": "new.user@example.com",
      "account_id": "ACC-99999"
    }
  }
}
```

### Flux RAG Attendu

```
1. Générer embedding pour: "Forgot password but email not arriving"

2. Chercher tickets similaires:
   - TICKET-001 (login issue, reset) → 0.82
   - TICKET-005 (password reset email) → 0.91 ✅ Très pertinent
   - TICKET-002 (billing) → 0.45
   
3. Injecter contexte RAG:
   "Similar past issue: TICKET-005 had password reset email not arriving.
    Resolution: Email service provider issue, manually reset password."

4. Support Expert Agent analyse avec contexte:
   - Catégorie suggérée: "account"
   - Confiance: 0.96
   - Suggestions:
     1. "Check if email is in spam folder"
     2. "Manually reset password via CRM"
   - Infos supplémentaires:
     - "Browser and OS for compatibility check"
     - "Current account status (active/locked)"
```

### Expected Output

```json
{
  "ticketId": "TICKET-NEW-001",
  "suggestedCategory": "account",
  "confidence": 0.96,
  "suggestions": [
    "Check spam/junk folder for reset email",
    "Manually reset password via customer account in CRM",
    "If repeated, check with email provider for delivery issues"
  ],
  "additionalInfoNeeded": [
    "Has customer checked spam folder?",
    "Is their email domain whitelisted?"
  ],
  "reasoning": "Similar to TICKET-005 which was email delivery issue. Pattern suggests email service problem.",
  "relatedPastTickets": [
    {
      "id": "TICKET-005",
      "title": "Password reset email not arriving",
      "similarity": 0.91,
      "resolution": "Email service provider had delivery issues"
    },
    {
      "id": "TICKET-001",
      "title": "Cannot login to account",
      "similarity": 0.82,
      "resolution": "Password reset token was expired"
    }
  ]
}
```

---

## 7️⃣ Checklist de Préparation des Données

- [ ] Créer répertoire `data/tickets/`
- [ ] Créer 5 fichiers YAML d'exemple (ticket-001 à 005)
- [ ] Créer fichier `data/context/support-expert.md`
- [ ] Valider syntaxe YAML (aucune erreur de parsing)
- [ ] Vérifier que chaque ticket a au minimum:
  - [ ] id (unique)
  - [ ] title
  - [ ] description
  - [ ] priority
  - [ ] tags
- [ ] Test d'indexation RAG:
  - [ ] Tous les tickets load correctement
  - [ ] Embeddings générés sans erreur
  - [ ] vectorStore.json créé
- [ ] Test de recherche RAG:
  - [ ] Similarity search retourne résultats
  - [ ] Threshold 0.4 fonctionne
  - [ ] Top-5 résultats correctement triés

---

## 📝 Notes Importantes

1. **Format YAML**:
   - Indentation: 2 espaces
   - Pas de tabs
   - Strings multilines avec `|` ou `>`

2. **IDs**:
   - `id`: Format lisible (TICKET-XXX)
   - Autres IDs UUID ou UUID pour embedding

3. **Timestamps**:
   - Format ISO 8601: "2026-04-28T10:30:00Z"
   - Utiliser UTC

4. **Résolutions**:
   - Garder historique (important pour apprentissage RAG)
   - Inclure root cause si trouvée

5. **Tags**:
   - Minuscules, kebab-case
   - Utiliser pour requêtes futures

---

**Document créé:** 2026-04-28  
**Version:** 1.0  
**Statut:** Prêt pour implémentation
