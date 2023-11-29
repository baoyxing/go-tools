// Copyright 2021 Tencent Inc. All rights reserved.

package validators

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/wechatpay-apiv3/wechatpay-go/core/auth"
)

// WechatPayNotifyValidator 微信支付 API v3 通知请求报文验证器
type WechatPayNotifyValidator struct {
	wechatPayValidator
}

// Validate 对接收到的微信支付 API v3 通知请求报文进行验证
func (v *WechatPayNotifyValidator) Validate(ctx context.Context, requestContext *app.RequestContext) error {
	return v.validateHTTPMessage(ctx, &requestContext.Request.Header, requestContext.Request.Body())
}

// NewWechatPayNotifyValidator 使用 auth.Verifier 初始化一个 WechatPayNotifyValidator
func NewWechatPayNotifyValidator(verifier auth.Verifier) *WechatPayNotifyValidator {
	return &WechatPayNotifyValidator{
		wechatPayValidator{verifier: verifier},
	}
}
