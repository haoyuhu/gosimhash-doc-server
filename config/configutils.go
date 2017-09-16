package config

import "github.com/HaoyuHu/gosimhash-doc-server/common"

const REDIS_CONFIG_PATH = "config/redis.json"
const COMMON_CONFIG_PATH = "config/common.json"

func GetRedisConfig() map[string]interface{} {
	ret, _ := common.ReadJsonFile(REDIS_CONFIG_PATH)
	return ret
}

func GetCommonConfig() map[string]interface{} {
	ret, _ := common.ReadJsonFile(COMMON_CONFIG_PATH)
	return ret
}
