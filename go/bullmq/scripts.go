package bullmq

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/taskforcesh/bullmq/go/luabundler"
)

var ErrScriptNotFound = errors.New("script not found")

type Scripts struct {
	scripts map[string]*redis.Script
}

func NewScripts(
	w *luabundler.Workspace,
) *Scripts {
	out := &Scripts{
		scripts: make(map[string]*redis.Script),
	}
	for k, v := range w.All() {
		out.scripts[k] = redis.NewScript(string(v))
	}
	return out
}

func (s *Scripts) ExecuteScript(ctx context.Context, c redis.Scripter, name string, keys []string, args ...any) (any, error) {
	script, ok := s.scripts[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrScriptNotFound, name)
	}
	return script.Eval(ctx, c, keys, args...).Result()
}
