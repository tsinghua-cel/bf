package staircase

import (
	"context"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/tsinghua-cel/strategy-gen/globalinfo"
	"github.com/tsinghua-cel/strategy-gen/types"
	"github.com/tsinghua-cel/strategy-gen/utils"
	"strconv"
	"time"
)

type Instance struct {
}

func (o *Instance) Name() string {
	return "staircase"
}

func (o *Instance) Description() string {
	desc_eng := `Staircase attack`
	return desc_eng
}

func (o *Instance) Run(ctx context.Context, params types.LibraryParams, feedbacker types.FeedBacker) {
	log.WithField("name", o.Name()).Info("start to run strategy")
	var latestEpoch int64 = -1
	ticker := time.NewTicker(time.Second * 3)
	slotTool := utils.SlotTool{SlotsPerEpoch: globalinfo.ChainBaseInfo().SlotsPerEpoch}
	for {
		select {
		case <-ctx.Done():
			log.WithField("name", o.Name()).Info("stop to run strategy")
			return
		case <-ticker.C:
			slot, err := utils.GetCurSlot(params.Attacker)
			if err != nil {
				log.WithField("error", err).Error("failed to get slot")
				continue
			}
			log.WithFields(log.Fields{
				"slot":      slot,
				"lastEpoch": latestEpoch,
			}).Info("get slot")
			epoch := slotTool.SlotToEpoch(int64(slot))
			// generate new strategy at the end of last epoch.
			if int64(slot) < slotTool.EpochEnd(epoch) {
				continue
			}
			if epoch == latestEpoch {
				continue
			}
			latestEpoch = epoch

			{
				nextEpoch := epoch + 1
				cas := 0

				nextDuties, err := utils.GetEpochDuties(params.Attacker, nextEpoch)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"epoch": nextEpoch,
					}).Error("failed to get duties")
					latestEpoch = epoch - 1
					continue
				}
				preDuties, err := utils.GetEpochDuties(params.Attacker, epoch-1)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"epoch": epoch - 1,
					}).Error("failed to get duties")
					latestEpoch = epoch - 1
					continue
				}
				curDuties, err := utils.GetEpochDuties(params.Attacker, epoch)
				if err != nil {
					log.WithFields(log.Fields{
						"error": err,
						"epoch": epoch,
					}).Error("failed to get duties")
					latestEpoch = epoch - 1
					continue
				}
				strategy := types.Strategy{}
				if checkFirstByzSlot(preDuties, params) &&
					checkFirstByzSlot(curDuties, params) &&
					!checkFirstByzSlot(nextDuties, params) {
					cas = 1
				}
				strategy.Uid = uuid.NewString()
				strategy.Slots = GenSlotStrategy(params.FillterHackDuties(nextDuties), cas)
				if err = utils.UpdateStrategy(params.Attacker, strategy); err != nil {
					log.WithField("error", err).Error("failed to update strategy")
				} else {
					log.WithFields(log.Fields{
						"epoch":    nextEpoch,
						"strategy": strategy,
					}).Info("update strategy successfully")
				}
			}
		}
	}
}

func getLatestHackerSlot(duties []types.ProposerDuty, param types.LibraryParams) int {
	latest, _ := strconv.Atoi(duties[0].Slot)
	for _, duty := range duties {
		idx, _ := strconv.Atoi(duty.ValidatorIndex)
		slot, _ := strconv.Atoi(duty.Slot)
		if !types.IsHackValidator(idx, param) {
			continue
		}
		if slot > latest {
			latest = slot
		}
	}
	return latest

}

func checkFirstByzSlot(duties []types.ProposerDuty, param types.LibraryParams) bool {
	firstproposerindex, _ := strconv.Atoi(duties[0].ValidatorIndex)
	if !types.IsHackValidator(firstproposerindex, param) {
		return false
	}
	return true
}
