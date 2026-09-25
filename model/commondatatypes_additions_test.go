package model

import (
	"encoding/json"
	"fmt"
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
		EndTime: NewAbsoluteOrRelativeTimeTypeFromDuration(time.Second * 3),
	}
	duration, err = tc.GetDuration()
	assert.Nil(t, err)
	assert.Equal(t, time.Second*3, duration)

	tc = NewTimePeriodTypeWithRelativeEndTime(time.Second * 3)

	duration, err = tc.GetDuration()
	assert.Nil(t, err)
	assert.Equal(t, time.Second*3, duration)

	data, err := json.Marshal(tc)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "{\"endTime\":\"PT3S\"}", string(data))

	var tp1 TimePeriodType
	err = json.Unmarshal(data, &tp1)
	assert.Nil(t, err)
	assert.Equal(t, *tc.EndTime, *tp1.EndTime)

	time.Sleep(time.Second * 1)

	duration, err = tc.GetDuration()
	assert.Nil(t, err)
	assert.Equal(t, time.Second*2, duration)

	data, err = json.Marshal(tc)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "{\"endTime\":\"PT2S\"}", string(data))

	time.Sleep(time.Second * 3)

	duration, err = tc.GetDuration()
	assert.Nil(t, err)
	assert.Equal(t, time.Second*0, duration)

	data, err = json.Marshal(tc)
	assert.Nil(t, err)
	assert.NotNil(t, data)
	assert.Equal(t, "{\"endTime\":\"P0D\"}", string(data))
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

func TestScaledNumberTypeUnmarshal(t *testing.T) {
	tc := []struct {
		json   string
		number int64
		scale  int16
	}{
		{`{"number":42}`, 42, 0},
		{`{"number":42,"scale":-2}`, 42, -2},
		{`{"number":-9007199254740990}`, -9007199254740990, 0},
		{`{"number":9007199254740990}`, 9007199254740990, 0},
		// sender encodes the number in scientific notation (float-like)
		{`{"number":9.00719925474099e+15}`, 9007199254740990, 0},
		{`{"number":-9.00719925474099e+15}`, -9007199254740990, 0},
		{`{"number":1E3}`, 1000, 0},
		{`{"number":1.5e1}`, 15, 0},
		{`{"number":42.0}`, 42, 0},
	}

	for _, tc := range tc {
		var got ScaledNumberType
		err := json.Unmarshal([]byte(tc.json), &got)
		assert.NoError(t, err, "input: %s", tc.json)
		if err != nil {
			continue
		}
		assert.Equal(t, NumberType(tc.number), *got.Number, "input: %s", tc.json)
		wantScale := ScaleType(tc.scale)
		if tc.scale == 0 {
			// scale is omitted when zero; accept both nil and explicit zero
			if got.Scale != nil {
				assert.Equal(t, wantScale, *got.Scale, "input: %s", tc.json)
			}
		} else {
			assert.Equal(t, wantScale, *got.Scale, "input: %s", tc.json)
		}
	}
}

func TestScaledNumberTypeUnmarshalInvalid(t *testing.T) {
	// only values decoding cleanly into an integer are accepted
	tc := []string{
		`{"number":1.5}`,
		`{"number":1.55e1}`,
		`{"number":9223372036854775808}`,
		`{"number":1e19}`,
		`{"number":"42"}`,
		`{"number":true}`,
	}

	for _, in := range tc {
		var got ScaledNumberType
		assert.Error(t, json.Unmarshal([]byte(in), &got), "input: %s", in)
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
		if got.Scale == nil {
			t.Errorf("NewScaledNumberType(%v): Scale may not be nil", tc.in)
		}

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

// TestDurationTypeIssue60 validates the fix for issue #60
// Ensures complex durations are formatted with preserved structure instead of seconds-only
func TestDurationTypeIssue60(t *testing.T) {
	// Test case from issue #60: complex duration should preserve structure
	duration := time.Duration(4357512417) * time.Second // Parsed P138Y1MT6H28M15S

	result := NewDurationType(duration)
	resultStr := string(*result)

	// Should NOT be "PT4357512417S" (old behavior)
	// Should use days: "P50434DT4H6M57S" (exact round-trip, no calendar dependency)
	assert.NotEqual(t, "PT4357512417S", resultStr, "Should not output seconds-only format")
	assert.Contains(t, resultStr, "D", "Should contain day component")

	// Verify it's still a valid duration that can be parsed back exactly
	parsedBack, err := result.GetTimeDuration()
	assert.NoError(t, err, "Result should be parseable")
	assert.Equal(t, duration, parsedBack, "Round-trip should be exact when using days")
}

// TestNewDurationTypeEdgeCases tests edge cases for the calendar-aware duration formatting
func TestNewDurationTypeEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		duration      time.Duration
		expectedRegex string
		description   string
	}{
		{
			name:          "zero duration",
			duration:      0,
			expectedRegex: "^P0D$",
			description:   "Zero duration should be P0D",
		},
		{
			name:          "exactly one second",
			duration:      1 * time.Second,
			expectedRegex: "^PT1S$",
			description:   "Should format single second",
		},
		{
			name:          "exactly one minute",
			duration:      1 * time.Minute,
			expectedRegex: "^PT1M$",
			description:   "Should format single minute",
		},
		{
			name:          "exactly one hour",
			duration:      1 * time.Hour,
			expectedRegex: "^PT1H$",
			description:   "Should format single hour",
		},
		{
			name:          "exactly 24 hours",
			duration:      24 * time.Hour,
			expectedRegex: "^P1D$",
			description:   "24 hours should become 1 day",
		},
		{
			name:          "just under 24 hours",
			duration:      23*time.Hour + 59*time.Minute + 59*time.Second,
			expectedRegex: "^PT23H59M59S$",
			description:   "Should not round up to days",
		},
		{
			name:          "exactly 7 days",
			duration:      7 * 24 * time.Hour,
			expectedRegex: "^P7D$",
			description:   "Should format as days, not weeks (calendar-aware)",
		},
		{
			name:          "complex time only",
			duration:      2*time.Hour + 30*time.Minute + 45*time.Second,
			expectedRegex: "^PT2H30M45S$",
			description:   "Should handle complex time components",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDurationType(tt.duration)
			resultStr := string(*result)

			assert.Regexp(t, tt.expectedRegex, resultStr, tt.description)

			// All durations should round-trip exactly (days-based, no calendar approximation)
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err, "Result should be parseable")
			assert.Equal(t, tt.duration, parsedBack, "Round-trip should be exact")
		})
	}
}

