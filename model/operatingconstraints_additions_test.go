package model

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func TestOperatingConstraintsInterruptListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: []OperatingConstraintsInterruptDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				IsPausable: util.Ptr(false),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				IsPausable: util.Ptr(false),
			},
		},
	}

	newData := OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: []OperatingConstraintsInterruptDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				IsPausable: util.Ptr(true),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.OperatingConstraintsInterruptData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, false, *item1.IsPausable)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, true, *item2.IsPausable)
}

func TestOperatingConstraintsDurationListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: []OperatingConstraintsDurationDataType{
			{
				SequenceId:        util.Ptr(PowerSequenceIdType(0)),
				ActiveDurationMin: NewDurationType(1 * time.Second),
			},
			{
				SequenceId:        util.Ptr(PowerSequenceIdType(1)),
				ActiveDurationMin: NewDurationType(1 * time.Second),
			},
		},
	}

	newData := OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: []OperatingConstraintsDurationDataType{
			{
				SequenceId:        util.Ptr(PowerSequenceIdType(1)),
				ActiveDurationMin: NewDurationType(10 * time.Second),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.OperatingConstraintsDurationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	duration, _ := item1.ActiveDurationMin.GetTimeDuration()
	assert.Equal(t, time.Duration(1*time.Second), duration)
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	duration, _ = item2.ActiveDurationMin.GetTimeDuration()
	assert.Equal(t, time.Duration(10*time.Second), duration)
}

func TestOperatingConstraintsPowerDescriptionListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: []OperatingConstraintsPowerDescriptionDataType{
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

	newData := OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: []OperatingConstraintsPowerDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeProduce),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.OperatingConstraintsPowerDescriptionData
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

func TestOperatingConstraintsPowerRangeListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: []OperatingConstraintsPowerRangeDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				PowerMin:   NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				PowerMin:   NewScaledNumberType(1),
			},
		},
	}

	newData := OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: []OperatingConstraintsPowerRangeDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				PowerMin:   NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.OperatingConstraintsPowerRangeData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.PowerMin.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.PowerMin.GetValue())
}

func TestOperatingConstraintsPowerLevelListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: []OperatingConstraintsPowerLevelDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(0)),
				Power:      NewScaledNumberType(1),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Power:      NewScaledNumberType(1),
			},
		},
	}

	newData := OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: []OperatingConstraintsPowerLevelDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Power:      NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.OperatingConstraintsPowerLevelData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.Power.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.Power.GetValue())
}

func TestOperatingConstraintsResumeImplicationListDataType_Update(t *testing.T) {
	sut := OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: []OperatingConstraintsResumeImplicationDataType{
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(0)),
				ResumeEnergyEstimated: NewScaledNumberType(1),
			},
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(1)),
				ResumeEnergyEstimated: NewScaledNumberType(1),
			},
		},
	}

	newData := OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: []OperatingConstraintsResumeImplicationDataType{
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(1)),
				ResumeEnergyEstimated: NewScaledNumberType(10),
			},
		},
	}

	// Act
	_, success := sut.UpdateList(false, true, &newData, NewFilterTypePartial(), nil, nil)
	assert.True(t, success)

	data := sut.OperatingConstraintsResumeImplicationData
	// check the non changing items
	assert.Equal(t, 2, len(data))
	item1 := data[0]
	assert.Equal(t, 0, int(*item1.SequenceId))
	assert.Equal(t, 1.0, item1.ResumeEnergyEstimated.GetValue())
	// check properties of updated item
	item2 := data[1]
	assert.Equal(t, 1, int(*item2.SequenceId))
	assert.Equal(t, 10.0, item2.ResumeEnergyEstimated.GetValue())
}

// ReadPartialData tests

