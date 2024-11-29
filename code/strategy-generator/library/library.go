package library

import (
	"context"
	"github.com/tsinghua-cel/strategy-gen/library/confuse"
	"github.com/tsinghua-cel/strategy-gen/library/exante"
	"github.com/tsinghua-cel/strategy-gen/library/ext_exante"
	"github.com/tsinghua-cel/strategy-gen/library/ext_sandwich"
	"github.com/tsinghua-cel/strategy-gen/library/ext_staircase"
	"github.com/tsinghua-cel/strategy-gen/library/ext_unrealized"
	"github.com/tsinghua-cel/strategy-gen/library/ext_withholding"
	"github.com/tsinghua-cel/strategy-gen/library/five"
	"github.com/tsinghua-cel/strategy-gen/library/one"
	"github.com/tsinghua-cel/strategy-gen/library/randomdelay"
	"github.com/tsinghua-cel/strategy-gen/library/sandwich"
	"github.com/tsinghua-cel/strategy-gen/library/staircase"
	"github.com/tsinghua-cel/strategy-gen/library/syncwrong"
	"github.com/tsinghua-cel/strategy-gen/library/three"
	"github.com/tsinghua-cel/strategy-gen/library/two"
	"github.com/tsinghua-cel/strategy-gen/library/unrealized"
	"github.com/tsinghua-cel/strategy-gen/library/withholding"
	"github.com/tsinghua-cel/strategy-gen/types"
	"sync"
)

type Strategy interface {
	Name() string
	Run(ctx context.Context, param types.LibraryParams, feedbacker types.FeedBacker)
	Description() string
}

var (
	allStrategies = sync.Map{}
)

func Init() {
	register(&one.Instance{})
	register(&two.Instance{})
	register(&three.Instance{})
	register(&syncwrong.Instance{})
	register(&five.Instance{})
	register(&exante.Instance{})
	register(&sandwich.Instance{})
	register(&withholding.Instance{})
	register(&unrealized.Instance{})
	register(&staircase.Instance{})
	register(&confuse.Instance{})
	register(&randomdelay.Instance{})

	register(&ext_sandwich.Instance{})
	register(&ext_exante.Instance{})
	register(&ext_staircase.Instance{})
	register(&ext_unrealized.Instance{})
	register(&ext_withholding.Instance{})

}

func register(ins Strategy) {
	allStrategies.Store(ins.Name(), ins)
}

func GetStrategy(name string) (Strategy, bool) {
	if v, ok := allStrategies.Load(name); ok {
		return v.(Strategy), true
	}
	return nil, false
}

func GetAllStrategies() map[string]Strategy {
	strategies := make(map[string]Strategy)
	allStrategies.Range(func(k, v interface{}) bool {
		strategies[k.(string)] = v.(Strategy)
		return true
	})
	return strategies
}
func GetStrategiesList() []Strategy {
	strategies := make([]Strategy, 0)
	allStrategies.Range(func(k, v interface{}) bool {
		strategies = append(strategies, v.(Strategy))
		return true
	})
	return strategies
}
