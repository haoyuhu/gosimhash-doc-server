package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	CODE_OK  = 0
	CODE_ERR = -1
)
const OK_MESSAGE = "ok"

type HttpBundle struct {
	Code       int    `json:"code"`
	ErrMessage string `json:"err_message"`
	Data interface{} `json:"data"`
}

func CreateOkResponse(data interface{}) *HttpBundle {
	bundle := &HttpBundle{}
	bundle.Code = CODE_OK
	bundle.ErrMessage = OK_MESSAGE
	if data != nil {
		bundle.Data = data
	}
	return bundle
}

func CreateErrResponse(errMessage string) *HttpBundle {
	bundle := &HttpBundle{}
	bundle.Code = CODE_ERR
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
