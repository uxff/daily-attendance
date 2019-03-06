package oncetoken

import (
	"math/rand"
	"time"
	"fmt"
	"sync"
	"github.com/astaxie/beego/logs"
)

type tokenVal struct {
	To *time.Timer
}
var tokens map[string]tokenVal

func init() {
	tokens = make(map[string]tokenVal, 0)
	rand.Seed(time.Now().UnixNano())
}

func GenToken() string {
	m := &sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	for {

		token := fmt.Sprintf("%d", rand.Int63())
		if _, exist := tokens[token]; exist {
			// 重复存在
			continue
		}

		to := time.AfterFunc(time.Second*300, func(){
			logs.Info("the oncetoken [%s] will be deleted", token)
			DeleteToken(token)
		})
		tokens[token] = tokenVal{To:to}

		return token
	}

	return ""
}

// verify ok will delete the token
func VerifyToken(t string) bool {
	m := &sync.Mutex{}
	m.Lock()
	defer m.Unlock()

	if tv, exist := tokens[t]; exist {
		tv.To.Stop()
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
