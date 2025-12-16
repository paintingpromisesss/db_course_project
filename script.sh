#!/usr/bin/env bash
set -euo pipefail

BASE="http://localhost:8000/api"
H="Content-Type: application/json"
TS=$(date +%s)

call() {
  local method=$1 url=$2 data=${3:-}
  echo "### $method $url"
  if [ -n "$data" ]; then
    curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X "$method" "$url" -H "$H" -d "$data"
  else
    curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X "$method" "$url"
  fi
  echo
}

# Health
call GET "http://localhost:8000/health"

# Disciplines
call POST "$BASE/batch-import/disciplines" "[
  {\"code\":\"disc-$TS\",\"name\":\"Discipline $TS\",\"description\":\"Test\"}
]"

# Teams (discipline_id=1 пример)
call POST "$BASE/batch-import/teams" "[
  {\"name\":\"Team$TS\",\"tag\":\"T$TS\",\"country_code\":\"US\",\"discipline_id\":1}
]"

# Players
call POST "$BASE/batch-import/players" "[
  {\"nickname\":\"Player$TS\",\"country_code\":\"US\"}
]"

# Tournaments
call POST "$BASE/batch-import/tournaments" "[
  {\"discipline_id\":1,\"name\":\"Cup $TS\",\"start_date\":\"2025-03-01\",\"end_date\":\"2025-03-05\",\"status\":\"upcoming\",\"is_online\":true}
]"

# Tournament registrations (нужны существующие tournament_id=1, team_id=1 — подправь при необходимости)
call POST "$BASE/batch-import/tournament-registrations" "[
  {\"tournament_id\":1,\"team_id\":1,\"status\":\"confirmed\"}
]"

# Matches (нужен tournament_id=1)
call POST "$BASE/batch-import/matches" "[
  {\"tournament_id\":1,\"start_time\":\"2025-04-01T12:00:00Z\",\"format\":\"bo3\"}
]"

# Match games (нужен match_id=1)
call POST "$BASE/batch-import/match-games" "[
  {\"match_id\":1,\"map_name\":\"Map$TS\",\"game_number\":1}
]"

# Game player stats (нужен game_id=1, player_id=1)
call POST "$BASE/batch-import/game-player-stats" "[
  {\"game_id\":1,\"player_id\":1,\"kills\":10,\"deaths\":5,\"assists\":3}
]"

# Squad members (нужен team_id=1, player_id=1)
call POST "$BASE/batch-import/squad-members" "[
  {\"team_id\":1,\"player_id\":1,\"role\":\"Tester\"}
]"

# Team profiles (нужен team_id=1)
call POST "$BASE/batch-import/team-profiles" "[
  {\"team_id\":1,\"website\":\"https://example.com/$TS\"}
]"