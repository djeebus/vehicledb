package db

import (
	"fmt"
	"strings"
)

var vehiclesTable = `
CREATE TABLE vehicles (
	"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"year" INTEGER NOT NULL,
	"make" STRING NOT NULL,
	"model" STRING NOT NULL
)`

type Vehicle struct {
	VehicleID RowID

	Year  Year
	Make  string
	Model string
}

func CreateVehicle(year Year, make, model string) (*Vehicle, error) {
	query := `INSERT INTO vehicles (year, make, model) VALUES (?, ?, ?)`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prep vehicle statement: %w", err)
	}

	result, err := stmt.Exec(year, make, model)
	if err != nil {
		return nil, fmt.Errorf("failed to exec create vehicle statement: %w", err)
	}

	defer stmt.Close()

	lastInserted, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve last inserted id: %w", err)
	}

	vehicle := Vehicle{
		VehicleID: RowID(lastInserted),
		Year:      year,
		Make:      make,
		Model:     model,
	}
	return &vehicle, nil
}

func GetVehicle(vehicleID RowID) (*Vehicle, error) {
	query := `SELECT year, make, model FROM vehicles WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare get vehicle query: %v", err)
	}

	rows, err := stmt.Query(vehicleID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get vehicle query: %v", err)
	}
	defer rows.Close()

	var year uint16
	var vehicleMake, model string

	for rows.Next() {
		err = rows.Scan(&year, &vehicleMake, &model)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		vehicle := Vehicle{
			VehicleID: vehicleID,
			Year:      Year(year),
			Make:      vehicleMake,
			Model:     model,
		}
		return &vehicle, nil
	}

	return nil, nil
}

func UpdateVehicle(vehicleID RowID, year *NullYear, vehicleMake, model *NullString) error {
	var values = make([]interface{}, 0, 0)
	sets := make([]string, 0, 0)

	if year.Valid {
		values = append(values, year.Year)
		sets = append(sets, "year = ?")
	}

	if vehicleMake.Valid {
		values = append(values, vehicleMake.String)
		sets = append(sets, "make = ?")
	}

	if model.Valid {
		values = append(values, model.String)
		sets = append(sets, "model = ?")
	}

	if len(values) == 0 {
		return nil
	}

	values = append(values, vehicleID)

	query := fmt.Sprintf(`UPDATE vehicles SET %s WHERE id = ?`, strings.Join(sets, ", "))
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update vehicle query: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(values...)
	if err != nil {
		return fmt.Errorf("failed to execute update vehicle query: %v", err)
	}

	return nil
}

func DeleteVehicle(vehicleID RowID) error {
	query := `DELETE FROM vehicles WHERE id = ?`
	stmt, err := sqlDb.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare delete vehicle query: %v", err)
	}

	_, err = stmt.Exec(vehicleID)
	if err != nil {
		return fmt.Errorf("failed to execute update vehicle query: %v", err)
	}

	return nil
}
