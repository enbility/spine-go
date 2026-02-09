package spine

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type DeviceRemote struct {
	*Device

	ski string

	entities      []api.EntityRemoteInterface
	entitiesMutex sync.Mutex

	sender api.SenderInterface

	localDevice api.DeviceLocalInterface

	// Version tracking fields
	supportedVersions []string      // List of versions supported by remote device
	negotiatedVersion string        // The negotiated common version
	
	// Version detection fields
	estimatedRemoteVersion string   // Our estimate of remote's version based on compatibility
	detectedRemoteVersion  string   // Actual version seen in remote's messages
	versionChanged         bool     // Track if version has changed
	
	versionsMutex     sync.RWMutex  // Mutex for thread-safe access to version fields
}

func NewDeviceRemote(localDevice api.DeviceLocalInterface, ski string, sender api.SenderInterface) *DeviceRemote {
	res := DeviceRemote{
		Device:      NewDevice(nil, nil, nil),
		ski:         ski,
		localDevice: localDevice,
		sender:      sender,
	}
	res.addNodeManagement()

	return &res
}

func (d *DeviceRemote) addNodeManagement() {
	deviceInformation := d.addNewEntity(model.EntityTypeTypeDeviceInformation, NewAddressEntityType([]uint{DeviceInformationEntityId}))
	nodeManagement := NewFeatureRemote(deviceInformation.NextFeatureId(), deviceInformation, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	deviceInformation.AddFeature(nodeManagement)
}

var _ shipapi.ShipConnectionDataReaderInterface = (*DeviceRemote)(nil)

/* ShipConnectionDataReaderInterface */

// processing incoming SPINE message from the associated SHIP connection
func (d *DeviceRemote) HandleShipPayloadMessage(message []byte) {
	if _, err := d.HandleSpineMesssage(message); err != nil {
		logging.Log().Errorf("error handling spine message", err)
	}
}

var _ api.DeviceRemoteInterface = (*DeviceRemote)(nil)

/* DeviceRemoteInterface */

// return the device SKI
func (d *DeviceRemote) Ski() string {
	return d.ski
}

func (d *DeviceRemote) AddEntity(entity api.EntityRemoteInterface) {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	d.entities = append(d.entities, entity)
}

func (d *DeviceRemote) addNewEntity(eType model.EntityTypeType, address []model.AddressEntityType) api.EntityRemoteInterface {
	newEntity := NewEntityRemote(d, eType, address)
	d.AddEntity(newEntity)
	return newEntity
}

// Remove an entity with a given address from this device
func (d *DeviceRemote) RemoveEntityByAddress(addr []model.AddressEntityType) api.EntityRemoteInterface {
	entityForRemoval := d.Entity(addr)
	if entityForRemoval == nil {
		return nil
	}

	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	var newEntities []api.EntityRemoteInterface
	for _, item := range d.entities {
		if !reflect.DeepEqual(item, entityForRemoval) {
			newEntities = append(newEntities, item)
		}
	}
	d.entities = newEntities

	return entityForRemoval
}

// Return an entity with a given address
func (d *DeviceRemote) Entity(id []model.AddressEntityType) api.EntityRemoteInterface {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	for _, e := range d.entities {
		if reflect.DeepEqual(id, e.Address().Entity) {
			return e
		}
	}
	return nil
}

// Return all entities of this device
func (d *DeviceRemote) Entities() []api.EntityRemoteInterface {
	d.entitiesMutex.Lock()
	defer d.entitiesMutex.Unlock()

	return d.entities
}

// Return the feature for a given address
func (d *DeviceRemote) FeatureByAddress(address *model.FeatureAddressType) api.FeatureRemoteInterface {
	entity := d.Entity(address.Entity)
	if entity != nil {
		return entity.FeatureOfAddress(address.Feature)
	}
	return nil
}

// Get the feature for a given entity, feature type and feature role
func (r *DeviceRemote) FeatureByEntityTypeAndRole(entity api.EntityRemoteInterface, featureType model.FeatureTypeType, role model.RoleType) api.FeatureRemoteInterface {
	if len(r.entities) < 1 {
		return nil
	}

	r.entitiesMutex.Lock()
	defer r.entitiesMutex.Unlock()

	for _, e := range r.entities {
		if entity != e {
			continue
		}
		for _, feature := range entity.Features() {
			if feature.Type() == featureType && feature.Role() == role {
				return feature
			}
		}
	}

	return nil
}

func (d *DeviceRemote) HandleSpineMesssage(message []byte) (*model.MsgCounterType, error) {
	datagram := model.Datagram{}
	if err := json.Unmarshal([]byte(message), &datagram); err != nil {
		return nil, err
	}

	// Check if this is a discovery message (which establishes version negotiation)
	isDiscoveryMessage := d.isDiscoveryMessage(&datagram.Datagram)
	
	// Validate protocol version
	if err := d.validateProtocolVersion(datagram.Datagram.Header.SpecificationVersion, isDiscoveryMessage); err != nil {
		// Send error response if appropriate
		d.sendVersionErrorResponse(&datagram.Datagram, err)
		return nil, err
	}

	if datagram.Datagram.Header.MsgCounterReference != nil {
		d.sender.ProcessResponseForMsgCounterReference(datagram.Datagram.Header.MsgCounterReference)
	}

	err := d.localDevice.ProcessCmd(datagram.Datagram, d)
	if err != nil {
		logging.Log().Trace(err)
		// Only propagate version incompatibility errors, preserve original behavior for others
		if IsVersionIncompatibilityError(err) {
			return datagram.Datagram.Header.MsgCounter, err
		}
	}

	return datagram.Datagram.Header.MsgCounter, nil
}

func (d *DeviceRemote) Sender() api.SenderInterface {
	return d.sender
}

func (d *DeviceRemote) UseCases() []model.UseCaseInformationDataType {
	entity := d.Entity(DeviceInformationAddressEntity)

	nodemgmt := d.FeatureByEntityTypeAndRole(entity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)

	data, ok := nodemgmt.DataCopy(model.FunctionTypeNodeManagementUseCaseData).(*model.NodeManagementUseCaseDataType)
	if ok && data != nil {
		return data.UseCaseInformation
	}

	return nil
}

func (d *DeviceRemote) UpdateDevice(description *model.NetworkManagementDeviceDescriptionDataType) {
	if description != nil {
		if description.DeviceAddress != nil && description.DeviceAddress.Device != nil {
			d.address = description.DeviceAddress.Device
		}
		if description.DeviceType != nil {
			d.dType = description.DeviceType
		}
		if description.NetworkFeatureSet != nil {
			d.featureSet = description.NetworkFeatureSet
		}
	}
}

func (d *DeviceRemote) AddEntityAndFeatures(
	initialData bool,
	data *model.NodeManagementDetailedDiscoveryDataType,
	entityAddressToAdd *model.EntityAddressType,
) ([]api.EntityRemoteInterface, error) {
	rEntites := make([]api.EntityRemoteInterface, 0)

	for _, ei := range data.EntityInformation {
		if err := d.CheckEntityInformation(initialData, ei); err != nil {
			return nil, err
		}

		entityAddress := ei.Description.EntityAddress.Entity
		// if entityAddressToAdd, make sure we are adding the correct entity
		if entityAddressToAdd != nil && !reflect.DeepEqual(entityAddress, entityAddressToAdd.Entity) {
			continue
		}

		entity := d.Entity(entityAddress)
		if entity == nil {
			entity = d.addNewEntity(*ei.Description.EntityType, entityAddress)
			rEntites = append(rEntites, entity)
		}

		// make sure the device address is set, which is not on entity 0 on startup !
		if entity.Address().Device == nil || len(*entity.Address().Device) == 0 {
			if data.DeviceInformation != nil &&
				data.DeviceInformation.Description != nil &&
				data.DeviceInformation.Description.DeviceAddress != nil &&
				data.DeviceInformation.Description.DeviceAddress.Device != nil {
				entity.UpdateDeviceAddress(*data.DeviceInformation.Description.DeviceAddress.Device)
			}
		}

		entity.SetDescription(ei.Description.Description)
		entity.RemoveAllFeatures()

		for _, fi := range data.FeatureInformation {
			if reflect.DeepEqual(fi.Description.FeatureAddress.Entity, entityAddress) {
				if f, ok := unmarshalFeature(entity, fi); ok {
					entity.AddFeature(f)
				}
			}
		}
	}

	return rEntites, nil
}

// check if the provided entity information is correct
// provide initialData to check if the entity is new and not an update
func (d *DeviceRemote) CheckEntityInformation(initialData bool, entity model.NodeManagementDetailedDiscoveryEntityInformationType) error {
	description := entity.Description
	if description == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description")
	}

	if description.EntityAddress == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress")
	}

	if description.EntityAddress.Entity == nil {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: invalid EntityInformation.Description.EntityAddress.Entity")
	}

	// Consider on initial NodeManagement Detailed Discovery, the device being empty as it is not yet known
	if initialData {
		return nil
	}

	address := d.Address()
	if description.EntityAddress.Device != nil && address != nil && *description.EntityAddress.Device != *address {
		return errors.New("nodemanagement.replyDetailedDiscoveryData: device address mismatch")
	}

	return nil
}

