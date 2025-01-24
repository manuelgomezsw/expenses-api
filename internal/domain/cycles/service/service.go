package service

import (
	"expenses-api/internal/domain/cycles"
	"expenses-api/internal/domain/cycles/repository"
	"fmt"
	"time"
)

var (
	MonthNamesShort = [...]string{"Ene", "Feb", "Mar", "Abr", "May", "Jun", "Jul", "Ago", "Sep", "Oct", "Nov", "Dic"}
)

func Get() ([]cycles.Cycle, error) {
	return repository.Get()
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

func checkDates(dateInit, dateEnd string) error {
	formatDate := "2006-01-02"

	parsedDateInit, err := time.Parse(formatDate, dateInit)
	if err != nil {
		return err
	}

	parsedDateEnd, err := time.Parse(formatDate, dateEnd)
	if err != nil {
		return err
	}

	if parsedDateEnd.After(parsedDateInit) {
		return nil
	}

	return fmt.Errorf("date_end %s must be grather than date_init %s", dateEnd, dateInit)
}

func getCycleName(dateInit, dateEnd string) string {
	formatDate := "2006-01-02"

	parsedDateInit, err := time.Parse(formatDate, dateInit)
	if err != nil {
		return ""
	}

	parsedDateEnd, err := time.Parse(formatDate, dateEnd)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%s %d", MonthNamesShort[int(parsedDateInit.Month())-1], parsedDateEnd.Year())
}
