package model

// Sources is an array of Source.
type Sources []*Source

// Source is the source configuration
type Source struct {
	// Name is the name of source
	Name string
	// Page is the page options
	Page PageOptions
	// Selector is the selector options
	Selector SelectorOptions
	// Category is the category options
	Category CategoryOptions
	// Debug determines whether to output debugging information
	Debug bool
}

// PageOptions is the page configuration
type PageOptions struct {
	// Entry is the first url to cralw
	Entry string
	// Template is the page template. e.g http:xxxx.xxx/proxy?page={page}
	Template string
	// From is the start page number
	From int
	// To is the end page number
	To int
}

// SelectorOptions is the selector configuration
type SelectorOptions struct {
	// Iterator is the iterable element of proxy items
	Iterator string
	// IP is the IP selector
	IP string
	// Port is the port selector
	Port string
}

// CategoryOptions is the category configuration
type CategoryOptions struct {
	// ParallelNumber is the number of parallels that the source crawls
	// e.g 10
	ParallelNumber int
	// DelayRange is the interval between crawls, random in this array range
	// e.g [0, 10]
	DelayRange []int
	// Interval is how long it takes to re-crawl from StarURL
	// e.g "@every 10m", "@every 10s", "@every 10h"
	Interval string
}
