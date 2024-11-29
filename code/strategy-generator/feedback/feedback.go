package feedback

import (
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"sync"
	"time"
)

type Feedbacker struct {
	url   string
	tasks sync.Map
	quit  chan struct{}
}

func NewFeedbacker(url string) *Feedbacker {
	return &Feedbacker{url: url, quit: make(chan struct{})}
}

func (f *Feedbacker) WaitFeedback(uid string, ch chan types.FeedBack) {
	f.tasks.Store(uid, ch)
}

func (f *Feedbacker) RemoveFeedback(uid string) {
	f.tasks.Delete(uid)
}

func (f *Feedbacker) Start() {
	go f.loop()
}

func (f *Feedbacker) Stop() {
	close(f.quit)
}

func (f *Feedbacker) loop() {
	tc := time.NewTicker(time.Second * 10)
	defer tc.Stop()
	for {
		select {
		case <-tc.C:
			f.tasks.Range(func(key, value interface{}) bool {
				uid := key.(string)
				ch := value.(chan types.FeedBack)
				feedback, err := utils.GetStrategyFeedback(f.url, uid)
				if err != nil {
					return true
				}
				ch <- types.FeedBack{
					Uid:  uid,
					Info: feedback,
				}
				f.tasks.Delete(uid)
				return true
			})
		case <-f.quit:
			f.tasks.Range(func(key, value interface{}) bool {
				ch := value.(chan types.FeedBackInfo)
				close(ch)
				return true
			})
		}
	}
}
