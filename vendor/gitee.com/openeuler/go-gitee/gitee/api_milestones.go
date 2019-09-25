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
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"fmt"
	"github.com/antihax/optional"
)

// Linger please
var (
	_ context.Context
)

type MilestonesApiService service

/* 
MilestonesApiService 删除仓库单个里程碑
删除仓库单个里程碑
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param owner 仓库所属空间地址(企业、组织或个人的地址path)
 * @param repo 仓库路径(path)
 * @param number 里程碑序号(id)
 * @param optional nil or *DeleteV5ReposOwnerRepoMilestonesNumberOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type DeleteV5ReposOwnerRepoMilestonesNumberOpts struct { 
	AccessToken optional.String
}

func (a *MilestonesApiService) DeleteV5ReposOwnerRepoMilestonesNumber(ctx context.Context, owner string, repo string, number int32, localVarOptionals *DeleteV5ReposOwnerRepoMilestonesNumberOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Delete")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/repos/{owner}/{repo}/milestones/{number}"
	localVarPath = strings.Replace(localVarPath, "{"+"owner"+"}", fmt.Sprintf("%v", owner), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repo"+"}", fmt.Sprintf("%v", repo), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"number"+"}", fmt.Sprintf("%v", number), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.AccessToken.IsSet() {
		localVarQueryParams.Add("access_token", parameterToString(localVarOptionals.AccessToken.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "multipart/form-data"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarHttpResponse, err
	}


	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		return localVarHttpResponse, newErr
	}

	return localVarHttpResponse, nil
}

/* 
MilestonesApiService 获取仓库所有里程碑
获取仓库所有里程碑
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param owner 仓库所属空间地址(企业、组织或个人的地址path)
 * @param repo 仓库路径(path)
 * @param optional nil or *GetV5ReposOwnerRepoMilestonesOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码
     * @param "State" (optional.String) -  里程碑状态: open, closed, all。默认: open
     * @param "Sort" (optional.String) -  排序方式: due_on
     * @param "Direction" (optional.String) -  升序(asc)或是降序(desc)。默认: asc
     * @param "Page" (optional.Int32) -  当前的页码
     * @param "PerPage" (optional.Int32) -  每页的数量，最大为 100

@return []Milestone
*/

type GetV5ReposOwnerRepoMilestonesOpts struct { 
	AccessToken optional.String
	State optional.String
	Sort optional.String
	Direction optional.String
	Page optional.Int32
	PerPage optional.Int32
}

