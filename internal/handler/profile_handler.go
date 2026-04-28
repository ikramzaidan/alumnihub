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

func NewProfileHandler(profileService *service.ProfileService) *Handler {
	return &Handler{ProfileService: profileService}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	profileData, err := h.ProfileService.GetByUsername(username)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, profileData)
}

func (h *Handler) MyProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserClaimsKey).(*auth.Claims)
	if !ok {
		_ = utils.ErrorJSON(w, errors.New("no claims in context"))
		return
	}

	profileData, err := h.ProfileService.GetMyProfile(claims)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, profileData)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserClaimsKey).(*auth.Claims)
	if !ok {
		_ = utils.ErrorJSON(w, errors.New("no claims in context"))
		return
	}

	var payload models.Profile
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ProfileService.UpdateProfile(claims, payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Profile updated"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) InsertAlumniEducation(w http.ResponseWriter, r *http.Request) {
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

	var education models.AlumniEducation
	if err := utils.ReadJSON(w, r, &education); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ProfileService.AddEducation(userID, education); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "New education section has been successfully added"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) DeleteAlumniEducation(w http.ResponseWriter, r *http.Request) {
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
	educationID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ProfileService.DeleteEducation(userID, educationID); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "This education section has been permanently deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) InsertAlumniJob(w http.ResponseWriter, r *http.Request) {
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

	var alumnijob models.AlumniJob
	if err := utils.ReadJSON(w, r, &alumnijob); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ProfileService.AddJob(userID, alumnijob); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "New job section has been successfully added"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) DeleteAlumniJob(w http.ResponseWriter, r *http.Request) {
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

	if err := h.ProfileService.DeleteJob(userID, jobID); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "This job section has been permanently deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}
