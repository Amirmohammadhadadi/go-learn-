package main

func splitIntSlice(s []int) ([]int, []int) {
	mid := len(s) / 2
	return s[:mid], s[mid:]
}
func splitStringSlice(s []string) ([]string, []string) {
	mid := len(s) / 2
	return s[:mid], s[mid:]
}
func splitAnySlice[T any](s []T) ([]T, []T) {
	mid := len(s) / 2
	return s[:mid], s[mid:]
}

type stringer interface {
	String() string
}

func concat[T stringer](vals []T, key K) string {
	result := ""
	for _, v := range vals {
		result += v.String()
	}
	return result
}

func t1[T ~int | ~string](val T) {

}

type Ordered interface {
	~int | ~string | ~float32
}

// -------------------------
type Store[P Product] interface {
	Sell(P)
}
type Product interface {
	Price() float32
	Name() string
}
type Book struct {
	title  string
	author string
	price  float32
}

func (b Book) Price() float32 {
	return b.price
}

func (b Book) Name() string {
	return b.title
}

type BookStore struct {
	booksSold []Book
}

func (s *BookStore) Sell(book Book) {
	s.booksSold = append(s.booksSold, book)
}
