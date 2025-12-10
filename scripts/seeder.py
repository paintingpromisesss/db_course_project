import asyncio
import random
import sys
import os
from datetime import date, datetime, timedelta
from decimal import Decimal

sys.path.append(os.getcwd())

from sqlalchemy import select, func
from sqlalchemy.ext.asyncio import AsyncSession
from faker import Faker
from tqdm import tqdm 

from app.db.session import AsyncSessionLocal
from app.models.models import (
    Discipline, Team, Player, SquadMember,
    Tournament, TournamentRegistration,
    Match, MatchGame, GamePlayerStats, AuditLog
)

fake = Faker()
Faker.seed(42)

NUM_PLAYERS = 600
NUM_TEAMS = 60
NUM_TOURNAMENTS = 100
TARGET_STATS_ROWS = 10000

DISCIPLINES_DATA = [
    ("Dota 2", "dota2", 5), ("Counter-Strike 2", "cs2", 5), 
    ("Valorant", "valorant", 5), ("League of Legends", "lol", 5),
    ("StarCraft II", "sc2", 1), ("Overwatch 2", "ow2", 5),
    ("Rainbow Six Siege", "r6", 5), ("Apex Legends", "apex", 3),
    ("Rocket League", "rl", 3), ("Street Fighter 6", "sf6", 1)
]

async def seed_disciplines(session: AsyncSession):
    print("🔹 Генерируем дисциплины...")
    disciplines = []
    for name, code, size in DISCIPLINES_DATA:
        disciplines.append(Discipline(
            name=name, code=code, team_size=size,
            description=fake.catch_phrase(),
            icon_url=f"https://icons.example.com/{code}.png",
            meta_info={"developer": fake.company()}
        ))
    session.add_all(disciplines)
    await session.commit()
    return await session.scalars(select(Discipline))

async def seed_players(session: AsyncSession):
    print(f"🔹 Генерируем {NUM_PLAYERS} игроков...")
    players = []
    for _ in tqdm(range(NUM_PLAYERS)):
        players.append(Player(
            nickname=fake.unique.user_name(),
            real_name=fake.name(),
            country_code=fake.country_code(),
            birth_date=fake.date_of_birth(minimum_age=16, maximum_age=35),
            steam_id=str(fake.unique.random_number(digits=17)),
            mmr_rating=Decimal(random.uniform(1000, 9000)),
            avatar_url=fake.image_url()
        ))
    session.add_all(players)
    await session.commit()
    return (await session.scalars(select(Player))).all()

async def seed_teams(session: AsyncSession, disciplines):
    print(f"🔹 Генерируем {NUM_TEAMS} команд...")
    teams = []
    disc_list = list(disciplines)
    
    for _ in tqdm(range(NUM_TEAMS)):
        disc = random.choice(disc_list)
        name = fake.company()
        tag = "".join([word[0] for word in name.split()[:3]]).upper()
        if len(tag) < 2: tag = name[:3].upper()
        
        teams.append(Team(
            name=name,
            tag=tag + str(random.randint(1,99)),
            country_code=fake.country_code(),
            discipline_id=disc.id,
            world_ranking=Decimal(random.uniform(0, 100)),
            is_verified=random.choice([True, False])
        ))
    session.add_all(teams)
    await session.commit()
    return (await session.scalars(select(Team))).all()

async def seed_squads(session: AsyncSession, teams, players):
    print("🔹 Формируем составы...")
    squads = []
    random.shuffle(players)
    player_idx = 0
    
    for team in tqdm(teams):
        team_size = 5 
        for _ in range(team_size):
            if player_idx >= len(players): break
            player = players[player_idx]
            
            squads.append(SquadMember(
                team_id=team.id,
                player_id=player.id,
                role=random.choice(["Captain", "Sniper", "Support", "Entry", "Lurker"]),
                join_date=fake.date_between(start_date="-2y", end_date="-1y"),
                contract_end_date=fake.date_between(start_date="+1y", end_date="+2y"),
                salary_monthly=round(Decimal(random.uniform(500, 15000)), 2)
            ))
            player_idx += 1
            
    session.add_all(squads)
    await session.commit()

