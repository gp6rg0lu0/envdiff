# envdiff

Compare `.env` files across environments and report missing, extra, or mismatched keys with optional redaction.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git
cd envdiff
go build -o envdiff .
```

---

## Usage

```bash
# Compare two .env files
envdiff .env.example .env.production

# Redact values in output
envdiff --redact .env.example .env.production

# Compare multiple environments against a base file
envdiff --base .env.example .env.staging .env.production
```

### Example Output

```
Missing keys in .env.production:
  - DATABASE_URL
  - REDIS_HOST

Extra keys in .env.production:
  + LEGACY_API_KEY

Mismatched keys:
  ~ APP_ENV  (expected: "development", got: [redacted])
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--redact` | Hide values in diff output |
| `--base` | Specify a base file to compare others against |
| `--json` | Output results as JSON |

---

## License

MIT © 2024 yourusername