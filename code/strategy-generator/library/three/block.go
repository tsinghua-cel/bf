package three

import (
	"fmt"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/types"
)

func getSlotStrategy(epoch int64, slot string, isLatestHackDuty bool) types.SlotStrategy {
	strategy := types.SlotStrategy{
		Slot:    slot,
		Level:   0,
		Actions: make(map[string]string),
	}
	secondsPerSlot := globalinfo.ChainBaseInfo().SecondsPerSlot
	slotsPerEpoch := globalinfo.ChainBaseInfo().SlotsPerEpoch
	switch epoch % 3 {
	case 0, 1:
		strategy.Actions["BlockBeforeSign"] = "return"
		strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")

	case 2:
		if isLatestHackDuty {

			strategy.Actions["AttestBeforeSign"] = fmt.Sprintf("return")

			strategy.Actions["BlockBeforeSign"] = "packPooledAttest"
			// delay half epoch first.
			strategy.Actions["BlockDelayForReceiveBlock"] = fmt.Sprintf("delayHalfEpoch")
			// and then delay 1.5 epoch and 1 slot.
			totalSlots := slotsPerEpoch*3/2 + 1
			totalSeconds := secondsPerSlot * totalSlots
			strategy.Actions["BlockBeforeBroadCast"] = fmt.Sprintf("delayWithSecond:%d", totalSeconds)

		} else {
			strategy.Actions["BlockBeforeSign"] = "return"
			strategy.Actions["AttestAfterSign"] = fmt.Sprintf("addAttestToPool")
			strategy.Actions["AttestBeforeBroadCast"] = fmt.Sprintf("return")
		}
	}
	return strategy
}

func GenSlotStrategy(allHackDuties []types.ProposerDuty, epoch int64) []types.SlotStrategy {
	strategys := make([]types.SlotStrategy, 0)
	for i := 0; i < len(allHackDuties); i++ {
		s := getSlotStrategy(epoch, allHackDuties[i].Slot, i == len(allHackDuties)-1)
		strategys = append(strategys, s)
	}

	return strategys

}
