package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

// PlayerImportInput represents a single player payload for batch import.
type PlayerImportInput struct {
	Nickname    string  `json:"nickname"`
	RealName    string  `json:"real_name"`
	CountryCode string  `json:"country_code"`
	BirthDate   string  `json:"birth_date"`
	SteamID     string  `json:"steam_id"`
	AvatarURL   string  `json:"avatar_url"`
	MMR         float64 `json:"mmr_rating"`
	IsRetired   *bool   `json:"is_retired"`
}

// ImportSummary reports results of batch import.
type ImportSummary struct {
	Inserted int      `json:"inserted"`
	Failed   int      `json:"failed"`
	Errors   []string `json:"errors"`
}

// ImportService handles batch imports and error logging.
type ImportService struct {
	db *sqlx.DB
}

func NewImportService(db *sqlx.DB) *ImportService {
	return &ImportService{db: db}
}

func (s *ImportService) ImportPlayers(ctx context.Context, source string, payload []PlayerImportInput) (ImportSummary, error) {
	summary := ImportSummary{}
	for _, row := range payload {
		if row.Nickname == "" {
			summary.Failed++
			summary.Errors = append(summary.Errors, "missing nickname")
			s.logError(ctx, source, row, errors.New("missing nickname"))
			continue
		}

		var birth *time.Time
		if row.BirthDate != "" {
			parsed, err := time.Parse("2006-01-02", row.BirthDate)
			if err != nil {
				summary.Failed++
				summary.Errors = append(summary.Errors, "invalid birth_date for "+row.Nickname)
				s.logError(ctx, source, row, err)
				continue
			}
			birth = &parsed
		}

		isRetired := false
		if row.IsRetired != nil {
			isRetired = *row.IsRetired
		}

		query := `INSERT INTO players (nickname, real_name, country_code, birth_date, steam_id, avatar_url, mmr_rating, is_retired)
				 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
		if _, err := s.db.ExecContext(ctx, query,
			row.Nickname,
			row.RealName,
			row.CountryCode,
			birth,
			row.SteamID,
			row.AvatarURL,
			row.MMR,
			isRetired,
		); err != nil {
			summary.Failed++
			summary.Errors = append(summary.Errors, err.Error())
			s.logError(ctx, source, row, err)
			continue
		}
		summary.Inserted++
	}
	return summary, nil
}

func (s *ImportService) logError(ctx context.Context, source string, row PlayerImportInput, logErr error) {
	rowData, _ := json.Marshal(row)
	_, _ = s.db.ExecContext(ctx,
		`INSERT INTO batch_import_errors (source, row_data, error_message) VALUES ($1, $2::jsonb, $3)`,
		source,
		string(rowData),
		logErr.Error(),
	)
}
