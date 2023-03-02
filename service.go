package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"k8s.io/klog/v2"
	"net/http"
)

type RestServer struct {
	sessions   *cache.Cache
	db         *DriveDb
	successUrl string
	errorUrl   string
	mailServer *MailServ
}

func (s *RestServer) runServer() {
	http.HandleFunc("/register", s.register)
	http.HandleFunc("/verify", s.verify)
	http.HandleFunc("/key", s.key)
	klog.Fatal(http.ListenAndServe(listenAddress, nil))
}

func (s *RestServer) verify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	session := r.URL.Query().Get("session")
	stored, b := s.sessions.Get(session)
	if !b {
		http.Redirect(w, r, s.errorUrl, http.StatusSeeOther)
		return
	} else {
		register := stored.(RegisterData)
		s.db.saveRegister(register)
		http.Redirect(w, r, s.successUrl, http.StatusSeeOther)
		return
	}
}

func (s *RestServer) key(w http.ResponseWriter, r *http.Request) { // check get params
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// get uuid
	mail := r.URL.Query().Get("mail")
	hash := r.URL.Query().Get("hash")

	data, err := s.db.loadRegister(mail)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if hash != data.PasswordHash {
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
	var register RegisterData
	err = json.Unmarshal(body, &register)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	session := uuid.New()
	sessions.Add(session.String(), register, cache.DefaultExpiration)
	s.mailServer.sendMail(register, session)

	klog.Info("Register: ", register.Email, " session: ", session.String())

	// Return a response indicating success.
	json.NewEncoder(w).Encode(sessions)
}
