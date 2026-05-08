package main

import (
	usecasecategory "budget-book-go/internal/application/usecase/category"
	"log"

	usecaseexpense "budget-book-go/internal/application/usecase/expense"
	usecaseincome "budget-book-go/internal/application/usecase/income"
	usecasesummary "budget-book-go/internal/application/usecase/summary"
	"budget-book-go/internal/infrastructure/config"
	"budget-book-go/internal/infrastructure/persistence/postgres"
	"budget-book-go/internal/presentation/handler"

	"github.com/gin-gonic/gin"
)

func main() {
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


	// Handler
	expenseHandler := handler.NewExpenseHandler(
		createExpenseUC,
		getExpenseUC,
		updateExpenseUC,
		deleteExpenseUC,
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

	// Router
	r := gin.Default()

	api := r.Group("/api")
	{
		expenses := api.Group("/expenses")
		{
			expenses.GET("",      expenseHandler.GetAll)
			expenses.GET("/date", expenseHandler.GetByDateRange)
			expenses.GET("/:id",  expenseHandler.GetByID)
			expenses.POST("",     expenseHandler.Create)
			expenses.PUT("/:id",  expenseHandler.Update)
			expenses.DELETE("/:id", expenseHandler.Delete)
		}
		incomes := api.Group("/incomes")
		{
			incomes.GET("",          incomeHandler.GetAll)
			incomes.GET("/planned",  incomeHandler.GetPlanned)
			incomes.GET("/date", incomeHandler.GetByDateRange)
			incomes.GET("/:id",      incomeHandler.GetByID)
			incomes.POST("",         incomeHandler.Create)
			incomes.PUT("/:id",      incomeHandler.Update)
			incomes.DELETE("/:id",   incomeHandler.Delete)
		}
		categories := api.Group("/categories")
		{
			categories.GET("",       categoryHandler.GetAllByUserID)
			categories.POST("",      categoryHandler.Create)
			categories.PUT("/:id",   categoryHandler.Update)
			categories.DELETE("/:id", categoryHandler.Delete)
		}
		summary := api.Group("/summary")
		{
			summary.GET("/forecast", summaryHandler.GetForecast)
		}
	}

	// サーバー起動
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}