#!/bin/bash

# Script to generate self-signed certificates for local HTTPS testing
# Run this script from the local/ directory

echo "Generating self-signed certificate for local HTTPS testing..."

# Generate private key
openssl genrsa -out server.key 2048

# Generate certificate signing request
openssl req -new -key server.key -out server.csr -subj "/C=FR/ST=IDF/L=Paris/O=Clariti Dev/OU=Development/CN=localhost"

# Generate self-signed certificate
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# Clean up CSR file
rm server.csr

echo "Certificate generated successfully!"
echo "- Private key: server.key"
echo "- Certificate: server.crt"
echo ""
echo "To use HTTPS, update your config.yaml with:"
echo "server:"
echo "  cert_file: \"local/server.crt\""
echo "  key_file: \"local/server.key\""
echo ""
echo "Note: This is a self-signed certificate for development only."
echo "Your browser will show a security warning that you can safely ignore for testing."
