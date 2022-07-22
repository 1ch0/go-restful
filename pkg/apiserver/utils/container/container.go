package container

import "github.com/barnettZQG/inject"

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
