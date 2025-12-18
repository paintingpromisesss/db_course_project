package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/brianvoe/gofakeit/v7"
)

type Discipline struct {
	ID          int
	Name        string
	Code        string
	Description string
	TeamSize    int
	Metadata    string
}

type Team struct {
	ID           int
	Name         string
	Tag          string
	Country      string
	DisciplineID int
	LogoURL      string
	WorldRanking float64
	IsVerified   bool
}

type TeamProfile struct {
	TeamID       int
	CoachName    string
	SponsorInfo  string
	Headquarters string
	Website      string
	ContactEmail string
}

type Player struct {
	ID        int
	Nickname  string
	RealName  string
	Country   string
	BirthDate time.Time
	SteamID   string
	AvatarURL string
	MMR       float64
	IsRetired bool
}

type SquadMember struct {
	ID       int
	TeamID   int
	PlayerID int
	Role     string
	JoinDate time.Time
	Salary   float64
}

type Tournament struct {
	ID            int
	DisciplineID  int
	Name          string
	StartDate     time.Time
	EndDate       time.Time
	PrizePool     float64
	Status        string
	IsOnline      bool
	BracketConfig string
}

type Registration struct {
	ID           int
	TournamentID int
	TeamID       int
	Seed         int
	Status       string
	ManagerEmail string
	RosterJSON   string
	IsInvited    bool
}

type Match struct {
	ID           int
	TournamentID int
	Team1ID      int
	Team2ID      int
	StartTime    time.Time
	Format       string
	Stage        string
	WinnerID     int
	IsForfeit    bool
}

type Game struct {
	ID           int
	MatchID      int
	MapName      string
	GameNumber   int
	DurationSec  int
	WinnerTeamID int
	ScoreTeam1   int
	ScoreTeam2   int
	StartedAt    time.Time
	TechPause    bool
}

type Stat struct {
	ID       int
	GameID   int
	PlayerID int
	TeamID   int
	Kills    int
	Deaths   int
	Assists  int
	Hero     string
	Damage   int
	Gold     int
	IsMVP    bool
}

func main() {
	var (
		outPath      string
		playersCount int
		teamsCount   int
		tournaments  int
	)

	flag.StringVar(&outPath, "out", "seeds/generated_seed.sql", "Path to write SQL seed file")
	flag.IntVar(&playersCount, "players", 1500, "Number of players")
	flag.IntVar(&teamsCount, "teams", 40, "Number of teams")
	flag.IntVar(&tournaments, "tournaments", 20, "Number of tournaments")
	flag.Parse()

	gofakeit.Seed(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())

	disciplines := buildDisciplines()
	teams := buildTeams(teamsCount, disciplines)
	profiles := buildTeamProfiles(teams)
	players := buildPlayers(playersCount)
	squad := buildSquadMembers(players, teams)
	tournamentsData := buildTournaments(tournaments, disciplines)
	registrations := buildRegistrations(tournamentsData, teams)
	matches := buildMatches(tournamentsData, registrations)
	games := buildGames(matches)
	stats := buildStats(games, matches, squad)

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		panic(err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	writeHeader(f)
	writeDisciplines(f, disciplines)
	writeTeams(f, teams)
	writeTeamProfiles(f, profiles)
	writePlayers(f, players)
	writeSquadMembers(f, squad)
	writeTournaments(f, tournamentsData)
	writeRegistrations(f, registrations)
	writeMatches(f, matches)
	writeGames(f, games)
	writeStats(f, stats)
	writeSequenceResets(f)
	writeFooter(f)

	fmt.Printf("Seed generated: %s\n", outPath)
}

func buildDisciplines() []Discipline {
	return []Discipline{
		{ID: 1, Name: "Counter-Strike 2", Code: "CS2", Description: "5v5 tactical FPS", TeamSize: 5, Metadata: `{"map_pool":["Inferno","Mirage","Nuke","Ancient"]}`},
		{ID: 2, Name: "Dota 2", Code: "DOTA2", Description: "5v5 MOBA", TeamSize: 5, Metadata: `{"map":"Ancient"}`},
		{ID: 3, Name: "Valorant", Code: "VAL", Description: "5v5 tac-shooter", TeamSize: 5, Metadata: `{"map_pool":["Ascent","Bind","Haven","Icebox"]}`},
	}
}

func buildTeams(count int, disciplines []Discipline) []Team {
	countries := []string{"US", "BR", "DE", "SE", "PL", "FR", "ES", "RU", "UA", "CN"}
	teams := make([]Team, 0, count)
	// Track used tags per discipline to avoid violating (tag, discipline_id) unique constraint
	usedTags := make(map[int]map[string]bool)
	for i := 1; i <= count; i++ {
		disc := disciplines[(i-1)%len(disciplines)]
		if usedTags[disc.ID] == nil {
			usedTags[disc.ID] = make(map[string]bool)
		}
		name := fmt.Sprintf("%s %s", gofakeit.Company(), gofakeit.JobDescriptor())
		tag := ""
		for {
			tag = strings.ToUpper(gofakeit.LetterN(3))
			if !usedTags[disc.ID][tag] {
				usedTags[disc.ID][tag] = true
				break
			}
		}
		teams = append(teams, Team{
			ID:           i,
			Name:         name,
			Tag:          tag,
			Country:      countries[rand.Intn(len(countries))],
			DisciplineID: disc.ID,
			LogoURL:      fmt.Sprintf("https://cdn.example.com/logos/%d.png", i),
			WorldRanking: round2(rand.Float64() * 100),
			IsVerified:   rand.Float64() < 0.35,
		})
	}
	return teams
}

func buildTeamProfiles(teams []Team) []TeamProfile {
	profiles := make([]TeamProfile, 0, len(teams))
	for _, t := range teams {
		profiles = append(profiles, TeamProfile{
			TeamID:       t.ID,
			CoachName:    fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.LastName()),
			SponsorInfo:  fmt.Sprintf("%s / %s", gofakeit.Company(), gofakeit.BeerName()),
			Headquarters: fmt.Sprintf("%s, %s", gofakeit.City(), gofakeit.Country()),
			Website:      fmt.Sprintf("https://%s.example.com", strings.ToLower(strings.ReplaceAll(t.Tag, " ", ""))),
			ContactEmail: fmt.Sprintf("contact+%d@%s.com", t.ID, strings.ToLower(t.Tag)),
		})
	}
	return profiles
}

