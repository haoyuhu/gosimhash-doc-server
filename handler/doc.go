package handler

import (
	"strconv"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/HaoyuHu/gosimhash-doc-server/common"
	"log"
)

func IdentifyDoc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	PrepareJsonHandler(w, r)

	docId := r.Form.Get("doc_id")
	docContent := r.Form.Get("doc")
	topNStr := r.Form.Get("top_n")
	ageStr := r.Form.Get("age")

	var bundle *HttpBundle
	if len(docId) == 0 || len(docContent) == 0 || len(topNStr) == 0 {
		bundle = CreateErrResponse("Empty doc_id or doc or top_n")
		SendJsonResponse(w, bundle)
		return
	}
	topN, err := strconv.Atoi(topNStr)
	if err != nil {
		bundle = CreateErrResponse("Incorrect top_n")
		SendJsonResponse(w, bundle)
		return
	}
	simhash := common.MakeSimhash(&docContent, topN)

	var age = 0
	if len(ageStr) != 0 {
		age, err = strconv.Atoi(ageStr)
		if err != nil {
			bundle = CreateErrResponse("Incorrect age")
			SendJsonResponse(w, bundle)
			return
		}
	}
	log.Printf("Start to find duplicated simhash in redis.\n")
	cache := common.GetCache()
	success, doc, err := cache.InsertIfNotDuplicated(docId, simhash, int64(age))
	if err != nil {
		bundle = CreateErrResponse(err.Error())
		SendJsonResponse(w, bundle)
		return
	}
	ret := map[string]interface{}{"has_similar_doc": !success}
	if !success {
		ret["similar_doc_id"] = doc.DocId
	}
	bundle = CreateOkResponse(ret)
	SendJsonResponse(w, bundle)
}
