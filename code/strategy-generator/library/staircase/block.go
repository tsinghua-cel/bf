package staircase

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/types"
	"strconv"
)

// 1. 每个epoch 的区块不出块，投票不广播
// 2. 每个epoch的最后一个作恶者出块，打包所有的 投票，并且延迟广播。
// 3.设block slot为t , receive block: (32 - t mod 32)*12s , broadcast block：在上一个的基础上12*12s

func getSlotStrategy(slot string, cas int, isLatestHackSlot bool) types.SlotStrategy {
	strategy := types.SlotStrategy{
		Slot:    slot,
		Level:   0,
		Actions: make(map[string]string),
	}
	secondsPerSlot := globalinfo.ChainBaseInfo().SecondsPerSlot
	slotsPerEpoch := globalinfo.ChainBaseInfo().SlotsPerEpoch
	switch cas {
	case 0:
		strategy.Actions["BlockBeforeSign"] = "return"
		strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")

	case 1:
		if isLatestHackSlot {
			islot, _ := strconv.Atoi(slot)
			stageI := (slotsPerEpoch - islot%slotsPerEpoch) * secondsPerSlot
			stageII := 12 * secondsPerSlot

			strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")

			strategy.Actions["BlockBeforeSign"] = "packPooledAttest"
			strategy.Actions["BlockDelayForReceiveBlock"] = fmt.Sprintf("%s:%d", "delayWithSecond", stageI)
			strategy.Actions["BlockBeforeBroadCast"] = fmt.Sprintf("%s:%d", "delayWithSecond", stageII)
		} else {
			strategy.Actions["BlockBeforeSign"] = "return"
			strategy.Actions["AttestAfterSign"] = fmt.Sprintf("addAttestToPool")
			strategy.Actions["AttestBeforeBroadCast"] = fmt.Sprintf("return")
		}
	}
	return strategy

}

func GenSlotStrategy(hackDuties []types.ProposerDuty, cas int) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for i := 0; i < len(hackDuties); i++ {
		s := getSlotStrategy(hackDuties[i].Slot, cas, i == len(hackDuties)-1)
		strategys = append(strategys, s)
	}
	return strategys
}
