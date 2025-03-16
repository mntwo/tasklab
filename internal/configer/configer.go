package configer

type Configer interface {
	Parse() error
}
