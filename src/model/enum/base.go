package enum

type IEnum interface {
	ToString() string
	Title(string) string
}
