package main

import (
	"fmt"
	"strings"
)

func generateBody(match Match, sbom SBOM) string {
	var res strings.Builder
	fmt.Fprintf(&res, "**Severity:** %s\n", match.Vulnerability.Severity)
	fmt.Fprintf(&res, "**Package:** %s %s %s\n", match.Artifact.Name, match.Artifact.Version, match.Artifact.Type)
	fmt.Fprintf(&res, "**PURL:** %s\n\n", match.Artifact.PURL)
	fmt.Fprintf(&res, "### Description\n%s\n\n", match.Vulnerability.Description)
	fmt.Fprintf(&res, "### Fix\n")
	if match.Vulnerability.Fix.State != "" {
		fmt.Fprintf(&res, "**State:** %s\n", match.Vulnerability.Fix.State)
	}else {
		fmt.Fprintf(&res, "**State:** not-fixed\n")
	}
	if len(match.Vulnerability.Fix.Versions) > 0 {
		fmt.Fprintf(&res, "**Fixed in:** %s\n\n", match.Vulnerability.Fix.Versions[0])
	}else {
		fmt.Fprintf(&res, "**Fixed in:** N/A\n\n")
	}
	fmt.Fprintf(&res, "### Scores\n")
	if len(match.Vulnerability.CVSS) > 0 {
		fmt.Fprintf(&res, "**CVSS Base Score:** %.1f\n", match.Vulnerability.CVSS[0].Metrics.BaseScore)
	}
	if len(match.Vulnerability.EPSS) > 0 {
		fmt.Fprintf(&res, "**EPSS:** %.3f%% (%.0fth percentile)\n\n", match.Vulnerability.EPSS[0].EPSS * 100, match.Vulnerability.EPSS[0].Percentile * 100)
	}
	if len(match.Vulnerability.URLs) > 0 {
		fmt.Fprintf(&res, "### References\n")
		fmt.Fprintf(&res, "- %s\n", match.Vulnerability.DataSource)
		for _, url := range match.Vulnerability.URLs {
			fmt.Fprintf(&res, "- %s\n", url)
		}
		fmt.Fprint(&res, "\n")
	}
	if sbom.Packages != nil {
		for _, rel := range sbom.Relationships {
			if rel.RelationshipType == "GENERATED_FROM" &&
				strings.Contains(rel.SpdxElementId, match.Artifact.Name) {
				for _, p := range sbom.Packages {
					if p.SPDXID == rel.RelatedSpdxElement &&
						p.DownloadLocation != "NOASSERTION" {
						fmt.Fprintf(&res, "### Source\n**Upstream:** %s\n", p.DownloadLocation)
					}
				}
			}
		}
	}
	return res.String()
}
