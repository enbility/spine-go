package api

import "github.com/enbility/spine-go/model"

/* Function */

type FunctionDataCmdInterface interface {
	FunctionDataInterface
	// Get the CmdType data for a read command
	ReadCmdType(partialSelector any, elements any) model.CmdType
	// Get the CmdType data for a reply command
	ReplyCmdType(partial bool) model.CmdType
	// Get the CmdType data for a notify or write command
	NotifyOrWriteCmdType(deleteSelector, partialSelector any, partialWithoutSelector bool, deleteElements any) model.CmdType
}

type FunctionDataInterface interface {
	// Get the function type
	FunctionType() model.FunctionType
	// Get a copy of the functions data
	DataCopyAny() any
	// Update the functions data
	UpdateDataAny(remoteWrite bool, data any, filterPartial *model.FilterType, filterDelete *model.FilterType) *model.ErrorType
}
