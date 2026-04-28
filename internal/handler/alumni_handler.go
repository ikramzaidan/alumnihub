package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"alumnihub/internal/models"
	"alumnihub/internal/service"
	"alumnihub/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"
)

func NewAlumniHandler(alumniService *service.AlumniService) *Handler {
	return &Handler{AlumniService: alumniService}
}

func (h *Handler) AllAlumni(w http.ResponseWriter, r *http.Request) {
	alumni, err := h.AlumniService.All()
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, alumni)
}

func (h *Handler) Alumni(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	alumniID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	alumni, err := h.AlumniService.Get(alumniID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, alumni)
}

func (h *Handler) InsertAlumni(w http.ResponseWriter, r *http.Request) {
	var alumni models.Alumni
	if err := utils.ReadJSON(w, r, &alumni); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.AlumniService.Create(alumni); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "New alumni has been successfully added"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) ImportAlumni(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	defer file.Close()

	f, err := excelize.OpenReader(file)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	var alumniList []models.Alumni
	for _, row := range rows[1:] {
		year, err := strconv.Atoi(row[5])
		if err != nil {
			log.Printf("Skipping row due to invalid graduation year: %v", row)
			continue
		}

		alumni := models.Alumni{
			NISN:   row[0],
			NIS:    row[1],
			Name:   row[2],
			Gender: row[3],
			Phone:  row[4],
			Year:   year,
			Class:  row[6],
		}

		alumniList = append(alumniList, alumni)
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, alumniList)
}

func (h *Handler) InsertImportAlumni(w http.ResponseWriter, r *http.Request) {
	var alumniList []models.Alumni
	if err := utils.ReadJSON(w, r, &alumniList); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.AlumniService.InsertBulk(alumniList); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: fmt.Sprintf("%d alumni berhasil ditambahkan", len(alumniList))}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateAlumni(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	alumniID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	var payload models.Alumni
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.AlumniService.Update(alumniID, payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Alumni has been successfully updated"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteAlumni(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.AlumniService.Delete(id); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Alumni has been permanently deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}
