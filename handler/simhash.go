package handler

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/HaoyuHu/gosimhash-doc-server/common"
	"strconv"
)

func Simhash(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	PrepareJsonHandler(w, r)

	docContent := r.Form.Get("doc")
	topNStr := r.Form.Get("top_n")

	var bundle *HttpBundle
	if len(docContent) == 0 || len(topNStr) == 0 {
		bundle = CreateErrResponse("Empty doc or top_n")
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
	bundle = CreateOkResponse(simhash)
	SendJsonResponse(w, bundle)
}

func Distance(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	PrepareJsonHandler(w, r)

	firstDoc := r.Form.Get("first_doc")
	secondDoc := r.Form.Get("second_doc")
	topNStr := r.Form.Get("top_n")
	var bundle *HttpBundle
	if len(firstDoc) == 0 || len(secondDoc) == 0 || len(topNStr) == 0 {
		bundle = CreateErrResponse("Empty first_doc or second_doc or top_n")
		SendJsonResponse(w, bundle)
		return
	}
	topN, err := strconv.Atoi(topNStr)
	if err != nil {
		bundle = CreateErrResponse("Incorrect top_n")
		SendJsonResponse(w, bundle)
		return
	}
	first := common.MakeSimhash(&firstDoc, topN)
	second := common.MakeSimhash(&secondDoc, topN)
	dist := common.Distance(first, second)

	ret := map[string]interface{}{"first_simhash": first, "second_simhash": second, "distance": dist}
	bundle = CreateOkResponse(ret)
	SendJsonResponse(w, bundle)
}
