package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, noSurf, app.authenticate)

	//mux := http.NewServeMux()
	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))

	mux.Get("/blog/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createBlogForm))
	mux.Post("/blog/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createBlog))
	mux.Get("/blog/:id", dynamicMiddleware.ThenFunc(app.showBlog))

	mux.Get("/user/register", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/register", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	//return app.recoverPanic(app.logRequest(secureHeaders(mux)))
	return standardMiddleware.Then(mux)
}
