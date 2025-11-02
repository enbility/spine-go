package model

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestPowerTimeSlotScheduleListDataType_Update(t *testing.T) {
	sut := PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: []PowerTimeSlotScheduleDataType{
			{
				SequenceId:    util.Ptr(PowerSequenceIdType(0)),
				SlotActivated: util.Ptr(false),
			},
			{
				SequenceId:    util.Ptr(PowerSequenceIdType(1)),
				SlotActivated: util.Ptr(false),
			},
		},
	}

	newData := PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: []PowerTimeSlotScheduleDataType{
			{
				SequenceId:    util.Ptr(PowerSequenceIdType(1)),
				SlotActivated: util.Ptr(true),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerTimeSlotScheduleData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, false, *item1.SlotActivated)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, true, *item2.SlotActivated)
}

func TestPowerTimeSlotValueListDataType_Update(t *testing.T) {
	sut := PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: []PowerTimeSlotValueDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Value:      NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Value:      NewScaledNumberType(1),
			},
		},
	}

	newData := PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: []PowerTimeSlotValueDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Value:      NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerTimeSlotValueData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.Value.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.Value.GetValue())
}

func TestPowerTimeSlotScheduleConstraintsListDataType_Update(t *testing.T) {
	sut := PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: []PowerTimeSlotScheduleConstraintsDataType{
			{
				SequenceId:  util.Ptr(PowerSequenceIdType(0)),
				MinDuration: NewDurationType(1 * time.Second),
			},
			{
				SequenceId:  util.Ptr(PowerSequenceIdType(1)),
				MinDuration: NewDurationType(1 * time.Second),
			},
		},
	}

	newData := PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: []PowerTimeSlotScheduleConstraintsDataType{
			{
				SequenceId:  util.Ptr(PowerSequenceIdType(1)),
				MinDuration: NewDurationType(10 * time.Second),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerTimeSlotScheduleConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	duration, _ := item1.MinDuration.GetTimeDuration()
	assert.Equal(t, time.Duration(1*time.Second), duration)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	duration, _ = item2.MinDuration.GetTimeDuration()
	assert.Equal(t, time.Duration(10*time.Second), duration)
}

func TestPowerSequenceAlternativesRelationListDataType_Update(t *testing.T) {
	sut := PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: []PowerSequenceAlternativesRelationDataType{
			{
				AlternativesId: util.Ptr(AlternativesIdType(0)),
				SequenceId:     []PowerSequenceIdType{0},
			},
			{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
				SequenceId:     []PowerSequenceIdType{0},
			},
		},
	}

	newData := PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: []PowerSequenceAlternativesRelationDataType{
			{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
				SequenceId:     []PowerSequenceIdType{1},
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerSequenceAlternativesRelationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.AlternativesId))
	assert.Equal(t, 0, int(item1.SequenceId[0]))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.AlternativesId))
	assert.Equal(t, 1, int(item2.SequenceId[0]))
}

func TestPowerSequenceDescriptionListDataType_Update(t *testing.T) {
	sut := PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: []PowerSequenceDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(0)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
			},
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
			},
		},
	}

	newData := PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: []PowerSequenceDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeProduce),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerSequenceDescriptionData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, EnergyDirectionTypeConsume, *item1.PositiveEnergyDirection)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, EnergyDirectionTypeProduce, *item2.PositiveEnergyDirection)
}

func TestPowerSequenceStateListDataType_Update(t *testing.T) {
	sut := PowerSequenceStateListDataType{
		PowerSequenceStateData: []PowerSequenceStateDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				State:      util.Ptr(PowerSequenceStateTypeRunning),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				State:      util.Ptr(PowerSequenceStateTypeRunning),
			},
		},
	}

	newData := PowerSequenceStateListDataType{
		PowerSequenceStateData: []PowerSequenceStateDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				State:      util.Ptr(PowerSequenceStateTypeCompleted),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerSequenceStateData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, PowerSequenceStateTypeRunning, *item1.State)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, PowerSequenceStateTypeCompleted, *item2.State)
}

