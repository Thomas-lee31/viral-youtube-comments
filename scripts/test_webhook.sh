#!/usr/bin/env bash
set -euo pipefail

# Load .env from project root
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ENV_FILE="$SCRIPT_DIR/../.env"

if [ -f "$ENV_FILE" ]; then
  export "$(grep -v '^#' "$ENV_FILE" | grep DISCORD_WEBHOOK_URL | xargs)"
fi

if [ -z "${DISCORD_WEBHOOK_URL:-}" ]; then
  echo "Error: DISCORD_WEBHOOK_URL is not set. Check your .env file."
  exit 1
fi

echo "Sending test message to Discord..."

HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
  -H "Content-Type: application/json" \
  -d '{
    "embeds": [{
      "title": "Webhook Test",
      "description": "If you see this, your Discord webhook is working correctly.",
      "color": 65280,
      "footer": {"text": "youtubeads - test ping"}
    }]
  }' \
  "$DISCORD_WEBHOOK_URL")

if [ "$HTTP_CODE" -ge 200 ] && [ "$HTTP_CODE" -lt 300 ]; then
  echo "Success! Discord returned HTTP $HTTP_CODE. Check your channel."
else
  echo "Failed. Discord returned HTTP $HTTP_CODE."
  exit 1
fi
