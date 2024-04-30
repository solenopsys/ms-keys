# Keys memorise service
save encrypted keys

## API
### Register new user
/api/register 

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

### Get key
/api/key 

### Check service health
/api/ok 



