package transport

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"user/pkg/app/model"
	"user/pkg/app/query"
	"user/pkg/app/service"
	"user/pkg/app/session"
)

const (
	sessionName = "session_id"
)

type Handler struct {
	userService      service.UserService
	userQueryService query.UserQueryService
	userSession      session.UserSession
	cookieStore      *sessions.CookieStore
}

func NewHandler(
	userService service.UserService,
	userQueryService query.UserQueryService,
	userSession session.UserSession,
	cookieStore *sessions.CookieStore,
) *Handler {
	return &Handler{
		userService:      userService,
		userQueryService: userQueryService,
		userSession:      userSession,
		cookieStore:      cookieStore,
	}
}

func (h *Handler) GetSignInPage(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./data/html/signin.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")
	if login == "" || pass == "" {
		http.Error(w, "login и password обязательны", http.StatusBadRequest)
		return
	}

	h.authenticate(w, r, login, pass)
}

func (h *Handler) GetSignUpPage(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("./data/html/signup.html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешён", http.StatusMethodNotAllowed)
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")
	if login == "" || pass == "" {
		http.Error(w, "login и password обязательны", http.StatusBadRequest)
		return
	}

	err := h.userService.Create(login, pass)
	if errors.Is(err, service.ErrUserAlreadyExists) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	h.authenticate(w, r, login, pass)
}

func (h *Handler) SignOut(w http.ResponseWriter, r *http.Request) {
	sess, err := h.cookieStore.Get(r, sessionName)
	if err != nil {
		http.Error(w, "Сессии недоступны", http.StatusInternalServerError)
		return
	}
	sess.Options.MaxAge = -1
	if err := sess.Save(r, w); err != nil {
		http.Error(w, "Не удалось завершить сессию", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	http.Redirect(w, r, "/user/signin", http.StatusSeeOther)
}

type AuthCheckResponse struct {
	Authenticated bool         `json:"authenticated"`
	UserID        model.UserID `json:"user_id,omitempty"`
}

func (h *Handler) CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := h.cookieStore.Get(r, sessionName)
	if err != nil {
		log.Printf("Session store error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	auth, ok := sess.Values["authenticated"].(bool)
	if !ok || !auth {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthCheckResponse{Authenticated: false})
		return
	}

	uidStr, ok := sess.Values["user_id"].(string)
	if !ok {
		sess.Options.MaxAge = -1
		_ = sess.Save(r, w)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthCheckResponse{Authenticated: false})
		return
	}

	user, err := h.userQueryService.FindByID(uidStr)
	if err != nil {
		log.Printf("UserQueryService error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if errors.Is(err, query.ErrUserNotFound) {
		sess.Options.MaxAge = -1
		_ = sess.Save(r, w)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AuthCheckResponse{Authenticated: false})
		return
	}
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthCheckResponse{
		Authenticated: true,
		UserID:        user.UserID,
	})
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request, login, password string) {
	user, err := h.userSession.Identify(login, password)
	if errors.Is(err, session.ErrInvalidCredentials) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	sess, err := h.cookieStore.Get(r, sessionName)
	if err != nil {
		log.Printf("Ошибка получения сессии: %v", err)
		http.Error(w, "Сессии недоступны", http.StatusInternalServerError)
		return
	}
	sess.Values["authenticated"] = true
	sess.Values["user_id"] = string(user.UserID)
	_ = sess.Save(r, w)

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"redirect":"/"}`))
}
