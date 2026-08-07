package core

// Formatter is implemented by any value that has a Format method.
type Formatter interface {
	Format() string
}

// Parser is implemented by any value that has a Parse method.
type Parser interface {
	Parse(value string) error
}
