package glhf

type opts struct {
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

func defaultOptions() *opts {
	return &opts{}
}
