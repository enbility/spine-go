package model

import (
	"encoding/json"
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestNodeManagementXSDComplianceSuite(t *testing.T) {
	suite.Run(t, new(NodeManagementXSDComplianceSuite))
}

type NodeManagementXSDComplianceSuite struct {
	suite.Suite
}

// Test that NewEntityInformationForNodeManagement creates XSD-compliant entity information
// This test should FAIL initially to drive the implementation
func (s *NodeManagementXSDComplianceSuite) Test_NewEntityInformationForNodeManagement_XSDCompliant() {
	// GIVEN: Valid entity address and type
	entityAddr := []AddressEntityType{0, 1}
	entityType := EntityTypeTypeDeviceInformation

	// WHEN: Creating entity information for node management using factory function
	info := NewEntityInformationForNodeManagement(entityAddr, entityType)

	// THEN: The result should be XSD compliant
	assert.NotNil(s.T(), info, "Entity information should not be nil")
	assert.NotNil(s.T(), info.Description, "Description should not be nil")
	assert.NotNil(s.T(), info.Description.EntityAddress, "EntityAddress should not be nil")

	// XSD compliance: Device field must be nil for NodeManagement context
	assert.Nil(s.T(), info.Description.EntityAddress.Device, "Device field must be nil for XSD compliance")

	// Entity field should be properly set
	assert.Equal(s.T(), entityAddr, info.Description.EntityAddress.Entity, "Entity field should match input")

	// EntityType should be properly set
	assert.NotNil(s.T(), info.Description.EntityType, "EntityType should not be nil")
	assert.Equal(s.T(), entityType, *info.Description.EntityType, "EntityType should match input")
}

// Test that JSON marshaling excludes the device field for XSD compliance
// This test should FAIL initially to drive the implementation
func (s *NodeManagementXSDComplianceSuite) Test_EntityInformation_JSONMarshal_NoDeviceField() {
	// GIVEN: Entity information created for node management
	entityAddr := []AddressEntityType{0, 1}
	entityType := EntityTypeTypeDeviceInformation
	info := NewEntityInformationForNodeManagement(entityAddr, entityType)

	// WHEN: Marshaling to JSON
	jsonData, err := json.Marshal(info)
	assert.NoError(s.T(), err, "JSON marshaling should not fail")

	// THEN: The JSON should not contain a device field
	jsonString := string(jsonData)
	assert.NotContains(s.T(), jsonString, "device", "JSON should not contain device field for XSD compliance")

	// BUT: Should contain entity field
	assert.Contains(s.T(), jsonString, "entity", "JSON should contain entity field")

	// Verify the structure by unmarshaling back
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	assert.NoError(s.T(), err, "JSON should be valid")

	description, ok := result["description"].(map[string]interface{})
	assert.True(s.T(), ok, "Description should be present")

	entityAddress, ok := description["entityAddress"].(map[string]interface{})
	assert.True(s.T(), ok, "EntityAddress should be present")

	// Critical XSD compliance check
	_, hasDevice := entityAddress["device"]
	assert.False(s.T(), hasDevice, "EntityAddress should NOT have device field")

	entity, hasEntity := entityAddress["entity"]
	assert.True(s.T(), hasEntity, "EntityAddress should have entity field")
	assert.NotNil(s.T(), entity, "Entity field should not be nil")
}

// Test XSD validation method
// This test should FAIL initially to drive the implementation
func (s *NodeManagementXSDComplianceSuite) Test_EntityInformation_ValidateXSD() {
	// GIVEN: XSD-compliant entity information
	validInfo := NewEntityInformationForNodeManagement([]AddressEntityType{1}, EntityTypeTypeCEM)

	// WHEN: Validating XSD compliance
	err := validInfo.ValidateXSD()

	// THEN: Should pass validation
	assert.NoError(s.T(), err, "XSD-compliant entity information should pass validation")

	// GIVEN: Entity information with Device field set (XSD violation)
	invalidInfo := &NodeManagementDetailedDiscoveryEntityInformationType{
		Description: &NetworkManagementEntityDescriptionDataType{
			EntityAddress: &EntityAddressType{
				Device: util.Ptr(AddressDeviceType("InvalidDevice")), // XSD violation
				Entity: []AddressEntityType{1},
			},
			EntityType: util.Ptr(EntityTypeTypeCEM),
		},
	}

	// WHEN: Validating XSD compliance
	err = invalidInfo.ValidateXSD()

	// THEN: Should fail validation
	assert.Error(s.T(), err, "Entity information with Device field should fail XSD validation")
	assert.Contains(s.T(), err.Error(), "XSD violation", "Error should mention XSD violation")
	assert.Contains(s.T(), err.Error(), "Device field", "Error should mention Device field")
}

// Test that factory function handles edge cases properly
func (s *NodeManagementXSDComplianceSuite) Test_NewEntityInformationForNodeManagement_EdgeCases() {
	// GIVEN: Empty entity address
	emptyEntityAddr := []AddressEntityType{}

	// WHEN: Creating entity information
	info := NewEntityInformationForNodeManagement(emptyEntityAddr, EntityTypeTypeDeviceInformation)

	// THEN: Should handle empty entity address gracefully
	assert.NotNil(s.T(), info.Description.EntityAddress, "EntityAddress should not be nil")
	assert.Nil(s.T(), info.Description.EntityAddress.Device, "Device should be nil")
	assert.Equal(s.T(), emptyEntityAddr, info.Description.EntityAddress.Entity, "Empty entity slice should be preserved")

	// Validate XSD compliance
	err := info.ValidateXSD()
	assert.NoError(s.T(), err, "Empty entity address should still be XSD compliant")
}

// Test comparison with manually created entity information to show the difference
func (s *NodeManagementXSDComplianceSuite) Test_ManualVsFactory_XSDCompliance_EntiyInformation() {
	entityAddr := []AddressEntityType{1, 2}
	entityType := EntityTypeTypeCEM

	// GIVEN: Manually created entity information (current approach)
	manualInfo := &NodeManagementDetailedDiscoveryEntityInformationType{
		Description: &NetworkManagementEntityDescriptionDataType{
			EntityAddress: &EntityAddressType{
				Device: util.Ptr(AddressDeviceType("SomeDevice")), // This violates XSD
				Entity: entityAddr,
			},
			EntityType: &entityType,
		},
	}

	// GIVEN: Factory-created entity information (XSD compliant)
	factoryInfo := NewEntityInformationForNodeManagement(entityAddr, entityType)

	// WHEN: Validating both
	manualErr := manualInfo.ValidateXSD()
	factoryErr := factoryInfo.ValidateXSD()

	// THEN: Manual creation should fail, factory should pass
	assert.Error(s.T(), manualErr, "Manually created info with Device field should fail XSD validation")
	assert.NoError(s.T(), factoryErr, "Factory-created info should pass XSD validation")
}

// Test that NewServerFeatureInformationForNodeManagement creates XSD-compliant feature information
// This test should FAIL initially to drive the implementation
func (s *NodeManagementXSDComplianceSuite) Test_NewServerFeatureInformationForNodeManagement_XSDCompliant() {
	// GIVEN: Valid feature address and type
	entityAddr := []AddressEntityType{0, 1}
	featureAddr := AddressFeatureType(1)
	featureType := FeatureTypeTypeMeasurement
	featureRole := RoleTypeServer
	featureDesc := DescriptionType("")
	supportedFuns := []FunctionPropertyType{}

	// WHEN: Creating feature information for node management using factory function
	info := NewFeatureInformationForNodeManagement(entityAddr, &featureAddr, &featureType, &featureRole, &featureDesc, supportedFuns)

	// THEN: The result should be XSD compliant
	assert.NotNil(s.T(), info, "Feature information should not be nil")
	assert.NotNil(s.T(), info.Description, "Description should not be nil")
	assert.NotNil(s.T(), info.Description.FeatureAddress, "FeatureAddress should not be nil")

	// XSD compliance: Device field must be nil for NodeManagement context
	assert.Nil(s.T(), info.Description.FeatureAddress.Device, "Device field must be nil for XSD compliance")

	// Entity field should be properly set
	assert.Equal(s.T(), entityAddr, info.Description.FeatureAddress.Entity, "Entity field should match input")
	// Feature field should be properly set
	assert.NotNil(s.T(), info.Description.FeatureAddress, "Feature field should not be nil")
	assert.Equal(s.T(), featureAddr, *info.Description.FeatureAddress.Feature, "Feature field should match input")

	// FeatureType should be properly set
	assert.NotNil(s.T(), info.Description.FeatureType, "FeatureType should not be nil")
	assert.Equal(s.T(), featureType, *info.Description.FeatureType, "FeatureType should match input")
	// Role should be properly set
	assert.NotNil(s.T(), info.Description.Role, "Role should not be nil")
	assert.Equal(s.T(), featureRole, *info.Description.Role, "Role should match input")
	// SupportedFunctions should be properly set
	assert.NotNil(s.T(), info.Description.SupportedFunction, "SupportedFunctions should not be nil")
	assert.Equal(s.T(), supportedFuns, info.Description.SupportedFunction, "SupportedFunctions should match input")
}

// Test that NewClientFeatureInformationForNodeManagement creates XSD-compliant feature information
// This test should FAIL initially to drive the implementation
func (s *NodeManagementXSDComplianceSuite) Test_NewClientFeatureInformationForNodeManagement_XSDCompliant() {
	// GIVEN: Valid feature address and type
	entityAddr := []AddressEntityType{0, 1}
	featureAddr := AddressFeatureType(1)
	featureType := FeatureTypeTypeMeasurement
	featureRole := RoleTypeClient
	featureDesc := DescriptionType("")

	// WHEN: Creating feature information for node management using factory function
	info := NewFeatureInformationForNodeManagement(entityAddr, &featureAddr, &featureType, &featureRole, &featureDesc, nil)

	// THEN: The result should be XSD compliant
	assert.NotNil(s.T(), info, "Feature information should not be nil")
	assert.NotNil(s.T(), info.Description, "Description should not be nil")
	assert.NotNil(s.T(), info.Description.FeatureAddress, "FeatureAddress should not be nil")

	// XSD compliance: Device field must be nil for NodeManagement context
	assert.Nil(s.T(), info.Description.FeatureAddress.Device, "Device field must be nil for XSD compliance")

	// Entity field should be properly set
	assert.Equal(s.T(), entityAddr, info.Description.FeatureAddress.Entity, "Entity field should match input")
	// Feature field should be properly set
	assert.NotNil(s.T(), info.Description.FeatureAddress, "Feature field should not be nil")
	assert.Equal(s.T(), featureAddr, *info.Description.FeatureAddress.Feature, "Feature field should match input")

	// FeatureType should be properly set
	assert.NotNil(s.T(), info.Description.FeatureType, "FeatureType should not be nil")
	assert.Equal(s.T(), featureType, *info.Description.FeatureType, "FeatureType should match input")
	// Role should be properly set
	assert.NotNil(s.T(), info.Description.Role, "Role should not be nil")
	assert.Equal(s.T(), featureRole, *info.Description.Role, "Role should match input")
	// SupportedFunctions should be properly set
	assert.Nil(s.T(), info.Description.SupportedFunction, "SupportedFunctions shall be nil")
}

// Test that JSON marshaling excludes the device field for XSD compliance
// This test should FAIL initially to drive the implementation
func (s *NodeManagementXSDComplianceSuite) Test_FeatureInformation_JSONMarshal_NoDeviceField() {
	// GIVEN: Valid feature address and type
	entityAddr := []AddressEntityType{0, 1}
	featureAddr := AddressFeatureType(1)
	featureType := FeatureTypeTypeMeasurement
	featureRole := RoleTypeClient
	featureDesc := DescriptionType("")

	// WHEN: Creating feature information for node management using factory function
	info := NewFeatureInformationForNodeManagement(entityAddr, &featureAddr, &featureType, &featureRole, &featureDesc, nil)

	// WHEN: Marshaling to JSON
	jsonData, err := json.Marshal(info)
	assert.NoError(s.T(), err, "JSON marshaling should not fail")

	// THEN: The JSON should not contain a device field
	jsonString := string(jsonData)
	assert.NotContains(s.T(), jsonString, "device", "JSON should not contain device field for XSD compliance")

	// BUT: Should contain entity field
	assert.Contains(s.T(), jsonString, "entity", "JSON should contain entity field")
	// BUT: Should contain feature field
	assert.Contains(s.T(), jsonString, "feature", "JSON should contain feature field")

	// Verify the structure by unmarshaling back
	var result map[string]interface{}
	err = json.Unmarshal(jsonData, &result)
	assert.NoError(s.T(), err, "JSON should be valid")

	description, ok := result["description"].(map[string]interface{})
	assert.True(s.T(), ok, "Description should be present")

	featureAddress, ok := description["featureAddress"].(map[string]interface{})
	assert.True(s.T(), ok, "FeatureAddress should be present")

	// Critical XSD compliance check
	_, hasDevice := featureAddress["device"]
	assert.False(s.T(), hasDevice, "FeatureAddress should NOT have device field")

	entity, hasEntity := featureAddress["entity"]
	assert.True(s.T(), hasEntity, "EntityAddress should have entity field")
	assert.NotNil(s.T(), entity, "Entity field should not be nil")

	feature, hasfeature := featureAddress["feature"]
	assert.True(s.T(), hasfeature, "FeatureAddress should have feature field")
	assert.NotNil(s.T(), feature, "Feature field should not be nil")
}

// Test XSD validation method
// This test should FAIL initially to drive the implementation
func (s *NodeManagementXSDComplianceSuite) Test_FeatureInformation_ValidateXSD() {
	// GIVEN: XSD-compliant feature information
	entityAddr := []AddressEntityType{1}
	featureAddr := AddressFeatureType(1)
	featureType := FeatureTypeTypeMeasurement
	featureRole := RoleTypeServer
	featureDesc := DescriptionType("")
	supportedFuns := []FunctionPropertyType{}
	validInfo := NewFeatureInformationForNodeManagement(entityAddr, &featureAddr, &featureType, &featureRole, &featureDesc, supportedFuns)

	// WHEN: Validating XSD compliance
	err := validInfo.ValidateXSD()

	// THEN: Should pass validation
	assert.NoError(s.T(), err, "XSD-compliant entity information should pass validation")

	// GIVEN: Feature information with Device field set (XSD violation)
	invalidInfo := &NodeManagementDetailedDiscoveryFeatureInformationType{
		Description: &NetworkManagementFeatureDescriptionDataType{
			FeatureAddress: &FeatureAddressType{
				Device:  util.Ptr(AddressDeviceType("InvalidDevice")), // XSD violation
				Entity:  []AddressEntityType{1},
				Feature: util.Ptr(AddressFeatureType(1)),
			},
			FeatureType: util.Ptr(FeatureTypeTypeMeasurement),
		},
	}

	// WHEN: Validating XSD compliance
	err = invalidInfo.ValidateXSD()

	// THEN: Should fail validation
	assert.Error(s.T(), err, "Feature information with Device field should fail XSD validation")
	assert.Contains(s.T(), err.Error(), "XSD violation", "Error should mention XSD violation")
	assert.Contains(s.T(), err.Error(), "Device field", "Error should mention Device field")
}

// Test comparison with manually created feature information to show the difference
func (s *NodeManagementXSDComplianceSuite) Test_ManualVsFactory_XSDCompliance_FeatureInformation() {
	entityAddr := []AddressEntityType{1, 2}
	featureAddr := AddressFeatureType(1)
	featureType := FeatureTypeTypeMeasurement
	featureRole := RoleTypeClient
	featureDesc := DescriptionType("")

	// GIVEN: Manually created feature information (current approach)
	manualInfo := &NodeManagementDetailedDiscoveryFeatureInformationType{
		Description: &NetworkManagementFeatureDescriptionDataType{
			FeatureAddress: &FeatureAddressType{
				Device:  util.Ptr(AddressDeviceType("SomeDevice")), // This violates XSD
				Entity:  entityAddr,
				Feature: &featureAddr,
			},
			FeatureType: &featureType,
			Role:        &featureRole,
			Description: &featureDesc,
		},
	}

	// GIVEN: Factory-created feature information (XSD compliant)
	factoryInfo := NewFeatureInformationForNodeManagement(entityAddr, &featureAddr, &featureType, &featureRole, &featureDesc, nil)

	// WHEN: Validating both
	manualErr := manualInfo.ValidateXSD()
	factoryErr := factoryInfo.ValidateXSD()

	// THEN: Manual creation should fail, factory should pass
	assert.Error(s.T(), manualErr, "Manually created info with Device field should fail XSD validation")
	assert.NoError(s.T(), factoryErr, "Factory-created info should pass XSD validation")
}
