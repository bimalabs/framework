package repositories

type (
	Transaction func(Repository) error

	Filter struct {
		Value    any
		Field    string
		Operator string
	}

	Repository interface {
		Model(model string)
		Transaction(Transaction) error
		Create(v any) error
		Update(v any) error
		Bind(v any, id string) error
		All(v any) error
		FindBy(v any, filters ...Filter) error
		Delete(v any, id string) error
	}
)