async def seed_tournaments_and_matches(session: AsyncSession, teams, players):
    print(f"🔹 Генерируем турниры, матчи и статистику (Цель: {TARGET_STATS_ROWS} записей)...")
    
    tournaments = []
    all_matches = []
    all_games = []
    all_stats = []
    
    teams_by_disc = {}
    for team in teams:
        if team.discipline_id not in teams_by_disc: teams_by_disc[team.discipline_id] = []
        teams_by_disc[team.discipline_id].append(team)

    total_stats_count = 0
    
    with tqdm(total=TARGET_STATS_ROWS) as pbar:
        for _ in range(NUM_TOURNAMENTS):
            disc_id = random.choice(list(teams_by_disc.keys()))
            available_teams = teams_by_disc[disc_id]
            
            if len(available_teams) < 2: continue
            
            start = fake.date_between(start_date="-1y", end_date="today")
            end = start + timedelta(days=random.randint(3, 14))
            
            tourn = Tournament(
                discipline_id=disc_id,
                name=f"{fake.city()} {random.choice(['Major', 'Cup', 'Invitational', 'Championship'])} 2025",
                start_date=start,
                end_date=end,
                prize_pool=Decimal(random.randint(10000, 1000000)),
                status="Finished",
                is_online=random.choice([True, False]),
                bracket_config={"type": "single_elimination"}
            )
            session.add(tourn)
            await session.flush()
            tournaments.append(tourn)
            
            num_teams = min(len(available_teams), random.randint(4, 16))
            tourn_teams = random.sample(available_teams, num_teams)
            
            regs = []
            for i, team in enumerate(tourn_teams):
                regs.append(TournamentRegistration(
                    tournament_id=tourn.id,
                    team_id=team.id,
                    seed_number=i+1,
                    status="Approved",
                    roster_snapshot={"player_ids": []}
                ))
            session.add_all(regs)
            
            match_pairs = []
            for i in range(0, len(tourn_teams) - 1, 2):
                match_pairs.append((tourn_teams[i], tourn_teams[i+1]))
                
            for t1, t2 in match_pairs:
                match = Match(
                    tournament_id=tourn.id,
                    team1_id=t1.id,
                    team2_id=t2.id,
                    start_time=datetime.combine(start, datetime.min.time()) + timedelta(hours=random.randint(10, 22)),
                    format="bo3",
                    stage="Group Stage",
                    winner_team_id=random.choice([t1.id, t2.id])
                )
                session.add(match)
                await session.flush()
                all_matches.append(match)
                
                num_games = random.randint(2, 3)
                for game_num in range(1, num_games + 1):
                    game = MatchGame(
                        match_id=match.id,
                        map_name=random.choice(["Dust 2", "Inferno", "Mirage", "Ancient"]),
                        game_number=game_num,
                        duration_seconds=random.randint(1200, 3600),
                        winner_team_id=match.winner_team_id if game_num == num_games else random.choice([t1.id, t2.id]),
                        score_team1=13, score_team2=random.randint(0, 11)
                    )
                    session.add(game)
                    await session.flush()
                    all_games.append(game)
                    
                    current_players = random.sample(players, 10)
                    
                    for pl in current_players:
                        stats = GamePlayerStats(
                            game_id=game.id,
                            player_id=pl.id,
                            team_id=random.choice([t1.id, t2.id]),
                            kills=random.randint(0, 30),
                            deaths=random.randint(0, 20),
                            assists=random.randint(0, 20),
                            hero_name=fake.first_name(),
                            damage_dealt=random.randint(5000, 30000),
                            gold_earned=random.randint(10000, 30000),
                            was_mvp=random.choice([True, False] + [False]*9)
                        )
                        all_stats.append(stats)
                        total_stats_count += 1
                        pbar.update(1)
            
            session.add_all(all_stats)
            await session.commit()
            all_stats = []

    print(f"✅ Успешно создано {total_stats_count} записей статистики!")


async def main():
    async with AsyncSessionLocal() as session:
        disciplines = await seed_disciplines(session)
        players = await seed_players(session)
        teams = await seed_teams(session, disciplines)
        await seed_squads(session, teams, players)
        await seed_tournaments_and_matches(session, teams, players)
        
        print("\n🎉 База данных успешно заполнена тестовыми данными!")
        print(f"Всего игроков: {len(players)}")
        print(f"Всего команд: {len(teams)}")

if __name__ == "__main__":
    if sys.platform == "win32":
        asyncio.set_event_loop_policy(asyncio.WindowsSelectorEventLoopPolicy())
    asyncio.run(main())
