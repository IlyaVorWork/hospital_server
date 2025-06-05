package main

import (
	"Server/internal/pkg/postgres"
	"fmt"
	"log"
	"os"

	hs "Server/internal/hospital_server"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	err := godotenv.Load("deploy_server/.env")
	if err != nil {
		log.Fatalf("cannot open .env file: %sv", err)
	}

	postgresDB, cancelPostgres, err := postgres.OpenPostgresDB(fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB")))
	if err != nil {
		log.Fatal("cannot open postgres connection ", err)
		return
	}
	defer cancelPostgres(postgresDB)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(CORSMiddleware())

	postgresRepo := hs.NewPostgresRepository(postgresDB)
	service := hs.NewService(postgresRepo)
	handler := hs.NewHandler(service)

	router.GET("/specs", handler.GetSpecs)
	router.GET("/doctors", handler.GetDoctors)
	router.GET("/doctor/:id/appointments", handler.GetSlots)

	router.PUT("/user/appointment", handler.MakeAppointment)
	router.DELETE("/user/appointment", handler.CancelAppointment)

	router.GET("/user/profile", handler.GetUserProfile)
	router.GET("/user/appointments", handler.GetUserAppointments)

	addr := ":" + os.Getenv("SERVICE_PORT")
	err = router.Run(addr)
	if err != nil {
		log.Fatal(err)
	}
}
