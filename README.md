# envdiff

> Diff and reconcile `.env` files across environments with secret masking.

---

## Installation

```bash
go install github.com/yourusername/envdiff@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envdiff.git && cd envdiff && go build -o envdiff .
```

---

## Usage

Compare two `.env` files and mask secret values:

```bash
envdiff .env.development .env.production
```

**Example output:**

```
~ DATABASE_URL  [secret changed]
+ NEW_FEATURE_FLAG=true
- DEPRECATED_KEY=old_value
= APP_NAME=myapp
```

### Flags

| Flag | Description |
|------|-------------|
| `--mask` | Mask secret values in output (default: true) |
| `--export` | Export reconciled output to a file |
| `--format` | Output format: `text`, `json` (default: `text`) |

### Reconcile

Merge differences and write a reconciled file:

```bash
envdiff --export .env.reconciled .env.staging .env.production
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

[MIT](LICENSE) © 2024 yourusername