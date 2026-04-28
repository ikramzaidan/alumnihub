package handler

import (
	"net/http"
	"strconv"

	"alumnihub/internal/models"
	"alumnihub/internal/service"
	"alumnihub/internal/utils"

	"github.com/go-chi/chi/v5"
)

func NewArticleHandler(articleService *service.ArticleService) *Handler {
	return &Handler{ArticleService: articleService}
}

func (h *Handler) AllArticles(w http.ResponseWriter, r *http.Request) {
	articles, err := h.ArticleService.All()
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, articles)
}

func (h *Handler) Article(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	article, err := h.ArticleService.GetBySlug(slug)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, article)
}

func (h *Handler) ShowArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	article, err := h.ArticleService.GetByID(articleID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, article)
}

func (h *Handler) InsertArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	if err := utils.ReadJSON(w, r, &article); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ArticleService.Create(article); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "New article has been successfully created"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	var payload models.Article
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ArticleService.Update(articleID, payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Article has been successfully updated"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.ArticleService.Delete(id); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Article has been permanently deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}
