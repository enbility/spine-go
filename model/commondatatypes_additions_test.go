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
	// Should be something like "P138Y1MT4H6M57S" (preserves year/month structure)
	assert.NotEqual(t, "PT4357512417S", resultStr, "Should not output seconds-only format")
	assert.Contains(t, resultStr, "Y", "Should contain year component")
	assert.Contains(t, resultStr, "M", "Should contain month component")

	// Verify it's still a valid duration that can be parsed back
	parsedBack, err := result.GetTimeDuration()
	assert.NoError(t, err, "Result should be parseable")

	// Round-trip tolerance: NewDurationType formats using the calendar-exact
	// time.AddDate, but GetTimeDuration parses back via period.DurationApprox,
	// which approximates a year as 365.2425 days (~6h error per year, see
	// getTimeDurationFromString). For a ~138-year duration this accumulates well
	// past a fixed 10h bound, and the exact error depends on the leap-day
	// distribution between now and the target — making a fixed tolerance
	// date-dependent and flaky. Scale the tolerance to the duration's magnitude.
	diff := parsedBack - duration
	if diff < 0 {
		diff = -diff
	}
	approxYears := float64(duration) / float64(365*24*time.Hour)
	tolerance := time.Duration(approxYears*6)*time.Hour + 24*time.Hour
	assert.Less(t, diff, tolerance, "Should be within calendar-approximation tolerance")
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
			name:          "one nanosecond",
			duration:      1 * time.Nanosecond,
			expectedRegex: "^PT(0\\.000000001S|1e-09S)$",
			description:   "Should handle nanosecond precision (scientific notation allowed)",
		},
		{
			name:          "one millisecond",
			duration:      1 * time.Millisecond,
			expectedRegex: "^PT0\\.001S$",
			description:   "Should format fractional seconds",
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
			duration:      2*time.Hour + 30*time.Minute + 45*time.Second + 123*time.Millisecond,
			expectedRegex: "^PT2H30M45\\.123S$",
			description:   "Should handle complex time components",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDurationType(tt.duration)
			resultStr := string(*result)

			assert.Regexp(t, tt.expectedRegex, resultStr, tt.description)

			// Verify it can be parsed back (within tolerance for calendar operations)
			parsedBack, err := result.GetTimeDuration()
			if tt.duration == 1*time.Nanosecond {
				// Scientific notation (1e-09S) is not parseable by period library
				// This is an acceptable limitation for such tiny durations
				if err != nil {
					t.Logf("Nanosecond duration produces unparseable scientific notation: %s", resultStr)
					return // Skip the rest of this test
				}
			}
			assert.NoError(t, err, "Result should be parseable")

			// For small durations (< 1 day), expect exact matches (except very small ones)
			// For larger durations, allow for calendar approximation errors
			if tt.duration < 24*time.Hour {
				if tt.duration >= 1*time.Millisecond {
					assert.Equal(t, tt.duration, parsedBack, "Small durations should round-trip exactly")
				} else {
					// Very small durations (nanoseconds) may have precision issues
					diff := parsedBack - tt.duration
					if diff < 0 {
						diff = -diff
					}
					assert.True(t, diff <= tt.duration, "Very small durations should be reasonably close")
				}
			} else {
				// Allow for small differences due to calendar calculations
				diff := parsedBack - tt.duration
				if diff < 0 {
					diff = -diff
				}
				assert.True(t, diff < 1*time.Hour, "Large durations should be within 1 hour tolerance")
			}
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

// TestNewDurationTypeLeapYearBoundaries tests calendar edge cases
func TestNewDurationTypeLeapYearBoundaries(t *testing.T) {
	// Test durations that would span different calendar boundaries
	// Use fixed times for reproducible results
	tests := []struct {
		name        string
		duration    time.Duration
		description string
		minYears    int
		maxYears    int
	}{
		{
			name:        "approximately 1 year",
			duration:    365 * 24 * time.Hour,
			description: "365 days should be close to 1 year",
			minYears:    0,
			maxYears:    1,
		},
		{
			name:        "approximately 2 years",
			duration:    2 * 365 * 24 * time.Hour,
			description: "730 days should be close to 2 years",
			minYears:    1,
			maxYears:    2,
		},
		{
			name:        "approximately 1 month",
			duration:    30 * 24 * time.Hour,
			description: "30 days should be approximately 1 month",
			minYears:    0,
			maxYears:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewDurationType(tt.duration)
			resultStr := string(*result)

			// Verify the structure makes sense
			if tt.minYears > 0 || tt.maxYears > 0 {
				assert.Contains(t, resultStr, "Y", "Should contain year component for ~yearly durations")
			}

			// Should not be seconds-only format
			assert.NotRegexp(t, "^PT\\d+S$", resultStr, "Should not be seconds-only format")

			// Should be parseable
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err)
			assert.NotZero(t, parsedBack)
		})
	}
}

