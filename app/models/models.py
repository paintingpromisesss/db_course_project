import datetime
from typing import Optional, List
from decimal import Decimal

from sqlalchemy import Integer, String, Date, DateTime, ForeignKey, Numeric, Text, CheckConstraint, BigInteger, func, Boolean, FetchedValue, Computed
from sqlalchemy.dialects.postgresql import JSONB
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.db.base_class import Base


class Discipline(Base):
    __tablename__ = "disciplines"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, index=True)
    name: Mapped[str] = mapped_column(String(100), unique=True, nullable=False)
    code: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    description: Mapped[Optional[str]] = mapped_column(Text)
    icon_url: Mapped[Optional[str]] = mapped_column(String(255))
    team_size: Mapped[int] = mapped_column(Integer, default=5)
    is_active: Mapped[bool] = mapped_column(default=True)
    meta_info: Mapped[Optional[dict]] = mapped_column("metadata", JSONB, default={})

    teams: Mapped[List["Team"]] = relationship(back_populates="discipline")
    tournaments: Mapped[List["Tournament"]] = relationship(back_populates="discipline")


class Team(Base):
    __tablename__ = "teams"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, index=True)
    name: Mapped[str] = mapped_column(String(100), nullable=False)
    tag: Mapped[str] = mapped_column(String(10), nullable=False)
    country_code: Mapped[str] = mapped_column(String(2), nullable=False)
    discipline_id: Mapped[int] = mapped_column(ForeignKey("disciplines.id"))
    created_at: Mapped[datetime.datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())
    logo_url: Mapped[Optional[str]] = mapped_column(String(255))
    world_ranking: Mapped[Decimal] = mapped_column(Numeric(5, 2), default=0.00)
    is_verified: Mapped[bool] = mapped_column(Boolean, default=False)

    discipline: Mapped["Discipline"] = relationship(back_populates="teams")
    squad_members: Mapped[List["SquadMember"]] = relationship(back_populates="team")
    registrations: Mapped[List["TournamentRegistration"]] = relationship(back_populates="team")

class Player(Base):
    __tablename__ = "players"
    id: Mapped[int] = mapped_column(Integer, primary_key=True, index=True)
    nickname: Mapped[str] = mapped_column(String(50), unique=True, nullable=False)
    real_name: Mapped[Optional[str]] = mapped_column(String(100))
    country_code: Mapped[Optional[str]] = mapped_column(String(2))
    birth_date: Mapped[Optional[datetime.date]] = mapped_column(Date)
    steam_id: Mapped[Optional[str]] = mapped_column(String(32), unique=True)
    avatar_url: Mapped[Optional[str]] = mapped_column(String(255))
    mmr_rating: Mapped[Decimal] = mapped_column(Numeric(7, 1), default=0.0)
    is_retired: Mapped[bool] = mapped_column(Boolean, default=False)
    created_at: Mapped[datetime.datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())

    squad_history: Mapped[List["SquadMember"]] = relationship(back_populates="player")

class SquadMember(Base):
    __tablename__ = "squad_members"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True)
    team_id: Mapped[int] = mapped_column(ForeignKey("teams.id"))
    player_id: Mapped[int] = mapped_column(ForeignKey("players.id"))
    role: Mapped[str] = mapped_column(String(50), default="Player")
    is_standin: Mapped[bool] = mapped_column(Boolean, default=False)
    join_date: Mapped[datetime.date] = mapped_column(Date, default=datetime.date.today)
    contract_end_date: Mapped[Optional[datetime.date]] = mapped_column(Date)
    leave_date: Mapped[Optional[datetime.date]] = mapped_column(Date)
    salary_monthly: Mapped[Optional[Decimal]] = mapped_column(Numeric(10, 2))

    team: Mapped["Team"] = relationship(back_populates="squad_members")
    player: Mapped["Player"] = relationship(back_populates="squad_history")

    __table_args__ = (
        CheckConstraint('leave_date IS NULL OR leave_date >= join_date', name='chk_dates'),
    )

class Tournament(Base):
    __tablename__ = "tournaments"

    id: Mapped[int] = mapped_column(Integer, primary_key=True, index=True)
    discipline_id: Mapped[int] = mapped_column(ForeignKey("disciplines.id"))
    name: Mapped[str] = mapped_column(String(200), nullable=False)
    start_date: Mapped[datetime.date] = mapped_column(Date, nullable=False)
    end_date: Mapped[datetime.date] = mapped_column(Date, nullable=False)
    prize_pool: Mapped[Optional[Decimal]] = mapped_column(Numeric(15, 2), default=0)
    currency: Mapped[str] = mapped_column(String(3), default="USD")
    status: Mapped[str] = mapped_column(String(20), default="Announced")
    is_online: Mapped[bool] = mapped_column(Boolean, default=False)
    bracket_config: Mapped[Optional[dict]] = mapped_column(JSONB)

    discipline: Mapped["Discipline"] = relationship(back_populates="tournaments")
    registrations: Mapped[List["TournamentRegistration"]] = relationship(back_populates="tournament")
    matches: Mapped[List["Match"]] = relationship(back_populates="tournament")

    __table_args__ = (
        CheckConstraint('end_date >= start_date', name='chk_tournament_dates'),
    )

