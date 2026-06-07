package main

type SBOM struct {
	Relationships []relationship `json:"relationships"`
	Packages []pkg `json:"packages"`
}

type relationship struct {
	SpdxElementId string `json:"spdxElementId"`
	RelationshipType string `json:"relationshipType"`
	RelatedSpdxElement string `json:"relatedSpdxElement"`
}

type pkg struct {
	SPDXID string `json:"SPDXID"`
	Name string `json:"name"`
	VersionInfo string `json:"versionInfo"`
	DownloadLocation string `json:"downloadLocation"`
	PrimaryPackagePurpose string `json:"primaryPackagePurpose"`
}
