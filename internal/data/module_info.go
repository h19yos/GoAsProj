package data

import (
	"context"
	"database/sql"
	"time"
)

type ModuleInfoModel struct {
	DB *sql.DB
}

func (mm ModuleInfoModel) Insert(mod *ModuleInfo) error {
	query := `INSERT INTO module_info (created_at, updated_at, module_name, module_duration, exam_type, version) VALUES ($1, $2, $3, $4, $5, $6)`
	args := []any{mod.CreatedAt, mod.UpdatedAt, mod.ModuleName, mod.ModuleDuration, mod.ExamType, mod.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return mm.DB.QueryRowContext(ctx, query, args...).Scan(&mod.ID, &mod.CreatedAt, &mod.Version)
}

func (mm ModuleInfoModel) Get(id int64) (ModuleInfo, error) {
	var mod ModuleInfo
	row := mm.DB.QueryRow("SELECT id, created_at, updated_at, module_name, module_duration, exam_type, version FROM module_info WHERE id = $1", id)
	err := row.Scan(&mod.ID, &mod.CreatedAt, &mod.UpdatedAt, &mod.ModuleName, &mod.ModuleDuration, &mod.ExamType, &mod.Version)
	if err != nil {
		return ModuleInfo{}, err
	}
	return mod, nil
}

func (mm ModuleInfoModel) Update(mod *ModuleInfo) error {
	_, err := mm.DB.Exec("UPDATE module_info SET updated_at = $1, module_name = $2, module_duration = $3, exam_type = $4, version = $5 WHERE id = $6",
		time.Now(), mod.ModuleName, mod.ModuleDuration, mod.ExamType, mod.Version, mod.ID)
	if err != nil {
		return err
	}
	return nil
}

func (mm ModuleInfoModel) Delete(id int64) error {
	_, err := mm.DB.Exec("DELETE FROM module_info WHERE id = $1", id)
	if err != nil {
		return err
	}
	return nil
}
