# ch — Companies House in your terminal

Fast, script-friendly CLI for the [Companies House API](https://developer.company-information.service.gov.uk/). JSON-first output, simple API key auth, and cross-platform binaries.

## Features

- **Company** — get company profiles and registered office addresses
- **Search** — search companies and officers
- **Officers** — list company officers (directors, secretaries, etc.)
- **Filing** — browse filing history, view individual filings
- **PSC** — persons with significant control
- **Charges** — mortgages and securities
- **Insolvency** — insolvency case information

## Install

### From source

```bash
go install github.com/anthonyencodeclub/ch/cmd/ch@latest
```

### From releases

Download the latest binary from [Releases](https://github.com/anthonyencodeclub/ch/releases).

## Quick start

```bash
# Set your API key (get one at https://developer.company-information.service.gov.uk/)
ch auth set-key YOUR_API_KEY

# Search for a company
ch search companies "OpenAI"

# Get a company profile
ch company get 00445790

# List officers
ch officers list 00445790

# Filing history
ch filing list 00445790

# JSON output (for scripting)
ch company get 00445790 --json
```

## Authentication

Get a free API key from the [Companies House Developer Hub](https://developer.company-information.service.gov.uk/).

```bash
# Option 1: Store in config
ch auth set-key YOUR_API_KEY

# Option 2: Environment variable
export CH_API_KEY=YOUR_API_KEY
```

## Output modes

| Flag | Description |
|------|-------------|
| (default) | Human-readable coloured output |
| `--json` | JSON to stdout |
| `--plain` | Stable, parseable text (no colours) |

## Environment variables

| Variable | Description |
|----------|-------------|
| `CH_API_KEY` | Companies House API key |
| `CH_JSON` | Default to JSON output (`1`/`true`) |
| `CH_PLAIN` | Default to plain output (`1`/`true`) |
| `CH_CONFIG_DIR` | Override config directory |

## License

MIT
