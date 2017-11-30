package config

import "github.com/HaoyuHu/gosimhash-doc-server/common"

const RedisConfigPath = "config/redis.json"
const CommonConfigPath = "config/common.json"

func GetRedisConfig() map[string]interface{} {
	ret, _ := common.ReadJsonFile(RedisConfigPath)
	return ret
}

func GetCommonConfig() map[string]interface{} {
	ret, _ := common.ReadJsonFile(CommonConfigPath)
	return ret
}
