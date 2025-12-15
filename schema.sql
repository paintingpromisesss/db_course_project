DROP VIEW IF EXISTS v_player_career_stats CASCADE;
DROP VIEW IF EXISTS v_match_results CASCADE;
DROP VIEW IF EXISTS v_active_rosters CASCADE;

DROP FUNCTION IF EXISTS fn_tournament_standings(INT) CASCADE;
DROP FUNCTION IF EXISTS fn_player_kda(INT) CASCADE;
DROP FUNCTION IF EXISTS refresh_team_rating(INT) CASCADE;
DROP FUNCTION IF EXISTS audit_log_changes() CASCADE;

DROP TABLE IF EXISTS batch_import_errors CASCADE;
DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS game_player_stats CASCADE;
DROP TABLE IF EXISTS match_games CASCADE;
DROP TABLE IF EXISTS matches CASCADE;
DROP TABLE IF EXISTS tournament_registrations CASCADE;
DROP TABLE IF EXISTS tournaments CASCADE;
DROP TABLE IF EXISTS squad_members CASCADE;
DROP TABLE IF EXISTS team_profiles CASCADE;
DROP TABLE IF EXISTS players CASCADE;
DROP TABLE IF EXISTS teams CASCADE;
DROP TABLE IF EXISTS disciplines CASCADE;

-- ==========================================
-- 1. disciplines
-- ==========================================
CREATE TABLE disciplines (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                    -- [INT]
    name VARCHAR(100) NOT NULL UNIQUE,                                  -- [VARCHAR]
    code VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,                                                   -- [TEXT]
    icon_url VARCHAR(255),
    team_size INT DEFAULT 5 CHECK (team_size > 0),
    is_active BOOLEAN DEFAULT TRUE,                                     -- [BOOLEAN]
    metadata JSONB DEFAULT '{}'::jsonb                                  -- [JSONB] (доп. настройки игры)
);

-- ==========================================
-- 2. teams
-- ==========================================
CREATE TABLE teams (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                    -- [INT]
    name VARCHAR(100) NOT NULL,                                          -- [VARCHAR]
    tag VARCHAR(10) NOT NULL,
    country_code CHAR(2) NOT NULL,                                       -- [CHAR]
    discipline_id INT NOT NULL REFERENCES disciplines(id) ON DELETE RESTRICT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,       -- [TIMESTAMP]
    logo_url VARCHAR(255),
    world_ranking DECIMAL(5,2) DEFAULT 0.00,                             -- [DECIMAL] (рейтинг команды 0-100)
    is_verified BOOLEAN DEFAULT FALSE,                                   -- [BOOLEAN] (верификация организации)
    
    CONSTRAINT uq_team_tag_discipline UNIQUE (tag, discipline_id)
);
CREATE INDEX idx_teams_discipline ON teams(discipline_id);

-- ==========================================
-- 2a. team_profiles (1:1 доп. сведения)
-- ==========================================
CREATE TABLE team_profiles (
    team_id INT PRIMARY KEY REFERENCES teams(id) ON DELETE CASCADE,      -- [INT] 1:1 с teams
    coach_name VARCHAR(100),                                             -- [VARCHAR]
    sponsor_info TEXT,                                                   -- [TEXT]
    headquarters VARCHAR(150),                                           -- [VARCHAR]
    website VARCHAR(255),                                                -- [VARCHAR]
    contact_email VARCHAR(150)                                           -- [VARCHAR]
);

-- ==========================================
-- 3. players
-- ==========================================
CREATE TABLE players (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                    -- [INT]
    nickname VARCHAR(50) NOT NULL UNIQUE,                                -- [VARCHAR]
    real_name VARCHAR(100),
    country_code CHAR(2),                                                -- [CHAR]
    birth_date DATE,                                                     -- [DATE]
    steam_id VARCHAR(32) UNIQUE,
    avatar_url VARCHAR(255),
    mmr_rating DECIMAL(7,1) DEFAULT 0.0,                                 -- [DECIMAL] (MMR/ELO рейтинг)
    is_retired BOOLEAN DEFAULT FALSE,                                    -- [BOOLEAN]
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP        -- [TIMESTAMP]
);

-- ==========================================
-- 4. squad_members
-- ==========================================
CREATE TABLE squad_members (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                 -- [BIGINT/INT]
    team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,        -- [INT]
    player_id INT NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL DEFAULT 'Player',                          -- [VARCHAR]
    is_standin BOOLEAN DEFAULT FALSE,                                    -- [BOOLEAN]
    join_date DATE NOT NULL DEFAULT CURRENT_DATE,                        -- [DATE]
    contract_end_date DATE,
    leave_date DATE,
    salary_monthly DECIMAL(10, 2),                                       -- [DECIMAL] (зарплата в USD)
    
    CONSTRAINT chk_dates CHECK (leave_date IS NULL OR leave_date >= join_date)
);
CREATE INDEX idx_squad_active ON squad_members(team_id) WHERE leave_date IS NULL;

