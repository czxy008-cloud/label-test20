package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"task-scheduler/internal/config"
	"task-scheduler/internal/database"
	"task-scheduler/internal/handler"
	"task-scheduler/internal/middleware"
	"task-scheduler/internal/repository"
	"task-scheduler/internal/service"
)

var (
	version   = "dev"
	buildDate = "unknown"
)

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	if err := database.Init(cfg.Database); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	defer database.Close()

	gin.SetMode(cfg.Server.Mode)
	router := gin.Default()

	userRepo := repository.NewUserRepository(database.DB)
	taskRepo := repository.NewTaskRepository(database.DB)
	execRepo := repository.NewExecutionRepository(database.DB)
	nodeRepo := repository.NewNodeRepository(database.DB)

	authService := service.NewAuthService(userRepo, cfg.Auth.Token.Secret)
	taskService := service.NewTaskService(taskRepo)
	execService := service.NewExecutionService(execRepo, taskRepo)
	nodeService := service.NewNodeService(nodeRepo)

	authHandler := handler.NewAuthHandler(authService)
	taskHandler := handler.NewTaskHandler(taskService)
	execHandler := handler.NewExecutionHandler(execService)
	nodeHandler := handler.NewNodeHandler(nodeService)

	authMiddleware := middleware.NewAuthMiddleware(authService, cfg.Auth.Token.Header)

	router.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		api.POST("/login", authHandler.Login)
	}

	auth := api.Group("")
	auth.Use(authMiddleware.Auth())
	{
		auth.GET("/user/me", authHandler.GetCurrentUser)

		tasks := auth.Group("/tasks")
		{
			tasks.POST("", taskHandler.Create)
			tasks.GET("", taskHandler.List)
			tasks.GET("/:id", taskHandler.GetByID)
			tasks.PUT("/:id", taskHandler.Update)
			tasks.DELETE("/:id", taskHandler.Delete)
		}

		executions := auth.Group("/executions")
		{
			executions.POST("/trigger", execHandler.TriggerTask)
			executions.GET("", execHandler.List)
			executions.GET("/:id", execHandler.GetByID)
		}

		nodes := auth.Group("/nodes")
		{
			nodes.GET("", nodeHandler.List)
			nodes.GET("/:id", nodeHandler.GetByID)
		}
	}

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动成功，监听端口: %d", cfg.Server.Port)
	if err := router.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
