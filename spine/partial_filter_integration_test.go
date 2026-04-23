package spine

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestPartialFilterIntegration(t *testing.T) {
	suite.Run(t, new(PartialFilterIntegrationTestSuite))
}

type PartialFilterIntegrationTestSuite struct {
	suite.Suite
	senderMock        *mocks.SenderInterface
	localDevice       *DeviceLocal
	localEntity       *EntityLocal
	localFeature      api.FeatureLocalInterface
	remoteDevice      *DeviceRemote
	remoteFeature     api.FeatureRemoteInterface
	serverFunction    model.FunctionType
	serverFeatureType model.FeatureTypeType
}

func (s *PartialFilterIntegrationTestSuite) BeforeTest(suiteName, testName string) {
	s.senderMock = mocks.NewSenderInterface(s.T())
	s.serverFunction = model.FunctionTypeLoadControlLimitListData
	s.serverFeatureType = model.FeatureTypeTypeLoadControl

	// Create local device and server feature
	s.localDevice, s.localEntity = createLocalDeviceAndEntity(1)
	_, s.localFeature = createLocalFeatures(s.localEntity, s.serverFeatureType, s.serverFunction)

	// Create remote device and client feature
	s.remoteDevice = createRemoteDevice(s.localDevice, "remotedevice", s.senderMock)
	s.remoteFeature, _ = createRemoteEntityAndFeature(s.remoteDevice, 1, s.serverFeatureType, s.serverFunction)
}

// Integration test: Complete message flow with partial filters
func (s *PartialFilterIntegrationTestSuite) Test_EndToEnd_PartialFilterIgnored() {
	// Setup: Add comprehensive test data to the local server feature
	testData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive:     util.Ptr(false),
				IsLimitChangeable: util.Ptr(true),
				Value:             model.NewScaledNumberType(1000),
				TimePeriod:        model.NewTimePeriodTypeWithRelativeEndTime(time.Minute * 30),
			},
			{
				LimitId:           util.Ptr(model.LoadControlLimitIdType(2)),
				IsLimitActive:     util.Ptr(true),
				IsLimitChangeable: util.Ptr(false),
				Value:             model.NewScaledNumberType(2000),
				TimePeriod:        model.NewTimePeriodTypeWithRelativeEndTime(time.Hour * 1),
			},
		},
	}
	s.localFeature.SetData(s.serverFunction, testData)

	// Create a complex partial filter that combines selectors and elements
	partialFilter := model.FilterType{
		CmdControl: &model.CmdControlType{
			Partial: &model.ElementTagType{},
		},
		// Selector: Only item with ID 1
		LoadControlLimitListDataSelectors: &model.LoadControlLimitListDataSelectorsType{
			LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
		},
		// Elements: Only LimitId and IsLimitActive
		LoadControlLimitDataElements: &model.LoadControlLimitDataElementsType{
			LimitId:       &model.ElementTagType{},
			IsLimitActive: &model.ElementTagType{},
		},
	}

	// Step 1: Create command with partial filter (simulating incoming read request)
	readCmd := model.CmdType{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		Filter:                   []model.FilterType{partialFilter},
	}

	// Step 2: Extract filters (simulating what ProcessCmd does)
	filterPartial, filterDelete := readCmd.ExtractFilter()
	assert.NotNil(s.T(), filterPartial)
	assert.Nil(s.T(), filterDelete)

	// Step 3: Create message with extracted filters
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(100)),
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		Cmd:           readCmd,
		FilterPartial: filterPartial,
		FilterDelete:  filterDelete,
		FeatureRemote: s.remoteFeature,
		EntityRemote:  s.remoteFeature.Entity(),
		DeviceRemote:  s.remoteFeature.Device(),
	}

	// Step 4: Setup expectation - reply should contain FULL data (ignoring filters)
	s.senderMock.EXPECT().Reply(
		mock.MatchedBy(func(header *model.HeaderType) bool {
			return header.MsgCounter != nil && *header.MsgCounter == model.MsgCounterType(100)
		}),
		s.localFeature.Address(),
		mock.MatchedBy(func(replyCmd model.CmdType) bool {
			// Verify the reply ignores the partial filter and returns full data
			if replyCmd.LoadControlLimitListData == nil {
				return false
			}

			data := replyCmd.LoadControlLimitListData

			// Should contain ALL entries (ignores selector for ID 1)
			if len(data.LoadControlLimitData) != 2 {
				return false
			}

			// Should contain ALL fields for each entry (ignores element filter)
			for _, entry := range data.LoadControlLimitData {
				if entry.LimitId == nil || entry.IsLimitActive == nil || entry.IsLimitChangeable == nil ||
					entry.Value == nil || entry.TimePeriod == nil {
					return false
				}
			}

			// Verify no filter is included in the reply
			return len(replyCmd.Filter) == 0
		}),
	).Return(nil)

	// Step 5: Process the message (this is the actual functionality being tested)
	err := s.localFeature.HandleMessage(msg)
	assert.Nil(s.T(), err)

	// Step 6: Verify the mocks were called as expected
	s.senderMock.AssertExpectations(s.T())
}

