package gitee

type ProjectHook struct {
	Id              int32      `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	Path            string     `json:"path,omitempty"`
	FullName        string     `json:"full_name,omitempty"`
	Owner           *UserBasic `json:"owner,omitempty"`
	Private         bool       `json:"private,omitempty"`
	HtmlUrl         string     `json:"html_url,omitempty"`
	Url             string     `json:"url,omitempty"`
	Description     string     `json:"description,omitempty"`
	Fork            bool       `json:"fork,omitempty"`
	PushedAt        string     `json:"pushed_at,omitempty"`
	CreatedAt       string     `json:"created_at,omitempty"`
	UpdatedAt       string     `json:"updated_at,omitempty"`
	SshUrl          string     `json:"ssh_url,omitempty"`
	GitUrl          string     `json:"git_url,omitempty"`
	CloneUrl        string     `json:"clone_url,omitempty"`
	SvnUrl          string     `json:"svn_url,omitempty"`
	GitHttpUrl      string     `json:"git_http_url,omitempty"`
	GitSshUrl       string     `json:"git_ssh_url,omitempty"`
	GitSvnUrl       string     `json:"git_svn_url,omitempty"`
	Homepage        string     `json:"homepage,omitempty"`
	StargazersCount int32      `json:"stargazers_count,omitempty"`
	WatchersCount   int32      `json:"watchers_count,omitempty"`
	ForksCount      int32      `json:"forks_count,omitempty"`
	Language        string     `json:"language,omitempty"`

	HasIssues bool   `json:"has_issues,omitempty"`
	HasWiki   bool   `json:"has_wiki,omitempty"`
	HasPage   bool   `json:"has_pages,omitempty"`
	License   string `json:"license,omitempty"`

	OpenIssuesCount int32  `json:"open_issues_count,omitempty"`
	DefaultBranch   string `json:"default_branch,omitempty"`
	Namespace       string `json:"namespace,omitempty"`

	NameWithNamespace string `json:"name_with_namespace,omitempty"`
	PathWithNamespace string `json:"path_with_namespace,omitempty"`
}
