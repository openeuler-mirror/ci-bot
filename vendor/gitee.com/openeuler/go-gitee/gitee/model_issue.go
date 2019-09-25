/*
 * 码云 Open API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 5.3.2
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package gitee

import (
	"time"
)

// 搜索 Issues
type Issue struct {
	Id int32 `json:"id,omitempty"`
	Url string `json:"url,omitempty"`
	RepositoryUrl string `json:"repository_url,omitempty"`
	LabelsUrl string `json:"labels_url,omitempty"`
	CommentsUrl string `json:"comments_url,omitempty"`
	HtmlUrl string `json:"html_url,omitempty"`
	ParentUrl string `json:"parent_url,omitempty"`
	Number string `json:"number,omitempty"`
	State string `json:"state,omitempty"`
	Title string `json:"title,omitempty"`
	Body string `json:"body,omitempty"`
	BodyHtml string `json:"body_html,omitempty"`
	User string `json:"user,omitempty"`
	Labels *Label `json:"labels,omitempty"`
	Assignee *UserBasic `json:"assignee,omitempty"`
	Collaborators *UserBasic `json:"collaborators,omitempty"`
	Repository string `json:"repository,omitempty"`
	Milestone *Milestone `json:"milestone,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	PlanStartedAt time.Time `json:"plan_started_at,omitempty"`
	Deadline time.Time `json:"deadline,omitempty"`
	FinishedAt time.Time `json:"finished_at,omitempty"`
	ScheduledTime string `json:"scheduled_time,omitempty"`
	Comments int32 `json:"comments,omitempty"`
	IssueType string `json:"issue_type,omitempty"`
	Program *ProgramBasic `json:"program,omitempty"`
}