// TestNewDurationTypeNegative tests negative duration handling
func TestNewDurationTypeNegative(t *testing.T) {
	tests := []struct {
		name         string
		duration     time.Duration
		expectedSign string
	}{
		{
			name:         "negative 1 hour",
			duration:     -1 * time.Hour,
			expectedSign: "-PT1H",
		},
		{
			name:         "negative 1 day",
			duration:     -24 * time.Hour,
			expectedSign: "-P1D",
		},
		{
			name:         "negative complex",
			duration:     -(2*time.Hour + 30*time.Minute),
			expectedSign: "-PT2H30M",
		},
		{
			name:         "negative zero",
			duration:     0,
			expectedSign: "P0D", // Zero is not negative
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDurationType(tt.duration)
			resultStr := string(*result)

			assert.Equal(t, tt.expectedSign, resultStr)

			// Verify parsing back gives the same duration
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err)
			assert.Equal(t, tt.duration, parsedBack)
		})
	}
}

// TestNewDurationTypeLargeDayBoundaries tests large day-based durations
func TestNewDurationTypeLargeDayBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "365 days",
			duration: 365 * 24 * time.Hour,
		},
		{
			name:     "730 days",
			duration: 2 * 365 * 24 * time.Hour,
		},
		{
			name:     "30 days",
			duration: 30 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDurationType(tt.duration)
			resultStr := string(*result)

			// Should use days, not seconds-only format
			assert.Contains(t, resultStr, "D", "Should contain day component")
			assert.NotRegexp(t, "^PT\\d+S$", resultStr, "Should not be seconds-only format")

			// Should round-trip exactly
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err)
			assert.Equal(t, tt.duration, parsedBack, "Round-trip should be exact")
		})
	}
}

// TestNewDurationTypeMonthBoundaries tests durations around month-length boundaries
func TestNewDurationTypeMonthBoundaries(t *testing.T) {
	monthLengths := []int{28, 29, 30, 31}

	for _, days := range monthLengths {
		t.Run(fmt.Sprintf("%d days", days), func(t *testing.T) {
			duration := time.Duration(days) * 24 * time.Hour
			result := NewDurationType(duration)
			resultStr := string(*result)

			// Should use days format (e.g. P28D, P30D)
			assert.Regexp(t, "^P\\d+D$", resultStr, "Should be days-only format")

			// Should round-trip exactly
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err)
			assert.Equal(t, duration, parsedBack, "Round-trip should be exact for day-based durations")
		})
	}
}

// TestNewDurationTypeLargeValues tests handling of large durations
func TestNewDurationTypeLargeValues(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
	}{
		{
			name:     "10 years in days",
			duration: 10 * 365 * 24 * time.Hour,
		},
		{
			name:     "100 years in days",
			duration: 100 * 365 * 24 * time.Hour,
		},
		{
			name:     "close to overflow",
			duration: 250 * 365 * 24 * time.Hour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.duration < 0 {
				t.Skip("Duration overflows time.Duration")
				return
			}

			result := NewDurationType(tt.duration)
			resultStr := string(*result)

			// Should use days, not seconds-only
			assert.Contains(t, resultStr, "D", "Should contain day component")
			assert.NotRegexp(t, "^PT\\d+S$", resultStr, "Large durations should not be seconds-only")

			// Should round-trip exactly
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err)
			assert.Equal(t, tt.duration, parsedBack, "Round-trip should be exact")
		})
	}
}

