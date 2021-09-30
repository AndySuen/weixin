package authorizer

import (
	"context"
)

const apiModifyDomain = "/wxa/modify_domain"
const apiSetWebViewDomain = "/wxa/setwebviewdomain"

/*
设置服务器域名
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/Server_Address_Configuration.html
授权给第三方的小程序，其服务器域名只可以为在第三方平台账号中配置的小程序服务器域名，当小程序通过第三方平台发布代码上线后，小程序原先自己配置的服务器域名将被删除，只保留第三方平台的域名，所以第三方平台在代替小程序发布代码之前，需要调用接口为小程序添加第三方平台自身的域名。
POST https://api.weixin.qq.com/wxa/modify_domain?access_token=ACCESS_TOKEN
{
  "action": "set",
  "requestdomain": ["https://www.qq.com", "https://www.qq.com"],
  "wsrequestdomain": ["wss://www.qq.com", "wss://www.qq.com"],
  "uploaddomain": ["https://www.qq.com", "https://www.qq.com"],
  "downloaddomain": ["https://www.qq.com", "https://www.qq.com"]
}
*/
func (api *AuthorizerApi) ModifyDomain(ctx context.Context, postData string) error {
	return api.Client.HTTPPostJson(ctx, apiModifyDomain, postData, nil)
}

/*
设置业务域名
授权给第三方的小程序，其业务域名只可以为在第三方平台账号中配置的小程序业务域名，当小程序通过第三方发布代码上线后，小程序原先自己配置的业务域名将被删除，只保留第三方平台的域名，所以第三方平台在代替小程序发布代码之前，需要调用接口为小程序添加业务域名。
https://developers.weixin.qq.com/doc/oplatform/Third-party_Platforms/2.0/api/Mini_Program_Basic_Info/setwebviewdomain.html
POST https://api.weixin.qq.com/wxa/setwebviewdomain?access_token=ACCESS_TOKEN
{
  "action": "add",
  "webviewdomain": ["https://www.qq.com", "https://m.qq.com"]
}
*/
func (api *AuthorizerApi) SetWebViewDomain(ctx context.Context, postData string) error {
	return api.Client.HTTPPostJson(ctx, apiSetWebViewDomain, postData, nil)
}
