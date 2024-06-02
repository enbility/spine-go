package spine

import (
	"sync"
	"time"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/rickb777/date/period"
)

const defaultMaxResponseDelay = time.Duration(time.Second * 10)

type FeatureRemote struct {
	*Feature

	entity           api.EntityRemoteInterface
	functionDataMap  map[model.FunctionType]api.FunctionDataInterface
	maxResponseDelay *time.Duration

	mux sync.Mutex
}

func NewFeatureRemote(id uint, entity api.EntityRemoteInterface, ftype model.FeatureTypeType, role model.RoleType) *FeatureRemote {
	res := &FeatureRemote{
		Feature: NewFeature(
			featureAddressType(id, entity.Address()),
			ftype,
			role),
		entity:          entity,
		functionDataMap: make(map[model.FunctionType]api.FunctionDataInterface),
	}
	for _, fd := range CreateFunctionData[api.FunctionDataInterface](ftype) {
		res.functionDataMap[fd.FunctionType()] = fd
	}

	res.operations = make(map[model.FunctionType]api.OperationsInterface)

	return res
}

var _ api.FeatureRemoteInterface = (*FeatureRemote)(nil)

/* FeatureRemoteInterface */

func (r *FeatureRemote) Device() api.DeviceRemoteInterface {
	return r.entity.Device()
}

func (r *FeatureRemote) Entity() api.EntityRemoteInterface {
	return r.entity
}
func (r *FeatureRemote) DataCopy(function model.FunctionType) any {
	r.mux.Lock()
	defer r.mux.Unlock()

	fd := r.functionData(function)
	if fd == nil {
		return nil
	}

	return r.functionData(function).DataCopyAny()
}

func (r *FeatureRemote) UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) *model.ErrorType {
	r.mux.Lock()
	defer r.mux.Unlock()

	fd := r.functionData(function)
	if fd == nil {
		return model.NewErrorTypeFromString("function data not found")
	}

	return fd.UpdateDataAny(false, data, filterPartial, filterDelete)
}

func (r *FeatureRemote) SetOperations(functions []model.FunctionPropertyType) {
	r.operations = make(map[model.FunctionType]api.OperationsInterface)
	for _, sf := range functions {
		if sf.PossibleOperations == nil {
			continue
		}
		r.operations[*sf.Function] = NewOperations(
			sf.PossibleOperations.Read != nil,
			sf.PossibleOperations.Read != nil && sf.PossibleOperations.Read.Partial != nil,
			sf.PossibleOperations.Write != nil,
			sf.PossibleOperations.Write != nil && sf.PossibleOperations.Write.Partial != nil,
		)
	}
}

func (r *FeatureRemote) SetMaxResponseDelay(delay *model.MaxResponseDelayType) {
	if delay == nil {
		return
	}
	p, err := period.Parse(string(*delay))
	if err != nil {
		r.maxResponseDelay = util.Ptr(p.DurationApprox())
	} else {
		logging.Log().Debug(err)
	}
}

func (r *FeatureRemote) MaxResponseDelayDuration() time.Duration {
	if r.maxResponseDelay != nil {
		return *r.maxResponseDelay
	}
	return defaultMaxResponseDelay
}

func (r *FeatureRemote) functionData(function model.FunctionType) api.FunctionDataInterface {
	fd, found := r.functionDataMap[function]
	if !found {
		logging.Log().Errorf("Data was not found for function '%s'", function)
		return nil
	}
	return fd
}
