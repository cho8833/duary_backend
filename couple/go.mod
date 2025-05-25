module github.com/cho8833/duary_lambda/couple

go 1.22.1

require (
	github.com/cho8833/duary_lambda/member v0.0.0
	github.com/cho8833/duary_lambda/shared v0.0.0
	github.com/google/uuid v1.6.0
)

replace (
	github.com/cho8833/duary_lambda/member => ../member
	github.com/cho8833/duary_lambda/shared => ../shared
)
