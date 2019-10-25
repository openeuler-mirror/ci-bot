package gitee

type NoteEvent struct {
	Action        *string          `json:"action,omitempty"`
	Comment       *Note            `json:"comment,omitempty"`
	Repository    *Project         `json:"repository,omitempty"`
	Project       *Project         `json:"project,omitempty"`
	Author        *User            `json:"author,omitempty"`
	Sender        *User            `json:"sender,omitempty"`
	URL           *string          `json:"url,omitempty"`
	Note          *string          `json:"note,omitempty"`
	NoteableType  *string          `json:"noteable_type,omitempty"`
	NoteableID    int64            `json:"noteable_id,omitempty"`
	Title         *string          `json:"title,omitempty"`
	PerIID        *string          `json:"per_iid,omitempty"`
	ShortCommitID *string          `json:"short_commit_id,omitempty"`
	Enterprise    *EnterpriseBasic `json:"enterprise,omitempty"`
	PullRequest   *PullRequest     `json:"pull_request,omitempty"`
	Issue         *Issue           `json:"issue,omitempty"`
	HookName      *string          `json:"hook_name,omitempty"`
	Password      *string          `json:"password,omitempty"`
}

type PushEvent struct {
	Action *string `json:"action,omitempty"`
}

type IssueEvent struct {
	Action      *string          `json:"action,omitempty"`
	Issue       *Issue           `json:"issue,omitempty"`
	Repository  *Project         `json:"repository,omitempty"`
	Project     *Project         `json:"project,omitempty"`
	Sender      *User            `json:"sender,omitempty"`
	TargetUser  *User            `json:"target_user,omitempty"`
	User        *User            `json:"user,omitempty"`
	Assignee    *UserBasic       `json:"assignee,omitempty"`
	UpdatedBy   *User            `json:"updated_by,omitempty"`
	IID         string           `json:"iid,omitempty"`
	Title       *string          `json:"title,omitempty"`
	Description *string          `json:"description,omitempty"`
	State       *string          `json:"state,omitempty"`
	Milestone   *Milestone       `json:"milestone,omitempty"`
	URL         *string          `json:"url,omitempty"`
	Enterprise  *EnterpriseBasic `json:"enterprise,omitempty"`
	HookName    *string          `json:"hook_name,omitempty"`
	Password    *string          `json:"password,omitempty"`
}

type RepoInfo struct {
	Project    *Project `json:"project,omitempty"`
	Repository *Project `json:"repository,omitempty"`
}

type PullRequestEvent struct {
	Action         *string          `json:"action,omitempty"`
	PullRequest    *PullRequest     `json:"pull_request,omitempty"`
	Number         int64            `json:"number,omitempty"`
	IID            int64            `json:"iid,omitempty"`
	Title          *string          `json:"title,omitempty"`
	Body           *string          `json:"body,omitempty"`
	State          *string          `json:"state,omitempty"`
	MergeStatus    *string          `json:"merge_status,omitempty"`
	MergeCommitSha *string          `json:"merge_commit_sha,omitempty"`
	URL            *string          `json:"url,omitempty"`
	SourceBranch   *string          `json:"source_branch,omitempty"`
	SourceRepo     *RepoInfo        `json:"source_repo,omitempty"`
	TargetBranch   *string          `json:"target_branch,omitempty"`
	TargetRepo     *RepoInfo        `json:"target_repo,omitempty"`
	Project        *Project         `json:"project,omitempty"`
	Repository     *Project         `json:"repository,omitempty"`
	Author         *User            `json:"author,omitempty"`
	UpdatedBy      *User            `json:"updated_by,omitempty"`
	Sender         *User            `json:"sender,omitempty"`
	TargetUser     *User            `json:"target_user,omitempty"`
	Enterprise     *EnterpriseBasic `json:"enterprise,omitempty"`
	HookName       *string          `json:"hook_name,omitempty"`
	Password       *string          `json:"password,omitempty"`
}

type TagPushEvent struct {
	Action *string `json:"action,omitempty"`
}
