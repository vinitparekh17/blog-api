package main

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	pb "github.com/jay-bhogayata/notifyHub/proto"
)

type config struct {
	port         string
	env          string
	sender_email string
	awsConfig    aws.Config
}

type application struct {
	config config
	pb.UnimplementedNotificationServiceServer
}

//	@title			NotifyHub API
//	@version		1.0.0
//	@description	This is a api for sending emails and sms using different providers.
//	@contact.name	Jay Bhogayata
//	@contact.url	https://github.com/jay-bhogayata/
//	@contact.email	jaybhogayata53@gmail.com

// @license.name    MIT
// @license.url     http://opensource.org/licenses/MIT
// @host.path	localhost:8080
// @BasePath	/api/v1
func main() {

	var cfg config

	app := &application{
		config: cfg,
	}

	app.LoggerInit()
	app.config.LoadConfig()
	app.config.ConfigAws()
	app.ServerInit()
}
