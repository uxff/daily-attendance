package oncetoken

import (
	"math/rand"
	"time"
	"fmt"
	"sync"
	"github.com/astaxie/beego/logs"
)

var tokens map[string]bool

func init() {
	tokens = make(map[string]bool, 0)
	rand.Seed(time.Now().UnixNano())
}

func GenToken() string {
	m := &sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	for {

		token := fmt.Sprintf("%d", rand.Int63())
		if _, exist := tokens[token]; exist {
			continue
		}

		tokens[token] = true
		time.AfterFunc(time.Second*300, func(){
			logs.Info("the oncetoken [%s] will be deleted", token)
			DeleteToken(token)
		})
		return token
	}

	return ""
}

// verify ok will delete the token
func VerifyToken(t string) bool {
	m := &sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	if _, exist := tokens[t]; exist {
		delete(tokens, t)
		return true
	}

	return false
}

func DeleteToken(t string) {
	m := &sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	delete(tokens, t)
}
