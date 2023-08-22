package main

import (
	"sync"

	"golang.org/x/net/html"
)

type DocPage struct {
	sync.RWMutex
	doc *html.Node
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
