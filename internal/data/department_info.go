package data

import (
	"database/sql"
	"errors"
)

type DepartmentInfoModel struct {
	DB *sql.DB
}

func (d DepartmentInfoModel) Insert(dep *DepartmentInfo) error {
	query := `
		INSERT INTO department_info (department_name, staff_quantity, department_director, module_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	row := d.DB.QueryRow(query, dep.DepartmentName, dep.StaffQuantity, dep.DepartmentDirector, dep.ModuleId)
	err := row.Scan(&dep.ID)
	if err != nil {
		return err
	}
	return nil
}

func (d DepartmentInfoModel) Get(id int64) (*DepartmentInfo, error) {
	query := `
SELECT id, department_name, department_director, staff_quantity, module_id
FROM department_info
WHERE id = $1`
	var dep DepartmentInfo
	err := d.DB.QueryRow(query, id).Scan(
		&dep.ID,
		&dep.DepartmentName,
		&dep.DepartmentDirector,
		&dep.StaffQuantity,
		&dep.ModuleId,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &dep, nil
}

func (d DepartmentInfoModel) GetAll() (*DepartmentInfo, error) {
	query := `
SELECT id, department_name, department_director, staff_quantity, module_id
FROM department_info`
	var dep DepartmentInfo
	err := d.DB.QueryRow(query).Scan(
		&dep.ID,
		&dep.DepartmentName,
		&dep.DepartmentDirector,
		&dep.StaffQuantity,
		&dep.ModuleId,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &dep, nil
}
