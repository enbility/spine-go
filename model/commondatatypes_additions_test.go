package model

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestTimePeriodType(t *testing.T) {
	tc := &TimePeriodType{}
	duration, err := tc.GetDuration()
	assert.NotNil(t, err)
	assert.Equal(t, time.Duration(0), duration)

	tc = &TimePeriodType{
		EndTime: NewAbsoluteOrRelativeTimeTypeFromDuration(time.Minute * 1),
	}
	duration, err = tc.GetDuration()
	assert.Nil(t, err)
	assert.Equal(t, time.Minute*1, duration)

	tc = NewTimePeriodTypeWithRelativeEndTime(time.Minute * 1)

	duration, err = tc.GetDuration()
	assert.Nil(t, err)
	assert.Equal(t, time.Minute*1, duration)

	data, err := json.Marshal(tc)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "{\"endTime\":\"PT1M\"}", string(data))

	var tp1 TimePeriodType
	err = json.Unmarshal(data, &tp1)
	assert.Nil(t, err)
	assert.Equal(t, *tc.EndTime, *tp1.EndTime)

	time.Sleep(time.Second * 1)

	duration, err = tc.GetDuration()
	assert.Nil(t, err)
	assert.Equal(t, time.Second*59, duration)

	data, err = json.Marshal(tc)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "{\"endTime\":\"PT59S\"}", string(data))
}

