module delete_event

go 1.22.1


require (
	github.com/cho8833/duary_lambda/event v0.0.0
	github.com/cho8833/duary_lambda/shared v0.0.0
)

replace (
	github.com/cho8833/duary_lambda/event => ../../event
	github.com/cho8833/duary_lambda/shared => ../../shared
)