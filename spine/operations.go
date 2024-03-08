package spine

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Operations struct {
	read, write bool
}

var _ api.OperationsInterface = (*Operations)(nil)

func NewOperations(read, write bool) *Operations {
	return &Operations{
		read:  read,
		write: write,
	}
}

func (r *Operations) Read() bool {
	return r.read
}

func (r *Operations) Write() bool {
	return r.write
}

func (r *Operations) String() string {
	switch {
	case r.read && !r.write:
		return "RO"
	case r.read && r.write:
		return "RW"
	default:
		return "--"
	}
}

func (r *Operations) Information() *model.PossibleOperationsType {
	res := new(model.PossibleOperationsType)
	if r.read {
		res.Read = &model.PossibleOperationsReadType{}
	}
	if r.write {
		res.Write = &model.PossibleOperationsWriteType{
			Partial: &model.ElementTagType{},
		}
	}

	return res
}
