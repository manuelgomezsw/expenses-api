package container

import (
	"expenses-api/internal/application/usecase"
	"expenses-api/internal/infrastructure/database"
	"expenses-api/internal/infrastructure/handler"
	"expenses-api/internal/infrastructure/repository"
	
	"gorm.io/gorm"
)

// Container holds all dependencies for dependency injection
type Container struct {
	// Database
	DB *gorm.DB

	// Repositories
	SalaryRepo            *repository.SalaryRepository
	PocketRepo            *repository.PocketRepository
	FixedExpenseRepo      *repository.FixedExpenseRepository
	DailyExpenseRepo      *repository.DailyExpenseRepository
	DailyExpenseConfigRepo *repository.DailyExpenseConfigRepository

	// Use Cases
	SalaryUseCase       *usecase.SalaryUseCase
	PocketUseCase       *usecase.PocketUseCase
	FixedExpenseUseCase *usecase.FixedExpenseUseCase
	DailyExpenseUseCase *usecase.DailyExpenseUseCase
	SummaryUseCase      *usecase.SummaryUseCase

	// Handlers
	ConfigHandler       *handler.ConfigHandler
	SummaryHandler      *handler.SummaryHandler
	FixedExpenseHandler *handler.FixedExpenseHandler
	DailyExpenseHandler *handler.DailyExpenseHandler
}

// NewContainer creates and initializes all dependencies
func NewContainer() (*Container, error) {
	container := &Container{}

	// Initialize database connection
	db := database.GetDB()
	container.DB = db

	// Initialize repositories
	container.SalaryRepo = repository.NewSalaryRepository(db)
	container.PocketRepo = repository.NewPocketRepository(db)
	container.FixedExpenseRepo = repository.NewFixedExpenseRepository(db)
	container.DailyExpenseRepo = repository.NewDailyExpenseRepository(db)
	container.DailyExpenseConfigRepo = repository.NewDailyExpenseConfigRepository(db)

	// Initialize use cases
	container.SalaryUseCase = usecase.NewSalaryUseCase(container.SalaryRepo)
	container.PocketUseCase = usecase.NewPocketUseCase(container.PocketRepo)
	container.FixedExpenseUseCase = usecase.NewFixedExpenseUseCase(container.FixedExpenseRepo)
	container.DailyExpenseUseCase = usecase.NewDailyExpenseUseCase(container.DailyExpenseRepo)
	
	// Summary use case needs multiple repositories
	container.SummaryUseCase = usecase.NewSummaryUseCase(
		container.SalaryRepo,
		container.FixedExpenseRepo,
		container.DailyExpenseRepo,
		container.DailyExpenseConfigRepo,
	)

	// Initialize handlers
	container.ConfigHandler = handler.NewConfigHandler(
		container.SalaryUseCase,
		container.PocketUseCase,
	)
	container.SummaryHandler = handler.NewSummaryHandler(container.SummaryUseCase)
	container.FixedExpenseHandler = handler.NewFixedExpenseHandler(container.FixedExpenseUseCase)
	container.DailyExpenseHandler = handler.NewDailyExpenseHandler(container.DailyExpenseUseCase)

	return container, nil
}
