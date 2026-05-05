# envault

A local secrets manager that encrypts .env files using age encryption and syncs them across team members via a shared backend.

## Installation

```bash
go install github.com/yourorg/envault@latest
```

## Usage

```bash
# Initialize envault in your project
envault init

# Encrypt and push your .env file
envault push

# Pull and decrypt the latest .env file
envault pull
```

## Requirements

- Go 1.21+
- A configured backend (S3, GCS, or local filesystem)
