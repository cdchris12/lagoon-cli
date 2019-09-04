package cmd

// ProjectByName struct.
type ProjectByName struct {
	ProjectByName Project `json:"projectByName"`
}

// WhatIsThere struct.
type WhatIsThere struct {
	AllProjects []Project `json:"allProjects"`
}

// Environments struct.
type Environments struct {
	Name            string `json:"name"`
	EnvironmentType string `json:"environmentType"`
	DeployType      string `json:"deployType"`
	Route           string `json:"route"`
}

// Project struct.
type Project struct {
	ID                           int            `json:"id"`
	GitURL                       string         `json:"gitUrl"`
	Subfolder                    string         `json:"subfolder"`
	Name                         string         `json:"name"`
	Branches                     string         `json:"branches"`
	Pullrequests                 string         `json:"pullrequests"`
	ProductionEnvironment        string         `json:"productionEnvironment"`
	Environments                 []Environments `json:"environments"`
	AutoIdle                     int            `json:"autoIdle"`
	DevelopmentEnvironmentsLimit int            `json:"developmentEnvironmentsLimit"`
}

type Environment struct {
	EnvironmentByOpenshiftProjectName Environments `json:"environmentByOpenshiftProjectName"`
}

type DeployResult struct {
	DeployEnvironmentBranch string `json:"deployEnvironmentBranch"`
}

type DeleteResult struct {
	DeleteEnvironment string `json:"deleteEnvironment"`
}
