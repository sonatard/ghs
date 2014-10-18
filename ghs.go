package main

var version = "0.0.1"

func main() {
	sort, order, query := GhsOptionParser()
	repos := NewRepositoryBySearch(sort, order, query)
	Print(repos)
}
