package options

type ComponentOptions struct {
	Definition ComponentDefinition
	Middleware []Middleware
}

type ComponentOption func(o *ComponentOptions)

func DefaultComponentOptions() ComponentOptions {
	return ComponentOptions{}
}

func (c *ComponentOptions) Validate() {
	if c.Definition == nil {
		panic("Component definition is required")
	}
}
