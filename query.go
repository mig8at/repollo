package repollo

import "sort"

// QueryResult estructura que permite realizar operaciones encadenadas en los resultados de una consulta.
type QueryResult[T any] struct {
	results []T
}

// Limit limita la cantidad de elementos en los resultados.
func (q *QueryResult[T]) Limit(limit int) *QueryResult[T] {
	if limit > len(q.results) {
		limit = len(q.results)
	}
	q.results = q.results[:limit]
	return q
}

// Offset aplica un desplazamiento a los resultados.
func (q *QueryResult[T]) Offset(offset int) *QueryResult[T] {
	if offset >= len(q.results) {
		q.results = []T{}
		return q
	}
	q.results = q.results[offset:]
	return q
}

func (q *QueryResult[T]) Sort(less func(a, b T) bool) *QueryResult[T] {
	sort.Slice(q.results, func(i, j int) bool {
		return less(q.results[i], q.results[j])
	})
	return q
}

// Results devuelve los resultados procesados.
func (q *QueryResult[T]) Results() []T {
	return q.results
}
