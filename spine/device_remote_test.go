package spine

import (
	"reflect"
	"testing"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	nm_detaileddiscovery_emptyarray_file_path          = "testdata/nm_detaileddiscovery_emptyarray.json"
	nm_usecaseinformationlistdata_recv_reply_file_path = "../spine/testdata/nm_usecaseinformationlistdata_recv_reply.json"
)

func TestDeviceRemoteSuite(t *testing.T) {
	suite.Run(t, new(DeviceRemoteSuite))
}

type DeviceRemoteSuite struct {
	suite.Suite

	localDevice  api.DeviceLocalInterface
	remoteDevice api.DeviceRemoteInterface
	remoteEntity api.EntityRemoteInterface
}

func (s *DeviceRemoteSuite) WriteShipMessageWithPayload([]byte) {}

func (s *DeviceRemoteSuite) SetupSuite() {}

func (s *DeviceRemoteSuite) BeforeTest(suiteName, testName string) {
	s.localDevice = NewDeviceLocal("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	ski := "test"
	sender := NewSender(s)
	s.remoteDevice = NewDeviceRemote(s.localDevice, ski, sender)
	desc := &model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{
			Device: util.Ptr(model.AddressDeviceType("test")),
		},
	}
	s.remoteDevice.UpdateDevice(desc)
	_ = s.localDevice.SetupRemoteDevice(ski, s)

	s.remoteEntity = NewEntityRemote(s.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})

	feature := NewFeatureRemote(0, s.remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	s.remoteEntity.AddFeature(feature)

	s.remoteDevice.AddEntity(s.remoteEntity)
}

func (s *DeviceRemoteSuite) Test_RemoveByAddress() {
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{2})
	assert.Equal(s.T(), 2, len(s.remoteDevice.Entities()))

	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))
}

func (s *DeviceRemoteSuite) Test_FeatureByEntityTypeAndRole() {
	entity := s.remoteDevice.Entity([]model.AddressEntityType{1})
	assert.NotNil(s.T(), entity)

	assert.Equal(s.T(), 1, len(entity.Features()))

	feature := s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	assert.Nil(s.T(), feature)

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), feature)

	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{1})
	assert.Equal(s.T(), 1, len(s.remoteDevice.Entities()))

	_ = s.remoteDevice.Entity([]model.AddressEntityType{0})
	s.remoteDevice.RemoveEntityByAddress([]model.AddressEntityType{0})
	assert.Equal(s.T(), 0, len(s.remoteDevice.Entities()))

	feature = s.remoteDevice.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.Nil(s.T(), feature)
}

func (s *DeviceRemoteSuite) Test_Usecases() {
	uc := s.remoteDevice.UseCases()
	assert.Nil(s.T(), uc)

	_, _ = s.remoteDevice.HandleSpineMesssage(loadFileData(s.T(), nm_usecaseinformationlistdata_recv_reply_file_path))

	uc = s.remoteDevice.UseCases()
	assert.NotNil(s.T(), uc)
}

// our simple EEBUS JSON to JSON conversion in ship is converting empty arrays to empty objects which will break unmarshalling
func (s *DeviceRemoteSuite) Test_EmptyArrayDataStructure() {
	_, err := s.remoteDevice.HandleSpineMesssage(loadFileData(s.T(), nm_detaileddiscovery_emptyarray_file_path))
	assert.Nil(s.T(), err)
}

func Test_findFieldTypeByJSONTag(t *testing.T) {
	// Test struct with various JSON tags
	type TestStruct struct {
		SimpleField    string  `json:"simpleField"`
		OmitEmptyField int     `json:"omitEmptyField,omitempty"`
		PointerField   *string `json:"pointerField"`
		NoTagField     string
		EmptyTagField  string `json:""`
		SliceField     []int  `json:"sliceField"`
	}

	type NestedStruct struct {
		Inner TestStruct `json:"inner"`
	}

	tests := []struct {
		name       string
		structType reflect.Type
		jsonName   string
		wantNil    bool
		wantKind   reflect.Kind
	}{
		{
			name:       "simple field found",
			structType: reflect.TypeOf(TestStruct{}),
			jsonName:   "simpleField",
			wantNil:    false,
			wantKind:   reflect.String,
		},
		{
			name:       "field with omitempty found",
			structType: reflect.TypeOf(TestStruct{}),
			jsonName:   "omitEmptyField",
			wantNil:    false,
			wantKind:   reflect.Int,
		},
		{
			name:       "pointer field found",
			structType: reflect.TypeOf(TestStruct{}),
			jsonName:   "pointerField",
			wantNil:    false,
			wantKind:   reflect.Ptr,
		},
		{
			name:       "slice field found",
			structType: reflect.TypeOf(TestStruct{}),
			jsonName:   "sliceField",
			wantNil:    false,
			wantKind:   reflect.Slice,
		},
		{
			name:       "field not found - wrong name",
			structType: reflect.TypeOf(TestStruct{}),
			jsonName:   "nonExistent",
			wantNil:    true,
		},
		{
			name:       "field without json tag not found",
			structType: reflect.TypeOf(TestStruct{}),
			jsonName:   "NoTagField",
			wantNil:    true,
		},
		{
			name:       "empty json tag not matched",
			structType: reflect.TypeOf(TestStruct{}),
			jsonName:   "EmptyTagField",
			wantNil:    true,
		},
		{
			name:       "pointer to struct works",
			structType: reflect.TypeOf(&TestStruct{}),
			jsonName:   "simpleField",
			wantNil:    false,
			wantKind:   reflect.String,
		},
		{
			name:       "nested struct field found",
			structType: reflect.TypeOf(NestedStruct{}),
			jsonName:   "inner",
			wantNil:    false,
			wantKind:   reflect.Struct,
		},
		{
			name:       "non-struct type returns nil",
			structType: reflect.TypeOf("string"),
			jsonName:   "anyField",
			wantNil:    true,
		},
		{
			name:       "slice type returns nil",
			structType: reflect.TypeOf([]int{}),
			jsonName:   "anyField",
			wantNil:    true,
		},
		{
			name:       "pointer to non-struct returns nil",
			structType: reflect.TypeOf(new(int)),
			jsonName:   "anyField",
			wantNil:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findFieldTypeByJSONTag(tt.structType, tt.jsonName)

			if tt.wantNil {
				assert.Nil(t, result, "expected nil result")
			} else {
				assert.NotNil(t, result, "expected non-nil result")
				if result != nil {
					assert.Equal(t, tt.wantKind, (*result).Kind(), "unexpected field kind")
				}
			}
		})
	}
}

