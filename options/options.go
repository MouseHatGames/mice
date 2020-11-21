package options

// Options holds the configuration for a service instance
type Options struct {
	Name string
}

// Option represents a function that can be used to mutate an Options object
type Option func(*Options)

// Name sets the name of the service
func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}
