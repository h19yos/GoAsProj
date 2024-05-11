package main

import (
	"errors"
	"fmt"
	"gaproject.terminator8000.net/internal/data"
	"net/http"
)

func (app *application) createDepInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DepartmentName     string `json:"departmentName"`
		StaffQuantity      int    `json:"staffQuantity"`
		DepartmentDirector string `json:"departmentDirector"`
		ModuleId           int    `json:"moduleID"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dep := &data.DepartmentInfo{
		DepartmentName:     input.DepartmentName,
		StaffQuantity:      input.StaffQuantity,
		DepartmentDirector: input.DepartmentDirector,
		ModuleId:           input.ModuleId,
	}

	err = app.models.DepartmentInfo.Insert(dep)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/departamentinfo/%d", dep.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"departament_info": dep}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getDepInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	dep, err := app.models.DepartmentInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"departament_info": dep}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getAllDepInfoHandler(w http.ResponseWriter, r *http.Request) {
	dep, err := app.models.DepartmentInfo.GetAll()
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"departament_info": dep}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
