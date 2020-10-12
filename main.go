package main

import (
	"CloudRestaurant/controller"
	"CloudRestaurant/tool"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := tool.ParseConfig("./config/app.json")

	if err != nil {
		panic(err.Error())
	}

	app := gin.Default()

	registerRouter(app)

//	fmt.Println(cfg.AppHost + ":" + cfg.AppPort)

	app.Run(cfg.AppHost + ":" + cfg.AppPort)

}

func registerRouter(router *gin.Engine) {
	new(controller.HelloController).Router(router)
}
