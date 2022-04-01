package main

import (
	"quik/internal/middleware"
	_mysqlPlayerRepo "quik/player/repository/mysql"
	_mysqlWalletRepo "quik/wallet/repository/mysql"
	_redisWalletRepo "quik/wallet/repository/redis"

	_playerService "quik/player/service"
	_walletService "quik/wallet/service"

	_playerHandler "quik/player/handler/http"
	_walletHandler "quik/wallet/handler/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func inject(d *DataSources) *gin.Engine {
	/*
	 * repository layer
	 */
	mysqlPlayerRepo := _mysqlPlayerRepo.NewMySqlPlayerRepository(d.MySQLDB)
	mysqlWalletRepo := _mysqlWalletRepo.NewMySqlWalletRepository(d.MySQLDB)
	redisWalletRepo := _redisWalletRepo.NewRedisInMemoryDB(d.RedisInMemoryDB)

	/*
	 * service layer
	 */
	playerService := _playerService.NewPlayerService(mysqlPlayerRepo)
	walletService := _walletService.NewWalletService(mysqlWalletRepo, redisWalletRepo)

	router := gin.Default()

	router.Use(middleware.LoggerToFile())
	router.Use(cors.Default())
	/*
	 * handler layer
	 */
	_playerHandler.NewPlayerHandler(router, playerService, walletService)
	_walletHandler.NewWalletHandler(router, walletService)

	return router
}
