package spine

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type FunctionDataCmd[T any] struct {
	*FunctionData[T]
}

func NewFunctionDataCmd[T any](function model.FunctionType) *FunctionDataCmd[T] {
	return &FunctionDataCmd[T]{
		FunctionData: NewFunctionData[T](function),
	}
}

var _ api.FunctionDataCmdInterface = (*FunctionDataCmd[int])(nil)

/* FunctionDataCmdInterface */

func (r *FunctionDataCmd[T]) ReadCmdType(partialSelector any, elements any) model.CmdType {
	cmd := createCmd[T](r.functionType, nil)

	var filters []model.FilterType
	filters = filtersForSelectorsElements(r.functionType, filters, nil, partialSelector, nil, elements)
	if len(filters) > 0 {
		cmd.Filter = filters
		cmd.Function = util.Ptr(model.FunctionType(""))
	}

	return cmd
}

func (r *FunctionDataCmd[T]) ReplyCmdType(partial bool) model.CmdType {
	data := r.DataCopy()
	cmd := createCmd(r.functionType, data)
	if partial {
		cmd.Filter = filterEmptyPartial()
		cmd.Function = util.Ptr(model.FunctionType(""))
	}
	return cmd
}

func (r *FunctionDataCmd[T]) NotifyOrWriteCmdType(deleteSelector, partialSelector any, partialWithoutSelector bool, deleteElements any) model.CmdType {
	data := r.DataCopy()
	cmd := createCmd(r.functionType, data)

	if partialWithoutSelector {
		cmd.Filter = filterEmptyPartial()
		cmd.Function = util.Ptr(model.FunctionType(r.functionType))
		return cmd
	}
	var filters []model.FilterType
	if filters := filtersForSelectorsElements(r.functionType, filters, deleteSelector, partialSelector, deleteElements, nil); len(filters) > 0 {
		cmd.Filter = filters
		cmd.Function = util.Ptr(model.FunctionType(r.functionType))
	}

	return cmd
}

func filtersForSelectorsElements(functionType model.FunctionType, filters []model.FilterType, deleteSelector, partialSelector any, deleteElements, readElements any) []model.FilterType {
	if deleteSelector != nil || deleteElements != nil {
		filter := model.FilterType{CmdControl: &model.CmdControlType{Delete: &model.ElementTagType{}}}
		if deleteSelector != nil {
			filter = addSelectorToFilter(filter, functionType, &deleteSelector)
		}
		if deleteElements != nil {
			filter = addElementToFilter(filter, functionType, &deleteElements)
		}
		filters = append(filters, filter)
	}

	if partialSelector != nil || readElements != nil {
		filter := model.FilterType{CmdControl: &model.CmdControlType{Partial: &model.ElementTagType{}}}
		if partialSelector != nil {
			filter = addSelectorToFilter(filter, functionType, &partialSelector)
		}
		if readElements != nil {
			filter = addElementToFilter(filter, functionType, &readElements)
		}
		filters = append(filters, filter)
	}

	return filters
}

// simple helper for adding a single filterType without any selectors
func filterEmptyPartial() []model.FilterType {
	return []model.FilterType{*model.NewFilterTypePartial()}
}

func addSelectorToFilter[T any](filter model.FilterType, function model.FunctionType, data *T) model.FilterType {
	result := filter

	result.SetDataForFunction(model.EEBusTagTypeTypeSelector, function, data)

	return result
}

func addElementToFilter[T any](filter model.FilterType, function model.FunctionType, data *T) model.FilterType {
	result := filter

	result.SetDataForFunction(model.EEbusTagTypeTypeElements, function, data)

	return result
}

func createCmd[T any](function model.FunctionType, data *T) model.CmdType {
	result := model.CmdType{}

	result.SetDataForFunction(function, data)

	return result
}
