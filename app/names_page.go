package app

// NamesPage represents list of users for a given namespace
// at a given page the page number and namespace are used
// to determine the path where the userIDs json file is stored
type NamesPage struct {
	PageNum   uint64
	UserIDs   []string `json:",flow"`
	Namespace string
	Count     int
}

// TODO fill out
func (np NamesPage) DataDir() string {
	return ""
}
