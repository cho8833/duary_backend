package shared

import "time"

func BoolPtr(value bool) *bool {
	return &value
}

func StringPtr(value string) *string {
	return &value
}

func TimePtr(value time.Time) *time.Time {
	return &value
}
