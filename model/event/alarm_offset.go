package event

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

type AlarmOffset string

const (
	None   AlarmOffset = "NONE"
	AtTime AlarmOffset = "AT_TIME"
	Min5   AlarmOffset = "MINUTE_5"
	Min10  AlarmOffset = "MINUTE_10"
	Min15  AlarmOffset = "MINUTE_15"
	Min30  AlarmOffset = "MINUTE_30"
	Hour1  AlarmOffset = "HOUR_1"
)

type AlarmOffsetDuration struct {
	Name     AlarmOffset
	Duration time.Duration
}

var AlarmOffsetMap = map[AlarmOffset]AlarmOffsetDuration{
	None: {
		Name:     None,
		Duration: time.Duration(0),
	},
	AtTime: {
		Name:     AtTime,
		Duration: time.Duration(0),
	},
	Min5: {
		Name:     Min5,
		Duration: time.Duration(-5) * time.Minute,
	},
	Min10: {
		Name:     Min10,
		Duration: time.Duration(-10) * time.Minute,
	},
	Min15: {
		Name:     Min15,
		Duration: time.Duration(-15) * time.Minute,
	},
	Min30: {
		Name:     Min30,
		Duration: time.Duration(-30) * time.Minute,
	},
	Hour1: {
		Name:     Hour1,
		Duration: time.Duration(-1) * time.Hour,
	},
}

func (af AlarmOffset) ToJSON() string {
	return string(af)
}

func FromJSON(input string) (AlarmOffset, error) {
	af := AlarmOffset(input)
	if _, ok := AlarmOffsetMap[af]; ok {
		return af, nil
	}
	return None, fmt.Errorf("unknown AlarmOffset: %s", input)
}

func (af AlarmOffset) MarshalJSON() ([]byte, error) {
	return json.Marshal(af.ToJSON())
}

func (af *AlarmOffset) UnmarshalJSON(data []byte) error {
	var name string
	if err := json.Unmarshal(data, &name); err != nil {
		return err
	}
	val, err := FromJSON(name)
	if err != nil {
		return err
	}
	*af = val
	return nil
}

func (af AlarmOffset) MarshalDynamoDBAttributeValue() (types.AttributeValue, error) {
	return &types.AttributeValueMemberS{
		Value: af.ToJSON(),
	}, nil
}

func (af *AlarmOffset) UnmarshalDynamoDBAttributeValue(av types.AttributeValue) error {
	s := av.(*types.AttributeValueMemberS).Value
	val, err := FromJSON(s)
	if err != nil {
		return err
	}
	*af = val
	return nil

}
