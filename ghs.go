package main

var version = "0.0.1"

func main() {
	args, opts := GhsOptionParser()
	repos := SearchRepository(opts.Sort, opts.Order, args[0])
	PrintRepository(repos)
}
