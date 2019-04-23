package main

import (
	"os"
	"testing"
)

func TestCompiler(t *testing.T) {
	defer os.Chdir("..")
	defer os.RemoveAll("components")

	os.Chdir("examples")
	main()
}

func BenchmarkCompiler(b *testing.B) {
	defer os.Chdir("..")
	defer os.RemoveAll("components")

	os.Chdir("examples")
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		main()
	}
}
