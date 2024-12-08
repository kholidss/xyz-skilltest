package http

import (
	"fmt"
	"github.com/kholidss/xyz-skilltest/pkg/app"
	"github.com/kholidss/xyz-skilltest/pkg/config"
	"github.com/kholidss/xyz-skilltest/pkg/logger"
)

func Start() {
	logger.SetJSONFormatter()
	cnf, err := config.LoadAllConfigs()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Failed to load configuration file: %v", err))
	}

	app.InitializeApp(cnf)
	application := app.GetServer()

	//go func() {
	//	StartServiceListener(cnf)
	//}()

	if err := application.StartServer(); err != nil {
		logger.Fatal(fmt.Sprintf("Failed to start server: %v", err))
	}
}
