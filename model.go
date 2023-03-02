package main

type RegisterData struct {
	Email        string `json:"email"`
	EncryptedKey string `json:"ekey"`
	PublicKey    string `json:"pkey"`
	PasswordHash string `json:"hash"`
}

type KeyRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"hash"`
}
