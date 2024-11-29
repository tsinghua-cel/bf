package confuse

import (
	"github.com/tsinghua-cel/strategy-gen/types"
)

func CheckDuties(param types.LibraryParams, duties []types.ProposerDuty) ([]types.ProposerDuty, bool) {
	result := param.FillterHackDuties(duties)
	return result, len(result) > 0
}