-- ==========================================
-- 5. tournaments
-- ==========================================
CREATE TABLE tournaments (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                    -- [INT]
    discipline_id INT NOT NULL REFERENCES disciplines(id) ON DELETE RESTRICT,
    name VARCHAR(200) NOT NULL,                                          -- [VARCHAR]
    start_date DATE NOT NULL,                                            -- [DATE]
    end_date DATE NOT NULL,
    prize_pool DECIMAL(15, 2) DEFAULT 0,                                 -- [DECIMAL]
    currency VARCHAR(3) DEFAULT 'USD',
    status VARCHAR(20) NOT NULL DEFAULT 'Announced',
    is_online BOOLEAN DEFAULT FALSE,                                     -- [BOOLEAN] (онлайн/оффлайн)
    bracket_config JSONB,                                                -- [JSONB] (конфиг сетки: single/double elim)
    
    CONSTRAINT chk_tournament_dates CHECK (end_date >= start_date)
);

-- ==========================================
-- 6. tournament_registrations
-- ==========================================
CREATE TABLE tournament_registrations (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                 -- [BIGINT/INT]
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE, -- [INT]
    team_id INT NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    seed_number INT,
    status VARCHAR(20) DEFAULT 'Pending',                                -- [VARCHAR]
    manager_contact VARCHAR(100),
    roster_snapshot JSONB,                                               -- [JSONB]
    is_invited BOOLEAN DEFAULT FALSE,                                    -- [BOOLEAN] (прямой инвайт vs квалификация)
    registered_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,    -- [TIMESTAMP]
    
    CONSTRAINT uq_tournament_team UNIQUE (tournament_id, team_id)
);
CREATE INDEX idx_registrations_tournament ON tournament_registrations(tournament_id);
CREATE INDEX idx_registrations_team ON tournament_registrations(team_id);

-- ==========================================
-- 7. matches
-- ==========================================
CREATE TABLE matches (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                 -- [BIGINT/INT]
    tournament_id INT NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE, -- [INT]
    team1_id INT REFERENCES teams(id) ON DELETE SET NULL,
    team2_id INT REFERENCES teams(id) ON DELETE SET NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,                        -- [TIMESTAMP]
    format VARCHAR(10) NOT NULL DEFAULT 'bo3',                           -- [VARCHAR]
    stage VARCHAR(50),
    winner_team_id INT REFERENCES teams(id),
    is_forfeit BOOLEAN DEFAULT FALSE,                                    -- [BOOLEAN] (техническое поражение)
    match_notes JSONB,                                                   -- [JSONB] (дополнительные данные: паузы, протесты)
    
    CONSTRAINT chk_different_teams CHECK (team1_id <> team2_id)
);
CREATE INDEX idx_matches_tournament ON matches(tournament_id);
CREATE INDEX idx_matches_start_time ON matches(start_time);

-- ==========================================
-- 8. match_games
-- ==========================================
CREATE TABLE match_games (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                 -- [BIGINT/INT]
    match_id BIGINT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,  -- [BIGINT]
    map_name VARCHAR(100) NOT NULL,                                      -- [VARCHAR]
    game_number INT NOT NULL,                                            -- [INT]
    duration_seconds INT,
    winner_team_id INT REFERENCES teams(id),
    score_team1 INT DEFAULT 0,
    score_team2 INT DEFAULT 0,
    started_at TIMESTAMP WITH TIME ZONE,                                 -- [TIMESTAMP] (реальное время начала)
    had_technical_pause BOOLEAN DEFAULT FALSE,                           -- [BOOLEAN]
    pick_ban_phase JSONB,                                                -- [JSONB] (пики/баны героев)
    
    CONSTRAINT uq_match_game_number UNIQUE (match_id, game_number)
);
CREATE INDEX idx_match_games_match ON match_games(match_id);

-- ==========================================
-- 9. game_player_stats
-- ==========================================
CREATE TABLE game_player_stats (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                 -- [BIGINT/INT]
    game_id BIGINT NOT NULL REFERENCES match_games(id) ON DELETE CASCADE, -- [BIGINT]
    player_id INT NOT NULL REFERENCES players(id) ON DELETE CASCADE,    -- [INT]
    team_id INT REFERENCES teams(id),
    kills INT DEFAULT 0,
    deaths INT DEFAULT 0,
    assists INT DEFAULT 0,
    hero_name VARCHAR(100),                                              -- [VARCHAR]
    damage_dealt INT DEFAULT 0,
    gold_earned INT DEFAULT 0,
    kda_ratio DECIMAL(5, 2) GENERATED ALWAYS AS                          -- [DECIMAL] (вычисляемое поле)
        (CASE WHEN deaths = 0 THEN (kills + assists)::decimal 
              ELSE (kills + assists)::decimal / deaths END) STORED,
    was_mvp BOOLEAN DEFAULT FALSE,                                       -- [BOOLEAN] (MVP карты)
    
    CONSTRAINT uq_game_player UNIQUE (game_id, player_id)
);
CREATE INDEX idx_stats_player ON game_player_stats(player_id);
CREATE INDEX idx_stats_game ON game_player_stats(game_id);

