package spine

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Operations struct {
	read, write               bool
	readPartial, writePartial bool
}

var _ api.OperationsInterface = (*Operations)(nil)

func NewOperations(read, readPartial, write, writePartial bool) *Operations {
	return &Operations{
		read:         read,
		readPartial:  readPartial,
		write:        write,
		writePartial: writePartial,
	}
}
func (r *Operations) Read() bool {
	return r.read
}

func (r *Operations) ReadPartial() bool {
	return r.readPartial
}

func (r *Operations) Write() bool {
	return r.write
}

func (r *Operations) WritePartial() bool {
	return r.writePartial
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
		if r.readPartial {
			res.Read = &model.PossibleOperationsReadType{
				Partial: &model.ElementTagType{},
			}
		}
	}
	if r.write {
		res.Write = &model.PossibleOperationsWriteType{}
		if r.writePartial {
			res.Write = &model.PossibleOperationsWriteType{
				Partial: &model.ElementTagType{},
			}
		}
	}

	return res
}
