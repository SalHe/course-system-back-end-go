package token

import (
	"encoding/json"
	"io"
	"os"
	"sync"
)

type InMemoryTokenStorage struct {
	tokenJwt map[string]string
	mux      *sync.Mutex
}

func NewInMemoryTokenStorage() *InMemoryTokenStorage {
	return &InMemoryTokenStorage{
		tokenJwt: make(map[string]string),
		mux:      &sync.Mutex{},
	}
}

func (i *InMemoryTokenStorage) Get(token string) (string, bool) {
	i.mux.Lock()
	defer i.mux.Unlock()
	s, ok := i.tokenJwt[token]
	return s, ok
}

func (i *InMemoryTokenStorage) Set(token string, jwt string) {
	i.mux.Lock()
	defer i.mux.Unlock()

	i.tokenJwt[token] = jwt
}

func (i InMemoryTokenStorage) SetExpire(token string, expire int64) {
	i.mux.Lock()
	defer i.mux.Unlock()

	// NOT SUPPORT
}

func (i *InMemoryTokenStorage) Delete(token string) {
	i.mux.Lock()
	defer i.mux.Unlock()
	delete(i.tokenJwt, token)
}

// Load 从文件加载token，当 overwrite 为 ture 时会覆盖已有的 token。
func (i *InMemoryTokenStorage) Load(file string, overwrite bool) {
	i.mux.Lock()
	defer i.mux.Unlock()
	f, err := os.OpenFile(file, os.O_RDONLY, 0666)
	if err == nil {
		defer f.Close()
		bytes, _ := io.ReadAll(f)
		tokensInFile := make(map[string]string)
		json.Unmarshal(bytes, &tokensInFile)
		for k, v := range tokensInFile {
			if _, ok := i.tokenJwt[k]; overwrite || !ok {
				i.tokenJwt[k] = v
			}
		}
		return
	}
}

// Save 将 token 写入文件。
func (i *InMemoryTokenStorage) Save(file string) {
	i.mux.Lock()
	defer i.mux.Unlock()
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err == nil {
		defer f.Close()
		bytes, _ := json.Marshal(i.tokenJwt)
		f.Write(bytes)
	}
}
