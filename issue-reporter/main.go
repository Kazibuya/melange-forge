package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"os"
	"strings"
)

func main() {
	grype := flag.String("grype", "grype-report.json", "Path to the report")
	sbom := flag.String("sbom", "", "Path to apko SBOM")
	severity := flag.String("severity", "high,critical", "Severity filter")
	flag.Parse()
	data, err := os.ReadFile(*grype)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1);
	}
	var report GrypeReport
	if err:= json.Unmarshal(data, &report); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	var s SBOM
	if *sbom != "" {
		data, err := os.ReadFile(*sbom)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1);
		}
		if err:= json.Unmarshal(data, &s); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	}
	var filter []string
	filter = strings.Split(*severity, ",")
	for _, match := range report.Matches {
		for _, check := range filter {
			if strings.EqualFold(match.Vulnerability.Severity, check) {
				body := generateBody(match, s)
				output := map[string]string{
					"title": fmt.Sprintf("%s in %s:%s", match.Vulnerability.ID, match.Artifact.Name, match.Artifact.Version),
					"body":  body,
					"severity": strings.ToLower(match.Vulnerability.Severity),
				}
				json.NewEncoder(os.Stdout).Encode(output)
			}
		}
	}
}
