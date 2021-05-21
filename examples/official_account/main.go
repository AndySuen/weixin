package main

import (
	"fmt"
	"net/url"

	"github.com/lixinio/weixin/test"
	"github.com/lixinio/weixin/utils/redis"
	"github.com/lixinio/weixin/weixin/official_account"
	"github.com/lixinio/weixin/weixin/user_api"
)

func main() {
	cache := redis.NewRedis(&redis.Config{RedisUrl: test.CacheUrl})
	officialAccount := official_account.New(cache, &official_account.Config{
		Appid:  test.OfficialAccountAppid,
		Secret: test.OfficialAccountSecret,
	})

	userApi := user_api.NewOfficialAccountApi(officialAccount)

	params := url.Values{}
	b, e := userApi.Get(params)
	fmt.Println(string(b), " ", e)

}
