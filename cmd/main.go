package main

import (
	"context"
	"os"

	"github.com/ansarctica/domashka4/internal/handlers"
	"github.com/ansarctica/domashka4/internal/postgres"
	"github.com/ansarctica/domashka4/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DATABASE_URL")
	dbPool, _ := pgxpool.New(context.Background(), dbURL)
	defer dbPool.Close()

	repo := postgres.NewRepository(dbPool)
	srv := service.NewService(repo)
	h := handlers.NewHandler(srv)

	e := echo.New()

	e.GET("/student/:id", h.GetStudent)
	e.GET("/all_class_schedule", h.GetAllSchedules)
	e.GET("/schedule/group/:id", h.GetGroupSchedule)
	e.POST("/attendance/subject", h.NewAttendance)
	e.GET("/attendanceBySubjectId/:id", h.GetAttendanceBySubject)
	e.GET("/attendanceByStudentId/:id", h.GetAttendanceByStudent)

	e.Logger.Fatal(e.Start(":8080"))
}
