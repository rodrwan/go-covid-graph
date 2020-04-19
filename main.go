package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/rodrwan/go-covid-graph/pkg/data"
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

type ErrorResponse struct {
	Code    int
	Message string
}

func homeHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			// Open the file
			cases, err := data.GetRegionalConfirmedCases()
			if err != nil {
				logger.Error(err)
				renderError(rw, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				return
			}

			t, err := template.ParseFiles("templates/layout.html", "templates/graph.html")
			if err != nil {
				logger.Error("Could not compile templates", err)
				http.Error(rw, "Internal server error", http.StatusInternalServerError)
				return
			}

			t.ExecuteTemplate(rw, "layout", cases[1:len(cases)-1])

			return
		}

		renderError(rw, http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
	}
}

func renderError(rw http.ResponseWriter, code int, message string) {
	logger.Errorf("%d - Error: %s", code, message)
	t, err := template.ParseFiles("templates/layout.html", "templates/error.html")
	if err != nil {
		logger.Error("Could not compile templates", err)
		http.Error(rw, "Internal server error", http.StatusInternalServerError)
		return
	}

	e := ErrorResponse{
		Code:    code,
		Message: message,
	}

	t.ExecuteTemplate(rw, "layout", e)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler())

	// serve static files.
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	logger = logrus.New()
	w := logger.Writer()
	defer w.Close()

	srv := http.Server{
		Addr:     ":5000",
		Handler:  mux,
		ErrorLog: log.New(w, "", 0),
	}

	logger.Info("Server is ready to handle requests at http://localhost:5000")

	srv.ListenAndServe()
}
