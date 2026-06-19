package httpapi

import (
	"html/template"
	"net/http"
	"strings"

	"example.com/pz6-web-security/internal/auth"
	"example.com/pz6-web-security/internal/store"
)

type Handler struct {
	store        *store.Store
	secureCookie bool
	templates    *template.Template
}

func NewHandler(s *store.Store, secureCookie bool) (*Handler, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}
	return &Handler{
		store:        s,
		secureCookie: secureCookie,
		templates:    tmpl,
	}, nil
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /login", h.login)
	mux.HandleFunc("GET /logout", h.logout)
	mux.HandleFunc("GET /profile", h.profileGet)
	mux.HandleFunc("POST /profile", h.profilePost)
	mux.HandleFunc("GET /hello", h.hello)
	mux.HandleFunc("GET /comments", h.commentsGet)
	mux.HandleFunc("POST /comments", h.commentsPost)
	mux.HandleFunc("GET /demo/unsafe", h.unsafeHello)
}

func (h *Handler) login(w http.ResponseWriter, r *http.Request) {
	sessionID, err := auth.RandomToken(16)
	if err != nil {
		http.Error(w, "failed to create session", http.StatusInternalServerError)
		return
	}
	csrfToken, err := auth.RandomToken(16)
	if err != nil {
		http.Error(w, "failed to create csrf token", http.StatusInternalServerError)
		return
	}

	h.store.Save(&store.UserProfile{
		SessionID: sessionID,
		Name:      "Student",
		CSRFToken: csrfToken,
		Comments:  nil,
	})
	auth.SetSessionCookie(w, sessionID, h.secureCookie)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *Handler) logout(w http.ResponseWriter, r *http.Request) {
	sessionID, err := auth.ReadSessionCookie(r)
	if err == nil {
		h.store.Delete(sessionID)
	}
	auth.ClearSessionCookie(w, h.secureCookie)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h *Handler) profileGet(w http.ResponseWriter, r *http.Request) {
	profile, ok := h.currentProfile(w, r)
	if !ok {
		return
	}
	h.render(w, "profile.html", profile)
}

func (h *Handler) profilePost(w http.ResponseWriter, r *http.Request) {
	profile, ok := h.currentProfile(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if r.FormValue("csrf_token") != profile.CSRFToken {
		http.Error(w, "invalid csrf token", http.StatusForbidden)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if !h.store.UpdateName(profile.SessionID, name) {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	newToken, err := auth.RandomToken(16)
	if err != nil {
		http.Error(w, "failed to rotate csrf token", http.StatusInternalServerError)
		return
	}
	if !h.store.UpdateCSRF(profile.SessionID, newToken) {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

func (h *Handler) hello(w http.ResponseWriter, r *http.Request) {
	profile, ok := h.currentProfile(w, r)
	if !ok {
		return
	}
	h.render(w, "hello.html", profile)
}

func (h *Handler) commentsGet(w http.ResponseWriter, r *http.Request) {
	profile, ok := h.currentProfile(w, r)
	if !ok {
		return
	}
	h.render(w, "comments.html", profile)
}

func (h *Handler) commentsPost(w http.ResponseWriter, r *http.Request) {
	profile, ok := h.currentProfile(w, r)
	if !ok {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if r.FormValue("csrf_token") != profile.CSRFToken {
		http.Error(w, "invalid csrf token", http.StatusForbidden)
		return
	}

	text := strings.TrimSpace(r.FormValue("text"))
	if text == "" {
		http.Error(w, "comment text is required", http.StatusBadRequest)
		return
	}
	if !h.store.AddComment(profile.SessionID, text) {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	newToken, err := auth.RandomToken(16)
	if err != nil {
		http.Error(w, "failed to rotate csrf token", http.StatusInternalServerError)
		return
	}
	if !h.store.UpdateCSRF(profile.SessionID, newToken) {
		http.Error(w, "session not found", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/comments", http.StatusSeeOther)
}

func (h *Handler) unsafeHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}
	// Намеренно небезопасный вывод для демонстрации XSS (только учебный пример).
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte("<html><body><h1>Hello, " + name + "!</h1></body></html>"))
}

func (h *Handler) currentProfile(w http.ResponseWriter, r *http.Request) (*store.UserProfile, bool) {
	sessionID, err := auth.ReadSessionCookie(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, false
	}
	profile, ok := h.store.Get(sessionID)
	if !ok {
		auth.ClearSessionCookie(w, h.secureCookie)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return nil, false
	}
	return profile, true
}

func (h *Handler) render(w http.ResponseWriter, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}
