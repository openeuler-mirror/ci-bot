/*
 * 码云 Open API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 5.3.2
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package gitee

// update pull request information
type PullRequestUpdateParam struct {
	// 用户授权码
	AccessToken string `json:"access_token,omitempty"`
	// 可选。Pull Request 标题
	Title string `json:"title,omitempty"`
	// 可选。Pull Request 内容
	Body string `json:"body,omitempty"`
	// 可选。Pull Request 状态
	State string `json:"state,omitempty"`
	// 可选。里程碑序号(id)
	MilestoneNumber int32 `json:"milestone_number,omitempty"`
	// 用逗号分开的标签，名称要求长度在 2-20 之间且非特殊字符。如: bug,performance
	Labels string `json:"labels,omitempty"`
	// 最少审查人数
	// change the int32 to *int32 manually, in order to pass 0
	AssigneesNumber *int32 `json:"assignees_number,omitempty"`
	// 最少测试人员
	// change the int32 to *int32 manually, in order to pass 0
	TestersNumber *int32 `json:"testers_number,omitempty"`
}
