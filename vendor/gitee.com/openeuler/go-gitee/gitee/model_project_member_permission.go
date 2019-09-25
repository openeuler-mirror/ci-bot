/*
 * 码云 Open API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 5.3.2
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package gitee

// 查看仓库成员的权限
type ProjectMemberPermission struct {
	Id int32 `json:"id,omitempty"`
	Login string `json:"login,omitempty"`
	Name string `json:"name,omitempty"`
	AvatarUrl string `json:"avatar_url,omitempty"`
	Url string `json:"url,omitempty"`
	HtmlUrl string `json:"html_url,omitempty"`
	FollowersUrl string `json:"followers_url,omitempty"`
	FollowingUrl string `json:"following_url,omitempty"`
	GistsUrl string `json:"gists_url,omitempty"`
	StarredUrl string `json:"starred_url,omitempty"`
	SubscriptionsUrl string `json:"subscriptions_url,omitempty"`
	OrganizationsUrl string `json:"organizations_url,omitempty"`
	ReposUrl string `json:"repos_url,omitempty"`
	EventsUrl string `json:"events_url,omitempty"`
	ReceivedEventsUrl string `json:"received_events_url,omitempty"`
	Type_ string `json:"type,omitempty"`
	SiteAdmin string `json:"site_admin,omitempty"`
	Permission string `json:"permission,omitempty"`
}
