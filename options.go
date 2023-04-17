package glhf

type opts struct {
	dsn         string
	urlParamKey any
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

func newFuncNodeOption(f func(*opts)) *funcOption {
	return &funcOption{
		f: f,
	}
}

func WithParametersKey(k any) Options {
	return newFuncNodeOption(func(o *opts) {
		o.urlParamKey = k
	})
}

func defaultOptions() *opts {
	return &opts{
		dsn: "",
	}
}
