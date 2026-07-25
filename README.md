# Spikeball League

:warning: This project is under heavy development

[![Test & Lint](https://github.com/Marcel2603/spikeball-league/actions/workflows/go-test.yml/badge.svg)](https://github.com/Marcel2603/spikeball-league/actions/workflows/go-test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Marcel2603/spikeball-league)](https://goreportcard.com/report/github.com/Marcel2603/spikeball-league)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/Marcel2603/spikeball-league)](https://github.com/Marcel2603/spikeball-league/releases/latest)

## Features

- tbd


## Docker

```bash
docker pull ghcr.io/marcel2603/spikeball-league/spikeball-league:latest

docker run \
  -p 3000:3000 \
  -v $PWD/app.yml:/app/app.yml \
  ghcr.io/marcel2603/spikeball-league/spikeball-league:latest
```

## Configuration

Copy and edit `cmd/config/app.default.yml`:

```yaml
```

Full reference → [docs/Configuration](https://marcel2603.github.io/spikeball-league/getting-started/configuration/)

## Contributing

See [CONTRIBUTING.md](.github/CONTRIBUTING.md) and the [docs](https://marcel2603.github.io/spikeball-league/development/).

```bash
# Install pre-commit hooks
make init-precommit

# Run tests
make test

# Run linter
make lint
```

## Documentation

Full documentation is available at: <https://marcel2603.github.io/spikeball-league/>

## Acknowledgements

Project default images are generated using Gemini 3 Pro.
