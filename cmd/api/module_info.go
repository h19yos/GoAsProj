package main

import (
	"errors"
	"fmt"
	"gaproject.terminator8000.net/internal/data"
	"net/http"
	"time"
)

func (app *application) createModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		ModuleName     string        `json:"moduleName"`
		ModuleDuration time.Duration `json:"moduleDuration"`
		ExamType       string        `json:"examType"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	mod := &data.ModuleInfo{
		ModuleName:     input.ModuleName,
		ModuleDuration: input.ModuleDuration,
		ExamType:       input.ExamType,
	}

	err = app.models.ModuleInfo.Insert(mod)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/moduleinfo/%d", mod.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"module_info": mod}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	mod, err := app.models.ModuleInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": mod}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	mod, err := app.models.ModuleInfo.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		ModuleName     string        `json:"moduleName"`
		ModuleDuration time.Duration `json:"moduleDuration"`
		ExamType       string        `json:"examType"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updated_mod := &data.ModuleInfo{
		ModuleName:     input.ModuleName,
		ModuleDuration: input.ModuleDuration,
		ExamType:       input.ExamType,
	}

	err = app.models.ModuleInfo.Update(updated_mod)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"module_info": mod}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteModuleInfoHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.ModuleInfo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "moduleinfo successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
