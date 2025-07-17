#!/bin/bash

# Script to test HTTPS with self-signed certificates

echo "🔧 Generating self-signed certificates..."
cd local/
./generate-certs.sh

echo ""
echo "🚀 Starting HTTPS server..."
echo "⚠️  WARNING: Your browser will show a security warning"
echo "   This is normal with self-signed certificates for testing."
echo ""
echo "🌐 HTTPS test endpoints:"
echo "   • https://localhost:8443/health"
echo "   • https://localhost:8443/api"
echo "   • https://localhost:8443/api/v1/weather"
echo ""
echo "📝 Credentials: admin / https123"
echo ""
echo "⏹️  Ctrl+C to stop server"
echo ""

cd ..
CONFIG_PATH="local/config.https.yaml" ./clariti-server
