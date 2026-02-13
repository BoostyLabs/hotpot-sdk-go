# HotPot Golang SDK Examples

This directory contains practical examples demonstrating how to use the HotPot SDK.

## Setup

Before running the examples, you can set your API key using environment variables:

```bash
export HOTPOT_API_KEY="your-api-key-here"
export HOTPOT_BASE_URL="https://api.hotpot.tech"  # Optional, defaults to this
```

## Running Examples

Run any example with:

```bash
go run examples/<example_name>/main.go
```

### Using .env file

Alternatively, you can use a `.env` file. Copy the provided `.env.example` to `<your_file>` and fill in the values:

```bash
cp examples/.env.example <your_file>
```

The examples will automatically load variables from the specified `.env` if `ENV_FILE` was set:

```bash
ENV_FILE=/path/to/your/.env go run examples/<example_name>/main.go
```

## API Documentation

For more details, see the [Hotpot API documentation](https://docs.hotpot.tech).