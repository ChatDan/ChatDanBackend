package main

import (
	"ChatDanBackend/bootstrap"
	_ "ChatDanBackend/docs"
	"ChatDanBackend/utils"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @title          ChatDan Backend
// @version        0.0.1
// @description    ChatDan, a message box and 'biaobai' platform for Fudaners.
// @termsOfService https://swagger.io/terms/

// @contact.name   JingYiJun
// @contact.url    https://www.jingyijun.xyz
// @contact.email  jingyijun3104@outlook.com

// @license.name  Apache 2.0
// @license.url   https://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	app := bootstrap.InitFiberApp()

	go func() {
		if innerErr := app.Listen("0.0.0.0:8000"); innerErr != nil {
			log.Println(innerErr)
		}
	}()

	interrupt := make(chan os.Signal, 1)

	// wait for CTRL-C interrupt
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	// close app
	err := app.Shutdown()
	if err != nil {
		utils.Logger.Error("app shutdown error", zap.Error(err))
	}

	_ = utils.Logger.Sync()
}