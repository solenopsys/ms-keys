package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/rs/cors"
	"io/ioutil"
	"k8s.io/klog/v2"
	"ms-keys/pkg"
	"net/http"
)

type RestServer struct {
	Sessions         *cache.Cache
	Db               PersistedData
	TransportService *Service
	SuccessUrl       string
	ErrorUrl         string
	ListenAddress    string
}

func (s *RestServer) Run() {
	mux := http.NewServeMux()
	prefix := "/api"
	mux.HandleFunc(prefix+"/register", s.register)
	mux.HandleFunc(prefix+"/verify", s.verify)
	mux.HandleFunc(prefix+"/key", s.getKey)
	mux.HandleFunc(prefix+"/ok", s.health)
	handler := cors.Default().Handler(mux)

	klog.Fatal(http.ListenAndServe(s.ListenAddress, handler))
}

func (s *RestServer) verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	session := r.URL.Query().Get("session")
	stored, b := s.Sessions.Get(session)
	if !b {
		http.Redirect(w, r, s.ErrorUrl, http.StatusSeeOther)
		return
	} else {
		register := stored.(pkg.RegisterData)
		s.Db.SaveRegister(register)
		http.Redirect(w, r, s.SuccessUrl, http.StatusSeeOther)
		return
	}
}

func (s *RestServer) health(w http.ResponseWriter, r *http.Request) { // check get params
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("AUTH WORKS"))
}

func (s *RestServer) getKey(w http.ResponseWriter, r *http.Request) { // check get params
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hashBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hash := string(hashBytes)
	data, err := s.Db.LoadRegister(hash)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if hash != data.Hash {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		marshal, err := json.Marshal(data)
		if err == nil {
			w.Write([]byte(marshal))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
}

func (s *RestServer) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Parse the JSON payload into a User struct.
	var register pkg.RegisterData
	err = json.Unmarshal(body, &register)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session := uuid.New()
	s.Sessions.Add(session.String(), register, cache.DefaultExpiration)
	s.TransportService.Send(register.Transport, register, session)

	klog.Info("Register: ", register.Login, " session: ", session.String())

	w.WriteHeader(http.StatusOK)
	// Return a response indicating success.
	json.NewEncoder(w).Encode(s.Sessions)
}

// list transport
func (s *RestServer) listTransport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Return a response indicating success.
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(s.TransportService.List())
}