func buildPlayers(count int) []Player {
	countries := []string{"US", "BR", "DE", "SE", "PL", "FR", "ES", "RU", "UA", "CN"}
	players := make([]Player, 0, count)
	for i := 1; i <= count; i++ {
		nick := fmt.Sprintf("%s_%d", sanitize(gofakeit.Gamertag()), i)
		players = append(players, Player{
			ID:        i,
			Nickname:  nick,
			RealName:  fmt.Sprintf("%s %s", gofakeit.FirstName(), gofakeit.LastName()),
			Country:   countries[rand.Intn(len(countries))],
			BirthDate: gofakeit.DateRange(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC), time.Date(2006, 12, 31, 0, 0, 0, 0, time.UTC)),
			SteamID:   fmt.Sprintf("STEAM_%X", rand.Uint64()),
			AvatarURL: fmt.Sprintf("https://cdn.example.com/avatars/%d.png", i),
			MMR:       round1(4000 + rand.Float64()*4000),
			IsRetired: rand.Float64() < 0.05,
		})
	}
	return players
}

func buildSquadMembers(players []Player, teams []Team) []SquadMember {
	roles := []string{"Player", "Support", "Carry", "Coach", "IGL", "Analyst"}
	squad := make([]SquadMember, 0, len(players))

	teamIndex := 0
	for i, p := range players {
		team := teams[teamIndex]
		squad = append(squad, SquadMember{
			ID:       i + 1,
			TeamID:   team.ID,
			PlayerID: p.ID,
			Role:     roles[rand.Intn(len(roles))],
			JoinDate: gofakeit.DateRange(time.Now().AddDate(-2, 0, 0), time.Now()),
			Salary:   round2(1500 + rand.Float64()*5000),
		})
		teamIndex = (teamIndex + 1) % len(teams)
	}
	return squad
}

func buildTournaments(count int, disciplines []Discipline) []Tournament {
	tournaments := make([]Tournament, 0, count)
	base := time.Now().AddDate(0, -2, 0)
	statuses := []string{"Announced", "Ongoing", "Completed"}
	for i := 1; i <= count; i++ {
		disc := disciplines[(i-1)%len(disciplines)]
		start := base.AddDate(0, 0, i*3)
		end := start.AddDate(0, 0, 4)
		tournaments = append(tournaments, Tournament{
			ID:            i,
			DisciplineID:  disc.ID,
			Name:          fmt.Sprintf("%s Masters %d", disc.Code, 2025+i%2),
			StartDate:     start,
			EndDate:       end,
			PrizePool:     round2(50000 + rand.Float64()*200000),
			Status:        statuses[rand.Intn(len(statuses))],
			IsOnline:      rand.Float64() < 0.6,
			BracketConfig: `{"type":"double_elim","format":"bo3"}`,
		})
	}
	return tournaments
}

