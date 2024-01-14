package spine

import (
	"fmt"
	"sync"
	"time"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/rickb777/date/period"
)

const defaultMaxResponseDelay = time.Duration(time.Second * 10)

var _ api.FeatureRemote = (*FeatureRemoteImpl)(nil)

type FeatureRemoteImpl struct {
	*FeatureImpl

	entity           api.EntityRemote
	functionDataMap  map[model.FunctionType]api.FunctionData
	maxResponseDelay *time.Duration

	mux sync.Mutex
}

func NewFeatureRemoteImpl(id uint, entity api.EntityRemote, ftype model.FeatureTypeType, role model.RoleType) *FeatureRemoteImpl {
	res := &FeatureRemoteImpl{
		FeatureImpl: NewFeatureImpl(
			featureAddressType(id, entity.Address()),
			ftype,
			role),
		entity:          entity,
		functionDataMap: make(map[model.FunctionType]api.FunctionData),
	}
	for _, fd := range CreateFunctionData[api.FunctionData](ftype) {
		res.functionDataMap[fd.Function()] = fd
	}

	res.operations = make(map[model.FunctionType]api.Operations)

	return res
}

func (r *FeatureRemoteImpl) DataCopy(function model.FunctionType) any {
	r.mux.Lock()
	defer r.mux.Unlock()

	return r.functionData(function).DataCopyAny()
}

func (r *FeatureRemoteImpl) SetData(function model.FunctionType, data any) {
	r.mux.Lock()

	fd := r.functionData(function)
	fd.UpdateDataAny(data, nil, nil)

	r.mux.Unlock()
}

func (r *FeatureRemoteImpl) UpdateData(function model.FunctionType, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.functionData(function).UpdateDataAny(data, filterPartial, filterDelete)
	// TODO: fire event
}

func (r *FeatureRemoteImpl) Sender() api.Sender {
	return r.Device().Sender()
}

func (r *FeatureRemoteImpl) Device() api.DeviceRemote {
	return r.entity.Device()
}

func (r *FeatureRemoteImpl) Entity() api.EntityRemote {
	return r.entity
}

func (r *FeatureRemoteImpl) SetOperations(functions []model.FunctionPropertyType) {
	r.operations = make(map[model.FunctionType]api.Operations)
	for _, sf := range functions {
		r.operations[*sf.Function] = NewOperations(sf.PossibleOperations.Read != nil, sf.PossibleOperations.Write != nil)
	}
}

func (r *FeatureRemoteImpl) SetMaxResponseDelay(delay *model.MaxResponseDelayType) {
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

func (r *FeatureRemoteImpl) MaxResponseDelayDuration() time.Duration {
	if r.maxResponseDelay != nil {
		return *r.maxResponseDelay
	}
	return defaultMaxResponseDelay
}

func (r *FeatureRemoteImpl) functionData(function model.FunctionType) api.FunctionData {
	fd, found := r.functionDataMap[function]
	if !found {
		panic(fmt.Errorf("Data was not found for function '%s'", function))
	}
	return fd
}
