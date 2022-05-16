package aurweb

type Result struct {
	Description    string  `json:"Description"`
	FirstSubmitted int     `json:"FirstSubmitted"`
	ID             int     `json:"ID"`
	LastModified   int     `json:"LastModified"`
	Maintainer     string  `json:"Maintainer"`
	Name           string  `json:"Name"`
	NumVotes       int     `json:"NumVotes"`
	OutOfDate      int     `json:"OutOfDate"`
	PackageBase    string  `json:"PackageBase"`
	PackageBaseID  int     `json:"PackageBaseID"`
	Popularity     float64 `json:"Popularity"`
	URL            string  `json:"URL"`
	URLPath        string  `json:"URLPath"`
	Version        string  `json:"Version"`
}

type SearchResponse struct {
	ResultCount int      `json:"resultcount"`
	Results     []Result `json:"results"`
	Type        string   `json:"type"`
	Version     int      `json:"version"`
	Error       string   `json:"error"`
}

type InfoResult struct {
	Result
	Conflicts   []string `json:"Conflicts"`
	Depends     []string `json:"Depends"`
	Keywords    []string `json:"Keywords"`
	License     []string `json:"License"`
	Maintainer  string   `json:"Maintainer"`
	MakeDepends []string `json:"MakeDepends"`
	OptDepends  []string `json:"OptDepends"`
	Provides    []string `json:"Provides"`
}

type InfoResponse struct {
	ResultCount int          `json:"resultcount"`
	Results     []InfoResult `json:"results"`
	Type        string       `json:"type"`
	Version     int          `json:"version"`
}
