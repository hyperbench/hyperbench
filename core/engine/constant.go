package engine

func newConstantEngine(base *baseEngine) *ConstantEngine {
	return &ConstantEngine{
		baseEngine: base,
	}
}

// ConstantEngine is control rate by constant.
type ConstantEngine struct {
	*baseEngine
}
