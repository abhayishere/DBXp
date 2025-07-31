package handlers

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgconn"
)

type Export struct {
	columns []string   // Column headers
	rows    [][]string // All the row data
}

func (e *Export) AddRow(row []any) {
	if len(row) != len(e.columns) {
		return // Ensure row matches column count
	}
	strValues := make([]string, len(row))
	for i, v := range row {
		if v != nil {
			strValues[i] = fmt.Sprintf("%v", v)
		} else {
			strValues[i] = "NULL"
		}
	}
	e.rows = append(e.rows, strValues) // Store row data for export
}

func (e *Export) AddColumns(columns []pgconn.FieldDescription) {
	if len(columns) == 0 {
		return // No columns to add
	}

	e.rows = nil
	e.columns = nil
	for _, col := range columns {
		e.columns = append(e.columns, col.Name)
	}
}

func (e *Export) ExportToCSV() error {
	file, err := os.Create("export.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.Write(e.columns) // Write header
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