func buildRegistrations(tournaments []Tournament, teams []Team) []Registration {
	regs := []Registration{}
	id := 1
	for _, t := range tournaments {
		sameDisc := []Team{}
		for _, tm := range teams {
			if tm.DisciplineID == t.DisciplineID {
				sameDisc = append(sameDisc, tm)
			}
		}
		rand.Shuffle(len(sameDisc), func(i, j int) { sameDisc[i], sameDisc[j] = sameDisc[j], sameDisc[i] })
		regCount := 12
		if len(sameDisc) < regCount {
			regCount = len(sameDisc)
		}
		rosterMap := map[int][]int{}
		for _, tm := range sameDisc[:regCount] {
			rosterMap[tm.ID] = []int{}
		}
		seed := 1
		for _, tm := range sameDisc[:regCount] {
			regs = append(regs, Registration{
				ID:           id,
				TournamentID: t.ID,
				TeamID:       tm.ID,
				Seed:         seed,
				Status:       "Confirmed",
				ManagerEmail: fmt.Sprintf("manager+%d@%s.com", tm.ID, strings.ToLower(tm.Tag)),
				RosterJSON:   fmt.Sprintf(`{"team_id":%d,"seed":%d}`, tm.ID, seed),
				IsInvited:    rand.Float64() < 0.4,
			})
			id++
			seed++
		}
	}
	return regs
}

func buildMatches(tournaments []Tournament, regs []Registration) []Match {
	matches := []Match{}
	id := 1
	regsByTournament := map[int][]Registration{}
	for _, r := range regs {
		regsByTournament[r.TournamentID] = append(regsByTournament[r.TournamentID], r)
	}
	for _, t := range tournaments {
		regsForT := regsByTournament[t.ID]
		if len(regsForT) < 2 {
			continue
		}
		rand.Shuffle(len(regsForT), func(i, j int) { regsForT[i], regsForT[j] = regsForT[j], regsForT[i] })
		matchCount := 8
		for i := 0; i < matchCount && i+1 < len(regsForT); i++ {
			team1 := regsForT[i%len(regsForT)].TeamID
			team2 := regsForT[(i+1)%len(regsForT)].TeamID
			if team1 == team2 {
				continue
			}
			start := t.StartDate.Add(time.Duration(i) * 3 * time.Hour)
			winner := team1
			if rand.Float64() < 0.5 {
				winner = team2
			}
			matches = append(matches, Match{
				ID:           id,
				TournamentID: t.ID,
				Team1ID:      team1,
				Team2ID:      team2,
				StartTime:    start,
				Format:       "bo3",
				Stage:        "Group Stage",
				WinnerID:     winner,
				IsForfeit:    rand.Float64() < 0.02,
			})
			id++
		}
	}
	return matches
}

func buildGames(matches []Match) []Game {
	games := []Game{}
	id := 1
	for _, m := range matches {
		gameCount := 2
		if rand.Float64() < 0.5 {
			gameCount = 3
		}
		for g := 1; g <= gameCount; g++ {
			games = append(games, Game{
				ID:           id,
				MatchID:      m.ID,
				MapName:      fmt.Sprintf("Map %d", g),
				GameNumber:   g,
				DurationSec:  1800 + rand.Intn(1200),
				WinnerTeamID: m.WinnerID,
				ScoreTeam1:   rand.Intn(16),
				ScoreTeam2:   rand.Intn(16),
				StartedAt:    m.StartTime.Add(time.Duration(g-1) * time.Hour),
				TechPause:    rand.Float64() < 0.08,
			})
			id++
		}
	}
	return games
}

