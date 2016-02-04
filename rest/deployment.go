package rest

type Deployment struct {
	Datacenter            string                `json:"datacenter"`
	Name                  string                `json:"name"`
	Applicances           []Appliance           `json:"appliances"`
	ApplicationDefinition ApplicationDefinition `json:"applicationDefintion"`
	Type                  int                   `json:"type"`
	Tenant                string                `json:"tenant"`
	Metadata              map[string]string     `json:"metadata"`
}

type Appliance struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type ApplicationDefinition struct {
	Steps []Step `json:"steps"`
}

type Step struct {
	Command       string `json:"command"`
	Role          string `json:"role"`
	Name          string `json:"name"`
	ApplianceName string `json:"applianceName"`
}
