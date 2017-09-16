package model

import (
	"encoding/json"
	"log"
)

type Document struct {
	DocId      string `json:"doc_id"`
	Simhash    uint64 `json:"simhash"`
	ExpireTime int64 `json:"expire_time"`
}

func (doc *Document) Doc2Json() (string, error) {
	ret, err := json.Marshal(doc)
	if err != nil {
		log.Printf("Cannot convert document to json, details: %s\n", err.Error())
		return "", err
	}
	return string(ret), nil
}

func Json2Doc(data string) (*Document, error) {
	ret := Document{}
	err := json.Unmarshal([]byte(data), &ret)
	if err != nil {
		log.Printf("Cannot convert json to document, details: %s\n", err.Error())
		return nil, err
	}
	return &ret, nil
}
