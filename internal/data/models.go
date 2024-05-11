package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
)

type Models struct {
	Movies         MovieModel
	ModuleInfo     ModuleInfoModel
	DepartmentInfo DepartmentInfoModel
	//Permissions     PermissionModel // Add a new Permissions field.
	//Users  UsersModel
	Tokens   TokenModel
	UserInfo UserInfoModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:         MovieModel{DB: db},
		ModuleInfo:     ModuleInfoModel{DB: db},
		DepartmentInfo: DepartmentInfoModel{DB: db},
		//Permissions:     PermissionModel{DB: db}, // Initialize a new PermissionModel instance.
		//Users:  UsersModel{DB: db},
		Tokens:   TokenModel{DB: db},
		UserInfo: UserInfoModel{DB: db},
	}
}

type ModuleInfo struct {
	ID             int           `json:"id"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	ModuleName     string        `json:"moduleName"`
	ModuleDuration time.Duration `json:"moduleDuration"`
	ExamType       string        `json:"examType"`
	Version        string        `json:"version"`
}

type DepartmentInfo struct {
	ID                 int    `json:"id"`
	DepartmentName     string `json:"departmentName"`
	StaffQuantity      int    `json:"staffQuantity"`
	DepartmentDirector string `json:"departmentDirector"`
	ModuleId           int    `json:"moduleId"`
}
