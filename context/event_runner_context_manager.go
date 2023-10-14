package context

import "github.com/aivyss/eventx/entity"

func manageEventRunnerContext(set entity.EventSet, literal func(set entity.EventSet)) {
	ctx := set.Context()
	ctx.Lock()
	defer ctx.Unlock()

	if ctx.IsRunnableInternal() {
		literal(set)
		ctx.Done = true
	}
}
