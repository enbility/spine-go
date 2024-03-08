package model

import (
	"fmt"
	"testing"

	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Id         *uint `eebus:"key"`
	Changeable *bool `eebus:"writecheck"`
	Active     *bool
	Data       *string
}

type testInvalidStruct struct {
	Id         *uint   `eebus:"key"`
	Changeable *string `eebus:"writecheck"`
	Data       *string
}

func (r testStruct) HashKey() string {
	return fmt.Sprintf("%d", r.Id)
}

func TestUnion_NewData(t *testing.T) {
	existingData := []testStruct{
		{Id: util.Ptr(uint(1)), Changeable: util.Ptr(true), Data: util.Ptr(string("data1"))},
	}

	newData := []testStruct{
		{Id: util.Ptr(uint(2)), Data: util.Ptr(string("data2"))},
	}

	// Act
	result := Merge(false, existingData, newData)

	if assert.Equal(t, 2, len(result)) {
		assert.Equal(t, 1, int(*result[0].Id))
		assert.Equal(t, "data1", string(*result[0].Data))
		assert.Equal(t, 2, int(*result[1].Id))
		assert.Equal(t, "data2", string(*result[1].Data))
	}
}

func TestUnion_NewAndUpdateData(t *testing.T) {
	existingData := []testStruct{
		{Id: util.Ptr(uint(0)), Changeable: util.Ptr(true), Active: util.Ptr(true), Data: util.Ptr(string("data1"))},
		{Id: util.Ptr(uint(1)), Changeable: util.Ptr(true), Active: util.Ptr(false), Data: util.Ptr(string("data2"))},
		{Id: util.Ptr(uint(2)), Changeable: util.Ptr(false), Active: util.Ptr(false), Data: util.Ptr(string("data3"))},
		{Id: util.Ptr(uint(3)), Active: util.Ptr(true), Data: util.Ptr(string("data4"))},
	}

	newData := []testStruct{
		{Id: util.Ptr(uint(1)), Changeable: util.Ptr(false), Active: util.Ptr(true), Data: util.Ptr(string("data22"))},
		{Id: util.Ptr(uint(2)), Data: util.Ptr(string("data33"))},
		{Id: util.Ptr(uint(3)), Data: util.Ptr(string("data44"))},
	}

	// Act
	result := Merge(false, existingData, newData)

	if assert.Equal(t, 4, len(result)) {
		assert.Equal(t, 0, int(*result[0].Id))
		assert.Equal(t, "data1", string(*result[0].Data))
		assert.Equal(t, true, bool(*result[0].Active))
		assert.Equal(t, 1, int(*result[1].Id))
		assert.Equal(t, "data22", string(*result[1].Data))
		assert.Equal(t, false, bool(*result[1].Changeable))
		assert.Equal(t, true, bool(*result[1].Active))
		assert.Equal(t, 2, int(*result[2].Id))
		assert.Equal(t, "data33", string(*result[2].Data))
		assert.Equal(t, false, bool(*result[2].Active))
		assert.Equal(t, 3, int(*result[3].Id))
		assert.Equal(t, "data44", string(*result[3].Data))
		assert.Equal(t, true, bool(*result[3].Active))
	}
}

func TestUnion_NewAndUpdateDataRemoteWrite(t *testing.T) {
	existingData := []testStruct{
		{Id: util.Ptr(uint(0)), Changeable: util.Ptr(true), Active: util.Ptr(true), Data: util.Ptr(string("data1"))},
		{Id: util.Ptr(uint(1)), Changeable: util.Ptr(true), Active: util.Ptr(false), Data: util.Ptr(string("data2"))},
		{Id: util.Ptr(uint(2)), Changeable: util.Ptr(false), Active: util.Ptr(false), Data: util.Ptr(string("data3"))},
		{Id: util.Ptr(uint(3)), Active: util.Ptr(true), Data: util.Ptr(string("data4"))},
	}

	newData := []testStruct{
		{Id: util.Ptr(uint(1)), Changeable: util.Ptr(false), Active: util.Ptr(true), Data: util.Ptr(string("data22"))},
		{Id: util.Ptr(uint(2)), Data: util.Ptr(string("data33"))},
		{Id: util.Ptr(uint(3)), Data: util.Ptr(string("data44"))},
	}

	// Act
	result := Merge(true, existingData, newData)

	if assert.Equal(t, 4, len(result)) {
		assert.Equal(t, 0, int(*result[0].Id))
		assert.Equal(t, "data1", string(*result[0].Data))
		assert.Equal(t, true, bool(*result[0].Active))
		assert.Equal(t, 1, int(*result[1].Id))
		assert.Equal(t, "data22", string(*result[1].Data))
		assert.Equal(t, true, bool(*result[1].Changeable))
		assert.Equal(t, true, bool(*result[1].Active))
		assert.Equal(t, 2, int(*result[2].Id))
		assert.Equal(t, "data3", string(*result[2].Data))
		assert.Equal(t, false, bool(*result[2].Active))
		assert.Equal(t, 3, int(*result[3].Id))
		assert.Equal(t, "data4", string(*result[3].Data))
		assert.Equal(t, true, bool(*result[3].Active))
	}
}

func TestUnion_InvalidData(t *testing.T) {
	existingData := []testInvalidStruct{
		{Id: util.Ptr(uint(0)), Changeable: util.Ptr("true"), Data: util.Ptr(string("data1"))},
		{Id: util.Ptr(uint(1)), Changeable: util.Ptr("true"), Data: util.Ptr(string("data2"))},
		{Id: util.Ptr(uint(2)), Changeable: util.Ptr("true"), Data: util.Ptr(string("data3"))},
		{Id: util.Ptr(uint(3)), Data: util.Ptr(string("data4"))},
	}

	newData := []testInvalidStruct{
		{Id: util.Ptr(uint(1)), Data: util.Ptr(string("data22"))},
		{Id: util.Ptr(uint(2)), Data: util.Ptr(string("data33"))},
		{Id: util.Ptr(uint(3)), Data: util.Ptr(string("data44"))},
	}

	// Act
	result := Merge(true, existingData, newData)

	if assert.Equal(t, 4, len(result)) {
		assert.Equal(t, 0, int(*result[0].Id))
		assert.Equal(t, "data1", string(*result[0].Data))
		assert.Equal(t, 1, int(*result[1].Id))
		assert.Equal(t, "data22", string(*result[1].Data))
		assert.Equal(t, 2, int(*result[2].Id))
		assert.Equal(t, "data33", string(*result[2].Data))
		assert.Equal(t, 3, int(*result[3].Id))
		assert.Equal(t, "data4", string(*result[3].Data))
	}
}
