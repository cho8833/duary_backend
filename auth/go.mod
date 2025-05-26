module github.com/cho8833/duary_lambda/auth

go 1.22.1

require (
	github.com/cho8833/duary_lambda/appjwt v0.0.0
	github.com/cho8833/duary_lambda/model v0.0.0
	github.com/cho8833/duary_lambda/shared v0.0.0
)

replace (
	github.com/cho8833/duary_lambda/appjwt => ../appjwt
	github.com/cho8833/duary_lambda/model => ./../model
	github.com/cho8833/duary_lambda/shared => ../shared
)