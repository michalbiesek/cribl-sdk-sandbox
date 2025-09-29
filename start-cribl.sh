#!/bin/bash

# Start Cribl Stream Leader
echo "🚀 Starting Cribl Stream Leader..."
docker compose up -d cribl-leader

echo "⏳ Waiting for Cribl Stream to be ready..."
sleep 10

# Check if Cribl is running
if curl -s -f http://localhost:19000/api/v1/system/health > /dev/null 2>&1; then
    echo "✅ Cribl Stream Leader is running!"
    echo "🌐 Access Cribl Stream UI at: http://localhost:19000"
    echo "📊 Default credentials: admin/admin (change on first login)"
else
    echo "⚠️  Cribl Stream may still be starting up..."
    echo "🌐 Try accessing: http://localhost:19000 in a few moments"
fi

echo ""
echo "📝 To stop Cribl Stream: ./stop-cribl.sh"
echo "📝 To view logs: docker compose logs -f cribl-leader"
