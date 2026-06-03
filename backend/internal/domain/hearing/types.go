package hearing

type HearingType string

const (
	HearingTypeConciliation HearingType = "conciliation"
	HearingTypeInstruction  HearingType = "instruction"
	HearingTypeJudgment     HearingType = "judgment"
	HearingTypeVirtual      HearingType = "virtual"
)
