#!/bin/bash

# Start Cribl Stream Leader
echo "🚀 Starting Cribl Stream Leader..."
docker compose up -d cribl-leader

echo "⏳ Waiting for Cribl Stream to be ready..."
sleep 10

PORT=19000
# Derive the forwarded URL if running in GitHub Codespaces
# GITHUB_CODESPACES_PORT_FORWARDING_DOMAIN is usually "app.github.dev"
if [[ -n "${CODESPACE_NAME}" ]]; then
  DOMAIN="${GITHUB_CODESPACES_PORT_FORWARDING_DOMAIN:-app.github.dev}"
  FORWARDED_URL="https://${CODESPACE_NAME}-${PORT}.${DOMAIN}/"
else
  FORWARDED_URL="http://localhost:${PORT}"
fi

echo -n "⏳ Waiting for Cribl Stream Leader to start "
while true; do
    if curl -s -f "http://localhost:${PORT}/api/v1/health" > /dev/null 2>&1; then
        echo ""
        echo "✅ Cribl Stream Leader is running!"
        echo "🌐 Access Cribl Stream UI at: ${FORWARDED_URL}"
        echo "📊 Default credentials: admin/admin (change on first login)"
        break
    else
        echo -n "."
        sleep 1
    fi
done

echo ""
echo "📝 To stop Cribl Stream: ./stop-cribl.sh"
echo "📝 To view logs: docker compose logs -f cribl-leader"
