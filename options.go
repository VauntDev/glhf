package glhf

type opts struct {
	defaultContentType string
}

type Options interface {
	Apply(*opts)
}

type funcOption struct {
	f func(*opts)
}

func (flo *funcOption) Apply(con *opts) {
	flo.f(con)
}

func newFuncOption(f func(*opts)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// Defines the Default Content-type to be used if one is not set by the user.
func WithDefaultContentType(contentType string) Options {
	return newFuncOption(func(o *opts) {
		o.defaultContentType = contentType
	})
}

func defaultOptions() *opts {
	return &opts{
		defaultContentType: ContentJSON,
	}
}
