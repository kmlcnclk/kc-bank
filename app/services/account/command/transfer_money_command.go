package command

type TransferMoneyCommand struct {
	Amount   float64
	FromIBAN string
	ToIBAN   string
}