func Test_findFieldTypeByJSONTag_WithModelTypes(t *testing.T) {
	// Test with actual SPINE model types to ensure compatibility
	datagramType := reflect.TypeOf(model.Datagram{})

	// The Datagram struct should have a "datagram" field
	result := findFieldTypeByJSONTag(datagramType, "datagram")
	assert.NotNil(t, result, "expected to find 'datagram' field in Datagram type")

	// Test with CmdType which has various slice fields
	cmdType := reflect.TypeOf(model.CmdType{})

	// Check for a known field in CmdType
	result = findFieldTypeByJSONTag(cmdType, "function")
	if result != nil {
		assert.Equal(t, reflect.Ptr, (*result).Kind())
	}

	// Test with non-existent field
	result = findFieldTypeByJSONTag(cmdType, "nonExistentField")
	assert.Nil(t, result, "expected nil for non-existent field")
}

func Test_fixupSliceFieldsRecursive(t *testing.T) {
	// Test struct definitions for type reference
	type InnerStruct struct {
		Name  string   `json:"name"`
		Items []string `json:"items"`
	}

	type TestStruct struct {
		StringField string        `json:"stringField"`
		IntField    int           `json:"intField"`
		SliceField  []string      `json:"sliceField"`
		NestedSlice []InnerStruct `json:"nestedSlice"`
		Inner       InnerStruct   `json:"inner"`
	}

	tests := []struct {
		name      string
		input     interface{}
		modelType reflect.Type
		validate  func(t *testing.T, result interface{})
	}{
		{
			name:      "nil value returns nil",
			input:     nil,
			modelType: reflect.TypeOf(TestStruct{}),
			validate: func(t *testing.T, result interface{}) {
				assert.Nil(t, result)
			},
		},
		{
			name:      "primitive string value unchanged",
			input:     "hello",
			modelType: reflect.TypeOf(""),
			validate: func(t *testing.T, result interface{}) {
				assert.Equal(t, "hello", result)
			},
		},
		{
			name:      "primitive int value unchanged",
			input:     42,
			modelType: reflect.TypeOf(0),
			validate: func(t *testing.T, result interface{}) {
				assert.Equal(t, 42, result)
			},
		},
		{
			name:      "primitive float value unchanged",
			input:     3.14,
			modelType: reflect.TypeOf(0.0),
			validate: func(t *testing.T, result interface{}) {
				assert.Equal(t, 3.14, result)
			},
		},
		{
			name:      "primitive bool value unchanged",
			input:     true,
			modelType: reflect.TypeOf(true),
			validate: func(t *testing.T, result interface{}) {
				assert.Equal(t, true, result)
			},
		},
		{
			name: "empty map for slice field converted to empty slice",
			input: map[string]interface{}{
				"sliceField": map[string]interface{}{},
			},
			modelType: reflect.TypeOf(TestStruct{}),
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				assert.True(t, ok)
				sliceVal, exists := resultMap["sliceField"]
				assert.True(t, exists)
				slice, ok := sliceVal.([]interface{})
				assert.True(t, ok, "expected slice type, got %T", sliceVal)
				assert.Empty(t, slice)
			},
		},
		{
			name: "non-empty map for non-slice field unchanged",
			input: map[string]interface{}{
				"inner": map[string]interface{}{
					"name": "test",
				},
			},
			modelType: reflect.TypeOf(TestStruct{}),
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				assert.True(t, ok)
				innerVal, exists := resultMap["inner"]
				assert.True(t, exists)
				innerMap, ok := innerVal.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "test", innerMap["name"])
			},
		},
		{
			name: "array with elements processed recursively",
			input: []interface{}{
				map[string]interface{}{"name": "item1", "items": map[string]interface{}{}},
				map[string]interface{}{"name": "item2", "items": []interface{}{"a", "b"}},
			},
			modelType: reflect.TypeOf([]InnerStruct{}),
			validate: func(t *testing.T, result interface{}) {
				resultSlice, ok := result.([]interface{})
				assert.True(t, ok)
				assert.Len(t, resultSlice, 2)

				// First element should have items converted to empty slice
				first, ok := resultSlice[0].(map[string]interface{})
				assert.True(t, ok)
				items1, ok := first["items"].([]interface{})
				assert.True(t, ok, "expected slice type for items, got %T", first["items"])
				assert.Empty(t, items1)

				// Second element should keep its array
				second, ok := resultSlice[1].(map[string]interface{})
				assert.True(t, ok)
				items2, ok := second["items"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, items2, 2)
			},
		},
		{
			name: "nested structure processed correctly",
			input: map[string]interface{}{
				"stringField": "test",
				"intField":    123,
				"sliceField":  map[string]interface{}{},
				"inner": map[string]interface{}{
					"name":  "nested",
					"items": map[string]interface{}{},
				},
			},
			modelType: reflect.TypeOf(TestStruct{}),
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				assert.True(t, ok)

				// String field unchanged
				assert.Equal(t, "test", resultMap["stringField"])

				// Int field unchanged
				assert.Equal(t, 123, resultMap["intField"])

				// Slice field converted
				sliceVal, ok := resultMap["sliceField"].([]interface{})
				assert.True(t, ok, "expected slice type for sliceField")
				assert.Empty(t, sliceVal)

				// Nested inner items also converted
				innerMap, ok := resultMap["inner"].(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "nested", innerMap["name"])
				innerItems, ok := innerMap["items"].([]interface{})
				assert.True(t, ok, "expected slice type for inner items")
				assert.Empty(t, innerItems)
			},
		},
		{
			name: "unknown field in map preserved",
			input: map[string]interface{}{
				"stringField":  "test",
				"unknownField": "unknown",
			},
			modelType: reflect.TypeOf(TestStruct{}),
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "test", resultMap["stringField"])
				assert.Equal(t, "unknown", resultMap["unknownField"])
			},
		},
		{
			name: "map with non-struct model type",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": map[string]interface{}{},
			},
			modelType: reflect.TypeOf((*interface{})(nil)).Elem(),
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				assert.True(t, ok)
				assert.Equal(t, "value1", resultMap["key1"])
				// Empty map stays as map when model type is interface{}
				_, ok = resultMap["key2"].(map[string]interface{})
				assert.True(t, ok)
			},
		},
		{
			name: "pointer model type dereferenced",
			input: map[string]interface{}{
				"sliceField": map[string]interface{}{},
			},
			modelType: reflect.TypeOf(&TestStruct{}),
			validate: func(t *testing.T, result interface{}) {
				resultMap, ok := result.(map[string]interface{})
				assert.True(t, ok)
				sliceVal, ok := resultMap["sliceField"].([]interface{})
				assert.True(t, ok, "expected slice type")
				assert.Empty(t, sliceVal)
			},
		},
		{
			name:      "empty array unchanged",
			input:     []interface{}{},
			modelType: reflect.TypeOf([]string{}),
			validate: func(t *testing.T, result interface{}) {
				resultSlice, ok := result.([]interface{})
				assert.True(t, ok)
				assert.Empty(t, resultSlice)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixupSliceFieldsRecursive(tt.input, tt.modelType)
			tt.validate(t, result)
		})
	}
}