func unmarshalFeature(entity api.EntityRemoteInterface,
	featureData model.NodeManagementDetailedDiscoveryFeatureInformationType,
) (api.FeatureRemoteInterface, bool) {
	var result api.FeatureRemoteInterface

	fid := featureData.Description

	if fid == nil {
		return nil, false
	}

	result = NewFeatureRemote(uint(*fid.FeatureAddress.Feature), entity, *fid.FeatureType, *fid.Role)

	result.SetDescription(fid.Description)
	result.SetMaxResponseDelay(fid.MaxResponseDelay)
	result.SetOperations(fid.SupportedFunction)

	return result, true
}

// isDiscoveryMessage checks if the datagram contains discovery data
func (d *DeviceRemote) isDiscoveryMessage(datagram *model.DatagramType) bool {
	if len(datagram.Payload.Cmd) == 0 {
		return false
	}
	
	// Check if any command contains NodeManagementDetailedDiscoveryData
	for _, cmd := range datagram.Payload.Cmd {
		if cmd.NodeManagementDetailedDiscoveryData != nil {
			return true
		}
	}
	return false
}

// validateProtocolVersion checks if the incoming message version is valid
// For discovery messages: uses compatibility checking
// For regular messages: detects and validates version
func (d *DeviceRemote) validateProtocolVersion(version *model.SpecificationVersionType, isDiscoveryMessage bool) error {
	// If no version provided, assume compatible (handle real-world devices)
	if version == nil {
		logging.Log().Debug("No protocol version in message from", d.address)
		// Update detected version for tracking
		if !isDiscoveryMessage {
			_ = d.UpdateDetectedVersion("")
		}
		return nil
	}

	versionStr := strings.TrimSpace(string(*version))
	
	// Check SPINE specification limit: max 128 characters
	if len(versionStr) > 128 {
		return fmt.Errorf("version string exceeds SPINE specification limit of 128 characters: %d", len(versionStr))
	}
	
	// Empty version - assume compatible
	if versionStr == "" {
		logging.Log().Debug("Empty protocol version from", d.address)
		// Update detected version for tracking
		if !isDiscoveryMessage {
			_ = d.UpdateDetectedVersion("")
		}
		return nil
	}

	// Discovery messages - accept any version, don't update detected version
	if isDiscoveryMessage {
		logging.Log().Trace("Discovery message - accepting any version:", versionStr)
		return nil
	}

	// Regular messages - detect and track the version
	if err := d.UpdateDetectedVersion(versionStr); err != nil {
		return err
	}

	// Log if detected version differs from estimate
	d.versionsMutex.RLock()
	estimated := d.estimatedRemoteVersion
	d.versionsMutex.RUnlock()
	
	if estimated != "" && estimated != versionStr {
		logging.Log().Debugf("Remote device %s using version %s (estimated: %s)", 
			d.ski, versionStr, estimated)
	}

	// Validate compatibility
	return d.validateVersionCompatibility(versionStr)
}

