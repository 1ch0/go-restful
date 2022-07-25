package container

import (
	"time"

	"github.com/1ch0/go-restful/pkg/apiserver/utils/log"
	"github.com/barnettZQG/inject"
)

type Container struct {
	graph inject.Graph
}

func NewContainer() *Container {
	return &Container{
		graph: inject.Graph{},
	}
}

func (c *Container) Provides(beans ...interface{}) error {
	for _, bean := range beans {
		if err := c.graph.Provide(&inject.Object{Value: bean}); err != nil {
			return err
		}
	}
	return nil
}

// ProvideWithName provide the bean with name
func (c *Container) ProvideWithName(name string, bean interface{}) error {
	return c.graph.Provide(&inject.Object{Name: name, Value: bean})
}

// Populate populate dependency fields for all beans.
func (c *Container) Populate() error {
	start := time.Now()
	defer func() {
		log.Logger.Infof("populate the bean container take time %s", time.Now().Sub(start))
	}()
	return c.graph.Populate()
}
