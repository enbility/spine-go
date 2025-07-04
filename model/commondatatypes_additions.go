package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/rickb777/period"
)

// TimePeriodType

func NewTimePeriodTypeWithRelativeEndTime(duration time.Duration) *TimePeriodType {
	now := time.Now().UTC()
	endTime := now.Add(duration)
	value := &TimePeriodType{
		EndTime: NewAbsoluteOrRelativeTimeTypeFromTime(endTime),
	}
	return value
}

// helper type to modify EndTime field value in json.Marshal and
// json.Unmarshal to allow provide accurate relative durations
//
// If only EndTime is provided and it is a duration, it has to
// decrease over time. To do this without actually changing
// the data, it will always be transformed into an absolute time
// in Marshal and returned as an up to date relative duration
// in Unmarshal
type tempTimePeriodType TimePeriodType

func setTimePeriodTypeEndTime(t *tempTimePeriodType) {
	if t.StartTime != nil || t.EndTime == nil {
		return
	}

	duration, err := t.EndTime.GetTimeDuration()
	if err != nil {
		return
	}

	time := time.Now().UTC().Add(duration)
	t.EndTime = NewAbsoluteOrRelativeTimeTypeFromTime(time)
}

func getTimePeriodTypeDuration(t *TimePeriodType) (time.Duration, error) {
	if t.StartTime != nil || t.EndTime == nil {
		return 0, errors.New("invalid data format")
	}

	if t.EndTime.IsRelativeTime() {
		return getTimeDurationFromString(string(*t.EndTime))
	}

	endTime, err := t.EndTime.GetTime()
	if err != nil {
		return 0, err
	}

	now := time.Now().UTC()
	duration := endTime.Sub(now)
	duration = duration.Round(time.Second)

	if duration < 0 {
		return 0, nil
	}

	return duration, nil
}

// when startTime is empty and endTime is an absolute time,
// then endTime should be returned as an relative timestamp
func (t TimePeriodType) MarshalJSON() ([]byte, error) {
	temp := tempTimePeriodType(t)

	if duration, err := getTimePeriodTypeDuration(&t); err == nil {
		temp.EndTime = NewAbsoluteOrRelativeTimeTypeFromDuration(duration)
	}

	return json.Marshal(temp)
}

// when startTime is empty and endTime is a relative time,
// then endTime should be written as an absolute timestamp
func (t *TimePeriodType) UnmarshalJSON(data []byte) error {
	var temp tempTimePeriodType
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	setTimePeriodTypeEndTime(&temp)

	*t = TimePeriodType(temp)

	return nil
}

// Return the current duration if StartTime is nil and EndTime is not
// otherwise returns an error
func (t *TimePeriodType) GetDuration() (time.Duration, error) {
	return getTimePeriodTypeDuration(t)
}

// TimeType xs:time

func NewTimeType(t string) *TimeType {
	value := TimeType(t)
	return &value
}

func (s *TimeType) GetTime() (time.Time, error) {
	allowedFormats := []string{
		"15:04:05.999999999",
		"15:04:05.999999999Z",
		"15:04:05",
		"15:04:05Z",
		"15:04:05+07:00",
		"15:04:05-07:00",
	}

	for _, format := range allowedFormats {
		if value, err := time.ParseInLocation(format, string(*s), time.UTC); err == nil {
			return value, nil
		}
	}

	return time.Time{}, errors.New("unsupported time format")
}

// DateType xs:date

func NewDateType(t string) *DateType {
	value := DateType(t)
	return &value
}

// 2001-10-26, 2001-10-26+02:00, 2001-10-26Z, 2001-10-26+00:00, -2001-10-26, or -20000-04-01
func (d *DateType) GetTime() (time.Time, error) {
	allowedFormats := []string{
		"2006-01-02",
		"2006-01-02Z",
		"2006-01-02+07:00",
	}

	for _, format := range allowedFormats {
		if value, err := time.ParseInLocation(format, string(*d), time.UTC); err == nil {
			return value, nil
		}
	}

	return time.Time{}, errors.New("unsupported date format")
}

// DateTimeType xs:datetime

func NewDateTimeType(t string) *DateTimeType {
	value := DateTimeType(t)
	return &value
}

func NewDateTimeTypeFromTime(t time.Time) *DateTimeType {
	s := t.Round(time.Second).UTC().Format("2006-01-02T15:04:05Z")
	return NewDateTimeType(s)
}

func (d *DateTimeType) GetTime() (time.Time, error) {
	allowedFormats := []string{
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.999999999Z",
		"2006-01-02T15:04:05",
		"2006-01-02T15:04:05Z",
	}

	for _, format := range allowedFormats {
		if value, err := time.ParseInLocation(format, string(*d), time.UTC); err == nil {
			return value, nil
		}
	}

	return time.Time{}, errors.New("unsupported datetime format")
}

