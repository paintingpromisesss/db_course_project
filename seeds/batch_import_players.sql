-- Batch import players from CSV with error logging.
-- Usage:
--   psql "postgresql://postgres:postgres@localhost:5432/cyber_tournament" \
--        -v csv_file='seeds/batch_import_players.csv' -f seeds/batch_import_players.sql

\set ON_ERROR_STOP on
\set csv_file :csv_file

BEGIN;

-- Staging table
CREATE TEMP TABLE stg_players (
    nickname TEXT,
    real_name TEXT,
    country_code CHAR(2),
    birth_date DATE,
    steam_id TEXT,
    avatar_url TEXT,
    mmr_rating NUMERIC,
    is_retired BOOLEAN
);

-- Client-side copy so local path works
\copy stg_players FROM :'csv_file' WITH (FORMAT csv, HEADER true);

DO $$
DECLARE
    r stg_players%ROWTYPE;
BEGIN
    FOR r IN SELECT * FROM stg_players LOOP
        BEGIN
            INSERT INTO players (nickname, real_name, country_code, birth_date, steam_id, avatar_url, mmr_rating, is_retired)
            VALUES (r.nickname, r.real_name, r.country_code, r.birth_date, r.steam_id, r.avatar_url, COALESCE(r.mmr_rating, 0), COALESCE(r.is_retired, false))
            ON CONFLICT (nickname) DO UPDATE
                SET real_name   = EXCLUDED.real_name,
                    country_code = EXCLUDED.country_code,
                    birth_date  = EXCLUDED.birth_date,
                    steam_id    = EXCLUDED.steam_id,
                    avatar_url  = EXCLUDED.avatar_url,
                    mmr_rating  = EXCLUDED.mmr_rating,
                    is_retired  = EXCLUDED.is_retired;
        EXCEPTION WHEN OTHERS THEN
            INSERT INTO batch_import_errors (source, row_data, error_message)
            VALUES ('players_csv', to_jsonb(r), SQLERRM);
        END;
    END LOOP;
END $$;

COMMIT;
