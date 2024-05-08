package main

import (
	"alumnihub/internal/models"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "active",
		Message: "Alumnihub",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) AllAlumni(w http.ResponseWriter, r *http.Request) {
	alumni, err := app.DB.AllAlumni()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, alumni)
}

func (app *application) Alumni(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	alumniID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	course, err := app.DB.Alumni(alumniID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, course)
}

// //////////////////
// Handler Articles
// //////////////////
func (app *application) allArticles(w http.ResponseWriter, r *http.Request) {
	article, err := app.DB.AllArticles()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, article)
}

func (app *application) article(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	article, err := app.DB.Article(articleID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, article)
}

func (app *application) insertArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article

	err := app.readJSON(w, r, &article)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()
	article.PublishedAt = time.Now()

	newID, err := app.DB.InsertArticle(article)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	message := fmt.Sprintf("Article added to %d", newID)

	resp := JSONResponse{
		Error:   false,
		Message: message,
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) updateArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload models.Article

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if articleID != payload.ID {
		app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}

	article, err := app.DB.Article(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	article.Title = payload.Title
	article.Slug = payload.Slug
	article.Body = payload.Body
	article.Status = payload.Status
	article.UpdatedAt = time.Now()
	article.PublishedAt = time.Now()

	err = app.DB.UpdateArticle(*article)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "article updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) deleteArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteArticle(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Article deleted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

// //////////////////
// Handler Forms
// //////////////////
func (app *application) allForms(w http.ResponseWriter, r *http.Request) {
	form, err := app.DB.AllForms()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, form)
}

func (app *application) form(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	form, err := app.DB.Form(formID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, form)
}

func (app *application) showForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	form, err := app.DB.ShowForm(formID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, form)
}

func (app *application) showFormAnswers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	form, err := app.DB.ShowFormAnswers(formID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, form)
}

func (app *application) insertForm(w http.ResponseWriter, r *http.Request) {
	var form models.Form

	err := app.readJSON(w, r, &form)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()

	newID, err := app.DB.InsertForm(form)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	message := fmt.Sprintf("Form added to %d", newID)

	resp := JSONResponse{
		Error:   false,
		Message: message,
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) updateForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload models.Form

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if formID != payload.ID {
		app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}

	form, err := app.DB.Form(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	form.Title = payload.Title
	form.Description = payload.Description
	form.HasTimeLimit = payload.HasTimeLimit
	form.StartDate = payload.StartDate
	form.UpdatedAt = time.Now()

	err = app.DB.UpdateForm(*form)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "article updated",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) deleteForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteForm(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Form deleted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) question(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	qID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	question, err := app.DB.Question(qID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, question)
}

func (app *application) insertQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Question

	err := app.readJSON(w, r, &question)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	question.CreatedAt = time.Now()
	question.UpdatedAt = time.Now()

	newID, err := app.DB.InsertQuestion(question)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// handle options when the type is multiple_choice
	if question.Type == "multiple_choice" {
		err = app.DB.UpdateQuestionOptions(newID, question.OptionsArray)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Question added",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) insertAnswers(w http.ResponseWriter, r *http.Request) {
	var answers []*models.Answer

	err := app.readJSON(w, r, &answers)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.InsertAnswers(answers)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Answers submitted",
	}

	app.writeJSON(w, http.StatusAccepted, resp)
}

func (app *application) uploadImage(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max file size
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Get the file from the form data
	file, handler, err := r.FormFile("image")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer file.Close()

	// Determine the file extension
	ext := filepath.Ext(handler.Filename)

	// Create a temporary file to store the uploaded image
	tempFile, err := os.CreateTemp("public", "upload-*"+ext)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer tempFile.Close()

	// Write the file content to the temporary file
	_, err = io.Copy(tempFile, file)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	fileName := filepath.Base(tempFile.Name())
	filePath := filepath.Join("public", fileName)
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	image := models.Image{
		FilePath: filePath,
		FileName: fileName,
	}

	app.writeJSON(w, http.StatusAccepted, image)
}

func (app *application) serveImage(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan path file dari URL
	imagePath := chi.URLParam(r, "image_path")

	// Gabungkan path file dengan direktori "public"
	imageFile := filepath.Join("public", imagePath)

	// Buka file gambar
	http.ServeFile(w, r, imageFile)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read JSON Payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid email"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid password"), http.StatusBadRequest)
		return
	}

	// create a jwt user
	u := jwtUser{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.IsAdmin,
	}

	// Generate token
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	refreshCookie := app.auth.GetRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)

	app.writeJSON(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.errorJSON(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.errorJSON(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			u := jwtUser{
				ID:       user.ID,
				Username: user.Username,
			}

			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.errorJSON(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.GetRefreshCookie(tokenPairs.RefreshToken))

			app.writeJSON(w, http.StatusOK, tokenPairs)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}
