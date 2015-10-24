package main

const Version string = "0.0.6"

func main() {
	args, opts := GhsOptionParser()
	query := ""

	for _, arg := range args {
		query += arg
	}

	if opts.Fields != "" {
		query += " in:" + opts.Fields
	}
	if opts.Language != "" {
		query += " language:" + opts.Language
	}
	if opts.User != "" {
		query += " user:" + opts.User
	}
	if opts.Repository != "" {
		query += " repo:" + opts.Repository
	}

	repos := SearchRepository(opts.Sort, opts.Order, opts.Max, opts.Enterprise, query)
	PrintRepository(repos)
}
