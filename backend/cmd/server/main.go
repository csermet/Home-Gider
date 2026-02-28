package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/caner/home-gider/internal/config"
	"github.com/caner/home-gider/internal/database"
	"github.com/caner/home-gider/internal/handlers"
	"github.com/caner/home-gider/internal/middleware"
	"github.com/caner/home-gider/internal/scheduler"
	"github.com/caner/home-gider/internal/services"
	"github.com/labstack/echo/v4"
	echomw "github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg := config.Load()

	// Veritabanı bağlantısı
	db := database.Connect(cfg)
	database.Seed(db)

	// Services
	authService := services.NewAuthService(db, cfg.JWTSecret)
	expenseService := services.NewExpenseService(db)
	recurringService := services.NewRecurringService(db)
	settlementService := services.NewSettlementService(db)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(db)
	expenseHandler := handlers.NewExpenseHandler(expenseService)
	recurringHandler := handlers.NewRecurringHandler(recurringService)
	summaryHandler := handlers.NewSummaryHandler(settlementService)
	settlementHandler := handlers.NewSettlementHandler(settlementService)

	// Scheduler
	cron := scheduler.Start(recurringService)
	defer cron.Stop()

	// Echo
	e := echo.New()
	e.Use(echomw.Logger())
	e.Use(echomw.Recover())
	// CORS sadece development'ta gerekli (production'da frontend aynı origin'den sunulur)
	if corsOrigin := os.Getenv("CORS_ORIGIN"); corsOrigin != "" {
		e.Use(echomw.CORSWithConfig(echomw.CORSConfig{
			AllowOrigins:     []string{corsOrigin},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: true,
		}))
	}

	// Static dosyalar (production'da frontend build) — SPA routing desteği
	if _, err := os.Stat("static"); err == nil {
		e.Static("/assets", "static/assets")
		e.File("/vite.svg", "static/vite.svg")
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				path := c.Request().URL.Path
				if !strings.HasPrefix(path, "/api") && !strings.HasPrefix(path, "/assets") && path != "/healthz" {
					// Dosya varsa onu servis et, yoksa index.html (SPA fallback)
					filePath := "static" + path
					if _, err := os.Stat(filePath); err == nil {
						return c.File(filePath)
					}
					return c.File("static/index.html")
				}
				return next(c)
			}
		})
	}

	// 404 handler for API routes
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, map[string]interface{}{"error": he.Message})
			return
		}
		c.JSON(http.StatusInternalServerError, map[string]string{"error": "Sunucu hatası"})
	}

	// Health check (K8s probe'ları için)
	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API Routes
	api := e.Group("/api")

	// Auth (public)
	api.POST("/auth/login", authHandler.Login)

	// Auth gerektiren rotalar
	auth := api.Group("", middleware.AuthMiddleware(cfg.JWTSecret))

	auth.POST("/auth/logout", authHandler.Logout)
	auth.GET("/auth/me", authHandler.Me)
	auth.POST("/auth/change-password", authHandler.ChangePassword)

	// Admin rotaları
	admin := auth.Group("/admin", middleware.AdminMiddleware())
	admin.GET("/users", authHandler.ListUsers)
	admin.POST("/users/:id/reset-password", authHandler.AdminResetPassword)

	// Kategoriler
	auth.GET("/categories", categoryHandler.List)
	auth.POST("/categories", categoryHandler.Create)

	// Giderler
	auth.GET("/expenses", expenseHandler.List)
	auth.POST("/expenses", expenseHandler.Create)
	auth.PUT("/expenses/:id", expenseHandler.Update)
	auth.DELETE("/expenses/:id", expenseHandler.Delete)
	auth.POST("/expenses/:id/approve", expenseHandler.Approve)
	auth.POST("/expenses/:id/reject", expenseHandler.Reject)

	// Sabit/Taksitli Giderler
	auth.GET("/recurring", recurringHandler.List)
	auth.POST("/recurring", recurringHandler.Create)
	auth.PUT("/recurring/:id", recurringHandler.Update)
	auth.DELETE("/recurring/:id", recurringHandler.Delete)
	auth.POST("/recurring/:id/approve", recurringHandler.Approve)
	auth.POST("/recurring/:id/reject", recurringHandler.Reject)

	// Özet & Hesaplaşma
	auth.GET("/summary", summaryHandler.GetSummary)
	auth.GET("/payments", settlementHandler.ListPayments)
	auth.POST("/payments", settlementHandler.AddPayment)
	auth.DELETE("/payments/:id", settlementHandler.DeletePayment)

	// Graceful shutdown
	go func() {
		if err := e.Start(":" + cfg.Port); err != nil {
			log.Printf("Sunucu kapatılıyor: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Sunucu kapatılıyor...")
}
