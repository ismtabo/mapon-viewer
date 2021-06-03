package service

import (
	"net/http"

	"github.com/kataras/go-sessions/v3"
)

const authAttribute = "authenticated"

// SessionsService implements sessions management.
type SessionsService interface {
	Login(rw http.ResponseWriter, r *http.Request)
	Logout(rw http.ResponseWriter, r *http.Request)
	IsAuthenticated(rw http.ResponseWriter, r *http.Request) bool
}

type sessionsSecurityService struct {
	ssns *sessions.Sessions
}

// NewSessionsService creates a new instance of SessionsService.
func NewSessionsService(ssns *sessions.Sessions) SessionsService {
	return &sessionsSecurityService{ssns: ssns}
}

// Login initialize an application session in the given HTTP response.
func (s sessionsSecurityService) Login(rw http.ResponseWriter, r *http.Request) {
	s.ssns.Start(rw, r).Set(authAttribute, true)
}

// Logout finalize an application session in the given HTTP response.
func (s sessionsSecurityService) Logout(rw http.ResponseWriter, r *http.Request) {
	s.ssns.Destroy(rw, r)
}

// IsAuthenticated validates that an application session exists in the given HTTP response.
func (s sessionsSecurityService) IsAuthenticated(rw http.ResponseWriter, r *http.Request) bool {
	return s.ssns.Start(rw, r).GetBooleanDefault(authAttribute, false)
}
