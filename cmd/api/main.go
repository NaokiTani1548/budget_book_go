package main

import (
    "os"
	usecasecategory "budget-book-go/internal/application/usecase/category"
	recurringexpense "budget-book-go/internal/application/usecase/recurring_expense"
	"log"

	usecaseexpense "budget-book-go/internal/application/usecase/expense"
	usecaseincome "budget-book-go/internal/application/usecase/income"
	usecasesummary "budget-book-go/internal/application/usecase/summary"
	usecaseauth "budget-book-go/internal/application/usecase/auth"
	usecaseocr "budget-book-go/internal/application/usecase/ocr"
	"budget-book-go/internal/infrastructure/config"
	"budget-book-go/internal/infrastructure/persistence/postgres"
	"budget-book-go/internal/presentation/handler"
	"budget-book-go/internal/presentation/middleware"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"
	"github.com/gin-contrib/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".envファイルが見つかりません。環境変数を使用します")
	}
    oauthConfig := config.NewGoogleOAuthConfig()
    jwtSecret := config.GetJWTSecret()
	// DB接続
	dbConfig := config.NewDBConfig()
	db, err := config.NewDBPool(dbConfig)
	if err != nil {
		log.Fatalf("DB接続に失敗しました: %v", err)
	}
	defer db.Close()

	// Repository
	expenseRepo := postgres.NewExpenseRepository(db)
	categoryRepo    := postgres.NewCategoryRepository(db)
	incomeRepo := postgres.NewIncomeRepository(db)
	summaryRepo := postgres.NewSummaryRepository(db)
	recurringRepo := postgres.NewRecurringExpenseRepository(db)
	userRepo := postgres.NewUserRepository(db)

	// UseCase
	createExpenseUC := usecaseexpense.NewCreateExpenseUseCase(expenseRepo)
	getExpenseUC    := usecaseexpense.NewGetExpenseUseCase(expenseRepo)
	updateExpenseUC := usecaseexpense.NewUpdateExpenseUseCase(expenseRepo)
	deleteExpenseUC := usecaseexpense.NewDeleteExpenseUseCase(expenseRepo)
	getCategoryUC   := usecasecategory.NewGetCategoryUseCase(categoryRepo)
	createCategoryUC := usecasecategory.NewCreateCategoryUseCase(categoryRepo)
	updateCategoryUC := usecasecategory.NewUpdateCategoryUseCase(categoryRepo)
	deleteCategoryUC := usecasecategory.NewDeleteCategoryUseCase(categoryRepo)
	createIncomeUC := usecaseincome.NewCreateIncomeUseCase(incomeRepo)
	getIncomeUC    := usecaseincome.NewGetIncomeUseCase(incomeRepo)
	updateIncomeUC := usecaseincome.NewUpdateIncomeUseCase(incomeRepo)
	deleteIncomeUC := usecaseincome.NewDeleteIncomeUseCase(incomeRepo)
	getForecastUC := usecasesummary.NewGetForecastUseCase(summaryRepo)
	getRecurringUC    := recurringexpense.NewGetRecurringExpenseUseCase(recurringRepo)
	createRecurringUC := recurringexpense.NewCreateRecurringExpenseUseCase(recurringRepo)
	updateRecurringUC := recurringexpense.NewUpdateRecurringExpenseUseCase(recurringRepo)
	deleteRecurringUC := recurringexpense.NewDeleteRecurringExpenseUseCase(recurringRepo)
	applyRecurringUC  := recurringexpense.NewApplyRecurringExpenseUseCase(recurringRepo, expenseRepo)
	googleAuthUC := usecaseauth.NewGoogleAuthUseCase(userRepo, oauthConfig, jwtSecret)
	geminiKey := config.GetGeminiAPIKey()
    analyzeReceiptUC := usecaseocr.NewAnalyzeReceiptUseCase(geminiKey)



	// Handler
	expenseHandler := handler.NewExpenseHandler(
		createExpenseUC,
		getExpenseUC,
		updateExpenseUC,
		deleteExpenseUC,
		applyRecurringUC,
		)
	categoryHandler := handler.NewCategoryHandler(
		getCategoryUC,
		createCategoryUC,
		updateCategoryUC,
		deleteCategoryUC,
		)
	incomeHandler := handler.NewIncomeHandler(
		createIncomeUC,
		getIncomeUC,
		updateIncomeUC,
		deleteIncomeUC,
	)
	summaryHandler := handler.NewSummaryHandler(getForecastUC)
	recurringHandler := handler.NewRecurringExpenseHandler(
		getRecurringUC,
		createRecurringUC,
		updateRecurringUC,
		deleteRecurringUC,
		applyRecurringUC,
	)
    authHandler := handler.NewAuthHandler(googleAuthUC)
    ocrHandler := handler.NewOCRHandler(analyzeReceiptUC)

	// Router
	r := gin.Default()

	r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000", os.Getenv("FRONTEND_URL")},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Authorization", "Content-Type"},
        AllowCredentials: true,
    }))

api := r.Group("/api")
{
	// 認証不要
	auth := api.Group("/auth")
	{
		auth.GET("/google/url", authHandler.GetGoogleAuthURL)
		auth.GET("/google/callback", authHandler.GoogleCallback)
	}

	// 認証必要（JWTミドルウェア適用）
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		expenses := protected.Group("/expenses")
		{
			expenses.GET("", expenseHandler.GetAllByUserID)
			expenses.GET("/planned", expenseHandler.GetPlanned)
			expenses.GET("/date", expenseHandler.GetByDateRange)
			expenses.GET("/search", expenseHandler.Search)
			expenses.GET("/:id", expenseHandler.GetByID)
			expenses.POST("", expenseHandler.Create)
			expenses.PUT("/:id", expenseHandler.Update)
			expenses.DELETE("/:id", expenseHandler.Delete)
		}

		incomes := protected.Group("/incomes")
		{
			incomes.GET("", incomeHandler.GetAllByUserID)
			incomes.GET("/planned", incomeHandler.GetPlanned)
			incomes.GET("/date", incomeHandler.GetByDateRange)
			incomes.GET("/:id", incomeHandler.GetByID)
			incomes.POST("", incomeHandler.Create)
			incomes.PUT("/:id", incomeHandler.Update)
			incomes.DELETE("/:id", incomeHandler.Delete)
		}

		categories := protected.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllByUserID)
			categories.POST("", categoryHandler.Create)
			categories.PUT("/:id", categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}

		recurring := protected.Group("/recurring-expenses")
		{
			recurring.GET("", recurringHandler.GetAll)
			recurring.GET("/:id", recurringHandler.GetByID)
			recurring.POST("", recurringHandler.Create)
			recurring.PUT("/:id", recurringHandler.Update)
			recurring.DELETE("/:id", recurringHandler.Delete)
		}

		summary := protected.Group("/summary")
		{
			summary.GET("/forecast", summaryHandler.GetForecast)
		}
        ocrGroup := protected.Group("/ocr")
        {
            ocrGroup.POST("/analyze", ocrHandler.Analyze)
        }
	}
}

	// サーバー起動
   port := config.GetPort()
   if err := r.Run(":" + port); err != nil {
	   log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}