// DurationType

// IMPORTANT: Duration Parsing Limitations
//
// The period library used for parsing ISO 8601 durations (getTimeDurationFromString)
// uses fixed approximations that introduce errors for month and year components:
//   - 1 year ≈ 365.2425 days (actual: 365 or 366)
//   - 1 month ≈ 30.4369 days (actual: 28-31)
//
// Error Magnitude:
//   - NO ERRORS: Durations using only weeks, days, hours, minutes, seconds
//     Examples: P1W, P7D, PT24H, P1W2DT3H4M5S
//   - SIGNIFICANT ERRORS: Durations using months or years
//     Examples: P1M (error: 11-33 hours), P1Y (error: ~6 hours)
//
// For SPINE use cases (typically seconds to hours), this is not a concern.
// However, for monthly/yearly scheduling, use calendar-based calculations instead.
//
// See: https://github.com/enbility/spine-go/issues/60

func NewDurationType(duration time.Duration) *DurationType {
	// Handle negative durations
	if duration < 0 {
		// For negative durations, we need to work backwards
		positiveDuration := -duration
		result := NewDurationType(positiveDuration)
		negativeResult := "-" + string(*result)
		value := DurationType(negativeResult)
		return &value
	}

	// For relative durations, always calculate from "now" to preserve calendar structure
	// This gives us accurate year/month representation instead of just seconds
	now := time.Now()
	target := now.Add(duration)

	// Calculate calendar units between now and target
	years := 0
	months := 0
	days := 0
	hours := 0
	minutes := 0
	seconds := 0

	// Calculate years first
	for now.AddDate(years+1, 0, 0).Before(target) || now.AddDate(years+1, 0, 0).Equal(target) {
		years++
	}

	// Then months
	tempTime := now.AddDate(years, 0, 0)
	for tempTime.AddDate(0, months+1, 0).Before(target) || tempTime.AddDate(0, months+1, 0).Equal(target) {
		months++
	}

	// Then days
	tempTime = now.AddDate(years, months, 0)
	for tempTime.AddDate(0, 0, days+1).Before(target) || tempTime.AddDate(0, 0, days+1).Equal(target) {
		days++
	}

	// Now handle time components
	tempTime = now.AddDate(years, months, days)
	remainingDuration := target.Sub(tempTime)

	// Extract hours, minutes, seconds from remaining duration
	totalSeconds := int64(remainingDuration.Seconds())
	hours = int(totalSeconds / 3600)
	totalSeconds %= 3600
	minutes = int(totalSeconds / 60)
	seconds = int(totalSeconds % 60)

	// Handle nanoseconds for sub-second precision
	nanos := remainingDuration.Nanoseconds() % 1e9

	// Build ISO 8601 duration string
	var result strings.Builder
	result.WriteString("P")

	// Date part
	if years > 0 {
		result.WriteString(fmt.Sprintf("%dY", years))
	}
	if months > 0 {
		result.WriteString(fmt.Sprintf("%dM", months))
	}
	if days > 0 {
		result.WriteString(fmt.Sprintf("%dD", days))
	}

	// Time part
	if hours > 0 || minutes > 0 || seconds > 0 || nanos > 0 {
		result.WriteString("T")
		if hours > 0 {
			result.WriteString(fmt.Sprintf("%dH", hours))
		}
		if minutes > 0 {
			result.WriteString(fmt.Sprintf("%dM", minutes))
		}
		if seconds > 0 || nanos > 0 {
			if nanos > 0 {
				// Format seconds with fractional part
				fractionalSeconds := float64(seconds) + float64(nanos)/1e9
				result.WriteString(fmt.Sprintf("%gS", fractionalSeconds))
			} else {
				result.WriteString(fmt.Sprintf("%dS", seconds))
			}
		}
	}

	// Handle edge case of zero duration
	if result.String() == "P" {
		// ISO 8601 specifies P0D for zero duration, though PT0S is also valid
		result.WriteString("0D")
	}

	value := DurationType(result.String())
	return &value
}

func (d *DurationType) GetTimeDuration() (time.Duration, error) {
	return getTimeDurationFromString(string(*d))
}

// helper for DurationType and AbsoluteOrRelativeTimeType
//
// WARNING: This function uses period.DurationApprox() which has limitations:
//   - EXACT for: weeks, days, hours, minutes, seconds (P1W, P7D, PT1H)
//   - APPROXIMATE for: years, months (P1Y ≈ 365.2425 days, P1M ≈ 30.4369 days)
//
// The approximation errors for month/year durations can be significant:
//   - P1M: 11-33 hours error depending on actual month
//   - P1Y: ~6 hours error
//
// For precise calendar operations with months/years, use time.AddDate() instead.
func getTimeDurationFromString(s string) (time.Duration, error) {
	p, err := period.Parse(string(s))
	if err != nil {
		return 0, err
	}

	return p.DurationApprox(), nil
}

