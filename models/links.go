package models

import (
	"os"
	"time"
	"math/rand"
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego/logs"
)

type Link struct{
	Name string `json:"name"`
	Url  string `json:"url"`
}

type IndexLinks []struct{
	Name string `json:"name"`
	Links []Link	`json:"links"`
}

type FriendlyLinks []Link

var
(
	// index links
	theIndexLinks IndexLinks

	// index links path
	theIndexLinksPath string

	// friendly links
	theFriendlyLinks FriendlyLinks

	// friendly links path
	theFriendlyLinksPath string
)


func GetIndexLinks() IndexLinks {
	return theIndexLinks
}

func SetLinksPath(p string) {
	theIndexLinksPath = p
	//LoadIndexLinksFromFile(theLinksPath)
}

func LoadIndexLinks() IndexLinks {
	return LoadIndexLinksFromFile(theIndexLinksPath)
}

func LoadIndexLinksFromFile(f string) IndexLinks {
	fhandle, err := os.Open(f)
	if err != nil {
		logs.Error("load index links from %s error:%v", f, err)
		return nil
	}

	defer fhandle.Close()

	content, err := ioutil.ReadAll(fhandle)
	if err != nil {
		logs.Error("load index links from %s error:%v", f, err)
		return nil
	}

	err = json.Unmarshal(content, &theIndexLinks)
	if err != nil {
		logs.Error("load index links from %s error:%v", f, err)
		return nil
	}

	logs.Info("load index links from file %s ok", f)

	return theIndexLinks
}

func GetFriendlyLinks() FriendlyLinks {
	return theFriendlyLinks
}

func SetFriendlyLinksPath(p string) {
	theFriendlyLinksPath = p
}

func LoadFriendlyLinks() FriendlyLinks {
	return LoadFriendlyLinksFromFile(theFriendlyLinksPath)
}

func LoadFriendlyLinksFromFile(f string) FriendlyLinks {
	fhandle, err := os.Open(f)
	if err != nil {
		logs.Error("load friendly links from %s error:%v", f, err)
		return nil
	}

	defer fhandle.Close()

	content, err := ioutil.ReadAll(fhandle)
	if err != nil {
		logs.Error("load friendly links from %s error:%v", f, err)
		return nil
	}

	err = json.Unmarshal(content, &theFriendlyLinks)
	if err != nil {
		logs.Error("load friendly links from %s error:%v", f, err)
		return nil
	}

	logs.Info("load friendly links from file %s ok", f)

	return theFriendlyLinks
}

func ShuffleLinks(links FriendlyLinks) FriendlyLinks {
	thelen := len(links)
	targetLinks := make(FriendlyLinks, 0, thelen)
	roundNum := time.Now().Unix()
	roundStart := rand.Int()%thelen

	switch true {
	case roundNum&1 == 0:
		// 正序
		for i := 0; i<thelen; i++ {
			targetLinks = append(targetLinks, links[(i+roundStart)%thelen])
		}
	case roundNum&1 == 1:
		// 倒叙
		for i := 0; i<thelen; i++ {
			targetLinks = append(targetLinks, links[(-i+roundStart+thelen)%thelen])
		}
	}

	return targetLinks
}


