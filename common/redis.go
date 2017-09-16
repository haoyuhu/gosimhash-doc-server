package common

import (
	"errors"
	"strconv"
	"log"
	"github.com/go-redis/redis"
	"github.com/HaoyuHu/gosimhash"
	"github.com/HaoyuHu/gosimhash-doc-server/model"
)

type RedisSimhashCache struct {
	operator *SimhashOperator
	client   *redis.Client
	limit    SimhashLimit
}

var cache *RedisSimhashCache

func InitializeRedisCache(limit SimhashLimit, host string, port int, passwd string) error {
	curr, err := NewRedisSimhashCache(limit, host, port, passwd)
	if err != nil {
		log.Printf("Cannot initialize redis cache, details: %s\n", err.Error())
		return err
	}
	if !curr.CheckConnection() {
		return errors.New("redis checkout failed")
	}
	cache = curr
	return nil
}

func GetCache() *RedisSimhashCache {
	return cache
}

func NewRedisSimhashCache(limit SimhashLimit, host string, port int, passwd string) (*RedisSimhashCache, error) {
	// init simhash operator
	op := NewSimhashOperator(int(limit) + 1)
	if !op.Check() {
		return nil, errors.New("incorrect limit for simhash operator")
	}
	// init redis server connection
	defaultDBIndex := 0
	endpoint := host + ":" + strconv.Itoa(port)
	client := redis.NewClient(&redis.Options{
		Addr:     endpoint,
		Password: passwd,
		DB:       defaultDBIndex})

	return &RedisSimhashCache{operator: op, client: client, limit: limit}, nil
}

func (redis *RedisSimhashCache) CheckConnection() bool {
	// test connection
	_, err := redis.client.Ping().Result()
	if err != nil {
		log.Printf("Cannot connect to redis server: %s\n", err.Error())
		return false
	}
	return true
}

func (redis *RedisSimhashCache) Init(docIds []string, simhashList []uint64, timeouts []int64) int {
	if len(docIds) != len(simhashList) || len(docIds) != len(timeouts) {
		return 0
	}
	count := 0
	for i, docId := range docIds {
		if success, _, err := redis.InsertIfNotDuplicated(docId, simhashList[i], timeouts[i]); err == nil && success {
			count += 1
		}
	}
	return count
}

func (redis *RedisSimhashCache) InsertIfNotDuplicated(docId string, simhash uint64, age int64) (bool, *model.Document, error) {
	var expireTime int64 = 0
	if age != 0 {
		expireTime = Now() + age
	}
	// find similar simhash in redis
	similarDoc, err := redis.similarDocExists(simhash)
	if err != nil {
		return false, nil, err
	}
	if similarDoc != nil {
		log.Printf("Find a document in redis which is similar with current document, doc_id is %s, simhash is %d\n", similarDoc.DocId, similarDoc.Simhash)
		return false, similarDoc, nil
	}

	// save document to redis if not exists
	doc := &model.Document{DocId: docId, Simhash: simhash, ExpireTime: expireTime}
	data, err := doc.Doc2Json()
	if err == nil {
		parts := redis.operator.Cut(simhash)
		for _, part := range parts {
			key := strconv.FormatUint(part, UINT_BASE)
			_, err := redis.client.LPush(key, data).Result()
			if err != nil {
				log.Printf("Cannot insert doc to redis, details: %s\n", err.Error())
			}
		}
	} else {
		log.Printf("Cannot convert doc to json string, details: %s\n", err.Error())
	}

	return true, nil, nil
}

func (redis *RedisSimhashCache) similarDocExists(simhash uint64) (*model.Document, error) {
	const LAST_ONE = -1

	parts := redis.operator.Cut(simhash)
	for _, part := range parts {
		key := strconv.FormatUint(part, UINT_BASE)
		// get all items from redis
		docs, err := redis.client.LRange(key, 0, LAST_ONE).Result()
		if err != nil {
			log.Printf("Cannot get doc list from redis, details: %s\n", err.Error())
			return nil, err
		}
		if len(docs) == 0 {
			continue
		}
		for i := len(docs) - 1; i >= 0; i-- {
			doc, err := model.Json2Doc(docs[i])
			if err != nil {
				log.Printf("Cannot convert json string to doc, details: %s\n", err.Error())
				continue
			}
			// has expired
			if doc.ExpireTime != 0 && doc.ExpireTime <= Now() {
				// remove last ONE item that equals to docs[i]
				_, err := redis.client.LRem(key, LAST_ONE, docs[i]).Result()
				if err != nil {
					log.Printf("Cannot remove expired item(%s), details: %s\n", docs[i], err.Error())
				}
				continue
			}
			if gosimhash.IsSimhashDuplicated(simhash, doc.Simhash, int(redis.limit)) {
				return doc, nil
			}
		}
	}
	return nil, nil
}
