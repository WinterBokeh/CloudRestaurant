package tool

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

var Gsession sessions.Session

func InitSession(engine *gin.Engine) {
	//参见文档
	cfg := GetConfig().Redis

	store, err := redis.NewStore(10, "tcp", cfg.Addr + ":" + cfg.Port, "", []byte("secret"))
	if err != nil {
		fmt.Println(err.Error())
	}

	engine.Use( sessions.Sessions("mysession", store) )
}

func SetSession(ctx *gin.Context, key interface{}, value interface{}) error {
	session := sessions.Default(ctx)
	Gsession = session

	if session == nil {
		return nil
	}
	session.Set(key, value)
	session.Set(6, 666)
//	fmt.Println("QAQAQAQ", session)
	return session.Save()
}

func GetSession(ctx *gin.Context, key interface{}) interface{} {
	//session := sessions.Default(gCtx)

//	fmt.Printf(" Get——%T ", key)
//	fmt.Println("阿巴", Gsession, "con", Gsession.Get(6) )
	return Gsession.Get(key)
}