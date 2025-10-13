module update_couple

go 1.22.1

require (
	github.com/cho8833/duary_lambda/model v0.0.0
)

replace (
	github.com/cho8833/duary_lambda/model => ../../model
)