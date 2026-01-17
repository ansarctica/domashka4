package main

import (
	"context"
	"os"

	"log"

	"github.com/ansarctica/domashka4/internal/handlers"
	"github.com/ansarctica/domashka4/internal/postgres"
	"github.com/ansarctica/domashka4/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const port = ":8080"

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Nid't find .env file", "error", err)
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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	auth := e.Group("/api/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	api := e.Group("/api", h.UserIdentity)
	api.GET("/users/me", h.GetMe)

	protected := e.Group("", h.UserIdentity)

	protected.GET("/student/:id", h.GetStudent)
	protected.GET("/all_class_schedule", h.GetAllSchedules)
	protected.GET("/schedule/group/:id", h.GetGroupSchedule)
	protected.POST("/attendance/subject", h.NewAttendance)
	protected.GET("/attendanceBySubjectId/:id", h.GetAttendanceBySubject)
	protected.GET("/attendanceByStudentId/:id", h.GetAttendanceByStudent)

	e.Logger.Fatal(e.Start(port))
}
