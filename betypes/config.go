package betypes

import "os"

var (
	// BOT_TOKEN given by BotFather to use telegram API CRM_BOT_TOKEN
	BOT_TOKEN = "1324365499:AAFjGjP0fKcq3JcnTaoAYikz0DFS0BUbpKQ"
	// SECRET_VASSAL_PASSWORD to register new user SECRET_CRM_BOT_PASSWORD_FOR_VASSALS
	SECRET_VASSAL_PASSWORD = os.Getenv("SECRET_CRM_BOT_PASSWORD_FOR_VASSALS")
	// SECRET_LORD_PASSWORD to register new admin SECRET_CRM_BOT_PASSWORD_FOR_LORDS
	SECRET_LORD_PASSWORD = os.Getenv("SECRET_CRM_BOT_PASSWORD_FOR_LORDS")
)

const (
	// MongoUsername ...
	MongoUsername = "haski0071"
	// MongoPassword ...
	MongoPassword = "Haski12345"
	// MongoHostname ...
	MongoHostname = "172.18.0.2"
	// MongoPort ...
	MongoPort = "27017"
)
