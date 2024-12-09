package repositories

const (
	//TableNameUsers is a const table name
	TableNameUsers = `users`
	//TableNameBuckets is a const table name
	TableNameBuckets = `buckets`
	//TableNameMerchants is a const table name
	TableNameMerchants = `merchants`
	//TableNameLimits is a const table name
	TableNameLimits = `limits`
	//TableNameTransactions is a const table name
	TableNameTransactions = `transactions`
	//TableNameTransactionCredits is a const table name
	TableNameTransactionCredits = `transaction_credits`
)

var (
	//DefaultQueryFindOne is a raw select query
	DefaultQueryFindOne = `SELECT %s FROM %s %s LIMIT 1;`
	//DefaultQueryFinds is a raw static select query
	DefaultQueryFinds = `SELECT %s FROM %s %s;`
)
