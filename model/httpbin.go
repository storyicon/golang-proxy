package model

// HTTPBinIP is the response struct of httpbin.org/ip
type HTTPBinIP struct {
	// Origin is the real ip returned from httpbin
	Origin string `json:"origin"`
}
