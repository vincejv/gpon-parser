package util

import (
	"sync"

	"github.com/vincejv/gpon-parser/model"
	"golang.org/x/net/html"
)

type DocPage struct {
	sync.RWMutex
	doc    *html.Node
	docStr string
}

func (docPage *DocPage) GetPage() *html.Node {
	docPage.RLock()
	defer docPage.RUnlock()
	return docPage.doc
}

func (docPage *DocPage) SetPage(doc *html.Node) {
	docPage.Lock()
	docPage.doc = doc
	docPage.Unlock()
}

func (docPage *DocPage) GetStrPage() string {
	docPage.RLock()
	defer docPage.RUnlock()
	return docPage.docStr
}

func (docPage *DocPage) SetStrPage(docStr string) {
	docPage.Lock()
	docPage.docStr = docStr
	docPage.Unlock()
}

type LoginCreds struct {
	sync.RWMutex
	username string
	password string
}

func (loginCreds *LoginCreds) GetCreds() (string, string) {
	loginCreds.RLock()
	defer loginCreds.RUnlock()
	return loginCreds.username, loginCreds.password
}

func (loginCreds *LoginCreds) SetCreds(username string, password string) {
	loginCreds.Lock()
	loginCreds.username = username
	loginCreds.password = password
	loginCreds.Unlock()
}

type GlobalFlag struct {
	sync.RWMutex
	flag bool
}

func (flag *GlobalFlag) GetFlag() bool {
	flag.RLock()
	defer flag.RUnlock()
	return flag.flag
}

func (flag *GlobalFlag) SetFlag(flagVal bool) {
	flag.Lock()
	flag.flag = flagVal
	flag.Unlock()
}

type CachedStat struct {
	sync.RWMutex
	stat *model.DeviceStats
}

func (stat *CachedStat) GetStat() *model.DeviceStats {
	stat.RLock()
	defer stat.RUnlock()
	return stat.stat
}

func (stat *CachedStat) SetStat(statVal *model.DeviceStats) {
	stat.Lock()
	stat.stat = statVal
	stat.Unlock()
}
