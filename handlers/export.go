package handlers

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Export struct {
	columns []string
	rows    [][]string
}

func (e *Export) AddRow(row []string) {
	if len(row) != len(e.columns) {
		return
	}
	e.rows = append(e.rows, row)
}

func (e *Export) AddColumns(columns []string) {
	if len(columns) == 0 {
		return
	}

	e.rows = nil
	e.columns = nil
	for _, col := range columns {
		e.columns = append(e.columns, col)
	}
}

func (e *Export) ExportToCSV() error {
	file, err := os.Create("export.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.Write(e.columns)
	if err != nil {
		return fmt.Errorf("error writing header to CSV: %w", err)
	}
	for _, row := range e.rows {
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("error writing row to CSV: %w", err)
		}
	}
	defer writer.Flush()
	return nil
}
