package main

import (
	"github.com/ahmetson/common-lib/data_type/key_value"
	"github.com/ahmetson/service-lib/communication/command"
	"github.com/ahmetson/service-lib/communication/message"
	"github.com/ahmetson/service-lib/configuration"
	"github.com/ahmetson/service-lib/controller"
	"github.com/ahmetson/service-lib/independent"
	"github.com/ahmetson/service-lib/log"
	"github.com/ahmetson/service-lib/remote"
)

func main() {
	// Any service starts with definition of the logger and configuration
	logger, _ := log.New("hello-service", true)
	appConfig, _ := configuration.New(logger)

	var onHello command.HandleFunc = func(request message.Request, _ *log.Logger, _ ...*remote.ClientSocket) message.Reply {
		name, _ := request.Parameters.GetString("name")
		replyParams := key_value.Empty().Set("message", "hello, "+name)

		return request.Ok(replyParams)
	}
	route := command.NewRoute("hello", onHello)

	replier, _ := controller.SyncReplier(logger)
	_ = replier.AddRoute(route)
	_ = controller.AnyRoute(replier) // add to the replier a route that returns nothing

	service, _ := independent.New(appConfig, logger)

	service.AddController("any-controller", replier)

	service.RequireProxy("github.com/ahmetson/web-proxy", configuration.DefaultContext)
	_ = service.Pipe("github.com/ahmetson/web-proxy", "any-controller")

	err := service.Prepare(configuration.IndependentType)
	if err != nil {
		logger.Fatal("service.Prepare", "error", err)
	}

	service.Run()
}
