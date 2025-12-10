DROP TABLE IF EXISTS audit_logs CASCADE;
DROP TABLE IF EXISTS game_player_stats CASCADE;
DROP TABLE IF EXISTS match_games CASCADE;
DROP TABLE IF EXISTS matches CASCADE;
DROP TABLE IF EXISTS tournament_registrations CASCADE;
DROP TABLE IF EXISTS tournaments CASCADE;
DROP TABLE IF EXISTS squad_members CASCADE;
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
