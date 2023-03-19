package pkg

type RegisterData struct {
	Transport    string `json:"transport"`
	Login        string `json:"login"`
	EncryptedKey string `json:"ekey"`
	PublicKey    string `json:"pkey"`
	Hash         string `json:"hash"`
}
