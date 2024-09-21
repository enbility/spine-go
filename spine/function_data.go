package spine

import (
	"fmt"
	"sync"

	"github.com/enbility/ship-go/logging"
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type FunctionData[T any] struct {
	functionType model.FunctionType
	data         *T

	mux sync.Mutex
}

func NewFunctionData[T any](function model.FunctionType) *FunctionData[T] {
	return &FunctionData[T]{
		functionType: function,
	}
}

var _ api.FunctionDataInterface = (*FunctionData[int])(nil)

/* FunctionDataInterface */

func (r *FunctionData[T]) FunctionType() model.FunctionType {
	return r.functionType
}

func (r *FunctionData[T]) SupportsPartialWrite() bool {
	return util.Implements[T, model.Updater]()
}

func (r *FunctionData[T]) DataCopy() *T {
	r.mux.Lock()
	defer r.mux.Unlock()

	// copy the data and return it as the data can be updated
	// and newly assigned at any time otherwise we run into panics
	// because of invalid memory address or nil pointer dereference
	var copiedData T
	if r.data == nil {
		return nil
	}

	copiedData = *r.data

	return &copiedData
}

func (r *FunctionData[T]) UpdateData(remoteWrite, persist bool, newData *T, filterPartial *model.FilterType, filterDelete *model.FilterType) (any, *model.ErrorType) {
	r.mux.Lock()
	defer r.mux.Unlock()

	if filterPartial == nil && filterDelete == nil && persist {
		// just set the data
		r.data = newData
		return r.data, nil
	}

	if !r.SupportsPartialWrite() {
		return nil, model.NewErrorTypeFromString(fmt.Sprintf("partial updates are not supported for type '%s'", util.Type[T]().Name()))
	}

	if r.data == nil {
		r.data = new(T)
	}

	updater := any(r.data).(model.Updater)
	data, success := updater.UpdateList(remoteWrite, persist, newData, filterPartial, filterDelete)
	if !success {
		return nil, model.NewErrorTypeFromString("update failed, likely not allowed to write")
	}

	return data, nil
}

func (r *FunctionData[T]) DataCopyAny() any {
	return r.DataCopy()
}

func (r *FunctionData[T]) UpdateDataAny(remoteWrite, persist bool, newData any, filterPartial *model.FilterType, filterDelete *model.FilterType) (any, *model.ErrorType) {
	data, err := r.UpdateData(remoteWrite, persist, newData.(*T), filterPartial, filterDelete)
	if err != nil {
		logging.Log().Debug(err.String())
	}

	return data, err
}