// TestNewDurationTypeRoundTrip tests round-trip consistency
func TestNewDurationTypeRoundTrip(t *testing.T) {
	// Test durations that should round-trip with high accuracy
	exactDurations := []time.Duration{
		1 * time.Second,
		30 * time.Second,
		5 * time.Minute,
		2 * time.Hour,
		6 * time.Hour,
		12 * time.Hour,
		1 * 24 * time.Hour,  // 1 day
		3 * 24 * time.Hour,  // 3 days
		7 * 24 * time.Hour,  // 1 week (in days)
		14 * 24 * time.Hour, // 2 weeks
		28 * 24 * time.Hour, // 4 weeks (close to month)
	}

	for _, original := range exactDurations {
		t.Run(fmt.Sprintf("round_trip_%v", original), func(t *testing.T) {
			// Format to ISO 8601
			durType := NewDurationType(original)
			isoStr := string(*durType)

			// Parse back
			parsed, err := durType.GetTimeDuration()
			assert.NoError(t, err)

			// All durations should round-trip exactly (only days/hours/minutes/seconds used)
			assert.Equal(t, original, parsed,
				"Round-trip should be exact for duration %v, got %v (iso: %s)",
				original, parsed, isoStr)
		})
	}
}

// TestNewDurationTypeStructurePreservation tests that structure is preserved vs old behavior
func TestNewDurationTypeStructurePreservation(t *testing.T) {
	// Test cases that would have been "PT...S" in the old implementation
	testCases := []struct {
		name           string
		inputSeconds   int64
		mustContain    []string
		mustNotContain []string
	}{
		{
			name:           "1 year in seconds",
			inputSeconds:   31556952, // ~1 year
			mustContain:    []string{"D"},
			mustNotContain: []string{"PT31556952S"},
		},
		{
			name:           "1 month in seconds",
			inputSeconds:   2629746, // ~1 month
			mustContain:    []string{"D"},
			mustNotContain: []string{"PT2629746S"},
		},
		{
			name:           "issue 60 duration",
			inputSeconds:   4357512417,
			mustContain:    []string{"D"},
			mustNotContain: []string{"PT4357512417S"},
		},
		{
			name:           "6 months in seconds",
			inputSeconds:   15778476, // ~6 months
			mustContain:    []string{"D"},
			mustNotContain: []string{"PT15778476S"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			duration := time.Duration(tc.inputSeconds) * time.Second
			result := NewDurationType(duration)
			resultStr := string(*result)

			for _, mustHave := range tc.mustContain {
				assert.Contains(t, resultStr, mustHave,
					"Result should contain %s: %s", mustHave, resultStr)
			}

			for _, mustNotHave := range tc.mustNotContain {
				assert.NotEqual(t, mustNotHave, resultStr,
					"Result should not be the old seconds-only format")
			}

			// Verify it's valid ISO 8601
			assert.Regexp(t, "^-?P", resultStr, "Should start with P (or -P)")

			// Verify it can be parsed
			parsed, err := result.GetTimeDuration()
			assert.NoError(t, err)
			assert.NotZero(t, parsed)
		})
	}
}

// TestNewDurationTypeSPINERealistic tests realistic SPINE protocol durations
func TestNewDurationTypeSPINERealistic(t *testing.T) {
	// Real-world SPINE durations from actual usage
	spineUseCases := []struct {
		name     string
		duration time.Duration
		context  string
	}{
		{
			name:     "heartbeat timeout",
			duration: 4 * time.Second,
			context:  "Device heartbeat interval",
		},
		{
			name:     "response timeout",
			duration: 30 * time.Second,
			context:  "Maximum response delay",
		},
		{
			name:     "measurement interval",
			duration: 5 * time.Minute,
			context:  "Measurement reporting interval",
		},
		{
			name:     "charging session",
			duration: 4 * time.Hour,
			context:  "EV charging duration",
		},
		{
			name:     "daily schedule",
			duration: 24 * time.Hour,
			context:  "Daily energy schedule",
		},
		{
			name:     "weekly pattern",
			duration: 7 * 24 * time.Hour,
			context:  "Weekly load pattern",
		},
		{
			name:     "maintenance window",
			duration: 30 * 24 * time.Hour,
			context:  "Monthly maintenance",
		},
	}

	for _, tc := range spineUseCases {
		t.Run(tc.name, func(t *testing.T) {
			result := NewDurationType(tc.duration)
			resultStr := string(*result)

			// Should produce human-readable format
			assert.NotRegexp(t, "^PT\\d{4,}S$", resultStr,
				"SPINE durations should not be large second counts")

			// Should be parseable with high accuracy (SPINE needs precision)
			parsed, err := result.GetTimeDuration()
			assert.NoError(t, err)

			// All SPINE durations should round-trip exactly (days-based, no calendar approximation)
			assert.Equal(t, tc.duration, parsed,
				"SPINE duration %s should round-trip exactly", tc.context)
		})
	}
}
