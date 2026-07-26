# Getting Started

## Installation

### Homebrew

```shell
brew tap Marcel2603/tap
brew install spikeball-league
```

### From Release

Download the latest binary from the [Releases](https://github.com/Marcel2603/spikeball-league/releases)
page and place it in your `$PATH`.

### Docker

The recommended way to run the service in production is via Docker:

```bash
docker pull ghcr.io/marcel2603/spikeball-league/spikeball-league:latest
```

Run with a custom config file:

```bash
docker run \
  -p 3000:3000 \
  -v $PWD/app.yml:/app/app.yml \
  ghcr.io/marcel2603/spikeball-league/spikeball-league:latest
```

Or with custom branding assets mounted:

```bash
docker run \
  -p 3000:3000 \
  -v $PWD/app.yml:/app/app.yml \
  -v $PWD/custom:/app/custom \
  ghcr.io/marcel2603/spikeball-league/spikeball-league:latest
```

### From Source

```bash
go install github.com/Marcel2603/spikeball-league@latest
```

This installs `spikeball-league` into your `$GOPATH/bin` or `$GOBIN`.

## First Run

### Start the service

```bash
go mod tidy
make generate   # only needed once, or when framework versions change
make run
```

The service is now available at [http://localhost:3000](http://localhost:3000).
