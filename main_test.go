package main

import (
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	defer os.Chdir("..")
	defer os.RemoveAll("components")

	os.Chdir("testdata")
	main()
}

func BenchmarkCompiler(b *testing.B) {
	defer os.Chdir("..")
	defer os.RemoveAll("components")

	os.Chdir("testdata")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		main()
	}
}
