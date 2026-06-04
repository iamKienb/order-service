package order

type QueryRepository interface {
}

type CommandRepository interface {
	CreateOrder()
	ConfirmO
}

type Repository interface {
	QueryRepository
	CommandRepository
}
