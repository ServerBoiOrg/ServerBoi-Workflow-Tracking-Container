module status

replace generalutils => ../../../ServerBoi-Lambdas-Go/Modules/GeneralUtils

replace discordhttpclient => ../../../ServerBoi-Lambdas-Go/Modules/DiscordHttpClient

replace responseutils => ../../../ServerBoi-Lambdas-Go/Modules/ResponseUtils

go 1.16

require (
	github.com/gin-gonic/gin v1.7.4
	github.com/joho/godotenv v1.3.0
	github.com/rumblefrog/go-a2s v1.0.1
)
