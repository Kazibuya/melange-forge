package main

type GrypeReport struct {
    Matches []Match `json:"matches"`
	Source	Source	`json:"source"`
}

type Source struct {
	Type	string	`json:"type"`
	Target	Target	`json:"target"`
}

type Target struct {
	UserInput string	`json:"userInput"`
	Tags	[]string	`json:"tags"`
}

type Match struct {
    Vulnerability Vulnerability `json:"vulnerability"`
    Artifact      Artifact      `json:"artifact"`
}

type Vulnerability struct {
    ID          string  `json:"id"`
    DataSource  string  `json:"dataSource"`
    Severity    string  `json:"severity"`
    Description string  `json:"description"`
    Fix         Fix     `json:"fix"`
    EPSS        []EPSS  `json:"epss"`
	CVSS	[]CVSS	`json:"cvss"`
	URLs	[]string	`json:"urls"`
}

type Fix struct {
    Versions []string `json:"versions"`
    State    string   `json:"state"`
}

type EPSS struct {
    EPSS       float64 `json:"epss"`
    Percentile float64 `json:"percentile"`
}

type Artifact struct {
    Name    string `json:"name"`
    Version string `json:"version"`
    Type    string `json:"type"`
	PURL	string	`json:"purl"`
	Licenses	[]string	`json:"licenses"`
}

type CVSS struct {
	Source string `json:"source"`
	Type string `json:"type"`
	Metrics Metrics `json:"metrics"`
}

type Metrics struct {
	BaseScore float64 `json:"baseScore"`
}
