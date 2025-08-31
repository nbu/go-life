package framework

type Help[T any] interface {
	DefineUsage() T
}

type GameLoop[T any] interface {
	Start(parameters T)
}

func Run[T any](
	help Help[T],
	gameLoop GameLoop[T]) {

	usage := help.DefineUsage()
	gameLoop.Start(usage)
}
