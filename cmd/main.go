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

	_ "github.com/ansarctica/domashka4/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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
	protected.POST("/assignments", h.NewAssignment)
	protected.POST("/grades", h.NewGrade)
	protected.GET("/students/:id/gpa", h.GetGPA)
	protected.GET("/students/:studentId/subjects/:subjectId/gpa", h.GetSubjectGPA)

	protected.GET("/groups/:id/ranking", h.GetGPARankingByGroup)
	protected.GET("/subjects/:id/ranking", h.GetGPARankingBySubject)
	protected.GET("/groups/:groupId/subjects/:subjectId/ranking", h.GetSubjectGPARankingByGroup)

	e.Logger.Fatal(e.Start(port))
}
