package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	_ "school-teacher-management/docs"
	"school-teacher-management/internal/config"
	"school-teacher-management/internal/handler"
	"school-teacher-management/internal/metrics"
	"school-teacher-management/internal/middleware"
	"school-teacher-management/internal/model"
	"school-teacher-management/internal/repository"
	"school-teacher-management/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// @title           School Teacher Management API
// @version         1.0
// @description     REST APIs for managing teachers and attendance records.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @host      localhost:8082
// @BasePath  /api/v1
// @schemes   http
func main() {

	// -------------------- DATABASE --------------------
	config.ConnectDatabase()

	if config.DB == nil {
		log.Fatal("Database connection failed")
	}

	config.DB.AutoMigrate(
		&model.Teacher{},
		&model.Attendance{},
	)

	// -------------------- REPOSITORIES --------------------
	teacherRepo := repository.NewTeacherRepository(config.DB)
	attendanceRepo := repository.NewAttendanceRepository(config.DB)

	// -------------------- SERVICES --------------------
	teacherService := service.NewTeacherService(teacherRepo)
	attendanceService := service.NewAttendanceService(attendanceRepo)

	// -------------------- HANDLERS --------------------
	teacherHandler := handler.NewTeacherHandler(teacherService)
	attendanceHandler := handler.NewAttendanceHandler(attendanceService)

	// -------------------- GIN SETUP --------------------
	r := gin.New()

	// Recovery + Metrics middleware
	r.Use(gin.Recovery())
	r.Use(middleware.MetricsMiddleware())

	// Trust proxies (fix warning)
	r.SetTrustedProxies(nil)

	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	// -------------------- METRICS --------------------
	metrics.Register()
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// -------------------- API ROUTES --------------------
	api := r.Group("/api/v1")
	{
		// Teachers
		api.POST("/teachers", teacherHandler.CreateTeacher)
		api.GET("/teachers", teacherHandler.SearchTeachers)
		api.GET("/teachers/:id", teacherHandler.GetTeacherByID)
		api.PUT("/teachers/:id", teacherHandler.UpdateTeacher)
		api.POST("/teachers/bulk", teacherHandler.CreateTeachers)

		// Attendance
		api.POST("/attendance", attendanceHandler.CreateAttendance)
		api.GET("/attendance", attendanceHandler.GetAttendances)
		api.GET("/attendance/:id", attendanceHandler.GetAttendanceByID)
		api.PUT("/attendance/:id", attendanceHandler.UpdateAttendance)
		api.DELETE("/attendance/:id", attendanceHandler.DeleteAttendance)
		api.GET("/attendanceByDate", attendanceHandler.GetAttendanceByDate)
		api.GET("/attendanceByFilterDate", attendanceHandler.GetAttendanceByFilterDate)
	}

	// -------------------- SWAGGER --------------------
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// -------------------- SERVER --------------------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Println("🚀 Server running on port", port)
	r.Run(":" + port)
}