// Integration test: Verify behavior across different function types
func (s *PartialFilterIntegrationTestSuite) Test_EndToEnd_DifferentFunctionTypes() {
	// Test with DeviceClassificationManufacturerData (different data type)
	manufacturerData := &model.DeviceClassificationManufacturerDataType{
		BrandName:    util.Ptr(model.DeviceClassificationStringType("Test Brand")),
		VendorName:   util.Ptr(model.DeviceClassificationStringType("Test Vendor")),
		DeviceName:   util.Ptr(model.DeviceClassificationStringType("Test Device")),
		DeviceCode:   util.Ptr(model.DeviceClassificationStringType("TEST001")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("SN123456")),
	}

	// Create a local feature for DeviceClassification
	dcFeatureType := model.FeatureTypeTypeDeviceClassification
	dcFunction := model.FunctionTypeDeviceClassificationManufacturerData
	_, dcLocalFeature := createLocalFeatures(s.localEntity, dcFeatureType, "")
	dcLocalFeature.SetData(dcFunction, manufacturerData)

	// Create remote feature for DeviceClassification
	dcRemoteFeature, _ := createRemoteEntityAndFeature(s.remoteDevice, 2, dcFeatureType, dcFunction)

	// Create partial filter for DeviceClassification data
	partialFilter := model.FilterType{
		CmdControl: &model.CmdControlType{
			Partial: &model.ElementTagType{},
		},
		DeviceClassificationManufacturerDataElements: &model.DeviceClassificationManufacturerDataElementsType{
			BrandName: &model.ElementTagType{},
			// Only requesting BrandName, not other fields
		},
	}

	// Create read message
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(200)),
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: &model.DeviceClassificationManufacturerDataType{},
			Filter:                               []model.FilterType{partialFilter},
		},
		FilterPartial: &partialFilter,
		FeatureRemote: dcRemoteFeature,
		EntityRemote:  dcRemoteFeature.Entity(),
		DeviceRemote:  dcRemoteFeature.Device(),
	}

	// Expect full data reply (all fields, not just BrandName)
	s.senderMock.EXPECT().Reply(
		mock.Anything,
		dcLocalFeature.Address(),
		mock.MatchedBy(func(replyCmd model.CmdType) bool {
			if replyCmd.DeviceClassificationManufacturerData == nil {
				return false
			}

			data := replyCmd.DeviceClassificationManufacturerData
			// Should contain ALL fields, not just BrandName
			return data.BrandName != nil && data.VendorName != nil && data.DeviceName != nil &&
				data.DeviceCode != nil && data.SerialNumber != nil
		}),
	).Return(nil)

	// Process the message
	err := dcLocalFeature.HandleMessage(msg)
	assert.Nil(s.T(), err)
}

// Integration test: Verify correct behavior when no filters are provided
func (s *PartialFilterIntegrationTestSuite) Test_EndToEnd_NoFilters_FullReply() {
	// Setup test data
	testData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				IsLimitActive: util.Ptr(false),
			},
		},
	}
	s.localFeature.SetData(s.serverFunction, testData)

	// Create read message WITHOUT any filters
	msg := &api.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(300)),
		},
		CmdClassifier: model.CmdClassifierTypeRead,
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
			// No Filter field
		},
		// No FilterPartial or FilterDelete
		FeatureRemote: s.remoteFeature,
		EntityRemote:  s.remoteFeature.Entity(),
		DeviceRemote:  s.remoteFeature.Device(),
	}

	// Expect full data reply
	s.senderMock.EXPECT().Reply(
		mock.Anything,
		s.localFeature.Address(),
		mock.MatchedBy(func(replyCmd model.CmdType) bool {
			return replyCmd.LoadControlLimitListData != nil &&
				len(replyCmd.LoadControlLimitListData.LoadControlLimitData) == 1
		}),
	).Return(nil)

	// Process the message
	err := s.localFeature.HandleMessage(msg)
	assert.Nil(s.T(), err)
}
