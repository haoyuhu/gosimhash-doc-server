package main

import (
	"log"
	"strings"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/HaoyuHu/gosimhash-doc-server/handler"
	"github.com/HaoyuHu/gosimhash-doc-server/config"
	"github.com/HaoyuHu/gosimhash-doc-server/common"
	"os"
)

func main() {
	commonConf := config.GetCommonConfig()
	redisConf := config.GetRedisConfig()
	if commonConf == nil || redisConf == nil {
		log.Fatal("Cannot initialize common and redis configurations")
	}
	host, hostExists := os.LookupEnv("HOST")
	port, portExists := os.LookupEnv("PORT")
	if !hostExists {
		value, exists := commonConf["host"]
		if exists {
			host = value.(string)
			hostExists = exists
		}
	}
	if !portExists {
		value, exists := commonConf["port"]
		if exists {
			port = value.(string)
			portExists = exists
		}
	}
	if !hostExists || !portExists {
		log.Fatal("Incorrect host or port in env or common.json")
	}
	addr := strings.Join([]string{host, port}, ":")

	router := httprouter.New()
	router.POST("/simhash", handler.Simhash)
	router.POST("/distance", handler.Distance)
	router.POST("/identify", handler.IdentifyDoc)

	hashType := commonConf["hash_type"]
	hmm := commonConf["hmm_dict"]
	idf := commonConf["idf_dict"]
	user := commonConf["user_dict"]
	stopWords := commonConf["stop_words"]
	limit := commonConf["simhash_limit"]
	if hashType == nil || hmm == nil || idf == nil || user == nil || stopWords == nil || limit == nil {
		log.Fatal("Incorrect hash_type or hmm_dict or idf_dict or user_dict or stop_words or limit in common.json")
	}
	common.InitializeSimhasher(common.HashType(uint8(hashType.(float64))), "", hmm.(string), user.(string), idf.(string), stopWords.(string))

	redisHost := redisConf["host"]
	redisPort := redisConf["port"]
	redisPasswd := redisConf["passwd"]
	if redisHost == nil || redisPort == nil || redisPasswd == nil {
		log.Fatal("Incorrect host or port or passwd in redis.json")
	}
	err := common.InitializeRedisCache(common.SimhashLimit(int(limit.(float64))), redisHost.(string), int(redisPort.(float64)), redisPasswd.(string))
	if err != nil {
		log.Fatal("Initialize redis cache failed.")
	}

	log.Fatal(http.ListenAndServe(addr, router))
}