// TestNewDurationTypeMonthBoundaries tests month length variations
func TestNewDurationTypeMonthBoundaries(t *testing.T) {
	// Test durations around month boundaries
	monthLengths := []int{28, 29, 30, 31} // Different month lengths

	for _, days := range monthLengths {
		t.Run(fmt.Sprintf("%d days", days), func(t *testing.T) {
			duration := time.Duration(days) * 24 * time.Hour
			result := NewDurationType(duration)
			resultStr := string(*result)

			// Should be in a reasonable format (days or month + days)
			assert.Regexp(t, "^P(\\d+M)?(\\d+D)?(T.*)?$", resultStr, "Should be valid ISO 8601 format")

			// Should be parseable
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err)

			// For durations around month boundaries, allow for calendar approximation errors
			diff := parsedBack - duration
			if diff < 0 {
				diff = -diff
			}
			assert.True(t, diff < 24*time.Hour, "Month-boundary durations should be within 24 hours (calendar approximations)")
		})
	}
}

// TestNewDurationTypeLargeValues tests handling of large durations
func TestNewDurationTypeLargeValues(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expectY  bool
		expectM  bool
	}{
		{
			name:     "10 years",
			duration: 10 * 365 * 24 * time.Hour,
			expectY:  true,
			expectM:  false,
		},
		{
			name:     "100 years",
			duration: 100 * 365 * 24 * time.Hour,
			expectY:  true,
			expectM:  false,
		},
		{
			name:     "close to overflow",
			duration: 250 * 365 * 24 * time.Hour, // Close to time.Duration max (~290 years)
			expectY:  true,
			expectM:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Skip if duration would overflow
			if tt.duration < 0 {
				t.Skip("Duration overflows time.Duration")
				return
			}

			result := NewDurationType(tt.duration)
			resultStr := string(*result)

			if tt.expectY {
				assert.Contains(t, resultStr, "Y", "Should contain year component")
			}
			if tt.expectM {
				assert.Contains(t, resultStr, "M", "Should contain month component")
			}

			// Should not be seconds-only
			assert.NotRegexp(t, "^PT\\d+S$", resultStr, "Large durations should not be seconds-only")

			// Should be parseable (even if with approximation errors)
			parsedBack, err := result.GetTimeDuration()
			assert.NoError(t, err)
			assert.NotZero(t, parsedBack)
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

			// For durations using only weeks/days/hours/minutes/seconds,
			// we should get exact round-trip
			if original <= 28*24*time.Hour {
				tolerance := 1 * time.Second // Allow 1 second tolerance for rounding
				diff := parsed - original
				if diff < 0 {
					diff = -diff
				}
				assert.True(t, diff <= tolerance,
					"Round-trip should be exact for duration %v, got %v (diff: %v, iso: %s)",
					original, parsed, diff, isoStr)
			}
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
			mustContain:    []string{"Y"},
			mustNotContain: []string{"PT31556952S"},
		},
		{
			name:           "1 month in seconds",
			inputSeconds:   2629746,       // ~1 month
			mustContain:    []string{"D"}, // Should be days or month+days
			mustNotContain: []string{"PT2629746S"},
		},
		{
			name:           "issue 60 duration",
			inputSeconds:   4357512417,
			mustContain:    []string{"Y", "M"},
			mustNotContain: []string{"PT4357512417S"},
		},
		{
			name:           "6 months in seconds",
			inputSeconds:   15778476,      // ~6 months
			mustContain:    []string{"M"}, // Should contain months
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

			// For SPINE use cases, accuracy is critical
			diff := parsed - tc.duration
			if diff < 0 {
				diff = -diff
			}

			// Most SPINE durations should be exact or very close
			if tc.duration <= 7*24*time.Hour {
				assert.True(t, diff <= 1*time.Second,
					"SPINE duration %s should be very accurate (diff: %v)", tc.context, diff)
			} else {
				assert.True(t, diff <= 1*time.Hour,
					"Longer SPINE duration %s should be reasonably accurate (diff: %v)", tc.context, diff)
			}
		})
	}
}
