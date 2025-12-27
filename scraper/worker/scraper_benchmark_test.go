package worker

import (
	"testing"
)

// just to make sure scraper has no memory leaks
// go test -bench=. -benchmem -benchtime=1x
func BenchmarkScrapeSkins(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _, err := ScrapeSkins()
		if err != nil {
			b.Fatal(err)
		}
	}
}
