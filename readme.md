# Keys memorise service
save encrypted keys

## API Keys
### Register new user
/api/register 

params:
```go
type RegisterData struct {
	Transport    string `json:"transport"`
	Login        string `json:"login"`
	EncryptedKey string `json:"encryptedKey"`
	PublicKey    string `json:"publicKey"`
	Hash         string `json:"hash"`
}
```

### Verify user
/api/verify 

params:
session - uuid for session (value: RegisterData)

### Get key
/api/key 

### Check service health
/api/ok 


## Google Auth

### Get auth token
/api/google/auth
- check auth token
- check register
- register if not exist
- return auth token


