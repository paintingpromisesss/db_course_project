from fastapi import FastAPI
from app.core.config import settings
from app.api.v1.endpoints import tournaments

app = FastAPI(
    title=settings.PROJECT_NAME,
    version=settings.VERSION,
    description="API for CyberTournament platform",
)

@app.get("/")
async def root():
    return {"message": "Welcome to the CyberTournament API!"}

app.include_router(tournaments.router, prefix="/api/v1", tags=["tournaments"])