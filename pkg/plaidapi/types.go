package plaidapi

type AccountType string

const (
	AccountTypeInvestment AccountType = "investment"
	AccountTypeCredit     AccountType = "credit"
	AccountTypeDepository AccountType = "depository"
	AccountTypeLoan       AccountType = "loan"
	AccountTypeOther      AccountType = "other"
)

type AccountSubtype string

const (
	//Loan
	AccountSubtype401a                           AccountSubtype = "401a"
	AccountSubtype401k                           AccountSubtype = "401k"
	AccountSubtype403B                           AccountSubtype = "403B"
	AccountSubtype457b                           AccountSubtype = "457b"
	AccountSubtype529                            AccountSubtype = "529"
	AccountSubtypeBrokerage                      AccountSubtype = "brokerage"
	AccountSubtypeCashISA                        AccountSubtype = "cash isa"
	AccountSubtypeEducationSavingsAccount        AccountSubtype = "education savings account"
	AccountSubtypeGic                            AccountSubtype = "gic"
	AccountSubtypeHealthReimbursementArrangement AccountSubtype = "health reimbursement arrangement"
	AccountSubtypeHsa                            AccountSubtype = "hsa"
	AccountSubtypeIsa                            AccountSubtype = "isa"
	AccountSubtypeIra                            AccountSubtype = "ira"
	AccountSubtypeLif                            AccountSubtype = "lif"
	AccountSubtypeLira                           AccountSubtype = "lira"
	AccountSubtypeLrif                           AccountSubtype = "lrif"
	AccountSubtypeLrsp                           AccountSubtype = "lrsp"
	AccountSubtypeNonTaxableBrokerageAccount     AccountSubtype = "non-taxable brokerage account"
	AccountSubtypeOther                          AccountSubtype = "other"
	AccountSubtypePrif                           AccountSubtype = "prif"
	AccountSubtypeRdsp                           AccountSubtype = "rdsp"
	AccountSubtypeResp                           AccountSubtype = "resp"
	AccountSubtypeRlif                           AccountSubtype = "rlif"
	AccountSubtypeRrif                           AccountSubtype = "rrif"
	AccountSubtypePension                        AccountSubtype = "pension"
	AccountSubtypeProfitSharingPlan              AccountSubtype = "profit sharing plan"
	AccountSubtypeRetirement                     AccountSubtype = "retirement"
	AccountSubtypeRoth                           AccountSubtype = "roth"
	AccountSubtypeRoth401k                       AccountSubtype = "roth 401k"
	AccountSubtypeRrsp                           AccountSubtype = "rrsp"
	AccountSubtypeSepIRA                         AccountSubtype = "sep ira"
	AccountSubtypeSimpleIRA                      AccountSubtype = "simple ira"
	AccountSubtypeSipp                           AccountSubtype = "sipp"
	AccountSubtypeStockPlan                      AccountSubtype = "stock plan"
	AccountSubtypeThriftSavingsPlan              AccountSubtype = "thrift savings plan"
	AccountSubtypeTfsa                           AccountSubtype = "tfsa"
	AccountSubtypeUgma                           AccountSubtype = "ugma"
	AccountSubtypeUtma                           AccountSubtype = "utma"
	AccountSubtypeVariableAnnuity                AccountSubtype = "variable annuity"

	//Credit
	AccountSubtypeCreditCard AccountSubtype = "credit card"
	AccountSubtypePaypal     AccountSubtype = "paypal"

	//Depository
	AccountSubtypeCD          AccountSubtype = "cd"
	AccountSubtypeChecking    AccountSubtype = "checking"
	AccountSubtypeSavings     AccountSubtype = "savings"
	AccountSubtypeMoneyMarket AccountSubtype = "money market"
	AccountSubtypePrepaid     AccountSubtype = "prepaid"
	// AccountSubtypePaypal is also permitted here

	//Loan
	AccountSubtypeAuto         AccountSubtype = "auto"
	AccountSubtypeCommercial   AccountSubtype = "commercial"
	AccountSubtypeConstruction AccountSubtype = "construction"
	AccountSubtypeConsumer     AccountSubtype = "consumer"
	AccountSubtypeHome         AccountSubtype = "home"
	AccountSubtypeHomeEquity   AccountSubtype = "home equity"
	AccountSubtypeLoan         AccountSubtype = "loan"
	AccountSubtypeMortgage     AccountSubtype = "mortgage"
	AccountSubtypeOverdraft    AccountSubtype = "overdraft"
	AccountSubtypeLineOfCredit AccountSubtype = "line of credit"
	AccountSubtypeStudent      AccountSubtype = "student"

	//Other
	AccountSubtypeCashManagement AccountSubtype = "cash management"
	AccountSubtypeKeogh          AccountSubtype = "keogh"
	AccountSubtypeMutualFund     AccountSubtype = "mutual fund"
	AccountSubtypeRecurring      AccountSubtype = "recurring"
	AccountSubtypeRewards        AccountSubtype = "rewards"
	AccountSubtypeSafeDeposit    AccountSubtype = "safe deposit"
	AccountSubtypeSarsep         AccountSubtype = "sarsep"
	//AccountSubtypePrepaid is also permitted here
	//AccountSubtypeOther is also permitted here
)
