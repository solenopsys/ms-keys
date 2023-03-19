package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
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
	prefix := "/api"
	http.HandleFunc(prefix+"/register", s.register)
	http.HandleFunc(prefix+"/verify", s.verify)
	http.HandleFunc(prefix+"/api", s.key)
	http.HandleFunc(prefix+"/ok", s.health)
	klog.Fatal(http.ListenAndServe(s.ListenAddress, nil))
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

func (s *RestServer) key(w http.ResponseWriter, r *http.Request) { // check get params
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	mail := r.URL.Query().Get("mail")
	hash := r.URL.Query().Get("hash")

	data, err := s.Db.LoadRegister(mail)

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
	json.NewEncoder(w).Encode(s.TransportService.List())
}