func (a *MilestonesApiService) GetV5ReposOwnerRepoMilestones(ctx context.Context, owner string, repo string, localVarOptionals *GetV5ReposOwnerRepoMilestonesOpts) ([]Milestone, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue []Milestone
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/repos/{owner}/{repo}/milestones"
	localVarPath = strings.Replace(localVarPath, "{"+"owner"+"}", fmt.Sprintf("%v", owner), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repo"+"}", fmt.Sprintf("%v", repo), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.AccessToken.IsSet() {
		localVarQueryParams.Add("access_token", parameterToString(localVarOptionals.AccessToken.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.State.IsSet() {
		localVarQueryParams.Add("state", parameterToString(localVarOptionals.State.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Sort.IsSet() {
		localVarQueryParams.Add("sort", parameterToString(localVarOptionals.Sort.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Direction.IsSet() {
		localVarQueryParams.Add("direction", parameterToString(localVarOptionals.Direction.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Page.IsSet() {
		localVarQueryParams.Add("page", parameterToString(localVarOptionals.Page.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.PerPage.IsSet() {
		localVarQueryParams.Add("per_page", parameterToString(localVarOptionals.PerPage.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "multipart/form-data"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v []Milestone
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/* 
MilestonesApiService 获取仓库单个里程碑
获取仓库单个里程碑
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param owner 仓库所属空间地址(企业、组织或个人的地址path)
 * @param repo 仓库路径(path)
 * @param number 里程碑序号(id)
 * @param optional nil or *GetV5ReposOwnerRepoMilestonesNumberOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码

@return Milestone
*/

type GetV5ReposOwnerRepoMilestonesNumberOpts struct { 
	AccessToken optional.String
}

func (a *MilestonesApiService) GetV5ReposOwnerRepoMilestonesNumber(ctx context.Context, owner string, repo string, number int32, localVarOptionals *GetV5ReposOwnerRepoMilestonesNumberOpts) (Milestone, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue Milestone
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/repos/{owner}/{repo}/milestones/{number}"
	localVarPath = strings.Replace(localVarPath, "{"+"owner"+"}", fmt.Sprintf("%v", owner), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repo"+"}", fmt.Sprintf("%v", repo), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"number"+"}", fmt.Sprintf("%v", number), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if localVarOptionals != nil && localVarOptionals.AccessToken.IsSet() {
		localVarQueryParams.Add("access_token", parameterToString(localVarOptionals.AccessToken.Value(), ""))
	}
	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "multipart/form-data"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v Milestone
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/* 
MilestonesApiService 更新仓库里程碑
更新仓库里程碑
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param owner 仓库所属空间地址(企业、组织或个人的地址path)
 * @param repo 仓库路径(path)
 * @param number 里程碑序号(id)
 * @param title 里程碑标题
 * @param dueOn 里程碑的截止日期
 * @param optional nil or *PatchV5ReposOwnerRepoMilestonesNumberOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码
     * @param "State" (optional.String) -  里程碑状态: open, closed, all。默认: open
     * @param "Description" (optional.String) -  里程碑具体描述

@return Milestone
*/

type PatchV5ReposOwnerRepoMilestonesNumberOpts struct { 
	AccessToken optional.String
	State optional.String
	Description optional.String
}

func (a *MilestonesApiService) PatchV5ReposOwnerRepoMilestonesNumber(ctx context.Context, owner string, repo string, number int32, title string, dueOn string, localVarOptionals *PatchV5ReposOwnerRepoMilestonesNumberOpts) (Milestone, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Patch")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue Milestone
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/repos/{owner}/{repo}/milestones/{number}"
	localVarPath = strings.Replace(localVarPath, "{"+"owner"+"}", fmt.Sprintf("%v", owner), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repo"+"}", fmt.Sprintf("%v", repo), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"number"+"}", fmt.Sprintf("%v", number), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "multipart/form-data"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if localVarOptionals != nil && localVarOptionals.AccessToken.IsSet() {
		localVarFormParams.Add("access_token", parameterToString(localVarOptionals.AccessToken.Value(), ""))
	}
	localVarFormParams.Add("title", parameterToString(title, ""))
	if localVarOptionals != nil && localVarOptionals.State.IsSet() {
		localVarFormParams.Add("state", parameterToString(localVarOptionals.State.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Description.IsSet() {
		localVarFormParams.Add("description", parameterToString(localVarOptionals.Description.Value(), ""))
	}
	localVarFormParams.Add("due_on", parameterToString(dueOn, ""))
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 200 {
			var v Milestone
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}

/* 
MilestonesApiService 创建仓库里程碑
创建仓库里程碑
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param owner 仓库所属空间地址(企业、组织或个人的地址path)
 * @param repo 仓库路径(path)
 * @param title 里程碑标题
 * @param dueOn 里程碑的截止日期
 * @param optional nil or *PostV5ReposOwnerRepoMilestonesOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码
     * @param "State" (optional.String) -  里程碑状态: open, closed, all。默认: open
     * @param "Description" (optional.String) -  里程碑具体描述

@return Milestone
*/

type PostV5ReposOwnerRepoMilestonesOpts struct { 
	AccessToken optional.String
	State optional.String
	Description optional.String
}

func (a *MilestonesApiService) PostV5ReposOwnerRepoMilestones(ctx context.Context, owner string, repo string, title string, dueOn string, localVarOptionals *PostV5ReposOwnerRepoMilestonesOpts) (Milestone, *http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		localVarReturnValue Milestone
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/repos/{owner}/{repo}/milestones"
	localVarPath = strings.Replace(localVarPath, "{"+"owner"+"}", fmt.Sprintf("%v", owner), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repo"+"}", fmt.Sprintf("%v", repo), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHttpContentTypes := []string{"application/json", "multipart/form-data"}

	// set Content-Type header
	localVarHttpContentType := selectHeaderContentType(localVarHttpContentTypes)
	if localVarHttpContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHttpContentType
	}

	// to determine the Accept header
	localVarHttpHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHttpHeaderAccept := selectHeaderAccept(localVarHttpHeaderAccepts)
	if localVarHttpHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHttpHeaderAccept
	}
	if localVarOptionals != nil && localVarOptionals.AccessToken.IsSet() {
		localVarFormParams.Add("access_token", parameterToString(localVarOptionals.AccessToken.Value(), ""))
	}
	localVarFormParams.Add("title", parameterToString(title, ""))
	if localVarOptionals != nil && localVarOptionals.State.IsSet() {
		localVarFormParams.Add("state", parameterToString(localVarOptionals.State.Value(), ""))
	}
	if localVarOptionals != nil && localVarOptionals.Description.IsSet() {
		localVarFormParams.Add("description", parameterToString(localVarOptionals.Description.Value(), ""))
	}
	localVarFormParams.Add("due_on", parameterToString(dueOn, ""))
	r, err := a.client.prepareRequest(ctx, localVarPath, localVarHttpMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, localVarFileName, localVarFileBytes)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHttpResponse, err := a.client.callAPI(r)
	if err != nil || localVarHttpResponse == nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	localVarBody, err := ioutil.ReadAll(localVarHttpResponse.Body)
	localVarHttpResponse.Body.Close()
	if err != nil {
		return localVarReturnValue, localVarHttpResponse, err
	}

	if localVarHttpResponse.StatusCode < 300 {
		// If we succeed, return the data, otherwise pass on to decode error.
		err = a.client.decode(&localVarReturnValue, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
		if err == nil { 
			return localVarReturnValue, localVarHttpResponse, err
		}
	}

	if localVarHttpResponse.StatusCode >= 300 {
		newErr := GenericSwaggerError{
			body: localVarBody,
			error: localVarHttpResponse.Status,
		}
		
		if localVarHttpResponse.StatusCode == 201 {
			var v Milestone
			err = a.client.decode(&v, localVarBody, localVarHttpResponse.Header.Get("Content-Type"));
				if err != nil {
					newErr.error = err.Error()
					return localVarReturnValue, localVarHttpResponse, newErr
				}
				newErr.model = v
				return localVarReturnValue, localVarHttpResponse, newErr
		}
		
		return localVarReturnValue, localVarHttpResponse, newErr
	}

	return localVarReturnValue, localVarHttpResponse, nil
}