func TestOperatingConstraintsInterruptListDataType_ReadPartialData(t *testing.T) {
	sut := &OperatingConstraintsInterruptListDataType{
		OperatingConstraintsInterruptData: []OperatingConstraintsInterruptDataType{
			{
				SequenceId:                  util.Ptr(PowerSequenceIdType(1)),
				IsPausable:                  util.Ptr(true),
				IsStoppable:                 util.Ptr(false),
				NotInterruptibleAtHighPower: util.Ptr(true),
				MaxCyclesPerDay:             util.Ptr(uint(5)),
			},
			{
				SequenceId:                  util.Ptr(PowerSequenceIdType(2)),
				IsPausable:                  util.Ptr(false),
				IsStoppable:                 util.Ptr(true),
				NotInterruptibleAtHighPower: util.Ptr(false),
				MaxCyclesPerDay:             util.Ptr(uint(10)),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*OperatingConstraintsInterruptListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsInterruptData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		OperatingConstraintsInterruptListDataSelectors: &OperatingConstraintsInterruptListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsInterruptListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsInterruptData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.OperatingConstraintsInterruptData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		OperatingConstraintsInterruptDataElements: &OperatingConstraintsInterruptDataElementsType{
			SequenceId: &ElementTagType{},
			IsPausable: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsInterruptListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsInterruptData))
	assert.NotNil(t, resultData.OperatingConstraintsInterruptData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsInterruptData[0].IsPausable)
	assert.Nil(t, resultData.OperatingConstraintsInterruptData[0].IsStoppable)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		OperatingConstraintsInterruptListDataSelectors: &OperatingConstraintsInterruptListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		OperatingConstraintsInterruptDataElements: &OperatingConstraintsInterruptDataElementsType{
			SequenceId:      &ElementTagType{},
			MaxCyclesPerDay: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsInterruptListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsInterruptData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.OperatingConstraintsInterruptData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsInterruptData[0].MaxCyclesPerDay)
	assert.Nil(t, resultData.OperatingConstraintsInterruptData[0].IsPausable)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		OperatingConstraintsInterruptListDataSelectors: &OperatingConstraintsInterruptListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsInterruptListDataType)
	assert.Equal(t, 0, len(resultData.OperatingConstraintsInterruptData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestOperatingConstraintsDurationListDataType_ReadPartialData(t *testing.T) {
	sut := &OperatingConstraintsDurationListDataType{
		OperatingConstraintsDurationData: []OperatingConstraintsDurationDataType{
			{
				SequenceId:           util.Ptr(PowerSequenceIdType(1)),
				ActiveDurationMin:    NewDurationType(time.Hour),
				ActiveDurationMax:    NewDurationType(2 * time.Hour),
				PauseDurationMin:     NewDurationType(30 * time.Minute),
				PauseDurationMax:     NewDurationType(time.Hour),
				ActiveDurationSumMin: NewDurationType(3 * time.Hour),
				ActiveDurationSumMax: NewDurationType(5 * time.Hour),
			},
			{
				SequenceId:           util.Ptr(PowerSequenceIdType(2)),
				ActiveDurationMin:    NewDurationType(2 * time.Hour),
				ActiveDurationMax:    NewDurationType(4 * time.Hour),
				PauseDurationMin:     NewDurationType(time.Hour),
				PauseDurationMax:     NewDurationType(2 * time.Hour),
				ActiveDurationSumMin: NewDurationType(6 * time.Hour),
				ActiveDurationSumMax: NewDurationType(10 * time.Hour),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*OperatingConstraintsDurationListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsDurationData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		OperatingConstraintsDurationListDataSelectors: &OperatingConstraintsDurationListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsDurationListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsDurationData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.OperatingConstraintsDurationData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		OperatingConstraintsDurationDataElements: &OperatingConstraintsDurationDataElementsType{
			SequenceId:        &ElementTagType{},
			ActiveDurationMin: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsDurationListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsDurationData))
	assert.NotNil(t, resultData.OperatingConstraintsDurationData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsDurationData[0].ActiveDurationMin)
	assert.Nil(t, resultData.OperatingConstraintsDurationData[0].ActiveDurationMax)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		OperatingConstraintsDurationListDataSelectors: &OperatingConstraintsDurationListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		OperatingConstraintsDurationDataElements: &OperatingConstraintsDurationDataElementsType{
			SequenceId:           &ElementTagType{},
			ActiveDurationSumMax: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsDurationListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsDurationData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.OperatingConstraintsDurationData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsDurationData[0].ActiveDurationSumMax)
	assert.Nil(t, resultData.OperatingConstraintsDurationData[0].ActiveDurationMin)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		OperatingConstraintsDurationListDataSelectors: &OperatingConstraintsDurationListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsDurationListDataType)
	assert.Equal(t, 0, len(resultData.OperatingConstraintsDurationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestOperatingConstraintsPowerDescriptionListDataType_ReadPartialData(t *testing.T) {
	sut := &OperatingConstraintsPowerDescriptionListDataType{
		OperatingConstraintsPowerDescriptionData: []OperatingConstraintsPowerDescriptionDataType{
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(1)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeConsume),
				PowerUnit:               util.Ptr(UnitOfMeasurementTypeW),
				EnergyUnit:              util.Ptr(UnitOfMeasurementTypeWh),
				Description:             util.Ptr(DescriptionType("Power Description 1")),
			},
			{
				SequenceId:              util.Ptr(PowerSequenceIdType(2)),
				PositiveEnergyDirection: util.Ptr(EnergyDirectionTypeProduce),
				PowerUnit:               util.Ptr(UnitOfMeasurementTypeV),
				EnergyUnit:              util.Ptr(UnitOfMeasurementTypeVAh),
				Description:             util.Ptr(DescriptionType("Power Description 2")),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*OperatingConstraintsPowerDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsPowerDescriptionData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		OperatingConstraintsPowerDescriptionListDataSelectors: &OperatingConstraintsPowerDescriptionListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsPowerDescriptionData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.OperatingConstraintsPowerDescriptionData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		OperatingConstraintsPowerDescriptionDataElements: &OperatingConstraintsPowerDescriptionDataElementsType{
			SequenceId:              &ElementTagType{},
			PositiveEnergyDirection: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerDescriptionListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsPowerDescriptionData))
	assert.NotNil(t, resultData.OperatingConstraintsPowerDescriptionData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsPowerDescriptionData[0].PositiveEnergyDirection)
	assert.Nil(t, resultData.OperatingConstraintsPowerDescriptionData[0].PowerUnit)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		OperatingConstraintsPowerDescriptionListDataSelectors: &OperatingConstraintsPowerDescriptionListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		OperatingConstraintsPowerDescriptionDataElements: &OperatingConstraintsPowerDescriptionDataElementsType{
			SequenceId:  &ElementTagType{},
			Description: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerDescriptionListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsPowerDescriptionData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.OperatingConstraintsPowerDescriptionData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsPowerDescriptionData[0].Description)
	assert.Nil(t, resultData.OperatingConstraintsPowerDescriptionData[0].PositiveEnergyDirection)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		OperatingConstraintsPowerDescriptionListDataSelectors: &OperatingConstraintsPowerDescriptionListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerDescriptionListDataType)
	assert.Equal(t, 0, len(resultData.OperatingConstraintsPowerDescriptionData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestOperatingConstraintsPowerRangeListDataType_ReadPartialData(t *testing.T) {
	sut := &OperatingConstraintsPowerRangeListDataType{
		OperatingConstraintsPowerRangeData: []OperatingConstraintsPowerRangeDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				PowerMin:   NewScaledNumberType(100.0),
				PowerMax:   NewScaledNumberType(500.0),
				EnergyMin:  NewScaledNumberType(1000.0),
				EnergyMax:  NewScaledNumberType(5000.0),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(2)),
				PowerMin:   NewScaledNumberType(200.0),
				PowerMax:   NewScaledNumberType(1000.0),
				EnergyMin:  NewScaledNumberType(2000.0),
				EnergyMax:  NewScaledNumberType(10000.0),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*OperatingConstraintsPowerRangeListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsPowerRangeData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		OperatingConstraintsPowerRangeListDataSelectors: &OperatingConstraintsPowerRangeListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerRangeListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsPowerRangeData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.OperatingConstraintsPowerRangeData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		OperatingConstraintsPowerRangeDataElements: &OperatingConstraintsPowerRangeDataElementsType{
			SequenceId: &ElementTagType{},
			PowerMin:   &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerRangeListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsPowerRangeData))
	assert.NotNil(t, resultData.OperatingConstraintsPowerRangeData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsPowerRangeData[0].PowerMin)
	assert.Nil(t, resultData.OperatingConstraintsPowerRangeData[0].PowerMax)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		OperatingConstraintsPowerRangeListDataSelectors: &OperatingConstraintsPowerRangeListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		OperatingConstraintsPowerRangeDataElements: &OperatingConstraintsPowerRangeDataElementsType{
			SequenceId: &ElementTagType{},
			EnergyMax:  &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerRangeListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsPowerRangeData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.OperatingConstraintsPowerRangeData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsPowerRangeData[0].EnergyMax)
	assert.Nil(t, resultData.OperatingConstraintsPowerRangeData[0].PowerMin)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		OperatingConstraintsPowerRangeListDataSelectors: &OperatingConstraintsPowerRangeListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerRangeListDataType)
	assert.Equal(t, 0, len(resultData.OperatingConstraintsPowerRangeData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestOperatingConstraintsPowerLevelListDataType_ReadPartialData(t *testing.T) {
	sut := &OperatingConstraintsPowerLevelListDataType{
		OperatingConstraintsPowerLevelData: []OperatingConstraintsPowerLevelDataType{
			{
				SequenceId: util.Ptr(PowerSequenceIdType(1)),
				Power:      NewScaledNumberType(250.0),
			},
			{
				SequenceId: util.Ptr(PowerSequenceIdType(2)),
				Power:      NewScaledNumberType(750.0),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*OperatingConstraintsPowerLevelListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsPowerLevelData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		OperatingConstraintsPowerLevelListDataSelectors: &OperatingConstraintsPowerLevelListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerLevelListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsPowerLevelData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.OperatingConstraintsPowerLevelData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		OperatingConstraintsPowerLevelDataElements: &OperatingConstraintsPowerLevelDataElementsType{
			SequenceId: &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerLevelListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsPowerLevelData))
	assert.NotNil(t, resultData.OperatingConstraintsPowerLevelData[0].SequenceId)
	assert.Nil(t, resultData.OperatingConstraintsPowerLevelData[0].Power)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		OperatingConstraintsPowerLevelListDataSelectors: &OperatingConstraintsPowerLevelListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		OperatingConstraintsPowerLevelDataElements: &OperatingConstraintsPowerLevelDataElementsType{
			SequenceId: &ElementTagType{},
			Power:      &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerLevelListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsPowerLevelData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.OperatingConstraintsPowerLevelData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsPowerLevelData[0].Power)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		OperatingConstraintsPowerLevelListDataSelectors: &OperatingConstraintsPowerLevelListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsPowerLevelListDataType)
	assert.Equal(t, 0, len(resultData.OperatingConstraintsPowerLevelData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}

func TestOperatingConstraintsResumeImplicationListDataType_ReadPartialData(t *testing.T) {
	sut := &OperatingConstraintsResumeImplicationListDataType{
		OperatingConstraintsResumeImplicationData: []OperatingConstraintsResumeImplicationDataType{
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(1)),
				ResumeEnergyEstimated: NewScaledNumberType(150.0),
				EnergyUnit:            util.Ptr(UnitOfMeasurementTypeWh),
				ResumeCostEstimated:   NewScaledNumberType(25.5),
				Currency:              util.Ptr(CurrencyType("EUR")),
			},
			{
				SequenceId:            util.Ptr(PowerSequenceIdType(2)),
				ResumeEnergyEstimated: NewScaledNumberType(300.0),
				EnergyUnit:            util.Ptr(UnitOfMeasurementTypeWh),
				ResumeCostEstimated:   NewScaledNumberType(45.75),
				Currency:              util.Ptr(CurrencyType("USD")),
			},
		},
	}

	// Test 1: nil filter should return original data
	result, success := sut.ReadPartialData(nil)
	assert.True(t, success)
	resultData := result.(*OperatingConstraintsResumeImplicationListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsResumeImplicationData))

	// Test 2: filter with selector should return filtered data
	filter := FilterType{
		OperatingConstraintsResumeImplicationListDataSelectors: &OperatingConstraintsResumeImplicationListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(1)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsResumeImplicationListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsResumeImplicationData))
	assert.Equal(t, PowerSequenceIdType(1), *resultData.OperatingConstraintsResumeImplicationData[0].SequenceId)

	// Test 3: filter with elements should return filtered data
	filter = FilterType{
		OperatingConstraintsResumeImplicationDataElements: &OperatingConstraintsResumeImplicationDataElementsType{
			SequenceId:            &ElementTagType{},
			ResumeEnergyEstimated: &ScaledNumberElementsType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsResumeImplicationListDataType)
	assert.Equal(t, 2, len(resultData.OperatingConstraintsResumeImplicationData))
	assert.NotNil(t, resultData.OperatingConstraintsResumeImplicationData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsResumeImplicationData[0].ResumeEnergyEstimated)
	assert.Nil(t, resultData.OperatingConstraintsResumeImplicationData[0].EnergyUnit)

	// Test 4: filter with both selector and elements
	filter = FilterType{
		OperatingConstraintsResumeImplicationListDataSelectors: &OperatingConstraintsResumeImplicationListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(2)),
		},
		OperatingConstraintsResumeImplicationDataElements: &OperatingConstraintsResumeImplicationDataElementsType{
			SequenceId: &ElementTagType{},
			Currency:   &ElementTagType{},
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsResumeImplicationListDataType)
	assert.Equal(t, 1, len(resultData.OperatingConstraintsResumeImplicationData))
	assert.Equal(t, PowerSequenceIdType(2), *resultData.OperatingConstraintsResumeImplicationData[0].SequenceId)
	assert.NotNil(t, resultData.OperatingConstraintsResumeImplicationData[0].Currency)
	assert.Nil(t, resultData.OperatingConstraintsResumeImplicationData[0].ResumeEnergyEstimated)

	// Test 5: non-matching selector should return empty data
	filter = FilterType{
		OperatingConstraintsResumeImplicationListDataSelectors: &OperatingConstraintsResumeImplicationListDataSelectorsType{
			SequenceId: util.Ptr(PowerSequenceIdType(99)),
		},
	}
	result, success = sut.ReadPartialData(&filter)
	assert.True(t, success)
	resultData = result.(*OperatingConstraintsResumeImplicationListDataType)
	assert.Equal(t, 0, len(resultData.OperatingConstraintsResumeImplicationData))

	// Test 6: ReadPartialData with empty filter
	emptyFilter := &FilterType{}
	_, success = sut.ReadPartialData(emptyFilter)
	assert.False(t, success)
}
