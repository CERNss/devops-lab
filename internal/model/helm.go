package model

type HelmRelease struct {
	Name      string `json:"name"`
	Chart     string `json:"chart"`
	Namespace string `json:"-"`
	Values    string `json:"values"`
	Timeout   string `json:"timeout"`
	Atomic    bool   `json:"atomic,omitempty"`
	Repo      *HelmRepo `json:"repo,omitempty"`
}

type HelmRepo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
