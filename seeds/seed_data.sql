-- Seed driver using gofakeit-generated SQL.
-- Usage:
--   go run ./seeds/generate_seeds.go --out seeds/generated_seed.sql
--   psql "postgresql://postgres:postgres@localhost:5432/cyber_tournament" -f seeds/generated_seed.sql

\echo 'Run go run ./seeds/generate_seeds.go --out seeds/generated_seed.sql before executing this file.'
\i seeds/generated_seed.sql
