package main

import (
	usecasecategory "budget-book-go/internal/application/usecase/category"
	"log"

	usecaseexpense "budget-book-go/internal/application/usecase/expense"
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

	// UseCase
	createExpenseUC := usecaseexpense.NewCreateExpenseUseCase(expenseRepo)
	getExpenseUC    := usecaseexpense.NewGetExpenseUseCase(expenseRepo)
	updateExpenseUC := usecaseexpense.NewUpdateExpenseUseCase(expenseRepo)
	deleteExpenseUC := usecaseexpense.NewDeleteExpenseUseCase(expenseRepo)
	getCategoryUC   := usecasecategory.NewGetCategoryUseCase(categoryRepo)


	// Handler
	expenseHandler := handler.NewExpenseHandler(
		createExpenseUC,
		getExpenseUC,
		updateExpenseUC,
		deleteExpenseUC,
	)
	categoryHandler := handler.NewCategoryHandler(getCategoryUC)

	// Router
	r := gin.Default()

	api := r.Group("/api")
	{
		expenses := api.Group("/expenses")
		{
			expenses.GET("",      expenseHandler.GetAll)
			expenses.GET("/:id",  expenseHandler.GetByID)
			expenses.POST("",     expenseHandler.Create)
			expenses.PUT("/:id",  expenseHandler.Update)
			expenses.DELETE("/:id", expenseHandler.Delete)
		}
		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetAllByUserID)
		}
	}

	// サーバー起動
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("サーバーの起動に失敗しました: %v", err)
	}
}