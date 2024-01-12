package main

import (
	"errors"
	"fmt"
	"gggrents/golangproject/pkg/forms"
	"gggrents/golangproject/pkg/models"
	"net/http"
	"strconv"
)
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// Because Pat matches the "/" path exactly, we can now remove the manual check
	// Because Pat matches the "/" path exactly, we can now remove the manual check
	// of r.URL.Path != "/" from this handler.
	s, err := app.snippets.Latest()
	if err != nil {
	app.serverError(w, err)
	return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
	Snippets: s,
	})
	}
	func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.URL.Query().Get(":id"))
		if err != nil || id < 1 {
		app.notFound(w)
		return
		}
		s, err := app.snippets.Get(id)
		if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		} else {
		app.serverError(w, err)
		}
		return
		}
		// Use the PopString() method to retrieve the value for the "flash" key.
		// PopString() also deletes the key and value from the session data, so it
		// acts like a one-time fetch. If there is no matching key in the session
		// data this will return the empty string.
		flash := app.session.PopString(r, "flash")
		// Pass the flash message to the template.
		app.render(w, r, "show.page.tmpl", &templateData{
		Flash: flash,
		Snippet: s,
		})
		}
	// Add a new createSnippetForm handler, which for now returns a placeholder response.
	func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
		app.render(w, r, "create.page.tmpl", &templateData{
		// Pass a new empty forms.Form object to the template.
		Form: forms.New(nil),
		})
		}
		func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
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
			id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))
			if err != nil {
			app.serverError(w, err)
			return
			}
			// Use the Put() method to add a string value ("Your snippet was saved
			// successfully!") and the corresponding key ("flash") to the session
			// data. Note that if there's no existing session for the current user
			// (or their session has expired) then a new, empty, session for them
			// will automatically be created by the session middleware.
			app.session.Put(r, "flash", "Snippet successfully created!")
			http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
			}

			func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
				app.render(w, r, "signup.page.tmpl", &templateData{
				Form: forms.New(nil),
				})
				}
				func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
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
					// Try to create a new user record in the database. If the email already exists
					// add an error message to the form and re-display it.
					err = app.users.Insert(form.Get("name"), form.Get("email"), form.Get("password"))
					if err != nil {
					if errors.Is(err, models.ErrDuplicateEmail) {
					form.Errors.Add("email", "Address is already in use")
					app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
					} else {
					app.serverError(w, err)
					}
					return
					}
					// Otherwise add a confirmation flash message to the session confirming that
					// their signup worked and asking them to log in.
					app.session.Put(r, "flash", "Your signup was successful. Please log in.")
					// And redirect the user to the login page.
					http.Redirect(w, r, "/user/login", http.StatusSeeOther)
					}
					
					func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
						app.render(w, r, "login.page.tmpl", &templateData{
						Form: forms.New(nil),
						})
						}
						
						func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
							err := r.ParseForm()
							if err != nil {
							app.clientError(w, http.StatusBadRequest)
							return
							}
							// Check whether the credentials are valid. If they're not, add a generic error
							// message to the form failures map and re-display the login page.
							form := forms.New(r.PostForm)
							id, err := app.users.Authenticate(form.Get("email"), form.Get("password"))
							if err != nil {
							if errors.Is(err, models.ErrInvalidCredentials) {
							form.Errors.Add("generic", "Email or Password is incorrect")
							app.render(w, r, "login.page.tmpl", &templateData{Form: form})
							} else {
							app.serverError(w, err)
							}
							return
							}
							// Add the ID of the current user to the session, so that they are now 'logged
							// in'.
							app.session.Put(r, "authenticatedUserID", id)
							// Redirect the user to the create snippet page.
							http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
							}
							
							func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
								// Remove the authenticatedUserID from the session data so that the user is
								// 'logged out'.
								app.session.Remove(r, "authenticatedUserID")
								// Add a flash message to the session to confirm to the user that they've been
								// logged out.
								app.session.Put(r, "flash", "You've been logged out successfully!")
								http.Redirect(w, r, "/", http.StatusSeeOther)
								}

								func ping(w http.ResponseWriter, r *http.Request) {
									w.Write([]byte("OK"))
									}