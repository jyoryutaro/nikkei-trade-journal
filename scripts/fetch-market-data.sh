#!/usr/bin/env bash
# Calls the backend fetch endpoint for each symbol so the local scheduler
# (launchd / cron) can keep the DB up to date. Idempotent (server upserts).
#
# Env overrides:
#   BASE_URL         backend base URL          (default http://localhost:8080)
#   INTERNAL_SECRET  bearer token for /fetch   (default dev-secret)
#   SYMBOLS          space-separated symbols   (default "^N225 NKD=F")
set -uo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
INTERNAL_SECRET="${INTERNAL_SECRET:-dev-secret}"
SYMBOLS="${SYMBOLS:-^N225 NKD=F}"

urlencode() {
  local s="$1"
  s="${s//%/%25}"; s="${s//&/%26}"; s="${s//=/%3D}"; s="${s//^/%5E}"; s="${s// /%20}"
  printf '%s' "$s"
}

for sym in $SYMBOLS; do
  enc="$(urlencode "$sym")"
  ts="$(date '+%Y-%m-%dT%H:%M:%S%z')"
  if resp="$(curl -fsS -m 30 -X POST "${BASE_URL}/api/market-data/fetch?symbol=${enc}" \
        -H "Authorization: Bearer ${INTERNAL_SECRET}" 2>&1)"; then
    echo "${ts} ${sym} OK ${resp}"
  else
    echo "${ts} ${sym} ERROR ${resp}"
  fi
done
