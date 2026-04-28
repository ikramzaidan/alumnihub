package handler

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"alumnihub/internal/models"
	"alumnihub/internal/utils"

	"github.com/go-chi/chi/v5"
)

func NewUploadHandler() *Handler {
	return &Handler{}
}

func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	file, handler, err := r.FormFile("image")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	defer file.Close()

	ext := filepath.Ext(handler.Filename)
	tempFile, err := os.CreateTemp("public", "upload-*"+ext)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	defer tempFile.Close()

	_, err = io.Copy(tempFile, file)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	err = os.Chmod(tempFile.Name(), 0644)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	fileName := filepath.Base(tempFile.Name())
	filePath := filepath.Join("public", fileName)
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	image := models.Image{
		FilePath: filePath,
		FileName: fileName,
	}

	_ = utils.WriteJSON(w, http.StatusAccepted, image)
}

func (h *Handler) ServeImage(w http.ResponseWriter, r *http.Request) {
	imagePath := chi.URLParam(r, "image_path")
	imageFile := filepath.Join("public", imagePath)
	http.ServeFile(w, r, imageFile)
}
