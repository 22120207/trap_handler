package services

var redisService *RedisService

func init() {
	var err error
	redisService, err = NewRedisService()

	if err != nil {
		panic(err)
	}
}