func TestPowerSequenceScheduleListDataType_Update(t *testing.T) {
	sut := PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: []PowerSequenceScheduleDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				EndTime:    NewAbsoluteOrRelativeTimeType("PT2H"),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				EndTime:    NewAbsoluteOrRelativeTimeType("PT2H"),
			},
		},
	}

	newData := PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: []PowerSequenceScheduleDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				EndTime:    NewAbsoluteOrRelativeTimeType("PT4H"),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerSequenceScheduleData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, "PT2H", string(*item1.EndTime))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, "PT4H", string(*item2.EndTime))
}

func TestPowerSequenceScheduleConstraintsListDataType_Update(t *testing.T) {
	sut := PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: []PowerSequenceScheduleConstraintsDataType{
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(0)),
				EarliestEndTime: NewAbsoluteOrRelativeTimeType("PT2H"),
			},
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(1)),
				EarliestEndTime: NewAbsoluteOrRelativeTimeType("PT2H"),
			},
		},
	}

	newData := PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: []PowerSequenceScheduleConstraintsDataType{
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(1)),
				EarliestEndTime: NewAbsoluteOrRelativeTimeType("PT4H"),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerSequenceScheduleConstraintsData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, "PT2H", string(*item1.EarliestEndTime))
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, "PT4H", string(*item2.EarliestEndTime))
}

func TestPowerSequencePriceListDataType_Update(t *testing.T) {
	sut := PowerSequencePriceListDataType{
		PowerSequencePriceData: []PowerSequencePriceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Price:      NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Price:      NewScaledNumberType(1),
			},
		},
	}

	newData := PowerSequencePriceListDataType{
		PowerSequencePriceData: []PowerSequencePriceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Price:      NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerSequencePriceData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.Price.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.Price.GetValue())
}

func TestPowerSequenceSchedulePreferenceListDataType_Update(t *testing.T) {
	sut := PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: []PowerSequenceSchedulePreferenceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Cheapest:   util.Ptr(false),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Cheapest:   util.Ptr(false),
			},
		},
	}

	newData := PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: []PowerSequenceSchedulePreferenceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Cheapest:   util.Ptr(true),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.PowerSequenceSchedulePreferenceData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, false, *item1.Cheapest)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, true, *item2.Cheapest)
}

func TestPowerTimeSlotScheduleListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerTimeSlotScheduleListDataType{
		PowerTimeSlotScheduleData: []PowerTimeSlotScheduleDataType{
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(1)),
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(1)),
				DefaultDuration: NewDurationType(time.Hour),
				SlotActivated:   util.Ptr(true),
				Description:     util.Ptr(DescriptionType("Schedule 1")),
			},
			{
				SequenceId:      util.Ptr(PowerSequenceIdType(2)),
				SlotNumber:      util.Ptr(PowerTimeSlotNumberType(2)),
				DefaultDuration: NewDurationType(2 * time.Hour),
				SlotActivated:   util.Ptr(false),
				Description:     util.Ptr(DescriptionType("Schedule 2")),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerTimeSlotScheduleListDataType)
	assert.Equal(t, 2, len(resultData.PowerTimeSlotScheduleData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerTimeSlotScheduleListDataSelectors: &PowerTimeSlotScheduleListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleListDataType)
	assert.Equal(t, 1, len(resultData.PowerTimeSlotScheduleData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerTimeSlotScheduleData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerTimeSlotScheduleDataElements: &PowerTimeSlotScheduleDataElementsType{
			SequenceId:    &ElementTagType{},
			SlotActivated: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleListDataType)
	assert.Equal(t, 2, len(resultData.PowerTimeSlotScheduleData))
	assert.NotNil(t, resultData.PowerTimeSlotScheduleData[0].SequenceId)
	assert.NotNil(t, resultData.PowerTimeSlotScheduleData[0].SlotActivated)
	assert.Nil(t, resultData.PowerTimeSlotScheduleData[0].SlotNumber)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerTimeSlotScheduleListDataSelectors: &PowerTimeSlotScheduleListDataSelectorsType{
			SlotNumber: util.Ptr(PowerTimeSlotNumberType(2)),
		},
		PowerTimeSlotScheduleDataElements: &PowerTimeSlotScheduleDataElementsType{
			SequenceId:  &ElementTagType{},
			Description: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleListDataType)
	assert.Equal(t, 1, len(resultData.PowerTimeSlotScheduleData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerTimeSlotScheduleData[0].SequenceId)
	assert.NotNil(t, resultData.PowerTimeSlotScheduleData[0].Description)
	assert.Nil(t, resultData.PowerTimeSlotScheduleData[0].SlotActivated)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerTimeSlotScheduleListDataSelectors: &PowerTimeSlotScheduleListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleListDataType)
	assert.Equal(t, 0, len(resultData.PowerTimeSlotScheduleData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerTimeSlotValueListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerTimeSlotValueListDataType{
		PowerTimeSlotValueData: []PowerTimeSlotValueDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				SlotNumber: util.Ptr(PowerTimeSlotNumberType(1)),
				ValueType:  util.Ptr(PowerTimeSlotValueTypeType("power")),
				Value:      NewScaledNumberType(100.0),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(2)),
				SlotNumber: util.Ptr(PowerTimeSlotNumberType(2)),
				ValueType:  util.Ptr(PowerTimeSlotValueTypeType("energy")),
				Value:      NewScaledNumberType(250.0),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerTimeSlotValueListDataType)
	assert.Equal(t, 2, len(resultData.PowerTimeSlotValueData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerTimeSlotValueListDataSelectors: &PowerTimeSlotValueListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotValueListDataType)
	assert.Equal(t, 1, len(resultData.PowerTimeSlotValueData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerTimeSlotValueData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerTimeSlotValueDataElements: &PowerTimeSlotValueDataElementsType{
			SequenceId: &ElementTagType{},
			ValueType:  &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotValueListDataType)
	assert.Equal(t, 2, len(resultData.PowerTimeSlotValueData))
	assert.NotNil(t, resultData.PowerTimeSlotValueData[0].SequenceId)
	assert.NotNil(t, resultData.PowerTimeSlotValueData[0].ValueType)
	assert.Nil(t, resultData.PowerTimeSlotValueData[0].SlotNumber)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerTimeSlotValueListDataSelectors: &PowerTimeSlotValueListDataSelectorsType{
			ValueType: util.Ptr(PowerTimeSlotValueTypeType("energy")),
		},
		PowerTimeSlotValueDataElements: &PowerTimeSlotValueDataElementsType{
			SequenceId: &ElementTagType{},
			Value:      &ScaledNumberElementsType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotValueListDataType)
	assert.Equal(t, 1, len(resultData.PowerTimeSlotValueData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerTimeSlotValueData[0].SequenceId)
	assert.NotNil(t, resultData.PowerTimeSlotValueData[0].Value)
	assert.Nil(t, resultData.PowerTimeSlotValueData[0].ValueType)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerTimeSlotValueListDataSelectors: &PowerTimeSlotValueListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotValueListDataType)
	assert.Equal(t, 0, len(resultData.PowerTimeSlotValueData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerTimeSlotScheduleConstraintsListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerTimeSlotScheduleConstraintsListDataType{
		PowerTimeSlotScheduleConstraintsData: []PowerTimeSlotScheduleConstraintsDataType{
			{
				SequenceId:   util.Ptr(PowerSequenceIdType(1)),
				SlotNumber:   util.Ptr(PowerTimeSlotNumberType(1)),
				MinDuration:  NewDurationType(time.Hour),
				MaxDuration:  NewDurationType(2 * time.Hour),
				OptionalSlot: util.Ptr(false),
			},
			{
				SequenceId:   util.Ptr(PowerSequenceIdType(2)),
				SlotNumber:   util.Ptr(PowerTimeSlotNumberType(2)),
				MinDuration:  NewDurationType(30 * time.Minute),
				MaxDuration:  NewDurationType(90 * time.Minute),
				OptionalSlot: util.Ptr(true),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerTimeSlotScheduleConstraintsListDataType)
	assert.Equal(t, 2, len(resultData.PowerTimeSlotScheduleConstraintsData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerTimeSlotScheduleConstraintsListDataSelectors: &PowerTimeSlotScheduleConstraintsListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleConstraintsListDataType)
	assert.Equal(t, 1, len(resultData.PowerTimeSlotScheduleConstraintsData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerTimeSlotScheduleConstraintsData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerTimeSlotScheduleConstraintsDataElements: &PowerTimeSlotScheduleConstraintsDataElementsType{
			SequenceId:   &ElementTagType{},
			OptionalSlot: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleConstraintsListDataType)
	assert.Equal(t, 2, len(resultData.PowerTimeSlotScheduleConstraintsData))
	assert.NotNil(t, resultData.PowerTimeSlotScheduleConstraintsData[0].SequenceId)
	assert.NotNil(t, resultData.PowerTimeSlotScheduleConstraintsData[0].OptionalSlot)
	assert.Nil(t, resultData.PowerTimeSlotScheduleConstraintsData[0].SlotNumber)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerTimeSlotScheduleConstraintsListDataSelectors: &PowerTimeSlotScheduleConstraintsListDataSelectorsType{
			SlotNumber: util.Ptr(PowerTimeSlotNumberType(2)),
		},
		PowerTimeSlotScheduleConstraintsDataElements: &PowerTimeSlotScheduleConstraintsDataElementsType{
			SequenceId:  &ElementTagType{},
			MinDuration: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleConstraintsListDataType)
	assert.Equal(t, 1, len(resultData.PowerTimeSlotScheduleConstraintsData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerTimeSlotScheduleConstraintsData[0].SequenceId)
	assert.NotNil(t, resultData.PowerTimeSlotScheduleConstraintsData[0].MinDuration)
	assert.Nil(t, resultData.PowerTimeSlotScheduleConstraintsData[0].OptionalSlot)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerTimeSlotScheduleConstraintsListDataSelectors: &PowerTimeSlotScheduleConstraintsListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerTimeSlotScheduleConstraintsListDataType)
	assert.Equal(t, 0, len(resultData.PowerTimeSlotScheduleConstraintsData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerSequenceAlternativesRelationListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerSequenceAlternativesRelationListDataType{
		PowerSequenceAlternativesRelationData: []PowerSequenceAlternativesRelationDataType{
			{
				AlternativesId: util.Ptr(AlternativesIdType(1)),
				SequenceId:     []PowerSequenceIdType{1, 2},
			},
			{
				AlternativesId: util.Ptr(AlternativesIdType(2)),
				SequenceId:     []PowerSequenceIdType{3, 4, 5},
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerSequenceAlternativesRelationListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceAlternativesRelationData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerSequenceAlternativesRelationListDataSelectors: &PowerSequenceAlternativesRelationListDataSelectorsType{
			AlternativesId: util.Ptr(AlternativesIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceAlternativesRelationListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceAlternativesRelationData))
	assert.Equal(t, AlternativesIdType(1), *resultData.PowerSequenceAlternativesRelationData[0].AlternativesId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerSequenceAlternativesRelationDataElements: &PowerSequenceAlternativesRelationDataElementsType{
			AlternativesId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceAlternativesRelationListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceAlternativesRelationData))
	assert.NotNil(t, resultData.PowerSequenceAlternativesRelationData[0].AlternativesId)
	assert.Nil(t, resultData.PowerSequenceAlternativesRelationData[0].SequenceId)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerSequenceAlternativesRelationListDataSelectors: &PowerSequenceAlternativesRelationListDataSelectorsType{
			AlternativesId: util.Ptr(AlternativesIdType(2)),
		},
		PowerSequenceAlternativesRelationDataElements: &PowerSequenceAlternativesRelationDataElementsType{
			AlternativesId: &ElementTagType{},
			SequenceId:     &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceAlternativesRelationListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceAlternativesRelationData))
	assert.Equal(t, AlternativesIdType(2), *resultData.PowerSequenceAlternativesRelationData[0].AlternativesId)
	assert.NotNil(t, resultData.PowerSequenceAlternativesRelationData[0].SequenceId)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerSequenceAlternativesRelationListDataSelectors: &PowerSequenceAlternativesRelationListDataSelectorsType{
			AlternativesId: util.Ptr(AlternativesIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceAlternativesRelationListDataType)
	assert.Equal(t, 0, len(resultData.PowerSequenceAlternativesRelationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerSequenceDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerSequenceDescriptionListDataType{
		PowerSequenceDescriptionData: []PowerSequenceDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				Description:             util.Ptr(DescriptionType("Power Sequence 1")),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
				PowerUnit:               util.Ptr(UnitOfMeasurementTypeW),
				EnergyUnit:              util.Ptr(UnitOfMeasurementTypeWh),
				TaskIdentifier:          util.Ptr(uint(101)),
				RepetitionsTotal:        util.Ptr(uint(5)),
			},
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(2)),
				Description:             util.Ptr(DescriptionType("Power Sequence 2")),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeProduce),
				PowerUnit:               util.Ptr(UnitOfMeasurementTypeV),
				EnergyUnit:              util.Ptr(UnitOfMeasurementTypeVAh),
				TaskIdentifier:          util.Ptr(uint(202)),
				RepetitionsTotal:        util.Ptr(uint(10)),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerSequenceDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceDescriptionData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerSequenceDescriptionListDataSelectors: &PowerSequenceDescriptionListDataSelectorsType{
			SequenceId: []PowerSequenceIdType{1},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceDescriptionData))
	if len(resultData.PowerSequenceDescriptionData) > 0 {
		assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerSequenceDescriptionData[0].SequenceId)
	}

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerSequenceDescriptionDataElements: &PowerSequenceDescriptionDataElementsType{
			SequenceId:  &ElementTagType{},
			Description: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceDescriptionData))
	assert.NotNil(t, resultData.PowerSequenceDescriptionData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequenceDescriptionData[0].Description)
	assert.Nil(t, resultData.PowerSequenceDescriptionData[0].PowerUnit)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerSequenceDescriptionListDataSelectors: &PowerSequenceDescriptionListDataSelectorsType{
			SequenceId: []PowerSequenceIdType{2},
		},
		PowerSequenceDescriptionDataElements: &PowerSequenceDescriptionDataElementsType{
			SequenceId:       &ElementTagType{},
			TaskIdentifier:   &ElementTagType{},
			RepetitionsTotal: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceDescriptionData))
	if len(resultData.PowerSequenceDescriptionData) > 0 {
		assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerSequenceDescriptionData[0].SequenceId)
		assert.NotNil(t, resultData.PowerSequenceDescriptionData[0].TaskIdentifier)
		assert.Nil(t, resultData.PowerSequenceDescriptionData[0].Description)
	}

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerSequenceDescriptionListDataSelectors: &PowerSequenceDescriptionListDataSelectorsType{
			SequenceId: []PowerSequenceIdType{99},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceDescriptionListDataType)
	assert.Equal(t, 0, len(resultData.PowerSequenceDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerSequenceStateListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerSequenceStateListDataType{
		PowerSequenceStateData: []PowerSequenceStateDataType{
			{
				SequenceId:                 util.Ptr(PowerSequenceIdType(1)),
				State:                      util.Ptr(PowerSequenceStateType("running")),
				ActiveSlotNumber:           util.Ptr(PowerTimeSlotNumberType(1)),
				ElapsedSlotTime:            NewDurationType(30 * time.Minute),
				RemainingSlotTime:          NewDurationType(90 * time.Minute),
				SequenceRemoteControllable: util.Ptr(true),
				ActiveRepetitionNumber:     util.Ptr(uint(2)),
			},
			{
				SequenceId:                 util.Ptr(PowerSequenceIdType(2)),
				State:                      util.Ptr(PowerSequenceStateType("paused")),
				ActiveSlotNumber:           util.Ptr(PowerTimeSlotNumberType(3)),
				ElapsedSlotTime:            NewDurationType(15 * time.Minute),
				RemainingSlotTime:          NewDurationType(45 * time.Minute),
				SequenceRemoteControllable: util.Ptr(false),
				ActiveRepetitionNumber:     util.Ptr(uint(1)),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerSequenceStateListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceStateData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerSequenceStateListDataSelectors: &PowerSequenceStateListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceStateListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceStateData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerSequenceStateData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerSequenceStateDataElements: &PowerSequenceStateDataElementsType{
			SequenceId: &ElementTagType{},
			State:      &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceStateListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceStateData))
	assert.NotNil(t, resultData.PowerSequenceStateData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequenceStateData[0].State)
	assert.Nil(t, resultData.PowerSequenceStateData[0].ActiveSlotNumber)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerSequenceStateListDataSelectors: &PowerSequenceStateListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		PowerSequenceStateDataElements: &PowerSequenceStateDataElementsType{
			SequenceId:                 &ElementTagType{},
			SequenceRemoteControllable: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceStateListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceStateData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerSequenceStateData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequenceStateData[0].SequenceRemoteControllable)
	assert.Nil(t, resultData.PowerSequenceStateData[0].State)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerSequenceStateListDataSelectors: &PowerSequenceStateListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceStateListDataType)
	assert.Equal(t, 0, len(resultData.PowerSequenceStateData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerSequenceScheduleListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerSequenceScheduleListDataType{
		PowerSequenceScheduleData: []PowerSequenceScheduleDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(2)),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerSequenceScheduleListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceScheduleData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerSequenceScheduleListDataSelectors: &PowerSequenceScheduleListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceScheduleData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerSequenceScheduleData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerSequenceScheduleDataElements: &PowerSequenceScheduleDataElementsType{
			SequenceId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceScheduleData))
	assert.NotNil(t, resultData.PowerSequenceScheduleData[0].SequenceId)
	assert.Nil(t, resultData.PowerSequenceScheduleData[0].StartTime)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerSequenceScheduleListDataSelectors: &PowerSequenceScheduleListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		PowerSequenceScheduleDataElements: &PowerSequenceScheduleDataElementsType{
			SequenceId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceScheduleData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerSequenceScheduleData[0].SequenceId)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerSequenceScheduleListDataSelectors: &PowerSequenceScheduleListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleListDataType)
	assert.Equal(t, 0, len(resultData.PowerSequenceScheduleData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerSequenceScheduleConstraintsListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerSequenceScheduleConstraintsListDataType{
		PowerSequenceScheduleConstraintsData: []PowerSequenceScheduleConstraintsDataType{
			{
				SequenceId:       util.Ptr(PowerSequenceIdType(1)),
				OptionalSequence: util.Ptr(false),
			},
			{
				SequenceId:       util.Ptr(PowerSequenceIdType(2)),
				OptionalSequence: util.Ptr(true),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerSequenceScheduleConstraintsListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceScheduleConstraintsData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerSequenceScheduleConstraintsListDataSelectors: &PowerSequenceScheduleConstraintsListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleConstraintsListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceScheduleConstraintsData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerSequenceScheduleConstraintsData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerSequenceScheduleConstraintsDataElements: &PowerSequenceScheduleConstraintsDataElementsType{
			SequenceId:       &ElementTagType{},
			OptionalSequence: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleConstraintsListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceScheduleConstraintsData))
	assert.NotNil(t, resultData.PowerSequenceScheduleConstraintsData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequenceScheduleConstraintsData[0].OptionalSequence)
	assert.Nil(t, resultData.PowerSequenceScheduleConstraintsData[0].EarliestStartTime)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerSequenceScheduleConstraintsListDataSelectors: &PowerSequenceScheduleConstraintsListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		PowerSequenceScheduleConstraintsDataElements: &PowerSequenceScheduleConstraintsDataElementsType{
			SequenceId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleConstraintsListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceScheduleConstraintsData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerSequenceScheduleConstraintsData[0].SequenceId)
	assert.Nil(t, resultData.PowerSequenceScheduleConstraintsData[0].OptionalSequence)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerSequenceScheduleConstraintsListDataSelectors: &PowerSequenceScheduleConstraintsListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceScheduleConstraintsListDataType)
	assert.Equal(t, 0, len(resultData.PowerSequenceScheduleConstraintsData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerSequencePriceListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerSequencePriceListDataType{
		PowerSequencePriceData: []PowerSequencePriceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Price:      NewScaledNumberType(25.5),
				Currency:   util.Ptr(CurrencyType("EUR")),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(2)),
				Price:      NewScaledNumberType(35.75),
				Currency:   util.Ptr(CurrencyType("USD")),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerSequencePriceListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequencePriceData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerSequencePriceListDataSelectors: &PowerSequencePriceListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequencePriceListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequencePriceData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerSequencePriceData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerSequencePriceDataElements: &PowerSequencePriceDataElementsType{
			SequenceId: &ElementTagType{},
			Currency:   &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequencePriceListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequencePriceData))
	assert.NotNil(t, resultData.PowerSequencePriceData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequencePriceData[0].Currency)
	assert.Nil(t, resultData.PowerSequencePriceData[0].Price)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerSequencePriceListDataSelectors: &PowerSequencePriceListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		PowerSequencePriceDataElements: &PowerSequencePriceDataElementsType{
			SequenceId: &ElementTagType{},
			Price:      &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequencePriceListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequencePriceData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerSequencePriceData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequencePriceData[0].Price)
	assert.Nil(t, resultData.PowerSequencePriceData[0].Currency)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerSequencePriceListDataSelectors: &PowerSequencePriceListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequencePriceListDataType)
	assert.Equal(t, 0, len(resultData.PowerSequencePriceData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestPowerSequenceSchedulePreferenceListDataType_ReadPartialData(t *testing.T) {
	sut := &PowerSequenceSchedulePreferenceListDataType{
		PowerSequenceSchedulePreferenceData: []PowerSequenceSchedulePreferenceDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Greenest:   util.Ptr(true),
				Cheapest:   util.Ptr(false),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(2)),
				Greenest:   util.Ptr(false),
				Cheapest:   util.Ptr(true),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*PowerSequenceSchedulePreferenceListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceSchedulePreferenceData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		PowerSequenceSchedulePreferenceListDataSelectors: &PowerSequenceSchedulePreferenceListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceSchedulePreferenceListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceSchedulePreferenceData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.PowerSequenceSchedulePreferenceData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		PowerSequenceSchedulePreferenceDataElements: &PowerSequenceSchedulePreferenceDataElementsType{
			SequenceId: &ElementTagType{},
			Greenest:   &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceSchedulePreferenceListDataType)
	assert.Equal(t, 2, len(resultData.PowerSequenceSchedulePreferenceData))
	assert.NotNil(t, resultData.PowerSequenceSchedulePreferenceData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequenceSchedulePreferenceData[0].Greenest)
	assert.Nil(t, resultData.PowerSequenceSchedulePreferenceData[0].Cheapest)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		PowerSequenceSchedulePreferenceListDataSelectors: &PowerSequenceSchedulePreferenceListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		PowerSequenceSchedulePreferenceDataElements: &PowerSequenceSchedulePreferenceDataElementsType{
			SequenceId: &ElementTagType{},
			Cheapest:   &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceSchedulePreferenceListDataType)
	assert.Equal(t, 1, len(resultData.PowerSequenceSchedulePreferenceData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.PowerSequenceSchedulePreferenceData[0].SequenceId)
	assert.NotNil(t, resultData.PowerSequenceSchedulePreferenceData[0].Cheapest)
	assert.Nil(t, resultData.PowerSequenceSchedulePreferenceData[0].Greenest)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		PowerSequenceSchedulePreferenceListDataSelectors: &PowerSequenceSchedulePreferenceListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*PowerSequenceSchedulePreferenceListDataType)
	assert.Equal(t, 0, len(resultData.PowerSequenceSchedulePreferenceData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
