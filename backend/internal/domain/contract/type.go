package contract

type ContractType string

const (
	ContractTypeFixedFee   ContractType = "fixed_fee"
	ContractTypeHourly     ContractType = "hourly"
	ContractTypeSuccessFee ContractType = "success_fee"
	ContractTypeMonthly    ContractType = "monthly"
)
