package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"db_course_project/internal/api"
)

func NewRouter(disciplineHandler *api.DisciplineHandler, teamHandler *api.TeamHandler, playerHandler *api.PlayerHandler, tournamentHandler *api.TournamentHandler, teamProfileHandler *api.TeamProfileHandler, squadMemberHandler *api.SquadMemberHandler, tournamentRegistrationHandler *api.TournamentRegistrationHandler, matchHandler *api.MatchHandler, matchGameHandler *api.MatchGameHandler, gamePlayerStatHandler *api.GamePlayerStatHandler, utilityHandler *api.UtilityHandler) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		api.RespondData(c, 200, gin.H{"status": "ok"}, nil)
	})

	apiGroup := r.Group("/api")
	disciplineHandler.Register(apiGroup)
	teamHandler.Register(apiGroup)
	playerHandler.Register(apiGroup)
	tournamentHandler.Register(apiGroup)
	teamProfileHandler.Register(apiGroup)
	squadMemberHandler.Register(apiGroup)
	tournamentRegistrationHandler.Register(apiGroup)
	matchHandler.Register(apiGroup)
	matchGameHandler.Register(apiGroup)
	gamePlayerStatHandler.Register(apiGroup)
	utilityHandler.Register(apiGroup)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