-- ==========================================
-- 10. audit_logs
-- ==========================================
CREATE TABLE audit_logs (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,                 -- [BIGINT/INT]
    table_name VARCHAR(50) NOT NULL,                                     -- [VARCHAR]
    record_id BIGINT NOT NULL,
    operation VARCHAR(10) NOT NULL,
    old_value JSONB,                                                     -- [JSONB]
    new_value JSONB,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,       -- [TIMESTAMP]
    changed_by VARCHAR(100),
    is_sensitive BOOLEAN DEFAULT FALSE                                   -- [BOOLEAN] (флаг чувствительных данных)
);
CREATE INDEX idx_audit_date ON audit_logs USING BRIN (changed_at);

-- ==========================================
-- 10b. batch_import_errors (логирование загрузок)
-- ==========================================
CREATE TABLE batch_import_errors (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    source VARCHAR(50) NOT NULL,
    row_data JSONB,
    error_message TEXT NOT NULL,
    occurred_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_import_errors_source ON batch_import_errors(source, occurred_at);

-- ==========================================
-- 11. Триггер для аудита
-- ==========================================
CREATE OR REPLACE FUNCTION audit_log_changes() RETURNS trigger AS $$
DECLARE
    v_new JSONB;
    v_old JSONB;
    v_pk  BIGINT;
BEGIN
    IF TG_OP IN ('INSERT', 'UPDATE') THEN
        v_new := to_jsonb(NEW);
    END IF;
    IF TG_OP IN ('UPDATE', 'DELETE') THEN
        v_old := to_jsonb(OLD);
    END IF;

    v_pk := COALESCE(
        (v_new ->> 'id')::BIGINT,
        (v_old ->> 'id')::BIGINT,
        (v_new ->> 'team_id')::BIGINT,
        (v_old ->> 'team_id')::BIGINT
    );

    IF TG_OP = 'INSERT' THEN
        INSERT INTO audit_logs(table_name, record_id, operation, new_value, changed_by)
        VALUES (TG_TABLE_NAME, v_pk, TG_OP, v_new, current_user);
        RETURN NEW;
    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO audit_logs(table_name, record_id, operation, old_value, new_value, changed_by)
        VALUES (TG_TABLE_NAME, v_pk, TG_OP, v_old, v_new, current_user);
        RETURN NEW;
    ELSE
        INSERT INTO audit_logs(table_name, record_id, operation, old_value, changed_by)
        VALUES (TG_TABLE_NAME, v_pk, TG_OP, v_old, current_user);
        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_disciplines_audit AFTER INSERT OR UPDATE OR DELETE ON disciplines
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_teams_audit AFTER INSERT OR UPDATE OR DELETE ON teams
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_team_profiles_audit AFTER INSERT OR UPDATE OR DELETE ON team_profiles
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_players_audit AFTER INSERT OR UPDATE OR DELETE ON players
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_squad_audit AFTER INSERT OR UPDATE OR DELETE ON squad_members
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_tournaments_audit AFTER INSERT OR UPDATE OR DELETE ON tournaments
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_registrations_audit AFTER INSERT OR UPDATE OR DELETE ON tournament_registrations
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_matches_audit AFTER INSERT OR UPDATE OR DELETE ON matches
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_match_games_audit AFTER INSERT OR UPDATE OR DELETE ON match_games
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();
CREATE TRIGGER trg_player_stats_audit AFTER INSERT OR UPDATE OR DELETE ON game_player_stats
FOR EACH ROW EXECUTE FUNCTION audit_log_changes();

-- ==========================================
-- 12. Агрегирующая функция рейтинга команды
-- ==========================================
CREATE OR REPLACE FUNCTION refresh_team_rating(p_team_id INT) RETURNS VOID AS $$
DECLARE
    v_avg_rating DECIMAL(7,2);
    v_scaled DECIMAL(7,2);
BEGIN
    SELECT AVG(p.mmr_rating) INTO v_avg_rating
    FROM squad_members sm
    JOIN players p ON p.id = sm.player_id
    WHERE sm.team_id = p_team_id
      AND sm.leave_date IS NULL;

    v_scaled := CASE
        WHEN v_avg_rating IS NULL THEN 0
        ELSE LEAST(ROUND(v_avg_rating / 100, 2), 1000)
    END;

    UPDATE teams
    SET world_ranking = v_scaled
    WHERE id = p_team_id;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_refresh_team_rating() RETURNS trigger AS $$
BEGIN
    PERFORM refresh_team_rating(COALESCE(NEW.team_id, OLD.team_id));
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION trg_refresh_team_rating_on_player() RETURNS trigger AS $$
DECLARE
    v_team_id INT;
BEGIN
    FOR v_team_id IN
        SELECT sm.team_id
        FROM squad_members sm
        WHERE sm.player_id = NEW.id
          AND sm.leave_date IS NULL
    LOOP
        PERFORM refresh_team_rating(v_team_id);
    END LOOP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_squad_rating_refresh
AFTER INSERT OR UPDATE OR DELETE ON squad_members
FOR EACH ROW EXECUTE FUNCTION trg_refresh_team_rating();

CREATE TRIGGER trg_player_rating_refresh
AFTER UPDATE OF mmr_rating ON players
FOR EACH ROW EXECUTE FUNCTION trg_refresh_team_rating_on_player();

-- ==========================================
-- 13. Функции и представления для отчетов
-- ==========================================
CREATE OR REPLACE FUNCTION fn_player_kda(p_player_id INT)
RETURNS DECIMAL(10,2) AS $$
DECLARE
    v_kills BIGINT;
    v_assists BIGINT;
    v_deaths BIGINT;
BEGIN
    SELECT COALESCE(SUM(kills),0), COALESCE(SUM(assists),0), COALESCE(SUM(deaths),0)
    INTO v_kills, v_assists, v_deaths
    FROM game_player_stats
    WHERE player_id = p_player_id;

    IF v_deaths = 0 THEN
        RETURN (v_kills + v_assists)::DECIMAL(10,2);
    END IF;
    RETURN ((v_kills + v_assists)::DECIMAL) / v_deaths;
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE FUNCTION fn_tournament_standings(p_tournament_id INT)
RETURNS TABLE (
    team_id INT,
    matches_played INT,
    wins INT,
    losses INT,
    forfeits INT
) AS $$
BEGIN
    RETURN QUERY
    WITH teams_in_matches AS (
        SELECT m.id AS match_id, m.tournament_id, m.team1_id AS team_id, m.winner_team_id, m.is_forfeit
        FROM matches m
        WHERE m.tournament_id = p_tournament_id
        UNION ALL
        SELECT m.id, m.tournament_id, m.team2_id AS team_id, m.winner_team_id, m.is_forfeit
        FROM matches m
        WHERE m.tournament_id = p_tournament_id
    )
    SELECT team_id,
           COUNT(*) AS matches_played,
           COUNT(*) FILTER (WHERE winner_team_id = team_id) AS wins,
           COUNT(*) FILTER (WHERE winner_team_id IS NOT NULL AND winner_team_id <> team_id) AS losses,
           COUNT(*) FILTER (WHERE is_forfeit) AS forfeits
    FROM teams_in_matches
    WHERE team_id IS NOT NULL
    GROUP BY team_id
    ORDER BY wins DESC, losses ASC;
END;
$$ LANGUAGE plpgsql STABLE;

CREATE OR REPLACE VIEW v_active_rosters AS
SELECT t.id AS team_id,
       t.name AS team_name,
       t.tag,
       p.id AS player_id,
       p.nickname,
       p.country_code,
       sm.role,
       sm.join_date
FROM squad_members sm
JOIN teams t ON t.id = sm.team_id
JOIN players p ON p.id = sm.player_id
WHERE sm.leave_date IS NULL;

CREATE OR REPLACE VIEW v_match_results AS
SELECT m.id AS match_id,
       m.tournament_id,
       m.start_time,
       m.stage,
       m.format,
       m.winner_team_id,
       COUNT(g.id) AS games_played,
       SUM(g.score_team1) AS total_score_team1,
       SUM(g.score_team2) AS total_score_team2
FROM matches m
LEFT JOIN match_games g ON g.match_id = m.id
GROUP BY m.id, m.tournament_id, m.start_time, m.stage, m.format, m.winner_team_id;

CREATE OR REPLACE VIEW v_player_career_stats AS
SELECT p.id AS player_id,
       p.nickname,
       COALESCE(SUM(s.kills),0) AS kills,
       COALESCE(SUM(s.deaths),0) AS deaths,
       COALESCE(SUM(s.assists),0) AS assists,
       COALESCE(SUM(s.damage_dealt),0) AS damage,
       COALESCE(SUM(s.gold_earned),0) AS gold,
       fn_player_kda(p.id) AS kda
FROM players p
LEFT JOIN game_player_stats s ON s.player_id = p.id
GROUP BY p.id, p.nickname;
