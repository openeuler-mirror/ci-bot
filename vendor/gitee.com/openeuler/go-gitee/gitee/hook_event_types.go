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
	Action *string `json:"action,omitempty"`
}

type PullRequestEvent struct {
	Action *string `json:"action,omitempty"`
}

type TagPushEvent struct {
	Action *string `json:"action,omitempty"`
}
