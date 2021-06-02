package service

import (
	"net/http"

	"github.com/kataras/go-sessions/v3"
)

const authAttribute = "authenticated"

type SecurityService interface {
	Login(rw http.ResponseWriter, r *http.Request)
	Logout(rw http.ResponseWriter, r *http.Request)
	IsAuthenticated(rw http.ResponseWriter, r *http.Request) bool
}

type sessionsSecurityService struct {
	ssns *sessions.Sessions
}

func NewSecurityService(ssns *sessions.Sessions) SecurityService {
	return &sessionsSecurityService{ssns: ssns}
}

func (s sessionsSecurityService) Login(rw http.ResponseWriter, r *http.Request) {
	s.ssns.Start(rw, r).Set(authAttribute, true)
}

func (s sessionsSecurityService) Logout(rw http.ResponseWriter, r *http.Request) {
	s.ssns.Destroy(rw, r)
}

func (s sessionsSecurityService) IsAuthenticated(rw http.ResponseWriter, r *http.Request) bool {
	return s.ssns.Start(rw, r).GetBooleanDefault(authAttribute, false)
}
