package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rcontrisha/api-workers-ewallet/internal/config"
	"github.com/rcontrisha/api-workers-ewallet/internal/router"
	"github.com/rcontrisha/api-workers-ewallet/pkg"
)

// @title           						The Workers E-wallet
// @version         						1.0
// @description     						Digital wallet for simplifying transfer and track your money
// @host      									localhost:8080
// @BasePath  									/
// @securityDefinitions.apikey 	BearerAuth
// @in 													header
// @name 												Authorization
// @description 								Bearer token used as identiity for accessing backend
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error oading env file")
	}

	pdb := config.NewPsqlDb(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	pool, err := pdb.Connect()
	if err != nil {
		log.Println("Cannot Connect to DB\nReason: ", err.Error())
		return
	}
	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Println("Database is not ready\nReason: ", err.Error())
		return
	}
	log.Println("Database is ready")

	rc := config.NewRedisClient(os.Getenv("REDIS_USER"), os.Getenv("REDIS_PASSWORD"), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	rdb := rc.Connect()
	defer rdb.Close()
	status, err := rdb.Ping(context.Background()).Result()

	if err != nil {
		fmt.Println("[failed to connect redis]:", err)
	}
	fmt.Println(status)

	fmt.Println("connected to redis")

	mailCfg := config.LoadMailConfig()
	mailClient, err := pkg.NewMailer(
		mailCfg.Host,
		mailCfg.Port,
		mailCfg.Username,
		mailCfg.Password,
		mailCfg.SenderEmail,
	)
	if err != nil {
		log.Fatalf("failed to initialize mailer: %v", err)
	}

	dokuClient := pkg.NewDokuClient()

	r := gin.Default()
	router.MainRouter(r, pool, rdb, mailClient, dokuClient)
	r.Run(fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
}
