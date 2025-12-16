package models

import "time"

type SquadMember struct {
	ID              int64      `db:"id" json:"id"`
	TeamID          int64      `db:"team_id" json:"team_id"`
	PlayerID        int64      `db:"player_id" json:"player_id"`
	Role            string     `db:"role" json:"role"`
	IsStandin       bool       `db:"is_standin" json:"is_standin"`
	JoinDate        time.Time  `db:"join_date" json:"join_date"`
	ContractEndDate *time.Time `db:"contract_end_date" json:"contract_end_date"`
	LeaveDate       *time.Time `db:"leave_date" json:"leave_date"`
	SalaryMonthly   *float64   `db:"salary_monthly" json:"salary_monthly"`
}

type SquadMemberFilter struct {
	TeamID     *int64
	PlayerID   *int64
	ActiveOnly bool
	Limit      int
	Offset     int
}
