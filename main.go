package main

import (
	"agnos_candidate_assignment/config"
	"agnos_candidate_assignment/database"
	"agnos_candidate_assignment/handlers"
	"agnos_candidate_assignment/middleware"
	"agnos_candidate_assignment/repositories"
	"agnos_candidate_assignment/services"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	conf := config.Load()

	if v, _ := os.LookupEnv("SILENCE_LOGS"); v == "true" || conf.GinMode == "release" {
		log.SetOutput(io.Discard)
		gin.DefaultWriter = io.Discard
		gin.DefaultErrorWriter = io.Discard
	}

	db, err := database.NewPostgresConnection(conf)

	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	hospitalRepo := repositories.NewHospitalRepository(db)
	staffRepo := repositories.NewStaffRepository(db)
	patientRepo := repositories.NewPatientRepository(db)

	authService := services.NewAuthService(staffRepo, hospitalRepo, conf)
	patientService := services.NewPatientService(patientRepo)

	hospitalHandler := handlers.NewHospitalHandler(hospitalRepo)
	staffHandler := handlers.NewStaffHandler(authService)
	patientHandler := handlers.NewPatientHandler(patientService)

	gin.SetMode(conf.GinMode)

	router := gin.New()
	router.Use(gin.RecoveryWithWriter(io.Discard))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api.POST("/hospital", hospitalHandler.Create)

	authMiddleWare := middleware.JWTAuth(conf, staffRepo)

	hospitalGroup := api.Group(":hospital")
	{
		hospitalGroup.POST("/staff/create", staffHandler.Register)
		hospitalGroup.POST("/staff/login", staffHandler.Login)

		hospitalGroup.GET("/patient/search/:id", func(c *gin.Context) {
			name := c.Param("hospital")
			hosp, err := hospitalRepo.FindByName(name)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "hospital not found"})
				return
			}
			c.Set("hospital_id", hosp.ID)
			patientHandler.GetByID(c)
		})
	}

	api.GET("/patient/search", authMiddleWare, func(c *gin.Context) {
		patientHandler.Search(c)
	})

	log.Printf("Starting server on port %s", conf.ServerPort)

	if err := router.Run(":" + conf.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
