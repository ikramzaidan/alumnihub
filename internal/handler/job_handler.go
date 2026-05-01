package handler

import (
	"errors"
	"net/http"
	"strconv"

	"alumnihub/internal/auth"
	"alumnihub/internal/models"
	"alumnihub/internal/service"
	"alumnihub/internal/utils"

	"github.com/go-chi/chi/v5"
)

func NewJobHandler(jobService *service.JobService) *Handler {
	return &Handler{JobService: jobService}
}

func (h *Handler) AllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.JobService.AllJobs()
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, jobs)
}

func (h *Handler) Job(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	job, err := h.JobService.GetJob(jobID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, job)
}

func (h *Handler) InsertJob(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserClaimsKey).(*auth.Claims)
	if !ok {
		_ = utils.ErrorJSON(w, errors.New("no claims in context"))
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	var job models.Job
	if err := utils.ReadJSON(w, r, &job); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.JobService.CreateJob(userID, job); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "New job has been successfully posted"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	var payload models.Job
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.JobService.UpdateJob(jobID, payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Job has been successfully updated"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserClaimsKey).(*auth.Claims)
	if !ok {
		_ = utils.ErrorJSON(w, errors.New("no claims in context"))
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.JobService.DeleteJob(userID, jobID); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Job has been successfully deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}
