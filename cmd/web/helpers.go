package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serveError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)
	http.Error(w, trace, http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	ts, ok := app.templateCache[name]
	if !ok {
		app.serveError(w, fmt.Errorf("the template %s does not exist", name))
		return
	}

	// Write to an in-memory buffer (instead of directly to `w`) so to
	// prevent incomplete HTML outputs in case of errors.
	buf := new(bytes.Buffer)

	err := ts.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serveError(w, err)
	}

	buf.WriteTo(w)
}

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}

	td.CurrentYear = time.Now().Year()
	td.Flash = app.session.PopString(r, "flash")
	td.IsAuthenticated = app.isAuthenticated(r)
	return td
}

func (app *application) isAuthenticated(r *http.Request) bool {
	return app.session.Exists(r, "authenticatedUserID")
}
