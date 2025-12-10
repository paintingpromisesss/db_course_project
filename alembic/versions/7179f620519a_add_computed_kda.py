"""Add computed KDA

Revision ID: 7179f620519a
Revises: 743f3022cd2f
Create Date: 2025-12-10 16:56:21.325592

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision: str = '7179f620519a'
down_revision: Union[str, Sequence[str], None] = '743f3022cd2f'
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None


def upgrade() -> None:
    # 1. Удаляем старую обычную колонку
    op.drop_column('game_player_stats', 'kda_ratio')
    
    # 2. Создаем новую как GENERATED ALWAYS AS (...) STORED
    op.add_column('game_player_stats', sa.Column(
        'kda_ratio', 
        sa.Numeric(5, 2), 
        sa.Computed('(CASE WHEN deaths = 0 THEN (kills + assists)::numeric ELSE (kills + assists)::numeric / deaths END)', persisted=True),
        nullable=True # Computed columns технически nullable, пока не вычислены
    ))

def downgrade() -> None:
    # Возвращаем всё назад (если нужно откатить)
    op.drop_column('game_player_stats', 'kda_ratio')
    op.add_column('game_player_stats', sa.Column(
        'kda_ratio', 
        sa.Numeric(5, 2), 
        server_default=sa.text('0'), # Или что там было
        nullable=False
    ))