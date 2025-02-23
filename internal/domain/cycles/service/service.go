package service

import (
	"expenses-api/internal/domain/budgets/service"
	conceptsRepository "expenses-api/internal/domain/concepts/repository"
	"expenses-api/internal/domain/cycles"
	"expenses-api/internal/domain/cycles/repository"
	expensesRepository "expenses-api/internal/domain/expenses/repository"
	"expenses-api/internal/util/customdate"
	"fmt"
	"time"
)

var (
	MonthNamesShort = [...]string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}
)

func GetAll() ([]cycles.Cycle, error) {
	return repository.GetAll()
}

func GetActive() ([]cycles.Cycle, error) {
	return repository.GetActive()
}

func GetByID(cycleID int) (cycles.Cycle, error) {
	return repository.GetByID(cycleID)
}

func Create(cycle *cycles.Cycle) error {
	if err := checkDates(cycle.DateInit, cycle.DateEnd); err != nil {
		return err
	}

	cycle.Name = getCycleName(cycle.DateInit, cycle.DateEnd)

	if err := repository.Create(cycle); err != nil {
		return err
	}

	return nil
}

func Update(cycleID int, currentCycle *cycles.Cycle) error {
	if err := checkDates(currentCycle.DateInit, currentCycle.DateEnd); err != nil {
		return err
	}

	if err := repository.Update(cycleID, currentCycle); err != nil {
		return err
	}

	return nil
}

func Delete(cycleID int) error {
	return repository.Delete(cycleID)
}

func Finish(cycleID int) error {
	cycle, err := repository.GetByID(cycleID)
	if err != nil {
		return err
	}

	cycleExpenses, err := expensesRepository.GetByCycleID(cycleID)
	if err != nil {
		return err
	}

	var cycleHistory cycles.History
	cycleHistory.PocketName = cycle.PocketName
	cycleHistory.CycleName = cycle.Name
	cycleHistory.Budget = cycle.Budget
	cycleHistory.Spent = service.SumExpenses(cycleExpenses)
	cycleHistory.SpentRatio = service.GetSpentRatio(float64(cycleHistory.Budget), float64(cycleHistory.Spent))
	cycleHistory.DateInit = cycle.DateInit
	cycleHistory.DateEnd = cycle.DateEnd

	if err = repository.CreateHistory(cycleHistory); err != nil {
		return err
	}

	if err = repository.Close(cycleID); err != nil {
		return err
	}

	if err = conceptsRepository.BulkUpdatePayed(cycle.PocketID, false); err != nil {
		return err
	}

	return nil
}

func checkDates(dateInit, dateEnd string) error {
	parsedDateInit, err := time.Parse(customdate.StandardDateTimeFormat, dateInit)
	if err != nil {
		return err
	}

	parsedDateEnd, err := time.Parse(customdate.StandardDateTimeFormat, dateEnd)
	if err != nil {
		return err
	}

	if parsedDateEnd.After(parsedDateInit) {
		return nil
	}

	return fmt.Errorf("date_end %s must be grather than date_init %s", dateEnd, dateInit)
}

func getCycleName(dateInit, dateEnd string) string {
	parsedDateInit, err := time.Parse(customdate.StandardDateTimeFormat, dateInit)
	if err != nil {
		return ""
	}

	parsedDateEnd, err := time.Parse(customdate.StandardDateTimeFormat, dateEnd)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s %d", MonthNamesShort[int(parsedDateInit.Month())-1], parsedDateEnd.Year())
}
