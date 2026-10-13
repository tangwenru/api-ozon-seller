package ozon

// 主要解决 ogl 仓库的 fbs 大包（发运 carriage / 货位 container）流程。
// 依据文档：docs.ozon.ru/api/seller/zh —— 信任验收（货位）发运：
//   /v1/posting/carriage-available/list 看板
//   → /v1/carriage/create 创建发运（所有 awaiting_deliver 货件自动加入）
//   → /v1/carriage/approve 确认（containers_count = 大包数）
//   → /v2/posting/fbs/act/get-container-labels 打印大包（货位）标签
//   → /v1/carriage/cancel 取消

import (
	"context"
	"net/http"
)

// CarriageCreate 创建发运：所有状态为 awaiting_deliver 的货件自动加入，状态 new。
func (c FBS) CarriageCreate(
	ctx context.Context,
	params *FbsCarriageCreateParams,
) (*FbsCarriageCreateResponse, error) {
	url := "/v1/carriage/create"

	resp := &FbsCarriageCreateResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// CarriageApprove 确认发运（发运 → 已生成）。
// 信任验收开启时须传 containers_count（大包数）。
func (c FBS) CarriageApprove(
	ctx context.Context,
	params *FbsCarriageApproveParams,
) (*FbsCarriageApproveResponse, error) {
	url := "/v1/carriage/approve"

	resp := &FbsCarriageApproveResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// CarriageGet 发运信息。
func (c FBS) CarriageGet(
	ctx context.Context,
	params *FbsCarriageGetParams,
) (*FbsCarriageGetResponse, error) {
	url := "/v1/carriage/get"

	resp := &FbsCarriageGetResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// CarriageCancel 删除发运，货件回到下一次发运。
func (c FBS) CarriageCancel(
	ctx context.Context,
	params *FbsCarriageCancelParams,
) (*FbsCarriageCancelResponse, error) {
	url := "/v1/carriage/cancel"

	resp := &FbsCarriageCancelResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// CarriageSetPostings 全量覆盖发运中的货件列表（仅限 awaiting_deliver 状态货件）。
func (c FBS) CarriageSetPostings(
	ctx context.Context,
	params *FbsCarriageSetPostingsParams,
) (*FbsCarriageSetPostingsResponse, error) {
	url := "/v1/carriage/set-postings"

	resp := &FbsCarriageSetPostingsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// ActGetPostings 单据（发运）中的货件列表。
func (c FBS) ActGetPostings(
	ctx context.Context,
	params *FbsActGetPostingsParams,
) (*FbsActGetPostingsResponse, error) {
	url := "/v2/posting/fbs/act/get-postings"

	resp := &FbsActGetPostingsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// CarriageDeliveryListV2 配送方式与发运看板（v2，v1 已于 2026-03-20 关闭）。
func (c FBS) CarriageDeliveryListV2(
	ctx context.Context,
	params *CarriageDeliveryListV2Params,
) (*CarriageDeliveryListV2Response, error) {
	url := "/v2/carriage/delivery/list"

	resp := &CarriageDeliveryListV2Response{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// CarriagePassCreate 创建进仓通行证（trusted acceptance dropoff 确认发运前需要）。
func (c FBS) CarriagePassCreate(
	ctx context.Context,
	params *FbsCarriagePassCreateParams,
) (*FbsCarriagePassCreateResponse, error) {
	url := "/v1/carriage/pass/create"

	resp := &FbsCarriagePassCreateResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
