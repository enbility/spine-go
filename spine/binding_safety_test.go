package spine

import (
	"testing"
	"time"

	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// BindingSafetyTestSuite demonstrates why single binding per server feature is a safety mechanism
type BindingSafetyTestSuite struct {
	suite.Suite

	// Devices that PROVIDE data/control points (have server features)
	evse       *DeviceLocal  // Wallbox with LoadControl server feature
	ev         *DeviceRemote // EV with Measurement server features
	smartMeter *DeviceRemote // Smart meter with Measurement server features

	// Devices that CONTROL/CONSUME (have client features)
	energyManager1 *DeviceRemote // Energy manager 1 wanting to control EVSE
	energyManager2 *DeviceRemote // Energy manager 2 wanting to control EVSE
	hems           *DeviceRemote // Home energy management system

	// Server features (on devices that provide data/control)
	evseLoadControl *FeatureLocal // Can be controlled
	evseEntity      *EntityLocal

	// Client features (on devices that control/consume)
	em1LoadControl *FeatureRemote // Wants to control EVSE
	em1Entity      *EntityRemote
	em2LoadControl *FeatureRemote // Also wants to control EVSE
	em2Entity      *EntityRemote

	// EV server features
	evMeasurement *FeatureRemote
	evEntity      *EntityRemote

	// Smart meter server features
	meterMeasurement *FeatureRemote
	meterEntity      *EntityRemote

	// HEMS client features for reading
	hemsMeasurement1 *FeatureRemote // For reading from EV
	hemsMeasurement2 *FeatureRemote // For reading from smart meter
	hemsEntity       *EntityRemote

	bindingManager *BindingManager
}

func (s *BindingSafetyTestSuite) SetupTest() {
	// Create EVSE (wallbox) - has SERVER features that can be controlled
	s.evse = NewDeviceLocal("ABB", "Terra AC", "12345", "TAC-22",
		"EVSE-Address", model.DeviceTypeTypeChargingStation, model.NetworkManagementFeatureSetTypeSmart)

	s.evseEntity = NewEntityLocal(s.evse, model.EntityTypeTypeEVSE, []model.AddressEntityType{1}, time.Second*4)
	s.evse.AddEntity(s.evseEntity)

	// EVSE has LoadControl SERVER feature - it can be controlled by energy managers
	s.evseLoadControl = s.evseEntity.GetOrAddFeature(model.FeatureTypeTypeLoadControl, model.RoleTypeServer).(*FeatureLocal)

	// Create Energy Manager 1 - has CLIENT features to control devices
	s.energyManager1 = NewDeviceRemote(s.evse, "em1-ski", nil)
	em1Address := model.AddressDeviceType("em1-address")
	s.energyManager1.UpdateDevice(&model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{Device: &em1Address},
	})

	s.em1Entity = NewEntityRemote(s.energyManager1, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	s.energyManager1.AddEntity(s.em1Entity)

	// Energy Manager 1 has LoadControl CLIENT feature to control EVSEs
	s.em1LoadControl = NewFeatureRemote(s.em1Entity.NextFeatureId(), s.em1Entity,
		model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	s.em1LoadControl.Address().Device = util.Ptr(em1Address)
	s.em1Entity.AddFeature(s.em1LoadControl)

	// Create Energy Manager 2 - also wants to control the EVSE
	s.energyManager2 = NewDeviceRemote(s.evse, "em2-ski", nil)
	em2Address := model.AddressDeviceType("em2-address")
	s.energyManager2.UpdateDevice(&model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{Device: &em2Address},
	})

	s.em2Entity = NewEntityRemote(s.energyManager2, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	s.energyManager2.AddEntity(s.em2Entity)

	s.em2LoadControl = NewFeatureRemote(s.em2Entity.NextFeatureId(), s.em2Entity,
		model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	s.em2LoadControl.Address().Device = util.Ptr(em2Address)
	s.em2Entity.AddFeature(s.em2LoadControl)

	// Create EV - has SERVER features providing measurement data
	s.ev = NewDeviceRemote(s.evse, "ev-ski", nil)
	evAddress := model.AddressDeviceType("ev-address")
	s.ev.UpdateDevice(&model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{Device: &evAddress},
	})

	s.evEntity = NewEntityRemote(s.ev, model.EntityTypeTypeEV, []model.AddressEntityType{1})
	s.ev.AddEntity(s.evEntity)

	// EV has Measurement SERVER feature - it provides its own measurement data
	s.evMeasurement = NewFeatureRemote(s.evEntity.NextFeatureId(), s.evEntity,
		model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	s.evMeasurement.Address().Device = util.Ptr(evAddress)
	s.evEntity.AddFeature(s.evMeasurement)

	// Create Smart Meter - has SERVER features providing grid data
	s.smartMeter = NewDeviceRemote(s.evse, "meter-ski", nil)
	meterAddress := model.AddressDeviceType("meter-address")
	s.smartMeter.UpdateDevice(&model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{Device: &meterAddress},
	})

	s.meterEntity = NewEntityRemote(s.smartMeter, model.EntityTypeTypeSubMeterElectricity, []model.AddressEntityType{1})
	s.smartMeter.AddEntity(s.meterEntity)

	s.meterMeasurement = NewFeatureRemote(s.meterEntity.NextFeatureId(), s.meterEntity,
		model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	s.meterMeasurement.Address().Device = util.Ptr(meterAddress)
	s.meterEntity.AddFeature(s.meterMeasurement)

	// Create HEMS - has CLIENT features to read from multiple devices
	s.hems = NewDeviceRemote(s.evse, "hems-ski", nil)
	hemsAddress := model.AddressDeviceType("hems-address")
	s.hems.UpdateDevice(&model.NetworkManagementDeviceDescriptionDataType{
		DeviceAddress: &model.DeviceAddressType{Device: &hemsAddress},
	})

	s.hemsEntity = NewEntityRemote(s.hems, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	s.hems.AddEntity(s.hemsEntity)

	// HEMS has multiple Measurement CLIENT features to read from different devices
	s.hemsMeasurement1 = NewFeatureRemote(s.hemsEntity.NextFeatureId(), s.hemsEntity,
		model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	s.hemsMeasurement1.Address().Device = util.Ptr(hemsAddress)
	s.hemsEntity.AddFeature(s.hemsMeasurement1)

	s.hemsMeasurement2 = NewFeatureRemote(s.hemsEntity.NextFeatureId(), s.hemsEntity,
		model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	s.hemsMeasurement2.Address().Device = util.Ptr(hemsAddress)
	s.hemsEntity.AddFeature(s.hemsMeasurement2)

	s.bindingManager = s.evse.BindingManager().(*BindingManager)
}

func (s *BindingSafetyTestSuite) Test_Control_Conflict_Prevention() {
	// This test demonstrates why allowing multiple energy managers to control
	// the same EVSE LoadControl feature would be dangerous

	// Energy Manager 1 binds to EVSE LoadControl
	binding1 := model.BindingManagementRequestCallType{
		ClientAddress:     s.em1LoadControl.Address(),  // EM1's client feature
		ServerAddress:     s.evseLoadControl.Address(), // EVSE's server feature
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	}
	err := s.bindingManager.AddBinding(s.energyManager1, binding1)
	assert.Nil(s.T(), err, "First energy manager should bind successfully")

	// Energy Manager 2 tries to bind to the SAME EVSE LoadControl
	// This would create conflicts: EM1 says "start charging", EM2 says "stop charging"
	binding2 := model.BindingManagementRequestCallType{
		ClientAddress:     s.em2LoadControl.Address(),  // EM2's client feature
		ServerAddress:     s.evseLoadControl.Address(), // Same EVSE server feature!
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	}
	err = s.bindingManager.AddBinding(s.energyManager2, binding2)
	assert.NotNil(s.T(), err, "Second energy manager should be prevented from binding")
	assert.Equal(s.T(), "the server feature already has a binding", err.Error())

	// This prevention is GOOD because it avoids:
	// 1. Conflicting commands (start vs stop)
	// 2. Notification loops (each change triggers the other to react)
	// 3. Unpredictable behavior (who wins?)
}

func (s *BindingSafetyTestSuite) Test_Sequential_Control_Transfer() {
	// This test shows how control can be transferred between energy managers
	// This is the safe way to handle multiple potential controllers

	// Energy Manager 1 takes control
	binding1 := model.BindingManagementRequestCallType{
		ClientAddress:     s.em1LoadControl.Address(),
		ServerAddress:     s.evseLoadControl.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	}
	err := s.bindingManager.AddBinding(s.energyManager1, binding1)
	assert.Nil(s.T(), err, "EM1 takes control")

	// Energy Manager 1 releases control
	unbinding1 := model.BindingManagementDeleteCallType{
		ClientAddress: s.em1LoadControl.Address(),
		ServerAddress: s.evseLoadControl.Address(),
	}
	err = s.bindingManager.RemoveBinding(s.energyManager1, unbinding1)
	assert.Nil(s.T(), err, "EM1 releases control")

	// Now Energy Manager 2 can take control
	binding2 := model.BindingManagementRequestCallType{
		ClientAddress:     s.em2LoadControl.Address(),
		ServerAddress:     s.evseLoadControl.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeLoadControl),
	}
	err = s.bindingManager.AddBinding(s.energyManager2, binding2)
	assert.Nil(s.T(), err, "EM2 can now take control after EM1 released it")

	// This sequential transfer prevents conflicts while allowing flexibility
}

func TestBindingSafetyTestSuite(t *testing.T) {
	suite.Run(t, new(BindingSafetyTestSuite))
}

// Key Insights from these tests:
//
// 1. Server features are on devices that PROVIDE data or control points:
//    - EVs have measurement server features (they measure their own data)
//    - EVSEs have loadcontrol server features (they can be controlled)
//    - Smart meters have measurement server features (they measure grid data)
//
// 2. Client features are on devices that CONSUME data or CONTROL:
//    - Energy managers have client features to control EVSEs
//    - HEMS have client features to read measurements
//
// 3. The single binding limitation prevents:
//    - Control conflicts (multiple managers sending conflicting commands)
//    - Notification loops (changes triggering endless reactions)
//    - Race conditions (unpredictable command ordering)
//
// 4. This is NOT a limitation but a SAFETY FEATURE:
//    - The spec doesn't define conflict resolution
//    - Multiple writers would create chaos
//    - One controller per feature ensures predictability