func buildStats(games []Game, matches []Match, squad []SquadMember) []Stat {
	stats := []Stat{}
	id := 1
	teamPlayers := map[int][]int{}
	for _, sm := range squad {
		teamPlayers[sm.TeamID] = append(teamPlayers[sm.TeamID], sm.PlayerID)
	}
	matchByID := map[int]Match{}
	for _, m := range matches {
		matchByID[m.ID] = m
	}
	for _, g := range games {
		m := matchByID[g.MatchID]
		candidates := append([]int{}, teamPlayers[m.Team1ID]...)
		candidates = append(candidates, teamPlayers[m.Team2ID]...)
		if len(candidates) == 0 {
			continue
		}
		rand.Shuffle(len(candidates), func(i, j int) { candidates[i], candidates[j] = candidates[j], candidates[i] })
		take := 10
		if len(candidates) < take {
			take = len(candidates)
		}
		for _, pid := range candidates[:take] {
			stats = append(stats, Stat{
				ID:       id,
				GameID:   g.ID,
				PlayerID: pid,
				TeamID:   pickTeam(pid, m, squad),
				Kills:    rand.Intn(25),
				Deaths:   rand.Intn(15),
				Assists:  rand.Intn(30),
				Hero:     fmt.Sprintf("Hero-%d", rand.Intn(120)),
				Damage:   10000 + rand.Intn(60000),
				Gold:     8000 + rand.Intn(40000),
				IsMVP:    rand.Float64() < 0.05,
			})
			id++
		}
	}
	return stats
}

func pickTeam(playerID int, m Match, squad []SquadMember) int {
	for _, sm := range squad {
		if sm.PlayerID == playerID {
			if sm.TeamID == m.Team1ID || sm.TeamID == m.Team2ID {
				return sm.TeamID
			}
		}
	}
	return m.Team1ID
}

func writeHeader(f *os.File) {
	f.WriteString("-- GENERATED BY gofakeit seed generator\n")
	f.WriteString("BEGIN;\n")
	f.WriteString("TRUNCATE TABLE game_player_stats, match_games, matches, tournament_registrations, tournaments, squad_members, team_profiles, players, teams, disciplines, audit_logs, batch_import_errors RESTART IDENTITY CASCADE;\n\n")
}

func writeFooter(f *os.File) {
	f.WriteString("\nCOMMIT;\n")
}

func writeSequenceResets(f *os.File) {
	f.WriteString("-- sync identity sequences to max(id)+1\n")
	seqs := []string{
		"disciplines_id_seq",
		"teams_id_seq",
		"players_id_seq",
		"squad_members_id_seq",
		"tournaments_id_seq",
		"tournament_registrations_id_seq",
		"matches_id_seq",
		"match_games_id_seq",
		"game_player_stats_id_seq",
		"audit_logs_id_seq",
		"batch_import_errors_id_seq",
	}
	for _, seq := range seqs {
		fmt.Fprintf(f, "SELECT setval('%s', (SELECT COALESCE(MAX(id),0)+1 FROM %s), false);\n", seq, strings.TrimSuffix(seq, "_id_seq"))
	}
	f.WriteString("\n")
}

func writeDisciplines(f *os.File, items []Discipline) {
	for _, d := range items {
		fmt.Fprintf(f, "INSERT INTO disciplines (id, name, code, description, team_size, metadata) OVERRIDING SYSTEM VALUE VALUES (%d, '%s', '%s', '%s', %d, '%s');\n",
			d.ID, esc(d.Name), esc(d.Code), esc(d.Description), d.TeamSize, esc(d.Metadata))
	}
	f.WriteString("\n")
}

func writeTeams(f *os.File, items []Team) {
	for _, t := range items {
		fmt.Fprintf(f, "INSERT INTO teams (id, name, tag, country_code, discipline_id, created_at, logo_url, world_ranking, is_verified) OVERRIDING SYSTEM VALUE VALUES (%d, '%s', '%s', '%s', %d, CURRENT_TIMESTAMP, '%s', %.2f, %t);\n",
			t.ID, esc(t.Name), esc(t.Tag), esc(t.Country), t.DisciplineID, esc(t.LogoURL), t.WorldRanking, t.IsVerified)
	}
	f.WriteString("\n")
}

func writeTeamProfiles(f *os.File, items []TeamProfile) {
	for _, p := range items {
		fmt.Fprintf(f, "INSERT INTO team_profiles (team_id, coach_name, sponsor_info, headquarters, website, contact_email) VALUES (%d, '%s', '%s', '%s', '%s', '%s');\n",
			p.TeamID, esc(p.CoachName), esc(p.SponsorInfo), esc(p.Headquarters), esc(p.Website), esc(p.ContactEmail))
	}
	f.WriteString("\n")
}

