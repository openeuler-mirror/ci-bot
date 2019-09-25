/*
 * 码云 Open API
 *
 * No description provided (generated by Swagger Codegen https://github.com/swagger-api/swagger-codegen)
 *
 * API version: 5.3.2
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package gitee

// 创建分支
type CompleteBranch struct {
	Name string `json:"name,omitempty"`
	Commit string `json:"commit,omitempty"`
	Links string `json:"_links,omitempty"`
	Protected string `json:"protected,omitempty"`
	ProtectionUrl string `json:"protection_url,omitempty"`
}
