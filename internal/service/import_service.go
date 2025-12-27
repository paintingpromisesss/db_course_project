package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"db_course_project/internal/models"
)

type PlayerImportInput struct {
	Nickname    string  `json:"nickname" csv:"nickname"`
	RealName    string  `json:"real_name" csv:"real_name"`
	CountryCode string  `json:"country_code" csv:"country_code"`
	BirthDate   string  `json:"birth_date" csv:"birth_date"`
	SteamID     string  `json:"steam_id" csv:"steam_id"`
	AvatarURL   string  `json:"avatar_url" csv:"avatar_url"`
	MMR         float64 `json:"mmr_rating" csv:"mmr_rating"`
	IsRetired   *bool   `json:"is_retired" csv:"is_retired"`
}

type DisciplineImportInput struct {
	Code        string          `json:"code" csv:"code"`
	Name        string          `json:"name" csv:"name"`
	Description string          `json:"description" csv:"description"`
	IconURL     *string         `json:"icon_url" csv:"icon_url"`
	TeamSize    *int            `json:"team_size" csv:"team_size"`
	Metadata    json.RawMessage `json:"metadata" swaggertype:"object" csv:"metadata"`
	IsActive    *bool           `json:"is_active" csv:"is_active"`
}

type TeamImportInput struct {
	Name         string   `json:"name" csv:"name"`
	Tag          string   `json:"tag" csv:"tag"`
	CountryCode  string   `json:"country_code" csv:"country_code"`
	DisciplineID int64    `json:"discipline_id" csv:"discipline_id"`
	LogoURL      *string  `json:"logo_url" csv:"logo_url"`
	WorldRanking *float64 `json:"world_ranking" csv:"world_ranking"`
	IsVerified   *bool    `json:"is_verified" csv:"is_verified"`
}

type TournamentImportInput struct {
	DisciplineID  int64           `json:"discipline_id" csv:"discipline_id"`
	Name          string          `json:"name" csv:"name"`
	StartDate     string          `json:"start_date" csv:"start_date"`
	EndDate       string          `json:"end_date" csv:"end_date"`
	PrizePool     float64         `json:"prize_pool" csv:"prize_pool"`
	Currency      string          `json:"currency" csv:"currency"`
	Status        string          `json:"status" csv:"status"`
	IsOnline      *bool           `json:"is_online" csv:"is_online"`
	BracketConfig json.RawMessage `json:"bracket_config" swaggertype:"object" csv:"bracket_config"`
}

type TournamentRegistrationImportInput struct {
	TournamentID   int64           `json:"tournament_id" csv:"tournament_id"`
	TeamID         int64           `json:"team_id" csv:"team_id"`
	SeedNumber     *int            `json:"seed_number" csv:"seed_number"`
	Status         string          `json:"status" csv:"status"`
	ManagerContact *string         `json:"manager_contact" csv:"manager_contact"`
	RosterSnapshot json.RawMessage `json:"roster_snapshot" swaggertype:"object" csv:"roster_snapshot"`
	IsInvited      *bool           `json:"is_invited" csv:"is_invited"`
}

type MatchImportInput struct {
	TournamentID int64           `json:"tournament_id" csv:"tournament_id"`
	Team1ID      *int64          `json:"team1_id" csv:"team1_id"`
	Team2ID      *int64          `json:"team2_id" csv:"team2_id"`
	StartTime    string          `json:"start_time" csv:"start_time"`
	Format       string          `json:"format" csv:"format"`
	Stage        *string         `json:"stage" csv:"stage"`
	WinnerTeamID *int64          `json:"winner_team_id" csv:"winner_team_id"`
	IsForfeit    *bool           `json:"is_forfeit" csv:"is_forfeit"`
	MatchNotes   json.RawMessage `json:"match_notes" swaggertype:"object" csv:"match_notes"`
}

