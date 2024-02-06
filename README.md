# triplebuffer

[![Go Reference](https://pkg.go.dev/badge/github.com/jalavosus/triplebuffer.svg)](https://pkg.go.dev/github.com/jalavosus/triplebuffer)

`go get github.com/jalavosus/triplebuffer`

An attempt at writing a triple buffer system in Go, mainly just to see if 
I could do it, make it work, AND test it.

There are GoDocs for all public things, a fairly extensive set of tests,
and a very basic usage example in [example_test.go](./example_test.go).

Don't use this in production without personally vetting it. Please. I beg you.

Inspired by the code in [this here article](https://brilliantsugar.github.io/posts/how-i-learned-to-stop-worrying-and-love-juggling-c++-atomics/),
all glory to brilliantsugar. 