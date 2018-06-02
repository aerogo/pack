package main

import (
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	os.Chdir("examples")

	main()

	os.RemoveAll("components")
	os.Chdir("..")
}

func BenchmarkCompiler(b *testing.B) {
	os.Chdir("examples")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		main()
	}

	os.RemoveAll("components")
	os.Chdir("..")
}
