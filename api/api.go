package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gitlab.com/rayone121/rayman/pacman"
	"net/http"
)

const authTokenLen = 10

type API struct {
	router *mux.Router
	token  string
	tasker *Tasker
}

func New() (*API, error) {
	var api API
	tokenBytes := make([]byte, authTokenLen)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return nil, err
	}

	api.token = base64.RawURLEncoding.EncodeToString(tokenBytes)

	api.tasker = NewTasker()

	return &api, nil
}

func (a *API) currentOperation(w http.ResponseWriter, _ *http.Request) {
	writeJson(w, a.tasker.GetCurrent())
}

func (a *API) search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	searchOp := pacman.NewSearchOperation(query.Get("keyword"), pacman.SearchField(query.Get("sort-by")))
	op := Operation{Operation: searchOp, ID: -1}
	op.Packages, op.Err = op.Operation.Execute()
	writeJson(w, op)
}

func (a *API) remove(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	packages := query["pkg"]

	a.schedule(w, pacman.NewRemoveOperation(packages, false))
}

func (a *API) install(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	packages := query["pkg"]

	a.schedule(w, pacman.NewInstallOperation(packages, false))
}

func (a *API) schedule(w http.ResponseWriter, op pacman.Operation) {
	id := a.tasker.Schedule(op)
	_, err := w.Write([]byte(fmt.Sprintf("{\"ID\":%d}", id)))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *API) completed(w http.ResponseWriter, _ *http.Request) {
	completed := a.tasker.GetCompleted()
	writeJson(w, completed)
}

func writeJson[T any](w http.ResponseWriter, anything T) {
	w.Header().Set("Content-Type", "application/json")
	bytes, err := json.Marshal(anything)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(bytes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (a *API) Listen() error {

	fmt.Printf("AUTH TOKEN: %s\n", a.token)

	a.router = mux.NewRouter()
	v1 := a.router.PathPrefix("/api/v1").Subrouter()
	v1.HandleFunc("/install", a.install)
	v1.HandleFunc("/remove", a.remove)
	v1.HandleFunc("/search", a.search)
	v1.HandleFunc("/current", a.currentOperation)
	v1.HandleFunc("/completed", a.completed)

	err := http.ListenAndServe(":8042", a)
	if err != nil {
		return err
	}
	return nil
}

func (a *API) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if request.URL.Query().Get("auth") != a.token {
		responseWriter.WriteHeader(http.StatusUnauthorized)
		_, err := responseWriter.Write([]byte("401 Unauthorized: Invalid Auth Token"))
		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
			return
		}
		return
	}

	a.router.ServeHTTP(responseWriter, request)
}
