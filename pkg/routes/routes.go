package routes

import (
	"net/http"

	"github.com/ismtabo/mapon-viewer/pkg/controller"
	"github.com/ismtabo/mapon-viewer/pkg/routes/mw"
	"github.com/ismtabo/mapon-viewer/pkg/service"
	"github.com/rs/zerolog"
)

// Routes implements an application HTTP router.
type Routes interface {
	AddRoutes()
	ServeHTTP(http.ResponseWriter, *http.Request)
}

type routes struct {
	maponCtrl controller.MaponController
	pagesCtrl controller.PagesController
	userCtrl  controller.UserController
	secSvc    service.SessionsService
	log       *zerolog.Logger
	mux       *http.ServeMux
}

// NewRoutes creates a new instance of Routes.
func NewRoutes(userCtrl controller.UserController, pagesCtrl controller.PagesController, maponCtrl controller.MaponController, secSvc service.SessionsService, log *zerolog.Logger) Routes {
	return &routes{maponCtrl: maponCtrl, pagesCtrl: pagesCtrl, userCtrl: userCtrl, secSvc: secSvc, log: log}
}

// AddRoutes initialize internal application routes handlers with HTTP middlewares.
func (r *routes) AddRoutes() {
	r.mux = http.NewServeMux()
	r.mux.HandleFunc("/static/", http.HandlerFunc(r.pagesCtrl.StaticFiles))
	r.mux.Handle("/", r.addMiddlewares([]string{"GET"}, "index-page", http.HandlerFunc(r.pagesCtrl.IndexPage)))
	r.mux.Handle("/login", r.addMiddlewares([]string{"GET"}, "login-page", http.HandlerFunc(r.pagesCtrl.LoginPage)))
	r.mux.Handle("/auth/login", r.addMiddlewares([]string{"POST"}, "login", http.HandlerFunc(r.userCtrl.LoginUser)))
	r.mux.Handle("/auth/logout", r.addMiddlewares([]string{"POST"}, "logout", http.HandlerFunc(r.userCtrl.LogoutUser)))
	r.mux.Handle("/auth/register", r.addMiddlewares([]string{"POST"}, "register", http.HandlerFunc(r.userCtrl.RegisterUser)))
	r.mux.Handle("/api/mapon", mw.SecurityHandler(r.secSvc)(r.addMiddlewares([]string{"GET"}, "get-mapon", http.HandlerFunc(r.maponCtrl.GetMaponInfo))))
}

func (r *routes) addMiddlewares(methods []string, op string, handleFunc http.HandlerFunc) http.Handler {
	return mw.MethodsHandler(methods...)(
		mw.InitAppCtxHandler(
			mw.InitLogCtxHandler(r.log)(
				mw.CorrelatorHandler(
					mw.LogHTTPHandler(
						mw.LogContextHandler(op)(handleFunc),
					),
				),
			),
		),
	)
}

// ServeHTTP implements http.Handler interface.
func (r *routes) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(rw, req)
}