func TestTimeType(t *testing.T) {
	tc := []struct {
		in    string
		parse string
	}{
		{"21:32:52.12679", "15:04:05.999999999"},
		{"21:32:52.12679Z", "15:04:05.999999999Z"},
		{"21:32:52", "15:04:05"},
		{"19:32:52Z", "15:04:05Z"},
		{"19:32:52+07:00", "15:04:05+07:00"},
		{"19:32:52-07:00", "15:04:05-07:00"},
	}

	for _, tc := range tc {
		got := NewTimeType(tc.in)
		expect, err := time.ParseInLocation(tc.parse, tc.in, time.UTC)
		if err != nil {
			t.Errorf("Parsing failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}
		value, err := got.GetTime()
		if err != nil {
			t.Errorf("Test Failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}

		if value.UTC() != expect.UTC() {
			t.Errorf("Test failure for %s, expected %s and got %s", tc.in, value.String(), expect.String())
		}
	}
}

func TestDateType(t *testing.T) {
	tc := []struct {
		in    string
		parse string
	}{
		{"2022-02-01", "2006-01-02"},
		{"2022-02-01Z", "2006-01-02Z"},
		{"2022-02-01+07:00", "2006-01-02+07:00"},
	}

	for _, tc := range tc {
		got := NewDateType(tc.in)
		expect, err := time.ParseInLocation(tc.parse, tc.in, time.UTC)
		if err != nil {
			t.Errorf("Parsing failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}
		value, err := got.GetTime()
		if err != nil {
			t.Errorf("Test Failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}

		if value.UTC() != expect.UTC() {
			t.Errorf("Test failure for %s, expected %s and got %s", tc.in, value.String(), expect.String())
		}
	}
}

func TestDateTimeType(t *testing.T) {
	tc := []struct {
		in    string
		parse string
	}{
		{"2022-02-01T21:32:52.12679", "2006-01-02T15:04:05.999999999"},
		{"2022-02-01T21:32:52.12679Z", "2006-01-02T15:04:05.999999999Z"},
		{"2022-02-01T21:32:52", "2006-01-02T15:04:05"},
		{"2022-02-01T19:32:52Z", "2006-01-02T15:04:05Z"},
	}

	for _, tc := range tc {
		got := NewDateTimeType(tc.in)
		expect, err := time.Parse(tc.parse, tc.in)
		if err != nil {
			t.Errorf("Parsing failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}
		value, err := got.GetTime()
		if err != nil {
			t.Errorf("Test Failure with %s and parser %s: %s", tc.in, tc.parse, err)
			continue
		}
		nDateTime := NewDateTimeTypeFromTime(value)
		nValue, err := nDateTime.GetTime()
		assert.Nil(t, err)
		assert.Equal(t, nValue.Hour(), value.Hour())
		assert.Equal(t, nValue.Minute(), value.Minute())
		assert.Equal(t, nValue.Second(), value.Second())
		if value.UTC() != expect.UTC() {
			t.Errorf("Test failure for %s, expected %s and got %s", tc.in, value.String(), expect.String())
		}
	}
}

func TestDurationType(t *testing.T) {
	tc := []struct {
		in  time.Duration
		out string
	}{
		{time.Duration(4) * time.Second, "PT4S"},
	}

	for _, tc := range tc {
		duration := NewDurationType(tc.in)
		got, err := duration.GetTimeDuration()
		if err != nil {
			t.Errorf("Test Failure with %s: %s", tc.in, err)
			continue
		}
		if got != tc.in {
			t.Errorf("Test failure for %d, got %d", tc.in, got)
		}
		if string(*duration) != tc.out {
			t.Errorf("Test failure for %d, expected %s got %s", tc.in, tc.out, string(*duration))
		}
	}
}

func TestAbsoluteOrRelativeTimeTypeAbsolute(t *testing.T) {
	tc := []struct {
		in       string
		dateTime time.Time
	}{
		{"2022-02-01T19:32:52Z", time.Date(2022, 02, 01, 19, 32, 52, 0, time.UTC)},
	}

	for _, tc := range tc {
		a := NewAbsoluteOrRelativeTimeType(tc.in)
		got, err := a.GetTime()
		if err != nil {
			t.Errorf("Test Failure with %s: %s", tc.in, err)
			continue
		}
		if got != tc.dateTime {
			t.Errorf("Test failure for %s, expected %s got %s", tc.in, tc.dateTime.String(), got.String())
		}

		d := a.GetDateTimeType()
		got, err = d.GetTime()
		if err != nil {
			t.Errorf("Test Failure with %s: %s", tc.in, err)
			continue
		}
		if got != tc.dateTime {
			t.Errorf("Test failure for %s, expected %s got %s", tc.in, tc.dateTime.String(), got.String())
		}
	}
}

func TestAbsoluteOrRelativeTimeTypeDuration(t *testing.T) {
	tc := []struct {
		in  time.Duration
		out string
	}{
		{time.Duration(4) * time.Second, "PT4S"},
	}

	for _, tc := range tc {
		a := NewAbsoluteOrRelativeTimeTypeFromDuration(tc.in)
		got, err := a.GetDurationType()
		if err != nil {
			t.Errorf("Test Failure with %d: %s", tc.in, err)
			continue
		}
		if string(*got) != tc.out {
			t.Errorf("Test failure for %d, expected %s got %s", tc.in, tc.out, string(*got))
		}

		d, err := a.GetTimeDuration()
		if err != nil {
			t.Errorf("Test Failure with %d: %s", tc.in, err)
			continue
		}
		got = NewDurationType(d)
		if string(*got) != tc.out {
			t.Errorf("Test failure for %d, expected %s got %s", tc.in, tc.out, string(*got))
		}
	}
}

func TestAbsoluteOrRelativeTimeTypeRelative(t *testing.T) {
	tc := []struct {
		in  string
		out time.Duration
	}{
		{"PT4S", time.Duration(4) * time.Second},
	}

	for _, tc := range tc {
		a := NewAbsoluteOrRelativeTimeType(tc.in)
		got, err := a.GetTimeDuration()
		if err != nil {
			t.Errorf("Test Failure with %s: %s", tc.in, err)
			continue
		}
		if got != tc.out {
			t.Errorf("Test failure for %s, expected %d got %d", tc.in, tc.out, got)
		}

		d, err := a.GetDurationType()
		if err != nil {
			t.Errorf("Test Failure with %s: %s", tc.in, err)
			continue
		}
		got, err = d.GetTimeDuration()
		if err != nil {
			t.Errorf("Test Failure with %s: %s", tc.in, err)
			continue
		}
		if got != tc.out {
			t.Errorf("Test failure for %s, expected %d got %d", tc.in, tc.out, got)
		}
	}
}

func TestNewScaledNumberType(t *testing.T) {
	tc := []struct {
		in     float64
		number int64
		scale  int
	}{
		{0, 0, 0},
		{0.1, 1, -1},
		{1.0, 1, 0},
		{6.25, 625, -2},
		{10, 10, 0},
		{12.5952, 125952, -4},
		{13.1637, 131637, -4},
	}

	for _, tc := range tc {
		got := NewScaledNumberType(tc.in)
		number := int64(*got.Number)
		scale := 0
		if got.Scale != nil {
			scale = int(*got.Scale)
		}
		if number != tc.number || scale != tc.scale {
			t.Errorf("NewScaledNumberType(%v) = %d %d, want %d %d", tc.in, got.Number, got.Scale, tc.number, tc.scale)
		}

		val := got.GetValue()
		if val != tc.in {
			t.Errorf("GetValue(%d %d) = %f, want %f", tc.number, tc.scale, val, tc.in)
		}
	}
}

func TestDeviceAddressTypeString(t *testing.T) {
	tc := []struct {
		device AddressDeviceType
		out    string
	}{
		{
			"Device 1",
			"Device 1",
		},
	}

	for _, tc := range tc {
		f := DeviceAddressType{
			Device: util.Ptr(tc.device),
		}

		got := f.String()
		if got != tc.out {
			t.Errorf("TestDeviceAddressTypeString(), got %s, expects %s", got, tc.out)
		}
	}
}

func TestEntityAddressTypeString(t *testing.T) {
	tc := []struct {
		device AddressDeviceType
		entity []AddressEntityType
		out    string
	}{
		{
			"Device",
			[]AddressEntityType{1, 1},
			"Device:[1,1]:",
		},
	}

	for _, tc := range tc {
		f := FeatureAddressType{
			Device: util.Ptr(tc.device),
			Entity: tc.entity,
		}

		got := f.String()
		if got != tc.out {
			t.Errorf("TestEntityAddressTypeString(), got %s, expects %s", got, tc.out)
		}
	}
}

func TestFeatureAddressTypeString(t *testing.T) {
	tc := []struct {
		device  AddressDeviceType
		entity  []AddressEntityType
		feature AddressFeatureType
		out     string
	}{
		{
			"Device",
			[]AddressEntityType{1, 1},
			0,
			"Device:[1,1]:0",
		},
	}

	for _, tc := range tc {
		f := FeatureAddressType{
			Device:  util.Ptr(tc.device),
			Entity:  tc.entity,
			Feature: util.Ptr(tc.feature),
		}

		got := f.String()
		if got != tc.out {
			t.Errorf("TestFeatureAddressTypeString(), got %s, expects %s", got, tc.out)
		}
	}
}
