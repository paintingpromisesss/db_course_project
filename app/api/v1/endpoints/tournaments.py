from typing import Any, Optional
from fastapi import APIRouter, Depends, HTTPException, Query
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy import select, func
from sqlalchemy.orm import selectinload

from app.db.session import get_db
from app.models.models import Tournament, Match, TournamentRegistration
from app.schemas import schemas

router = APIRouter()

@router.get("/tournaments", response_model=schemas.TournamentListResponse)
async def get_tournaments(
    db: AsyncSession = Depends(get_db),
    skip: int = 0,
    limit: int = 20,
    status: Optional[str] = None,
    discipline_id: Optional[int] = None
):
    query = select(Tournament).options(
        selectinload(Tournament.discipline)
    )

    if status:
        query = query.where(Tournament.status == status)
    if discipline_id:
        query = query.where(Tournament.discipline_id == discipline_id)

    count_query = select(func.count()).select_from(query.subquery())
    total = await db.scalar(count_query)
    
    query = query.order_by(Tournament.start_date.desc()).offset(skip).limit(limit)
    
    result = await db.execute(query)
    tournaments = result.scalars().all()
    
    return {"total": total, "items": tournaments}

@router.get("/tournaments/{tournament_id}", response_model=schemas.TournamentDetail)
async def get_tournament_detail(
    tournament_id: int,
    db: AsyncSession = Depends(get_db)
):
    query = select(Tournament).where(Tournament.id == tournament_id).options(
        selectinload(Tournament.discipline),
        selectinload(Tournament.registrations).selectinload(TournamentRegistration.team)
    )
    
    result = await db.execute(query)
    tournament = result.scalar_one_or_none()
    
    if not tournament:
        raise HTTPException(status_code=404, detail="Tournament not found")
        
    return tournament


@router.get("/matches", response_model=schemas.MatchListResponse)
async def get_matches(
    db: AsyncSession = Depends(get_db),
    skip: int = 0,
    limit: int = 50,
    tournament_id: Optional[int] = None
):
    query = select(Match).options(
        selectinload(Match.team1),
        selectinload(Match.team2),
        selectinload(Match.tournament)
    )
    
    if tournament_id:
        query = query.where(Match.tournament_id == tournament_id)
        
    query = query.order_by(Match.start_time.desc()).offset(skip).limit(limit)
    
    if tournament_id:
         count_q = select(func.count()).where(Match.tournament_id == tournament_id)
    else:
         count_q = select(func.count(Match.id))
         
    total = await db.scalar(count_q)
    
    result = await db.execute(query)
    matches = result.scalars().all()
    
    return {"total": total, "items": matches}