// validateVersionCompatibility performs compatibility checking for version strings
func (d *DeviceRemote) validateVersionCompatibility(versionStr string) error {
	// Try to parse as semantic version (major.minor.patch)
	major, minor, patch, valid := parseSemanticVersion(versionStr)
	
	if !valid {
		// Not a valid semantic version format - log and assume compatible
		logging.Log().Debug("Non-compliant protocol version format from", d.address, ":", versionStr)
		return nil
	}

	// Valid semantic version - check major version compatibility  
	if major >= 2 {
		return fmt.Errorf("incompatible major version: %s", versionStr)
	}

	// Log successful validation of compliant version
	logging.Log().Trace("Valid protocol version from", d.address, ":", major, ".", minor, ".", patch)

	// Major version 0 or 1 - compatible
	return nil
}

// parseSemanticVersion attempts to parse a version string as major.minor.patch
// Returns the parsed values and whether the format is valid
func parseSemanticVersion(version string) (major, minor, patch int, valid bool) {
	// Split by dots
	parts := strings.Split(version, ".")
	if len(parts) != 3 {
		return 0, 0, 0, false
	}

	// Parse major version
	majorVal, err := parseVersionNumber(parts[0])
	if err != nil || majorVal < 0 {
		return 0, 0, 0, false
	}

	// Parse minor version
	minorVal, err := parseVersionNumber(parts[1])
	if err != nil || minorVal < 0 {
		return 0, 0, 0, false
	}

	// Parse patch version
	patchVal, err := parseVersionNumber(parts[2])
	if err != nil || patchVal < 0 {
		return 0, 0, 0, false
	}

	return majorVal, minorVal, patchVal, true
}