class TournamentRegistration(Base):
    __tablename__ = "tournament_registrations"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True)
    tournament_id: Mapped[int] = mapped_column(ForeignKey("tournaments.id"))
    team_id: Mapped[int] = mapped_column(ForeignKey("teams.id"))
    seed_number: Mapped[Optional[int]] = mapped_column(Integer)
    status: Mapped[str] = mapped_column(String(20), default="Pending")
    manager_contact: Mapped[Optional[str]] = mapped_column(String(100))
    roster_snapshot: Mapped[Optional[dict]] = mapped_column(JSONB)
    is_invited: Mapped[bool] = mapped_column(Boolean, default=False)
    registered_at: Mapped[datetime.datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())

    tournament: Mapped["Tournament"] = relationship(back_populates="registrations")
    team: Mapped["Team"] = relationship(back_populates="registrations")

class Match(Base):
    __tablename__ = "matches"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True)
    tournament_id: Mapped[int] = mapped_column(ForeignKey("tournaments.id"))
    team1_id: Mapped[Optional[int]] = mapped_column(ForeignKey("teams.id"))
    team2_id: Mapped[Optional[int]] = mapped_column(ForeignKey("teams.id"))
    start_time: Mapped[datetime.datetime] = mapped_column(DateTime(timezone=True))
    format: Mapped[str] = mapped_column(String(10), default="bo3")
    stage: Mapped[Optional[str]] = mapped_column(String(50))
    winner_team_id: Mapped[Optional[int]] = mapped_column(ForeignKey("teams.id"))
    is_forfeit: Mapped[bool] = mapped_column(Boolean, default=False)
    match_notes: Mapped[Optional[dict]] = mapped_column(JSONB)

    tournament: Mapped["Tournament"] = relationship(back_populates="matches")
    games: Mapped[List["MatchGame"]] = relationship(back_populates="match")
    team1: Mapped["Team"] = relationship(foreign_keys=[team1_id])
    team2: Mapped["Team"] = relationship(foreign_keys=[team2_id])

class MatchGame(Base):
    __tablename__ = "match_games"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True)
    match_id: Mapped[int] = mapped_column(ForeignKey("matches.id"))
    map_name: Mapped[str] = mapped_column(String(100))
    game_number: Mapped[int] = mapped_column(Integer)
    duration_seconds: Mapped[Optional[int]] = mapped_column(Integer)
    winner_team_id: Mapped[Optional[int]] = mapped_column(ForeignKey("teams.id"))
    score_team1: Mapped[int] = mapped_column(Integer, default=0)
    score_team2: Mapped[int] = mapped_column(Integer, default=0)
    started_at: Mapped[Optional[datetime.datetime]] = mapped_column(DateTime(timezone=True))
    had_technical_pause: Mapped[bool] = mapped_column(Boolean, default=False)
    pick_ban_phase: Mapped[Optional[dict]] = mapped_column(JSONB)

    match: Mapped["Match"] = relationship(back_populates="games")
    player_stats: Mapped[List["GamePlayerStats"]] = relationship(back_populates="game")

class GamePlayerStats(Base):
    __tablename__ = "game_player_stats"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True)
    game_id: Mapped[int] = mapped_column(ForeignKey("match_games.id"))
    player_id: Mapped[int] = mapped_column(ForeignKey("players.id"))
    team_id: Mapped[Optional[int]] = mapped_column(ForeignKey("teams.id"))
    
    kills: Mapped[int] = mapped_column(Integer, default=0)
    deaths: Mapped[int] = mapped_column(Integer, default=0)
    assists: Mapped[int] = mapped_column(Integer, default=0)
    hero_name: Mapped[Optional[str]] = mapped_column(String(100))
    damage_dealt: Mapped[int] = mapped_column(Integer, default=0)
    gold_earned: Mapped[int] = mapped_column(Integer, default=0)

    kda_ratio: Mapped[Decimal] = mapped_column(
        Numeric(5, 2),
        Computed(
            "(CASE WHEN deaths = 0 THEN kills + assists ELSE (kills + assists)::decimal / deaths END)"
        ),
    )

    was_mvp: Mapped[bool] = mapped_column(Boolean, default=False)

    game: Mapped["MatchGame"] = relationship(back_populates="player_stats")
    player: Mapped["Player"] = relationship()

class AuditLog(Base):
    __tablename__ = "audit_logs"

    id: Mapped[int] = mapped_column(BigInteger, primary_key=True)
    table_name: Mapped[str] = mapped_column(String(50))
    record_id: Mapped[int] = mapped_column(BigInteger)
    operation: Mapped[str] = mapped_column(String(10))
    old_value: Mapped[Optional[dict]] = mapped_column(JSONB)
    new_value: Mapped[Optional[dict]] = mapped_column(JSONB)
    changed_at: Mapped[datetime.datetime] = mapped_column(DateTime(timezone=True), server_default=func.now())
    changed_by: Mapped[Optional[str]] = mapped_column(String(100))
    is_sensitive: Mapped[bool] = mapped_column(Boolean, default=False)