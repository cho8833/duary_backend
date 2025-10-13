module send_notification

go 1.22.1

require (
	github.com/cho8833/duary_lambda/fcm v0.0.0
	github.com/cho8833/duary_lambda/shared v0.0.0

)

replace (
	github.com/cho8833/duary_lambda/fcm => ../../fcm
	github.com/cho8833/duary_lambda/shared => ../../shared

)