type MatchGameImportInput struct {
	MatchID           int64           `json:"match_id" csv:"match_id"`
	MapName           string          `json:"map_name" csv:"map_name"`
	GameNumber        int             `json:"game_number" csv:"game_number"`
	DurationSeconds   *int            `json:"duration_seconds" csv:"duration_seconds"`
	WinnerTeamID      *int64          `json:"winner_team_id" csv:"winner_team_id"`
	ScoreTeam1        *int            `json:"score_team1" csv:"score_team1"`
	ScoreTeam2        *int            `json:"score_team2" csv:"score_team2"`
	StartedAt         *string         `json:"started_at" csv:"started_at"`
	HadTechnicalPause *bool           `json:"had_technical_pause" csv:"had_technical_pause"`
	PickBanPhase      json.RawMessage `json:"pick_ban_phase" swaggertype:"object" csv:"pick_ban_phase"`
}

type GamePlayerStatImportInput struct {
	GameID      int64   `json:"game_id" csv:"game_id"`
	PlayerID    int64   `json:"player_id" csv:"player_id"`
	TeamID      *int64  `json:"team_id" csv:"team_id"`
	Kills       int     `json:"kills" csv:"kills"`
	Deaths      int     `json:"deaths" csv:"deaths"`
	Assists     int     `json:"assists" csv:"assists"`
	HeroName    *string `json:"hero_name" csv:"hero_name"`
	DamageDealt int     `json:"damage_dealt" csv:"damage_dealt"`
	GoldEarned  int     `json:"gold_earned" csv:"gold_earned"`
	WasMVP      *bool   `json:"was_mvp" csv:"was_mvp"`
}

type SquadMemberImportInput struct {
	TeamID          int64    `json:"team_id" csv:"team_id"`
	PlayerID        int64    `json:"player_id" csv:"player_id"`
	Role            string   `json:"role" csv:"role"`
	IsStandin       *bool    `json:"is_standin" csv:"is_standin"`
	JoinDate        *string  `json:"join_date" csv:"join_date"`
	ContractEndDate *string  `json:"contract_end_date" csv:"contract_end_date"`
	LeaveDate       *string  `json:"leave_date" csv:"leave_date"`
	SalaryMonthly   *float64 `json:"salary_monthly" csv:"salary_monthly"`
}

type TeamProfileImportInput struct {
	TeamID       int64   `json:"team_id" csv:"team_id"`
	CoachName    *string `json:"coach_name" csv:"coach_name"`
	SponsorInfo  *string `json:"sponsor_info" csv:"sponsor_info"`
	Headquarters *string `json:"headquarters" csv:"headquarters"`
	Website      *string `json:"website" csv:"website"`
	ContactEmail *string `json:"contact_email" csv:"contact_email"`
}

