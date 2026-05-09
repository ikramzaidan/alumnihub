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

	// Prevent directory traversal
	cleanedPath := filepath.Clean(imageFile)
	if !strings.HasPrefix(cleanedPath, filepath.Clean("public")) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	http.ServeFile(w, r, cleanedPath)
}

func (h *Handler) DownloadExcel(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "filename")

	// Prevent directory traversal - only allow filenames without path separators
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") || filename == ".." || strings.HasPrefix(filename, ".") {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	excelFile := filepath.Join("public", "excel", filename)

	// Verify the file exists before serving
	if _, err := os.Stat(excelFile); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	http.ServeFile(w, r, excelFile)
}
