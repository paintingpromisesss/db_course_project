package main

// @title DB Course Project API
// @version 1.0
// @description REST API for disciplines, teams, tournaments, matches, and reports.
// @BasePath /api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"db_course_project/docs"
	"db_course_project/internal/api"
	"db_course_project/internal/config"
	"db_course_project/internal/db"
	"db_course_project/internal/repository"
	"db_course_project/internal/server"
	"db_course_project/internal/service"
)

func main() {
	cfg := config.New()

	docs.SwaggerInfo.Title = "DB Course Project API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/api"
	docs.SwaggerInfo.Description = "REST API for managing esports data with reports and imports."

	sqlxDB, err := db.Connect(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer sqlxDB.Close()

	// Wiring repositories and services.
	disciplineRepo := repository.NewDisciplineRepository(sqlxDB)
	teamRepo := repository.NewTeamRepository(sqlxDB)
	playerRepo := repository.NewPlayerRepository(sqlxDB)
	reportRepo := repository.NewReportRepository(sqlxDB)
	tournamentRepo := repository.NewTournamentRepository(sqlxDB)
	teamProfileRepo := repository.NewTeamProfileRepository(sqlxDB)
	squadMemberRepo := repository.NewSquadMemberRepository(sqlxDB)
	tournamentRegistrationRepo := repository.NewTournamentRegistrationRepository(sqlxDB)
	matchRepo := repository.NewMatchRepository(sqlxDB)
	matchGameRepo := repository.NewMatchGameRepository(sqlxDB)
	gamePlayerStatRepo := repository.NewGamePlayerStatRepository(sqlxDB)

	disciplineSvc := service.NewDisciplineService(disciplineRepo)
	teamSvc := service.NewTeamService(teamRepo)
	playerSvc := service.NewPlayerService(playerRepo)
	reportSvc := service.NewReportService(reportRepo)
	importSvc := service.NewImportService(sqlxDB)
	tournamentSvc := service.NewTournamentService(tournamentRepo)
	teamProfileSvc := service.NewTeamProfileService(teamProfileRepo)
	squadMemberSvc := service.NewSquadMemberService(squadMemberRepo)
	tournamentRegistrationSvc := service.NewTournamentRegistrationService(tournamentRegistrationRepo)
	matchSvc := service.NewMatchService(matchRepo)
	matchGameSvc := service.NewMatchGameService(matchGameRepo)
	gamePlayerStatSvc := service.NewGamePlayerStatService(gamePlayerStatRepo)

	disciplineHandler := api.NewDisciplineHandler(disciplineSvc)
	teamHandler := api.NewTeamHandler(teamSvc)
	playerHandler := api.NewPlayerHandler(playerSvc)
	tournamentHandler := api.NewTournamentHandler(tournamentSvc)
	teamProfileHandler := api.NewTeamProfileHandler(teamProfileSvc)
	squadMemberHandler := api.NewSquadMemberHandler(squadMemberSvc)
	tournamentRegistrationHandler := api.NewTournamentRegistrationHandler(tournamentRegistrationSvc)
	matchHandler := api.NewMatchHandler(matchSvc)
	matchGameHandler := api.NewMatchGameHandler(matchGameSvc)
	gamePlayerStatHandler := api.NewGamePlayerStatHandler(gamePlayerStatSvc)
	utilityHandler := api.NewUtilityHandler(reportSvc, importSvc)

	router := server.NewRouter(disciplineHandler, teamHandler, playerHandler, tournamentHandler, teamProfileHandler, squadMemberHandler, tournamentRegistrationHandler, matchHandler, matchGameHandler, gamePlayerStatHandler, utilityHandler)

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: router,
	}

	go func() {
		log.Printf("starting http server on %s", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
	log.Println("server exited")
}
