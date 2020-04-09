package usecase

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"no_homomorphism/internal/pkg/csrf"
	"time"
)

//expireTime in seconds
type CryptToken struct {
	Secret     []byte
	ExpireTime int64
	BlackList  csrf.Repository
}

type TokenData struct {
	SessionID string
	TimeStamp int64
}

func NewAesCryptHashToken(secret string, expireTime int64, blackList csrf.Repository) (CryptToken, error) {
	key := []byte(secret)
	_, err := aes.NewCipher(key)
	if err != nil {
		return CryptToken{}, fmt.Errorf("cipher problem %v", err)
	}
	return CryptToken{Secret: key, ExpireTime: expireTime, BlackList: blackList}, nil
}

func (tk *CryptToken) Create(sid string, timeStamp int64) (string, error) {
	block, err := aes.NewCipher(tk.Secret)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	td := &TokenData{SessionID: sid, TimeStamp: timeStamp}
	data, _ := json.Marshal(td)
	ciphertext := aesgcm.Seal(nil, nonce, data, nil)

	res := append([]byte(nil), nonce...)
	res = append(res, ciphertext...)

	token := base64.StdEncoding.EncodeToString(res)

	return token, nil
}

func (tk *CryptToken) Check(sid string, inputToken string) (bool, error) {
	block, err := aes.NewCipher(tk.Secret)
	if err != nil {
		return false, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return false, err
	}
	ciphertext, err := base64.StdEncoding.DecodeString(inputToken)
	if err != nil {
		return false, err
	}
	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return false, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return false, fmt.Errorf("decrypt fail: %v", err)
	}

	td := TokenData{}
	err = json.Unmarshal(plaintext, &td)
	if err != nil {
		return false, fmt.Errorf("bad json: %v", err)
	}

	if time.Now().Unix()-td.TimeStamp > tk.ExpireTime {
		return false, fmt.Errorf("token expired")
	}

	expected := TokenData{SessionID: sid, TimeStamp: td.TimeStamp}

	err = tk.BlackList.Check(inputToken)

	if td != expected || err != nil {
		return false, nil
	}

	err = tk.BlackList.Add(inputToken, tk.ExpireTime)
	if err != nil {
		return false, err
	}
	return true, nil
}
