from typing import Optional, List, Dict, Any
from datetime import date, datetime
from decimal import Decimal
from pydantic import BaseModel, ConfigDict, Field, AliasPath

class TeamBase(BaseModel):
    id: int
    name: str
    tag: str
    logo_url: Optional[str] = None
    model_config = ConfigDict(from_attributes=True)

class DisciplineBase(BaseModel):
    id: int
    name: str
    code: str
    model_config = ConfigDict(from_attributes=True)

class TournamentBase(BaseModel):
    name: str
    start_date: date
    end_date: date
    prize_pool: Optional[Decimal]
    currency: str
    status: str
    is_online: bool

class TournamentListItem(TournamentBase):
    id: int
    discipline: DisciplineBase
    model_config = ConfigDict(from_attributes=True)

class TournamentListResponse(BaseModel):
    total: int
    items: List[TournamentListItem]

class TournamentRegistrationSchema(BaseModel):
    team: TeamBase
    status: str
    seed_number: Optional[int]
    model_config = ConfigDict(from_attributes=True)

class TournamentDetail(TournamentBase):
    id: int
    discipline: DisciplineBase
    registrations: List[TournamentRegistrationSchema] = []
    model_config = ConfigDict(from_attributes=True)

class MatchBase(BaseModel):
    id: int
    start_time: datetime
    format: str
    stage: Optional[str]
    team1: Optional[TeamBase] = None
    team2: Optional[TeamBase] = None
    winner_team_id: Optional[int]
    tournament_name: str = Field(validation_alias=AliasPath("tournament", "name"))

    model_config = ConfigDict(from_attributes=True)

class MatchListResponse(BaseModel):
    total: int
    items: List[MatchBase]