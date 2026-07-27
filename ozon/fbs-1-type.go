package ozon

import ozonCore "github.com/tangwenru/api-ozon-seller"

type FbsCarriageCreateParams struct {
	// 表示需要创建包含可追溯商品的发运, 一般为 false
	AllBLRTraceable  bool   `json:"all_blr_traceable,omitempty"`
	DeliveryMethodId int64  `json:"delivery_method_id"`
	DepartureDate    string `json:"departure_date,omitempty"` // 格式: 2006-01-02T15:04:05Z
}

type FbsCarriageCreateResponse struct {
	ozonCore.CommonResponse

	// Array of shipments
	Result GFbsCarriageCreateResult `json:"result"`
}

type GFbsCarriageCreateResult struct {
	CarriageId int64 `json:"carriage_id"`
}
