module kakao_sign_in

go 1.22.1

require (
	github.com/cho8833/duary_lambda/appjwt v0.0.0
	github.com/cho8833/duary_lambda/auth v0.0.0
	github.com/cho8833/duary_lambda/member v0.0.0
	github.com/cho8833/duary_lambda/shared v0.0.0
)

replace (
	github.com/cho8833/duary_lambda/appjwt => ../../appjwt
	github.com/cho8833/duary_lambda/auth => ../../auth
	github.com/cho8833/duary_lambda/member => ../../member
	github.com/cho8833/duary_lambda/shared => ../../shared
)
