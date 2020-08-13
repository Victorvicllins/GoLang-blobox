package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Victorvicllins/pkg/forms"
	"github.com/Victorvicllins/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// if r.URL.Path != "/" {
	// 	app.notFound(w)
	// 	return
	// }

	s, err := app.blog.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		BlogList: s,
	})
	// for _, blog := range s {
	// 	fmt.Fprintf(w, "%v\n", blog)
	// }
	// data := &templateData{BlogList: s}
	//
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }
	//
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	//app.errorLog.Println(err.Error())
	// 	app.serverError(w, err)
	// 	// http.Error(w, "Internal Server Error", 500)
	// 	return
	// }
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	//app.errorLog.Println(err.Error())
	// 	//http.Error(w, "Internal Server Error", 500)
	// }
}

func (app *application) showBlog(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.blog.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}
	// flash := app.session.PopString(r, "flash")

	app.render(w, r, "show.page.tmpl", &templateData{
		Viccy: s,
	})
	// data := &templateData{Viccy: s}
	//
	// files := []string{
	// 	"./ui/html/show.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }
	// //fmt.Fprintf(w, "%v", s)
	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// err = ts.Execute(w, data)
	// if err != nil {
	// 	app.serverError(w, err)
	// }
}

func (app *application) createBlogForm(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("Create new blog.."))
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createBlog(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	w.Header().Set("Allow", http.MethodPost)
	// 	app.clientError(w, http.StatusMethodNotAllowed)
	// 	//http.Error(w, "Method Not Allowed", 405)
	// 	return
	// }

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	id, err := app.blog.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
	if err != nil {
		app.serverError(w, err)
		return
	}
	app.session.Put(r, "flash", "Blog item created successfully.!")
	// w.Write([]byte("Create new Blog..."))
	// ## Redirect user to the relevant page of the Blog
	http.Redirect(w, r, fmt.Sprintf("/blog/%d", id), http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Display signup/Registeration page/Form..")
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Create new user..")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)

	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.Errors.Add("email", "Email address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}

	// fmt.Fprintln(w, "Create a new user..")
	app.session.Put(r, "flash", "Your Sign-Up was successful. Please loggin..")
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "User login page/Form..")
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Authenticate user / Login user...")
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.Errors.Add("generic", "Email or Password incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, err)
		}
		return
	}
	app.session.Put(r, "authenticatedUserID", id)
	http.Redirect(w, r, "/blog/create", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Logout user.")
	app.session.Remove(r, "authenticatedUserID")
	app.session.Put(r, "flash", "You've been logged out successfully.!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
