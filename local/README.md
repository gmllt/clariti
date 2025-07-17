# Local Development & Testing

This directory contains files for testing and local development that should not be committed.

## Configuration

1. Copy `config.example.yaml` to `config.yaml` and customize as needed
2. To enable HTTPS, generate certificates with the provided script

## HTTPS certificate generation for testing

```bash
cd local/
./generate-certs.sh
```

This will generate:
- `server.key` - Private key
- `server.crt` - Self-signed certificate

Then in your `config.yaml`:

```yaml
server:
  host: "localhost"
  port: "8443"  # Alternative HTTPS standard port
  cert_file: "local/server.crt"
  key_file: "local/server.key"
```

## File structure

- `config.example.yaml` - Example configuration for development
- `generate-certs.sh` - Script to generate test certificates
- `server.crt` / `server.key` - Generated certificates (not committed)
- `config.yaml` - Your local configuration (not committed)

## Notes

- All files in this directory are ignored by Git (`.gitignore`)
- Generated certificates are self-signed and for testing only
- Use valid certificates in production