func writePlayers(f *os.File, items []Player) {
	for _, p := range items {
		fmt.Fprintf(f, "INSERT INTO players (id, nickname, real_name, country_code, birth_date, steam_id, avatar_url, mmr_rating, is_retired, created_at) OVERRIDING SYSTEM VALUE VALUES (%d, '%s', '%s', '%s', '%s', '%s', '%s', %.1f, %t, CURRENT_TIMESTAMP);\n",
			p.ID, esc(p.Nickname), esc(p.RealName), esc(p.Country), p.BirthDate.Format("2006-01-02"), esc(p.SteamID), esc(p.AvatarURL), p.MMR, p.IsRetired)
	}
	f.WriteString("\n")
}

func writeSquadMembers(f *os.File, items []SquadMember) {
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	for _, sm := range items {
		fmt.Fprintf(f, "INSERT INTO squad_members (id, team_id, player_id, role, join_date, salary_monthly) OVERRIDING SYSTEM VALUE VALUES (%d, %d, %d, '%s', '%s', %.2f);\n",
			sm.ID, sm.TeamID, sm.PlayerID, esc(sm.Role), sm.JoinDate.Format("2006-01-02"), sm.Salary)
	}
	f.WriteString("\n")
}

func writeTournaments(f *os.File, items []Tournament) {
	for _, t := range items {
		fmt.Fprintf(f, "INSERT INTO tournaments (id, discipline_id, name, start_date, end_date, prize_pool, currency, status, is_online, bracket_config) OVERRIDING SYSTEM VALUE VALUES (%d, %d, '%s', '%s', '%s', %.2f, 'USD', '%s', %t, '%s');\n",
			t.ID, t.DisciplineID, esc(t.Name), t.StartDate.Format("2006-01-02"), t.EndDate.Format("2006-01-02"), t.PrizePool, esc(t.Status), t.IsOnline, esc(t.BracketConfig))
	}
	f.WriteString("\n")
}

func writeRegistrations(f *os.File, items []Registration) {
	for _, r := range items {
		fmt.Fprintf(f, "INSERT INTO tournament_registrations (id, tournament_id, team_id, seed_number, status, manager_contact, roster_snapshot, is_invited, registered_at) OVERRIDING SYSTEM VALUE VALUES (%d, %d, %d, %d, '%s', '%s', '%s', %t, CURRENT_TIMESTAMP);\n",
			r.ID, r.TournamentID, r.TeamID, r.Seed, esc(r.Status), esc(r.ManagerEmail), esc(r.RosterJSON), r.IsInvited)
	}
	f.WriteString("\n")
}

func writeMatches(f *os.File, items []Match) {
	for _, m := range items {
		fmt.Fprintf(f, "INSERT INTO matches (id, tournament_id, team1_id, team2_id, start_time, format, stage, winner_team_id, is_forfeit) OVERRIDING SYSTEM VALUE VALUES (%d, %d, %d, %d, '%s', '%s', '%s', %d, %t);\n",
			m.ID, m.TournamentID, m.Team1ID, m.Team2ID, m.StartTime.Format(time.RFC3339), esc(m.Format), esc(m.Stage), m.WinnerID, m.IsForfeit)
	}
	f.WriteString("\n")
}

func writeGames(f *os.File, items []Game) {
	for _, g := range items {
		fmt.Fprintf(f, "INSERT INTO match_games (id, match_id, map_name, game_number, duration_seconds, winner_team_id, score_team1, score_team2, started_at, had_technical_pause, pick_ban_phase) OVERRIDING SYSTEM VALUE VALUES (%d, %d, '%s', %d, %d, %d, %d, %d, '%s', %t, '%s');\n",
			g.ID, g.MatchID, esc(g.MapName), g.GameNumber, g.DurationSec, g.WinnerTeamID, g.ScoreTeam1, g.ScoreTeam2, g.StartedAt.Format(time.RFC3339), g.TechPause, esc(`{"note":"auto"}`))
	}
	f.WriteString("\n")
}

func writeStats(f *os.File, items []Stat) {
	for _, s := range items {
		fmt.Fprintf(f, "INSERT INTO game_player_stats (id, game_id, player_id, team_id, kills, deaths, assists, hero_name, damage_dealt, gold_earned, was_mvp) OVERRIDING SYSTEM VALUE VALUES (%d, %d, %d, %d, %d, %d, %d, '%s', %d, %d, %t);\n",
			s.ID, s.GameID, s.PlayerID, s.TeamID, s.Kills, s.Deaths, s.Assists, esc(s.Hero), s.Damage, s.Gold, s.IsMVP)
	}
	f.WriteString("\n")
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "'", "")
	return s
}

func esc(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func round2(v float64) float64 { return float64(int(v*100)) / 100 }
func round1(v float64) float64 { return float64(int(v*10)) / 10 }
