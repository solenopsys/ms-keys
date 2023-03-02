package pkg

type RegisterData struct {
	Email        string `json:"email"`
	EncryptedKey string `json:"ekey"`
	PublicKey    string `json:"pkey"`
	PasswordHash string `json:"hash"`
}
