package model

import (
	"fmt"
	"time"
)

var tabelNamePrefix = "1ch0_"

var registeredModels = map[string]Interface{}

type Interface interface {
	TableName() string
	ShortTableName() string
}

// RegisterModel register model
func RegisterModel(models ...Interface) {
	for _, model := range models {
		if _, exist := registeredModels[model.TableName()]; exist {
			panic(fmt.Errorf("model table name %s conflict", model.TableName()))
		}
		registeredModels[model.TableName()] = model
	}
}

type BaseModel struct {
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// SetCreateTime set create time
func (m *BaseModel) SetCreateTime(time time.Time) {
	m.CreateTime = time
}

// SetUpdateTime set update time
func (m *BaseModel) SetUpdateTime(time time.Time) {
	m.UpdateTime = time
}
