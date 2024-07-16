package main

import (
	"alumnihub/internal/models"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
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

func (app *application) Dashboard(w http.ResponseWriter, r *http.Request) {
	type Payload struct {
		CountAlumni        int               `json:"count_alumni"`
		CountAlumniAccount int               `json:"count_alumni_account"`
		Profiles           []*models.Profile `json:"profiles,omitempty"`
	}

	countAlumni, err := app.DB.CountAlumni()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	countAlumniAccount, err := app.DB.CountAlumniAccount()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	profiles, err := app.DB.GetProfiles()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	payload := Payload{
		CountAlumni:        countAlumni,
		CountAlumniAccount: countAlumniAccount,
		Profiles:           profiles,
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

	alumni, err := app.DB.Alumni(alumniID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	profile, err := app.DB.GetProfileByAlumniID(alumniID)
	if err != nil && err != sql.ErrNoRows {
		app.errorJSON(w, err)
		return
	}

	if profile != nil {
		user, err := app.DB.GetUserByID(profile.UserID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		alumni.UserUsername = user.Username
	}

	_ = app.writeJSON(w, http.StatusOK, alumni)
}

func (app *application) insertAlumni(w http.ResponseWriter, r *http.Request) {
	var alumni models.Alumni

	err := app.readJSON(w, r, &alumni)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.InsertAlumni(alumni)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	message := "New alumni has been successfully added"

	resp := JSONResponse{
		Error:   false,
		Message: message,
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) importAlumni(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("file")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	defer file.Close()

	// Read the Excel file
	f, err := excelize.OpenReader(file)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	rows, err := f.GetRows("Sheet1") // Assuming data is in Sheet1
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var alumniList []models.Alumni

	for _, row := range rows[1:] { // Skipping header row
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

	app.writeJSON(w, http.StatusAccepted, alumniList)
}

func (app *application) insertImportAlumni(w http.ResponseWriter, r *http.Request) {
	var alumniList []models.Alumni

	err := app.readJSON(w, r, &alumniList)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	count := 0

	for _, alumni := range alumniList {
		err = app.DB.InsertAlumni(alumni)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		count++
	}

	message := fmt.Sprintf("%d alumni berhasil ditambahkan", count)

	resp := JSONResponse{
		Error:   false,
		Message: message,
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) updateAlumni(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	alumniID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload models.Alumni

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	alumni, err := app.DB.Alumni(alumniID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	alumni.Name = payload.Name
	alumni.Gender = payload.Gender
	alumni.Phone = payload.Phone
	alumni.Year = payload.Year
	alumni.Class = payload.Class
	alumni.NISN = payload.NISN
	alumni.NIS = payload.NIS

	err = app.DB.UpdateAlumni(*alumni)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Alumni has been successfully updated",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) deleteAlumni(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteAlumni(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Alumni has been permanently deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) profile(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	userID, err := app.DB.GetUserIDByUsername(username)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	profile, err := app.DB.GetProfileByUserID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	alumniName, err := app.DB.GetAlumniNameByID(profile.AlumniID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	userUsername, err := app.DB.GetUserUsernameByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	userPhoto, err := app.DB.GetUserPhotoByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	profile.UserName = alumniName
	profile.UserUsername = userUsername
	profile.Photo = userPhoto

	app.writeJSON(w, http.StatusOK, profile)
}

func (app *application) myProfile(w http.ResponseWriter, r *http.Request) {
	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	var profile models.Profile

	if !claims.IsAdmin {
		profileAlumni, err := app.DB.GetProfileByUserID(userID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		alumniName, err := app.DB.GetAlumniNameByID(profileAlumni.AlumniID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		profile = *profileAlumni
		profile.UserName = alumniName
	} else {
		profileAdmin, err := app.DB.GetAdminProfileByUserID(userID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}

		profile = *profileAdmin
	}

	userUsername, err := app.DB.GetUserUsernameByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	userPhoto, err := app.DB.GetUserPhotoByID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	profile.UserUsername = userUsername
	profile.Photo = userPhoto

	app.writeJSON(w, http.StatusOK, profile)
}

func (app *application) updateProfile(w http.ResponseWriter, r *http.Request) {
	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	var payload models.Profile

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	profile, err := app.DB.GetProfileByUserID(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	profile.Bio = payload.Bio
	profile.Location = payload.Location
	profile.Facebook = payload.Facebook
	profile.Instagram = payload.Instagram
	profile.Twitter = payload.Twitter
	profile.Tiktok = payload.Tiktok
	profile.Photo = payload.Photo

	err = app.DB.UpdateProfile(*profile)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Profile updated",
	}

	app.writeJSON(w, http.StatusOK, resp)
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
	slug := chi.URLParam(r, "slug")

	article, err := app.DB.ArticleBySlug(slug)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, article)
}

func (app *application) showArticle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	articleId, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	article, err := app.DB.Article(articleId)
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

	imgSrc, err := app.getFirstImageFromHtml(article.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")

	article.Image = imgSrc
	article.CreatedAt = time.Now().In(loc)
	article.UpdatedAt = time.Now().In(loc)
	article.PublishedAt = time.Now().In(loc)

	_, err = app.DB.InsertArticle(article)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "New article has been successfully created",
	}

	app.writeJSON(w, http.StatusCreated, resp)
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
		app.errorJSON(w, errors.New("invalid request"))
		return
	}

	article, err := app.DB.Article(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	imgSrc, err := app.getFirstImageFromHtml(payload.Body)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	article.Image = imgSrc
	article.Title = payload.Title
	article.Slug = payload.Slug
	article.Body = payload.Body
	article.Status = payload.Status
	article.UpdatedAt = time.Now()
	if article.Status != "published" && payload.Status == "published" {
		article.PublishedAt = time.Now()
	}

	err = app.DB.UpdateArticle(*article)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Article has been successfully updated",
	}

	app.writeJSON(w, http.StatusOK, resp)
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
		Message: "Article has been permanently deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)
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

func (app *application) insertForm(w http.ResponseWriter, r *http.Request) {
	var form models.Form

	err := app.readJSON(w, r, &form)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	form.CreatedAt = time.Now()
	form.UpdatedAt = time.Now()

	_, err = app.DB.InsertForm(form)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "New survey has been successfully created",
	}

	app.writeJSON(w, http.StatusCreated, resp)
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
		Message: "Survey has been successfully updated",
	}

	app.writeJSON(w, http.StatusOK, resp)
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
		Message: "Survey has been permanently deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)
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

	groupAnswers, err := app.DB.GroupAnswersByQuestion(question.FormID, qID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	question.GroupAnswer = groupAnswers

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
		Message: "Question has been successfully created",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) updateQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	questionID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload models.Question

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if questionID != payload.ID {
		app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}

	question, err := app.DB.Question(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	question.Question = payload.Question
	question.Type = payload.Type
	question.Extension = payload.Extension
	question.UpdatedAt = time.Now()
	question.ID = payload.ID
	question.OptionsArray = payload.OptionsArray

	err = app.DB.UpdateQuestion(*question)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// handle options when the type is multiple_choice
	if question.Type == "multiple_choice" {
		err = app.DB.UpdateQuestionOptions(payload.ID, question.OptionsArray)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.DB.DeleteQuestionOptions(payload.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	if question.Extension && payload.QuestionExtension != nil {
		err = app.DB.UpdateQuestionExtension(payload.QuestionExtension)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.DB.DeleteQuestionExtension(payload.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Question has been successfully updated",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) deleteQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteQuestion(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Question has been permanently deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)
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
		Message: "Answers has been successfully saved",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) showQuestionAnswers(w http.ResponseWriter, r *http.Request) {
	fID := chi.URLParam(r, "fid")
	formID, err := strconv.Atoi(fID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	qID := chi.URLParam(r, "qid")
	questionID, err := strconv.Atoi(qID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	groupAnswers, err := app.DB.GroupAnswersByQuestion(formID, questionID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, groupAnswers)
}

func (app *application) userAnswers(w http.ResponseWriter, r *http.Request) {
	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	answers, err := app.DB.GetAnswersByUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, answers)
}

// //////////////////
// Handler Forums
// //////////////////
func (app *application) allForums(w http.ResponseWriter, r *http.Request) {
	forums, err := app.DB.AllForums()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, forums)
}

func (app *application) allUserForums(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")

	userID, err := app.DB.GetUserIDByUsername(username)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	forums, err := app.DB.GetForumsByUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, forums)
}

func (app *application) forum(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	form, err := app.DB.Forum(forumID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, form)
}

func (app *application) insertForum(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Forum string `json:"forum_text"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	var forum models.Forum

	forum.Forum = payload.Forum
	forum.UserID = userID
	forum.PublishedAt = time.Now()

	_, err = app.DB.InsertForum(forum)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "New forum has been successfully created",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) deleteForum(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteForum(forumID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Forum has been successfully deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) insertComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload struct {
		ForumID int    `json:"forum_id"`
		Comment string `json:"reply_text"`
	}

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if forumID != payload.ForumID {
		app.errorJSON(w, errors.New("invalid request"), http.StatusBadRequest)
		return
	}

	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	var comment models.Comment

	comment.Comment = payload.Comment
	comment.UserID = userID
	comment.ForumID = payload.ForumID
	comment.PublishedAt = time.Now()

	_, err = app.DB.InsertComment(comment)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "New comment has been succesfully created",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) userLikes(w http.ResponseWriter, r *http.Request) {
	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	likes, err := app.DB.GetLikesByUser(userID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, likes)
}

func (app *application) insertLike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	var like models.Like

	like.ForumID = forumID
	like.UserID = userID
	like.CreatedAt = time.Now()

	err = app.DB.InsertLike(like)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Like has been succesfully added",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) deleteLike(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	forumID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	err = app.DB.DeleteLike(userID, forumID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Like has been successfully deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

// //////////////////
// Handler Jobs
// //////////////////
func (app *application) allJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := app.DB.AllJobs()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, jobs)
}

func (app *application) job(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	job, err := app.DB.Job(jobID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	_ = app.writeJSON(w, http.StatusOK, job)
}

func (app *application) insertJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job

	err := app.readJSON(w, r, &job)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Ambil klaim dari konteks menggunakan tipe kunci khusus
	claims, ok := r.Context().Value(userClaimsKey).(*Claims)
	if !ok {
		app.errorJSON(w, errors.New("no claims in context"))
		return
	}

	// Konversi userID dari string ke int
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		app.errorJSON(w, errors.New("invalid user ID in token"))
		return
	}

	job.UserID = userID
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	_, err = app.DB.InsertJob(job)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "New job has been successfully posted",
	}

	app.writeJSON(w, http.StatusCreated, resp)
}

func (app *application) updateJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var payload models.Job

	err = app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jobID != payload.ID {
		app.errorJSON(w, errors.New("invalid request"))
		return
	}

	job, err := app.DB.Job(payload.ID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	job.JobPosition = payload.JobPosition
	job.Company = payload.Company
	job.JobLocation = payload.JobLocation
	job.JobType = payload.JobType
	job.MinSalary = payload.MinSalary
	job.MaxSalary = payload.MaxSalary
	job.Description = payload.Description
	job.Closed = payload.Closed
	job.UpdatedAt = time.Now()

	err = app.DB.UpdateJob(*job)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Job has been successfully updated",
	}

	app.writeJSON(w, http.StatusOK, resp)
}

func (app *application) deleteJob(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	jobID, err := strconv.Atoi(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.DB.DeleteJob(jobID)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Job has been successfully deleted",
	}

	app.writeJSON(w, http.StatusOK, resp)
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

// Matching NISN before register
func (app *application) registerCheck(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		NISN string `json:"nisn"`
	}

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Validate NISN
	alumni, err := app.DB.GetAlumniByNISN(requestPayload.NISN)
	if err != nil {
		app.errorJSON(w, errors.New("nisn doesn't match any record"), http.StatusBadRequest)
		return
	}

	// Check if user already exist
	_, err = app.DB.GetProfileByAlumniID(alumni.ID)
	if err == nil {
		app.errorJSON(w, errors.New("account already registered"), http.StatusBadRequest)
		return
	}

	// if other error occured, return internal server error
	if err != nil && err != sql.ErrNoRows {
		app.errorJSON(w, errors.New("error checking profile"), http.StatusInternalServerError)
		return
	}

	app.writeJSON(w, http.StatusOK, alumni)

}

func (app *application) register(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		AlumniID int    `json:"alumni_id"`
	}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Validasi payload untuk nilai kosong
	if payload.Username == "" || payload.Email == "" || payload.Password == "" || payload.AlumniID == 0 {
		app.errorJSON(w, errors.New("all fields are required"), http.StatusBadRequest)
		return
	}

	// Cek apakah alumni sudah terdaftar?
	_, err = app.DB.GetProfileByAlumniID(payload.AlumniID)
	if err == nil {
		app.errorJSON(w, errors.New("account already registered"), http.StatusBadRequest)
		return
	}

	var user models.User

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	user.Username = payload.Username
	user.Email = payload.Email
	user.Password = string(hashedPassword)
	user.IsAdmin = false
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	userID, err := app.DB.InsertUser(user)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var profile models.Profile
	profile.AlumniID = payload.AlumniID
	profile.UserID = userID

	_, err = app.DB.InsertProfile(profile)
	if err != nil {
		_ = app.DB.DeleteUser(userID)
		app.errorJSON(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Register success",
	}

	app.writeJSON(w, http.StatusCreated, resp)

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

	app.writeJSON(w, http.StatusOK, tokens)
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
				Role:     user.IsAdmin,
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
	w.WriteHeader(http.StatusNoContent)
}
