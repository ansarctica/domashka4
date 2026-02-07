package main

import (
	"context"
	"log"
	"os"

	"github.com/ansarctica/domashka4/internal/handlers"
	"github.com/ansarctica/domashka4/internal/postgres"
	"github.com/ansarctica/domashka4/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/ansarctica/domashka4/docs"
)

// @title Student Management API
// @version 1.0
// @description API for managing students, schedules, attendance, and grades.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Did't find .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	} else if port[0] != ':' {
		port = ":" + port
	}

	dbURL := os.Getenv("DATABASE_URL")
	dbPool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatal("Didn't connect to db", "error", err)
	}
	defer dbPool.Close()

	repo := postgres.NewRepository(dbPool)
	srv := service.NewService(repo)
	h := handlers.NewHandler(srv)

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "https://front-for-hw4.onrender.com"},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{
			echo.GET,
			echo.PUT,
			echo.POST,
			echo.DELETE,
			echo.PATCH,
			echo.OPTIONS,
		},
		MaxAge: 86400,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Swagger Endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := e.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	protected := e.Group("", h.UserIdentity)
	protected.GET("/users/me", h.GetMe)
	protected.GET("/students", h.GetStudents)
	protected.GET("/students/:id", h.GetStudent)
	protected.POST("/students", h.CreateStudent)
	protected.PATCH("/students/:id", h.UpdateStudent)
	protected.DELETE("/students/:id", h.DeleteStudent)
	protected.GET("/students/:id/gpa", h.GetStudentGPA)
	protected.GET("/groups", h.GetGroups)
	protected.GET("/schedules", h.GetSchedules)
	protected.POST("/schedules", h.CreateSchedule)
	protected.PATCH("/schedules/:id", h.UpdateSchedule)
	protected.DELETE("/schedules/:id", h.DeleteSchedule)
	protected.GET("/attendance", h.GetAttendance)
	protected.POST("/attendance", h.CreateAttendance)
	protected.PATCH("/attendance/:id", h.UpdateAttendance)
	protected.DELETE("/attendance/:id", h.DeleteAttendance)
	protected.GET("/assignments", h.GetAssignments)
	protected.POST("/assignments", h.CreateAssignment)
	protected.POST("/grades", h.CreateGrade)
	protected.GET("/rankings", h.GetRankings)
	protected.GET("/subjects", h.GetSubjects)

	e.Logger.Fatal(e.Start(port))
}
