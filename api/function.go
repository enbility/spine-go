package api

import "github.com/enbility/spine-go/model"

/* Function */

type FunctionDataCmdInterface interface {
	FunctionDataInterface
	ReadCmdType(partialSelector any, elements any) model.CmdType
	ReplyCmdType(partial bool) model.CmdType
	NotifyOrWriteCmdType(deleteSelector, partialSelector any, partialWithoutSelector bool, deleteElements any) model.CmdType
}

type FunctionDataInterface interface {
	Function() model.FunctionType
	DataCopyAny() any
	UpdateDataAny(data any, filterPartial *model.FilterType, filterDelete *model.FilterType)
}