// AbsoluteOrRelativeTimeType
// can be of type TimeType or DurationType

func NewAbsoluteOrRelativeTimeType(s string) *AbsoluteOrRelativeTimeType {
	value := AbsoluteOrRelativeTimeType(s)
	return &value
}

func NewAbsoluteOrRelativeTimeTypeFromDuration(t time.Duration) *AbsoluteOrRelativeTimeType {
	s := NewDurationType(t)
	value := AbsoluteOrRelativeTimeType(*s)
	return &value
}

func NewAbsoluteOrRelativeTimeTypeFromTime(t time.Time) *AbsoluteOrRelativeTimeType {
	s := NewDateTimeTypeFromTime(t)
	value := AbsoluteOrRelativeTimeType(*s)
	return &value
}

func (a *AbsoluteOrRelativeTimeType) GetDateTimeType() *DateTimeType {
	value := NewDateTimeType(string(*a))
	return value
}

func (a *AbsoluteOrRelativeTimeType) GetTime() (time.Time, error) {
	value := NewDateTimeType(string(*a))
	t, err := value.GetTime()
	if err == nil {
		return t, nil
	}

	// Check if this is a relative time
	d, err := getTimeDurationFromString(string(*a))
	if err != nil {
		return time.Time{}, err
	}
	r := time.Now().Add(d)
	return r, nil
}

func (a *AbsoluteOrRelativeTimeType) IsRelativeTime() bool {
	_, err := getTimeDurationFromString(string(*a))
	return err == nil
}

func (a *AbsoluteOrRelativeTimeType) GetDurationType() (*DurationType, error) {
	value, err := a.GetTimeDuration()
	if err != nil {
		return nil, err
	}

	return NewDurationType(value), nil
}

func (a *AbsoluteOrRelativeTimeType) GetTimeDuration() (time.Duration, error) {
	return getTimeDurationFromString(string(*a))
}

// ScaledNumberType

func (m *ScaledNumberType) GetValue() float64 {
	if m.Number == nil {
		return 0
	}
	var scale float64 = 0
	if m.Scale != nil {
		scale = float64(*m.Scale)
	}
	return float64(*m.Number) * math.Pow(10, scale)
}

func NewScaledNumberType(value float64) *ScaledNumberType {
	m := &ScaledNumberType{}

	numberOfDecimals := 0
	temp := strconv.FormatFloat(value, 'f', -1, 64)
	index := strings.IndexByte(temp, '.')
	if index > -1 {
		numberOfDecimals = len(temp) - index - 1
	}

	// We limit this to 4 digits for now
	if numberOfDecimals > 4 {
		numberOfDecimals = 4
	}

	numberValue := NumberType(math.Trunc(value * math.Pow(10, float64(numberOfDecimals))))
	m.Number = &numberValue

	var scaleValue ScaleType
	scaleValue = ScaleType(0)
	if numberValue != 0 && -numberOfDecimals >= math.MinInt8 && -numberOfDecimals <= math.MaxInt8 {
		scaleValue = ScaleType(int8(-numberOfDecimals))
	}
	m.Scale = &scaleValue

	return m
}

// DeviceAddressType

var _ UpdateHelper = (*DeviceAddressType)(nil)

func (r *DeviceAddressType) String() string {
	if r == nil {
		return ""
	}

	var result = ""
	if r.Device != nil {
		result += string(*r.Device)
	}

	return result
}

// EntityAddressType

var _ UpdateHelper = (*EntityAddressType)(nil)

func (r *EntityAddressType) String() string {
	if r == nil {
		return ""
	}

	var result = ""
	if r.Device != nil {
		result += string(*r.Device)
	}
	result += ":["
	for index, id := range r.Entity {
		if index > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", id)
	}
	result += "]:"
	return result
}

// FeatureAddressType

var _ UpdateHelper = (*FeatureAddressType)(nil)

func (r *FeatureAddressType) String() string {
	if r == nil {
		return ""
	}

	var result = ""
	if r.Device != nil {
		result += string(*r.Device)
	}
	result += ":["
	for index, id := range r.Entity {
		if index > 0 {
			result += ","
		}
		result += fmt.Sprintf("%d", id)
	}
	result += "]:"
	if r.Feature != nil {
		result += fmt.Sprintf("%d", *r.Feature)
	}
	return result
}