type ImportSummary struct {
	Inserted int      `json:"inserted"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}
type ImportService struct {
	db              *sqlx.DB
	disciplineSvc   *DisciplineService
	teamSvc         *TeamService
	playerSvc       *PlayerService
	tournamentSvc   *TournamentService
	registrationSvc *TournamentRegistrationService
	matchSvc        *MatchService
	matchGameSvc    *MatchGameService
	statSvc         *GamePlayerStatService
	squadSvc        *SquadMemberService
	teamProfileSvc  *TeamProfileService
}

func NewImportService(db *sqlx.DB, disciplineSvc *DisciplineService, teamSvc *TeamService, playerSvc *PlayerService, tournamentSvc *TournamentService, registrationSvc *TournamentRegistrationService, matchSvc *MatchService, matchGameSvc *MatchGameService, statSvc *GamePlayerStatService, squadSvc *SquadMemberService, teamProfileSvc *TeamProfileService) *ImportService {
	return &ImportService{
		db:              db,
		disciplineSvc:   disciplineSvc,
		teamSvc:         teamSvc,
		playerSvc:       playerSvc,
		tournamentSvc:   tournamentSvc,
		registrationSvc: registrationSvc,
		matchSvc:        matchSvc,
		matchGameSvc:    matchGameSvc,
		statSvc:         statSvc,
		squadSvc:        squadSvc,
		teamProfileSvc:  teamProfileSvc,
	}
}

func (s *ImportService) ImportPlayers(ctx context.Context, source string, payload []PlayerImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		player, err := s.toPlayer(row)
		if err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		if err := s.playerSvc.Create(ctx, player); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportDisciplines(ctx context.Context, source string, payload []DisciplineImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		d := &models.Discipline{
			Code:        row.Code,
			Name:        row.Name,
			Description: row.Description,
			IconURL:     row.IconURL,
			TeamSize:    row.TeamSize,
			Metadata:    row.Metadata,
		}
		if row.IsActive != nil {
			d.IsActive = *row.IsActive
		} else {
			d.IsActive = true
		}
		if len(d.Metadata) == 0 {
			d.Metadata = json.RawMessage(`{}`)
		}
		if err := s.disciplineSvc.Create(ctx, d); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportTeams(ctx context.Context, source string, payload []TeamImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		team := &models.Team{
			Name:         row.Name,
			Tag:          row.Tag,
			CountryCode:  row.CountryCode,
			DisciplineID: row.DisciplineID,
			LogoURL:      row.LogoURL,
			WorldRanking: 0,
			IsVerified:   false,
		}
		if row.WorldRanking != nil {
			team.WorldRanking = *row.WorldRanking
		}
		if row.IsVerified != nil {
			team.IsVerified = *row.IsVerified
		}
		if err := s.teamSvc.Create(ctx, team); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportTournaments(ctx context.Context, source string, payload []TournamentImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		start, err := parseDate(row.StartDate)
		if err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		end, err := parseDate(row.EndDate)
		if err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		isOnline := false
		if row.IsOnline != nil {
			isOnline = *row.IsOnline
		}
		t := &models.Tournament{
			DisciplineID:  row.DisciplineID,
			Name:          row.Name,
			StartDate:     *start,
			EndDate:       *end,
			PrizePool:     row.PrizePool,
			Currency:      row.Currency,
			Status:        row.Status,
			IsOnline:      isOnline,
			BracketConfig: row.BracketConfig,
		}
		if err := s.tournamentSvc.Create(ctx, t); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportTournamentRegistrations(ctx context.Context, source string, payload []TournamentRegistrationImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		status := row.Status
		isInvited := false
		if row.IsInvited != nil {
			isInvited = *row.IsInvited
		}
		reg := &models.TournamentRegistration{
			TournamentID:   row.TournamentID,
			TeamID:         row.TeamID,
			SeedNumber:     row.SeedNumber,
			Status:         status,
			ManagerContact: row.ManagerContact,
			RosterSnapshot: row.RosterSnapshot,
			IsInvited:      isInvited,
		}
		if err := s.registrationSvc.Create(ctx, reg); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportMatches(ctx context.Context, source string, payload []MatchImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		start, err := parseDateTime(row.StartTime)
		if err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		isForfeit := false
		if row.IsForfeit != nil {
			isForfeit = *row.IsForfeit
		}
		var notes *json.RawMessage
		if row.MatchNotes != nil {
			n := row.MatchNotes
			notes = &n
		}
		m := &models.Match{
			TournamentID: row.TournamentID,
			Team1ID:      row.Team1ID,
			Team2ID:      row.Team2ID,
			StartTime:    *start,
			Format:       row.Format,
			Stage:        row.Stage,
			WinnerTeamID: row.WinnerTeamID,
			IsForfeit:    isForfeit,
			MatchNotes:   notes,
		}
		if err := s.matchSvc.Create(ctx, m); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportMatchGames(ctx context.Context, source string, payload []MatchGameImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		var startedAt *time.Time
		if row.StartedAt != nil && *row.StartedAt != "" {
			parsed, err := parseDateTime(*row.StartedAt)
			if err != nil {
				s.recordError(ctx, source, row, err, &summary)
				continue
			}
			startedAt = parsed
		}
		hasTech := false
		if row.HadTechnicalPause != nil {
			hasTech = *row.HadTechnicalPause
		}
		g := &models.MatchGame{
			MatchID:           row.MatchID,
			MapName:           row.MapName,
			GameNumber:        row.GameNumber,
			DurationSeconds:   row.DurationSeconds,
			WinnerTeamID:      row.WinnerTeamID,
			ScoreTeam1:        row.ScoreTeam1,
			ScoreTeam2:        row.ScoreTeam2,
			StartedAt:         startedAt,
			HadTechnicalPause: hasTech,
			PickBanPhase:      row.PickBanPhase,
		}
		if err := s.matchGameSvc.Create(ctx, g); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportGamePlayerStats(ctx context.Context, source string, payload []GamePlayerStatImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		stat := &models.GamePlayerStat{
			GameID:      row.GameID,
			PlayerID:    row.PlayerID,
			TeamID:      row.TeamID,
			Kills:       row.Kills,
			Deaths:      row.Deaths,
			Assists:     row.Assists,
			HeroName:    row.HeroName,
			DamageDealt: row.DamageDealt,
			GoldEarned:  row.GoldEarned,
			WasMVP:      false,
		}
		if row.WasMVP != nil {
			stat.WasMVP = *row.WasMVP
		}
		if err := s.statSvc.Create(ctx, stat); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportSquadMembers(ctx context.Context, source string, payload []SquadMemberImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		joinDate := time.Time{}
		if row.JoinDate != nil && *row.JoinDate != "" {
			parsed, err := parseDate(*row.JoinDate)
			if err != nil {
				s.recordError(ctx, source, row, err, &summary)
				continue
			}
			joinDate = *parsed
		}
		var contract *time.Time
		if row.ContractEndDate != nil && *row.ContractEndDate != "" {
			parsed, err := parseDate(*row.ContractEndDate)
			if err != nil {
				s.recordError(ctx, source, row, err, &summary)
				continue
			}
			contract = parsed
		}
		var leave *time.Time
		if row.LeaveDate != nil && *row.LeaveDate != "" {
			parsed, err := parseDate(*row.LeaveDate)
			if err != nil {
				s.recordError(ctx, source, row, err, &summary)
				continue
			}
			leave = parsed
		}
		isStandin := false
		if row.IsStandin != nil {
			isStandin = *row.IsStandin
		}
		m := &models.SquadMember{
			TeamID:          row.TeamID,
			PlayerID:        row.PlayerID,
			Role:            row.Role,
			IsStandin:       isStandin,
			JoinDate:        joinDate,
			ContractEndDate: contract,
			LeaveDate:       leave,
			SalaryMonthly:   row.SalaryMonthly,
		}
		if err := s.squadSvc.Create(ctx, m); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) ImportTeamProfiles(ctx context.Context, source string, payload []TeamProfileImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		profile := &models.TeamProfile{
			TeamID:       row.TeamID,
			CoachName:    row.CoachName,
			SponsorInfo:  row.SponsorInfo,
			Headquarters: row.Headquarters,
			Website:      row.Website,
			ContactEmail: row.ContactEmail,
		}
		if err := s.teamProfileSvc.Create(ctx, profile); err != nil {
			s.recordError(ctx, source, row, err, &summary)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}
func (s *ImportService) recordError(ctx context.Context, source string, row any, err error, summary *ImportSummary) {
	summary.Failed++
	summary.Errors = append(summary.Errors, err.Error())
	s.logError(ctx, source, row, err)
}

func (s *ImportService) logError(ctx context.Context, source string, row any, logErr error) {
	rowData, _ := json.Marshal(row)
	_, _ = s.db.ExecContext(ctx,
		`INSERT INTO batch_import_errors (source, row_data, error_message) VALUES ($1, $2::jsonb, $3)`,
		source,
		string(rowData),
		logErr.Error(),
	)
}

func parseDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, fmt.Errorf("date is required")
	}
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func parseDateTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, fmt.Errorf("datetime is required")
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (s *ImportService) toPlayer(row PlayerImportInput) (*models.Player, error) {
	var birth *time.Time
	if row.BirthDate != "" {
		parsed, err := parseDate(row.BirthDate)
		if err != nil {
			return nil, fmt.Errorf("invalid birth_date for %s", row.Nickname)
		}
		birth = parsed
	}
	isRetired := false
	if row.IsRetired != nil {
		isRetired = *row.IsRetired
	}
	return &models.Player{
		Nickname:    row.Nickname,
		RealName:    stringPtr(row.RealName),
		CountryCode: stringPtr(row.CountryCode),
		BirthDate:   birth,
		SteamID:     stringPtr(row.SteamID),
		AvatarURL:   stringPtr(row.AvatarURL),
		MMRRating:   row.MMR,
		IsRetired:   isRetired,
	}, nil
}

func stringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}
