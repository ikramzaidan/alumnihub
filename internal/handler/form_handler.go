package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"alumnihub/internal/auth"
	"alumnihub/internal/models"
	"alumnihub/internal/service"
	"alumnihub/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/xuri/excelize/v2"
)

func NewFormHandler(formService *service.FormService) *Handler {
	return &Handler{FormService: formService}
}

func (h *Handler) AllForms(w http.ResponseWriter, r *http.Request) {
	forms, err := h.FormService.AllForms()
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}
	_ = utils.WriteJSON(w, http.StatusOK, forms)
}

func (h *Handler) Form(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	form, err := h.FormService.GetForm(formID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, form)
}

func (h *Handler) ShowForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	form, err := h.FormService.ShowForm(formID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, form)
}

func (h *Handler) InsertForm(w http.ResponseWriter, r *http.Request) {
	var form models.Form
	if err := utils.ReadJSON(w, r, &form); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	formID, err := h.FormService.CreateForm(form)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: fmt.Sprintf("New survey has been successfully created with id %d", formID)}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateForm(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	var payload models.Form
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.FormService.UpdateForm(formID, payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Survey has been successfully updated"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.FormService.DeleteForm(id); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Survey has been permanently deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) Question(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	qID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	question, err := h.FormService.GetQuestion(qID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, question)
}

func (h *Handler) InsertQuestion(w http.ResponseWriter, r *http.Request) {
	var question models.Question
	if err := utils.ReadJSON(w, r, &question); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.FormService.CreateQuestion(question); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Question has been successfully created"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	questionID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	var payload models.Question
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.FormService.UpdateQuestion(questionID, payload); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Question has been successfully updated"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.FormService.DeleteQuestion(id); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Question has been permanently deleted"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) ShowFormAnswers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	form, err := h.FormService.ShowFormAnswers(formID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, form)
}

func (h *Handler) InsertAnswers(w http.ResponseWriter, r *http.Request) {
	var answers []*models.Answer
	if err := utils.ReadJSON(w, r, &answers); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	if err := h.FormService.InsertAnswers(answers); err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Answers has been successfully saved"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) ShowQuestionAnswers(w http.ResponseWriter, r *http.Request) {
	fID := chi.URLParam(r, "fid")
	formID, err := strconv.Atoi(fID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	qID := chi.URLParam(r, "qid")
	questionID, err := strconv.Atoi(qID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	answers, err := h.FormService.ShowQuestionAnswers(formID, questionID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, answers)
}

func (h *Handler) UserAnswers(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(auth.UserClaimsKey).(*auth.Claims)
	if !ok {
		_ = utils.ErrorJSON(w, errors.New("failed to retrieve user claims"))
		return
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	answers, err := h.FormService.GetAnswersByUser(userID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, answers)
}

func (h *Handler) ExportAnswers(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	formID, err := strconv.Atoi(id)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	form, err := h.FormService.GetForm(formID)
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	xlsx := excelize.NewFile()
	for _, formQuestion := range form.Questions {
		sheetName := fmt.Sprintf("Question %d", formQuestion.ID)
		index, err := xlsx.NewSheet(sheetName)
		if err != nil {
			_ = utils.ErrorJSON(w, err)
			return
		}

		question, err := h.FormService.GetQuestion(formQuestion.ID)
		if err != nil {
			_ = utils.ErrorJSON(w, err)
			return
		}

		xlsx.MergeCell(sheetName, "A1", "G1")
		xlsx.SetCellValue(sheetName, "A1", question.Question)
		style, err := xlsx.NewStyle(&excelize.Style{Font: &excelize.Font{Bold: true, Size: 11, Color: "000000"}})
		if err != nil {
			_ = utils.ErrorJSON(w, err)
			return
		}
		xlsx.SetCellStyle(sheetName, "A1", "A1", style)

		xlsx.SetCellValue(sheetName, "A3", "No")
		xlsx.SetCellValue(sheetName, "B3", "User ID")
		xlsx.SetCellValue(sheetName, "C3", "Jawaban")
		style, err = xlsx.NewStyle(&excelize.Style{
			Font: &excelize.Font{Bold: true, Size: 11, Color: "000000"},
			Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#60A5FA"}},
			Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			},
		})
		if err != nil {
			_ = utils.ErrorJSON(w, err)
			return
		}
		xlsx.SetCellStyle(sheetName, "A3", "C3", style)
		xlsx.SetColWidth(sheetName, "A", "A", 4)
		xlsx.SetColWidth(sheetName, "B", "B", 12)
		xlsx.SetColWidth(sheetName, "C", "C", 32)

		for i, answer := range question.Answers {
			cellRow := i + 4
			xlsx.SetCellValue(sheetName, fmt.Sprintf("A%d", cellRow), i+1)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("B%d", cellRow), answer.UserID)
			xlsx.SetCellValue(sheetName, fmt.Sprintf("C%d", cellRow), answer.Answer)
			style, err := xlsx.NewStyle(&excelize.Style{Border: []excelize.Border{
				{Type: "top", Color: "000000", Style: 1},
				{Type: "left", Color: "000000", Style: 1},
				{Type: "bottom", Color: "000000", Style: 1},
				{Type: "right", Color: "000000", Style: 1},
			}})
			if err != nil {
				_ = utils.ErrorJSON(w, err)
				return
			}
			xlsx.SetCellStyle(sheetName, fmt.Sprintf("A%d", cellRow), fmt.Sprintf("C%d", cellRow), style)
		}

		xlsx.SetActiveSheet(index)
	}

	currentTime := time.Now().Format("2006-01-02_15-04-05")
	fileName := fmt.Sprintf("survei_%s_%s.xlsx", utils.SanitizeFileName(form.Title), currentTime)

	err = xlsx.SaveAs("/home/ikramzaidann/alumnihub/public/excel/workbook.xlsx")
	if err != nil {
		_ = utils.ErrorJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	_ = xlsx.Write(w)
}
