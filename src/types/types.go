package types

type SubnetsTagsFields struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

type SubnetTags struct {
	Tags []SubnetsTagsFields `json:"Tags"`
}
