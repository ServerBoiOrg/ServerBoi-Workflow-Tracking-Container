module workflow-tracking

require (
	discordhttpclient v0.0.0
	generalutils v0.0.0
	github.com/awlsring/discordtypes v0.1.6
	responseutils v0.0.0
)

replace generalutils => ../../../ServerBoi-Lambdas-Go/Modules/GeneralUtils

replace discordhttpclient => ../../../ServerBoi-Lambdas-Go/Modules/DiscordHttpClient

replace responseutils => ../../../ServerBoi-Lambdas-Go/Modules/ResponseUtils

go 1.16