// parseVersionNumber parses a single version number component
func parseVersionNumber(s string) (int, error) {
	// Must be non-empty
	if s == "" {
		return 0, fmt.Errorf("empty version component")
	}

	// Must contain only digits
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("non-numeric character in version component")
		}
	}

	// Convert to int - handle potential overflow by checking length
	if len(s) > 9 { // Prevent integer overflow
		return 0, fmt.Errorf("version number too large")
	}

	val := 0
	for _, r := range s {
		val = val*10 + int(r-'0')
	}

	return val, nil
}

// sendVersionErrorResponse sends an error response for version incompatibility when appropriate
func (d *DeviceRemote) sendVersionErrorResponse(datagram *model.DatagramType, versionErr error) {
	// Only send error response if the message expects a response
	if datagram.Header.CmdClassifier == nil {
		return
	}

	// Don't send error responses to REPLY or RESULT messages (avoid loops)
	if *datagram.Header.CmdClassifier == model.CmdClassifierTypeReply ||
		*datagram.Header.CmdClassifier == model.CmdClassifierTypeResult {
		return
	}

	// Send error for READ, WRITE, or messages with ackRequest
	shouldSendError := false
	if *datagram.Header.CmdClassifier == model.CmdClassifierTypeRead ||
		*datagram.Header.CmdClassifier == model.CmdClassifierTypeWrite {
		shouldSendError = true
	} else if datagram.Header.AckRequest != nil && *datagram.Header.AckRequest {
		shouldSendError = true
	}

	if !shouldSendError {
		return
	}

	// Create appropriate sender address
	senderAddress := &model.FeatureAddressType{
		Device:  d.address,
		Entity:  []model.AddressEntityType{0},
		Feature: util.Ptr(model.AddressFeatureType(0)),
	}

	// Send error response
	errorType := model.NewErrorType(
		model.ErrorNumberTypeGeneralError,
		fmt.Sprintf("Protocol version incompatibility: %s", versionErr.Error()),
	)

	_ = d.sender.ResultError(&datagram.Header, senderAddress, errorType)
}

