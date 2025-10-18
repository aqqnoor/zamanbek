#!/usr/bin/env bash
set -euo pipefail

HOST="${HOST:-http://localhost:8080}"

echo "== health =="
curl -fsS -i "$HOST/healthz" | head -n 1

echo "== new session =="
curl -fsS -c cookies.txt "$HOST/session/new" | jq .

echo "== ping =="
curl -fsS "$HOST/llm/ping" | jq .

echo "== chat (RU) =="
curl -fsS -b cookies.txt -X POST "$HOST/llm/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"Объясните, что такое мурабаха и максимальная сумма"}' | jq .

echo "== chat (KZ) =="
curl -fsS -b cookies.txt -X POST "$HOST/llm/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"Бұл карта халал ма?"}' | jq .

echo "== chat (EN) =="
curl -fsS -b cookies.txt -X POST "$HOST/llm/chat" \
  -H "Content-Type: application/json" \
  -d '{"message":"Is your deposit halal?"}' | jq .
