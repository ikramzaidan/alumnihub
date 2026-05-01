package handler

import (
	"net/http"

	"alumnihub/internal/auth"
	"alumnihub/internal/service"
	"alumnihub/internal/utils"
)

type authPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	AlumniID int    `json:"alumni_id"`
}

type registerCheckPayload struct {
	NISN string `json:"nisn"`
}

type forumPayload struct {
	Forum string `json:"forum_text"`
}

type forgotPasswordPayload struct {
	Email string `json:"email"`
}

type resetPasswordPayload struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func NewAuthHandler(authService *service.AuthService, authSvc *auth.Auth) *Handler {
	return &Handler{AuthService: authService, Auth: authSvc}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	payload := struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}{
		Status:  "active",
		Message: "Alumnihub",
	}

	_ = utils.WriteJSON(w, http.StatusOK, payload)
}

func (h *Handler) RegisterCheck(w http.ResponseWriter, r *http.Request) {
	var payload registerCheckPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	alumni, err := h.AuthService.RegisterCheck(payload.NISN)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	_ = utils.WriteJSON(w, http.StatusOK, alumni)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var payload registerPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err := h.AuthService.Register(service.RegisterPayload{
		Username: payload.Username,
		Email:    payload.Email,
		Password: payload.Password,
		AlumniID: payload.AlumniID,
	})
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Register success"}
	_ = utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h *Handler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var payload authPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	tokens, err := h.AuthService.Authenticate(payload.Email, payload.Password)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	http.SetCookie(w, h.Auth.GetRefreshCookie(tokens.RefreshToken))
	_ = utils.WriteJSON(w, http.StatusOK, tokens)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == h.Auth.CookieName {
			tokenPairs, err := h.AuthService.RefreshToken(cookie.Value)
			if err != nil {
				_ = utils.ErrorJSON(w, err, http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, h.Auth.GetRefreshCookie(tokenPairs.RefreshToken))
			_ = utils.WriteJSON(w, http.StatusOK, tokenPairs)
			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, h.Auth.GetExpiredRefreshCookie())
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var payload forgotPasswordPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	// Always return generic success message for security
	// Do not reveal whether email exists
	_ = h.AuthService.ForgotPassword(payload.Email)

	resp := utils.JSONResponse{Error: false, Message: "If the email exists, a reset link has been sent."}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var payload resetPasswordPayload
	if err := utils.ReadJSON(w, r, &payload); err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	err := h.AuthService.ResetPassword(payload.Token, payload.Password)
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}

	resp := utils.JSONResponse{Error: false, Message: "Password has been reset successfully"}
	_ = utils.WriteJSON(w, http.StatusOK, resp)
}
