package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	CodeOk  = 0
	CodeErr = -1
)
const OkMessage = "ok"

type HttpBundle struct {
	Code       int         `json:"code"`
	ErrMessage string      `json:"err_message"`
	Data       interface{} `json:"data"`
}

func CreateOkResponse(data interface{}) *HttpBundle {
	bundle := &HttpBundle{}
	bundle.Code = CodeOk
	bundle.ErrMessage = OkMessage
	if data != nil {
		bundle.Data = data
	}
	return bundle
}

func CreateErrResponse(errMessage string) *HttpBundle {
	bundle := &HttpBundle{}
	bundle.Code = CodeErr
	bundle.ErrMessage = errMessage
	return bundle
}

func PrepareJsonHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	r.ParseForm()
}

func SendJsonResponse(w http.ResponseWriter, bundle *HttpBundle) {
	ret, _ := json.Marshal(bundle)
	w.Write(ret)
}

func parseHttpBundle(bytes []byte) (*HttpBundle, error) {
	bundle := &HttpBundle{}
	err := json.Unmarshal(bytes, bundle)
	if err != nil {
		log.Printf("Cannot convert bytes to HttpBundle, details: %s\n", err.Error())
		return nil, err
	}
	return bundle, nil
}
