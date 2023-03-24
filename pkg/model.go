package pkg

type RegisterData struct {
	Transport    string `json:"transport"`
	Login        string `json:"login"`
	EncryptedKey string `json:"encryptedKey"`
	PublicKey    string `json:"publicKey"`
	Hash         string `json:"hash"`
}
