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

type MiscellaneousApiService service

/* 
MiscellaneousApiService 列出可使用的 Emoji
列出可使用的 Emoji
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param optional nil or *GetV5EmojisOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5EmojisOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5Emojis(ctx context.Context, localVarOptionals *GetV5EmojisOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/emojis"

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
MiscellaneousApiService 列出可使用的 .gitignore 模板
列出可使用的 .gitignore 模板
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param optional nil or *GetV5GitignoreTemplatesOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5GitignoreTemplatesOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5GitignoreTemplates(ctx context.Context, localVarOptionals *GetV5GitignoreTemplatesOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/gitignore/templates"

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
MiscellaneousApiService 获取一个 .gitignore 模板
获取一个 .gitignore 模板
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param name .gitignore 模板名
 * @param optional nil or *GetV5GitignoreTemplatesNameOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5GitignoreTemplatesNameOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5GitignoreTemplatesName(ctx context.Context, name string, localVarOptionals *GetV5GitignoreTemplatesNameOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/gitignore/templates/{name}"
	localVarPath = strings.Replace(localVarPath, "{"+"name"+"}", fmt.Sprintf("%v", name), -1)

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
MiscellaneousApiService 获取一个 .gitignore 模板原始文件
获取一个 .gitignore 模板原始文件
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param name .gitignore 模板名
 * @param optional nil or *GetV5GitignoreTemplatesNameRawOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5GitignoreTemplatesNameRawOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5GitignoreTemplatesNameRaw(ctx context.Context, name string, localVarOptionals *GetV5GitignoreTemplatesNameRawOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/gitignore/templates/{name}/raw"
	localVarPath = strings.Replace(localVarPath, "{"+"name"+"}", fmt.Sprintf("%v", name), -1)

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
MiscellaneousApiService 列出可使用的开源许可协议
列出可使用的开源许可协议
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param optional nil or *GetV5LicensesOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5LicensesOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5Licenses(ctx context.Context, localVarOptionals *GetV5LicensesOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/licenses"

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
MiscellaneousApiService 获取一个开源许可协议
获取一个开源许可协议
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param license 协议名称
 * @param optional nil or *GetV5LicensesLicenseOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5LicensesLicenseOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5LicensesLicense(ctx context.Context, license string, localVarOptionals *GetV5LicensesLicenseOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/licenses/{license}"
	localVarPath = strings.Replace(localVarPath, "{"+"license"+"}", fmt.Sprintf("%v", license), -1)

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
MiscellaneousApiService 获取一个开源许可协议原始文件
获取一个开源许可协议原始文件
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param license 协议名称
 * @param optional nil or *GetV5LicensesLicenseRawOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5LicensesLicenseRawOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5LicensesLicenseRaw(ctx context.Context, license string, localVarOptionals *GetV5LicensesLicenseRawOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/licenses/{license}/raw"
	localVarPath = strings.Replace(localVarPath, "{"+"license"+"}", fmt.Sprintf("%v", license), -1)

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
MiscellaneousApiService 获取一个仓库使用的开源许可协议
获取一个仓库使用的开源许可协议
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param owner 仓库所属空间地址(企业、组织或个人的地址path)
 * @param repo 仓库路径(path)
 * @param optional nil or *GetV5ReposOwnerRepoLicenseOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type GetV5ReposOwnerRepoLicenseOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) GetV5ReposOwnerRepoLicense(ctx context.Context, owner string, repo string, localVarOptionals *GetV5ReposOwnerRepoLicenseOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Get")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/repos/{owner}/{repo}/license"
	localVarPath = strings.Replace(localVarPath, "{"+"owner"+"}", fmt.Sprintf("%v", owner), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"repo"+"}", fmt.Sprintf("%v", repo), -1)

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
MiscellaneousApiService 渲染 Markdown 文本
渲染 Markdown 文本
 * @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 * @param text Markdown 文本
 * @param optional nil or *PostV5MarkdownOpts - Optional Parameters:
     * @param "AccessToken" (optional.String) -  用户授权码


*/

type PostV5MarkdownOpts struct { 
	AccessToken optional.String
}

func (a *MiscellaneousApiService) PostV5Markdown(ctx context.Context, text string, localVarOptionals *PostV5MarkdownOpts) (*http.Response, error) {
	var (
		localVarHttpMethod = strings.ToUpper("Post")
		localVarPostBody   interface{}
		localVarFileName   string
		localVarFileBytes  []byte
		
	)

	// create path and map variables
	localVarPath := a.client.cfg.BasePath + "/v5/markdown"

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
	localVarFormParams.Add("text", parameterToString(text, ""))
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