// UpdateEstimatedRemoteVersion calculates what version the remote device will likely use
// based on its supported versions and compatibility groups
func (d *DeviceRemote) UpdateEstimatedRemoteVersion() {
	d.versionsMutex.Lock()
	defer d.versionsMutex.Unlock()
	
	if len(d.supportedVersions) == 0 {
		d.estimatedRemoteVersion = ""
		return
	}
	
	// Get our local version's major number for compatibility group
	localVersion := string(SpecificationVersion)
	localMajor, _, _, _ := parseSemanticVersion(localVersion)
	
	// Find highest version in compatible group
	// Major versions 0 and 1 are compatible
	highestCompatible := ""
	for _, v := range d.supportedVersions {
		major, _, _, valid := parseSemanticVersion(v)
		if valid && (major == localMajor || (localMajor <= 1 && major <= 1)) {
			if highestCompatible == "" || compareVersions(v, highestCompatible) > 0 {
				highestCompatible = v
			}
		}
	}
	
	d.estimatedRemoteVersion = highestCompatible
}

// EstimatedRemoteVersion returns the estimated version the remote will use
func (d *DeviceRemote) EstimatedRemoteVersion() string {
	d.versionsMutex.RLock()
	defer d.versionsMutex.RUnlock()
	return d.estimatedRemoteVersion
}

// UpdateDetectedVersion updates the detected version from actual messages
func (d *DeviceRemote) UpdateDetectedVersion(version string) error {
	d.versionsMutex.Lock()
	defer d.versionsMutex.Unlock()
	
	// Check if this is an incompatible version
	if version != "" && version != "..." && version != "draft" {
		major, _, _, valid := parseSemanticVersion(version)
		if valid && major >= 2 {
			return fmt.Errorf("incompatible major version: %s", version)
		}
	}
	
	// Track version changes
	if d.detectedRemoteVersion != "" && d.detectedRemoteVersion != version {
		d.versionChanged = true
		logging.Log().Debugf("Remote device %s changed version from %s to %s", 
			d.ski, d.detectedRemoteVersion, version)
	}
	
	d.detectedRemoteVersion = version
	return nil
}

// DetectedRemoteVersion returns the actual version seen in messages
func (d *DeviceRemote) DetectedRemoteVersion() string {
	d.versionsMutex.RLock()
	defer d.versionsMutex.RUnlock()
	return d.detectedRemoteVersion
}

// HasVersionChanged returns whether the remote has changed versions
func (d *DeviceRemote) HasVersionChanged() bool {
	d.versionsMutex.RLock()
	defer d.versionsMutex.RUnlock()
	return d.versionChanged
}

// ValidateDatagramVersion validates the version in a datagram
func (d *DeviceRemote) ValidateDatagramVersion(datagram *model.DatagramType) error {
	// Handle nil datagram
	if datagram == nil {
		return nil
	}
	
	// Discovery messages are never rejected
	if d.isDiscoveryMessage(datagram) {
		return nil
	}
	
	// Extract version
	versionStr := ""
	if datagram.Header.SpecificationVersion != nil {
		versionStr = string(*datagram.Header.SpecificationVersion)
	}
	
	// Update detected version
	if versionStr != "" {
		if err := d.UpdateDetectedVersion(versionStr); err != nil {
			return err
		}
	}
	
	return nil
}

// SetSupportedProtocolVersions stores the supported versions from remote device
func (d *DeviceRemote) SetSupportedProtocolVersions(versions []string) {
	d.versionsMutex.Lock()
	defer d.versionsMutex.Unlock()
	d.supportedVersions = make([]string, len(versions))
	copy(d.supportedVersions, versions)
}

// SupportedProtocolVersions returns the supported versions (thread-safe)
func (d *DeviceRemote) SupportedProtocolVersions() []string {
	d.versionsMutex.RLock()
	defer d.versionsMutex.RUnlock()
	if d.supportedVersions == nil {
		return nil
	}
	result := make([]string, len(d.supportedVersions))
	copy(result, d.supportedVersions)
	return result
}

// SetNegotiatedProtocolVersion stores the negotiated version
func (d *DeviceRemote) SetNegotiatedProtocolVersion(version string) {
	d.versionsMutex.Lock()
	defer d.versionsMutex.Unlock()
	d.negotiatedVersion = version
}

// NegotiatedProtocolVersion returns the negotiated version (thread-safe)
func (d *DeviceRemote) NegotiatedProtocolVersion() string {
	d.versionsMutex.RLock()
	defer d.versionsMutex.RUnlock()
	return d.negotiatedVersion
}

