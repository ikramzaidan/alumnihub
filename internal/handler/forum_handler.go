package handler

import (
	"errors"
	"net/http"
	"strconv"

	"alumnihub/internal/auth"
	"alumnihub/internal/service"
	"alumnihub/internal/utils"

	"github.com/go-chi/chi/v5"
)

func NewForumHandler(forumService *service.ForumService) *Handler {
	return &Handler{ForumService: forumService}
}

func (h *Handler) AllForums(w http.ResponseWriter, r *http.Request) {
	forums, err := h.ForumService.AllForums()
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, forums)
}

func (h *Handler) AllUserForums(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	forums, err := h.ForumService.GetForumsByUser(username)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, forums)
}

func (h *Handler) Forum(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	forum, err := h.ForumService.GetForum(forumID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, forum)
}

func (h *Handler) InsertForum(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Forum string `json:"forum_text"`
	}
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

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

	if err := h.ForumService.CreateForum(userID, payload.Forum); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "New forum has been successfully created"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) DeleteForum(w http.ResponseWriter, r *http.Request) {
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
	forumID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ForumService.DeleteForum(userID, forumID); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Forum post has been successfully deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) InsertComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	var payload struct {
		ForumID int    `json:"forum_id"`
		Comment string `json:"reply_text"`
	}
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if forumID != payload.ForumID {
		_ = utils.ErrorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}

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

	if err := h.ForumService.InsertComment(userID, forumID, payload.Comment); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "New comment has been succesfully created"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UserLikes(w http.ResponseWriter, r *http.Request) {
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

	likes, err := h.ForumService.GetLikesByUser(userID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, likes)
}

func (h *Handler) InsertLike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

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

	if err := h.ForumService.AddLike(userID, forumID); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Like has been succesfully added"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) DeleteLike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

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

	if err := h.ForumService.RemoveLike(userID, forumID); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Like has been successfully deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}
