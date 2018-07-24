package model

type Sources []Source
type Source struct {
	Name     string
	Page     PageOptions
	Selector SelectorOptions
	Category CategoryOptions
	Debug    bool
}

type PageOptions struct {
	Entry    string
	Template string
	From     int
	To       int
}

type SelectorOptions struct {
	Iterator string
	IP       string
	Port     string
	Scheme   string
	Filter   string
}

type CategoryOptions struct {
	ParallelNumber int
	DelayRange     []int
	Interval       string
}