func Test_fixupSliceFieldsRecursive_WithModelTypes(t *testing.T) {
	// Test with actual SPINE model types
	datagramType := reflect.TypeOf(model.Datagram{})

	// Simulate JSON-parsed data with empty object where array should be
	input := map[string]interface{}{
		"datagram": map[string]interface{}{
			"header": map[string]interface{}{
				"specificationVersion": "1.3.0",
			},
			"payload": map[string]interface{}{
				"cmd": map[string]interface{}{}, // This should be converted to []
			},
		},
	}

	result := fixupSliceFieldsRecursive(input, datagramType)

	resultMap, ok := result.(map[string]interface{})
	assert.True(t, ok)
	assert.NotNil(t, resultMap["datagram"])
}

func Test_fixupSliceFields(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains string
	}{
		{
			name:     "no empty objects - unchanged",
			input:    `{"datagram":{"header":{"specificationVersion":"1.3.0"}}}`,
			contains: `"specificationVersion":"1.3.0"`,
		},
		{
			name:     "invalid JSON returns original",
			input:    `{invalid json}`,
			contains: `{invalid json}`,
		},
		{
			name:     "empty object pattern triggers fixup",
			input:    `{"datagram":{"payload":{"cmd":{}}}}`,
			contains: `"cmd"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fixupSliceFields([]byte(tt.input))
			assert.Contains(t, string(result), tt.contains)
		})
	}
}
