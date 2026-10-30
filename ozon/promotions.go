package ozon

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

type Promotions struct {
	client *ozonCore.Client
}

type GetAvailablePromotionsResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result []GetAvailablePromotionsResult `json:"result"`
}

type GetAvailablePromotionsResult struct {
	// Promotion identifier
	Id float64 `json:"id"`

	// Promotion name
	Title string `json:"title"`

	// Promotion type
	ActionType string `json:"action_type"`

	// Promotion description
	Description string `json:"description"`

	// Promotion start date
	DateStart string `json:"date_start"`

	// Promotion end date
	DateEnd string `json:"date_end"`

	// 商品自动加入活动的时间点（统一价格逻辑后新增）
	AutoAddDates []string `json:"auto_add_dates"`

	// Promotion freeze date.
	//
	// If the field is filled, the seller can't increase prices, change the list of products, or decrease the number of product units in the promotion.
	//
	// The seller can lower prices and increase the product units number in the promotion
	FreezeDate string `json:"freeze_date"`

	// Number of products that can participate in the promotion
	PotentialProductsCount float64 `json:"potential_products_count"`

	// Number of products that participate in the promotion
	ParticipatingProductsCount float64 `json:"participating_products_count"`

	// Whether or not you participate in the promotion
	IsParticipating bool `json:"is_participating"`

	// Indication that customers need a promo code to participate in the promotion
	IsVoucherAction bool `json:"is_voucher_action"`

	// Number of blocked products
	BannedProductsCount float64 `json:"banned_products_count"`

	// Indication of the promotion is with the target audience
	WithTargeting bool `json:"with_targeting"`

	// Order amount
	OrderAmount float64 `json:"order_amount"`

	// Discount type
	DiscountType string `json:"discount_type"`

	// Discount size
	DiscountValue float64 `json:"discount_value"`
}

// A method for getting a list of promotions that you can participate in
func (c Promotions) GetAvailablePromotions(ctx context.Context) (*GetAvailablePromotionsResponse, error) {
	url := "/v1/actions"

	resp := &GetAvailablePromotionsResponse{}

	response, err := c.client.Request(ctx, http.MethodGet, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type AddProductToPromotionParams struct {
	// Promotion identifier
	ActionId float64 `json:"action_id"`

	// Products list
	Products []AddProductToPromotionProduct `json:"products"`
}

// PromotionMoney 统一价格逻辑后的金额对象。
// action_price 从「活动价」变成了「价格上限(предельная цена)」。
type PromotionMoney struct {
	// Amount 金额，字符串，例如 "88" / "88.5"
	Amount string `json:"amount"`

	// Currency 币种，Ozon Seller API 固定 RUB
	Currency string `json:"currency"`
}

// PromotionProductId 兼容 Ozon 把 uint64 商品 ID 序列化成「数字」或「字符串」两种形态
type PromotionProductId uint64

func (p *PromotionProductId) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		*p = 0
		return nil
	}
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	*p = PromotionProductId(v)
	return nil
}

type AddProductToPromotionProduct struct {
	// Product identifier
	ProductId float64 `json:"product_id"`

	// 价格上限(предельная цена)：<= 活动限额则加入，> 限额则移出；同时写入商品卡价格
	ActionPrice PromotionMoney `json:"action_price"`

	// 仅「Промокоды」活动需要；「Эластичный бустинг」「Максимальный бустинг」不要传
	Stock float64 `json:"stock,omitempty"`
}

type AddProductToPromotionResponse struct {
	ozonCore.CommonResponse

	// 已加入/更新的商品 ID（Ozon 文档标 string，实际可能返回 number）
	ActiveProductIds []PromotionProductId `json:"active_product_ids"`

	// 被移出活动的商品 ID
	DeactivatedProductIds []PromotionProductId `json:"deactivated_product_ids"`

	// 未成功加入活动的商品
	Rejected []AddProductToPromotionResultRejected `json:"rejected"`

	// 被移出活动的原因
	Warnings []AddProductToPromotionResultRejected `json:"warnings"`
}

type AddProductToPromotionResultRejected struct {
	// Product identifier
	ProductId float64 `json:"product_id"`

	// Reason why the product wasn't added to the promotion
	Reason string `json:"reason"`
}

// AddToPromotion 参加活动 / 更新活动商品价格上限。
//
// 2026-10-13 起 Ozon 停用 /v1/actions/products/activate（统一价格逻辑），
// 该接口返回 400 {"code":9,"message":"the method is blocked due to unified price logic..."}。
// 改调 /v1/actions/products/update：action_price 是价格上限，会写入商品卡。
func (c Promotions) AddToPromotion(ctx context.Context, params *AddProductToPromotionParams) (*AddProductToPromotionResponse, error) {
	url := "/v1/actions/products/update"

	resp := &AddProductToPromotionResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ProductsAvailableForPromotionParams struct {
	// Promotion identifier
	ActionId float64 `json:"action_id"`

	// Number of values in the response. The default value is 100
	Limit float64 `json:"limit"`

	// Number of elements that will be skipped in the response.
	// For example, if offset=10, the response will start with the 11th element found
	Offset float64 `json:"offset,omitempty"`
}

type ProductsAvailableForPromotionResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result ProductsAvailableForPromotionResult `json:"result"`
}

type ProductsAvailableForPromotionResult struct {
	// Products list
	Products []PromotionProduct `json:"products"`

	// Total number of products that can participate in the promotion
	Total float64 `json:"total"`
}

type PromotionProduct struct {
	// Product identifier
	Id float64 `json:"id"`

	// Current product price without a discount
	Price float64 `json:"price"`

	// Promotional product price
	ActionPrice float64 `json:"action_price"`

	// Maximum possible promotional product price
	MaxActionPrice float64 `json:"max_action_price"`

	// Type of adding a product to the promotion: automatically or manually by the seller
	AddMode string `json:"add_mode"`

	// Minimum number of product units in a stock discount type promotion
	MinStock float64 `json:"min_stock"`

	// Number of product units in a stock discount type promotion
	Stock float64 `json:"stock"`
}

// A method for getting a list of products that can participate in the promotion by the promotion identifier
func (c Promotions) ProductsAvailableForPromotion(ctx context.Context, params *ProductsAvailableForPromotionParams) (*ProductsAvailableForPromotionResponse, error) {
	url := "/v1/actions/candidates"

	resp := &ProductsAvailableForPromotionResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ProductsInPromotionParams struct {
	// Promotion identifier
	ActionId float64 `json:"action_id"`

	// Number of values in the response. The default value is 100
	Limit float64 `json:"limit"`

	// Number of elements that will be skipped in the response. For example, if offset=10, the response will start with the 11th element found
	Offset float64 `json:"offset,omitempty"`
}

type ProductsInPromotionResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result ProductsInPromotionResult `json:"result"`
}

type ProductsInPromotionResult struct {
	// Products list
	Products []PromotionProduct `json:"products"`

	// Total number of products that can participate in the promotion
	Total float64 `json:"total"`
}

// A method for getting the list of products participating in the promotion by its identifier
func (c Promotions) ProductsInPromotion(ctx context.Context, params *ProductsInPromotionParams) (*ProductsInPromotionResponse, error) {
	url := "/v1/actions/products"

	resp := &ProductsInPromotionResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type RemoveProductFromPromotionParams struct {
	// Promotion identifier
	ActionId float64 `json:"action_id"`

	// List of products identifiers
	ProductIds []float64 `json:"product_ids"`
}

type RemoveProductFromPromotionResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result RemoveProductFromPromotionResult `json:"result"`
}

type RemoveProductFromPromotionResult struct {
	// List of product identifiers that were removed from the promotion
	ProductIds []float64 `json:"product_ids"`

	// List of product identifiers that weren't removed from the promotion
	Rejected []RemoveProductFromPromotionResultRejected `json:"rejected"`
}

type RemoveProductFromPromotionResultRejected struct {
	// Product identifier
	ProductId float64 `json:"product_id"`

	// Reason why the product wasn't removed from the promotion
	Reason string `json:"reason"`
}

// A method for removing products from the promotion
func (c Promotions) RemoveProduct(ctx context.Context, params *RemoveProductFromPromotionParams) (*RemoveProductFromPromotionResponse, error) {
	url := "/v1/actions/products/deactivate"

	resp := &RemoveProductFromPromotionResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ListDiscountRequestsParams struct {
	// Discount request status: ALL / NEW / APPROVED / DECLINED
	Status ListDiscountRequestsStatus `json:"status" default:"ALL"`

	// Identifier of the last value on the page. Leave empty for the first request
	LastId int64 `json:"last_id,omitempty"`

	// The maximum number of requests on a page: 5 / 10 / 15 / 20 / 30 / 50 (default 50)
	Limit uint64 `json:"limit" default:"50"`
}

type ListDiscountRequestsResponse struct {
	ozonCore.CommonResponse

	// List of requests
	Tasks []ListDiscountRequestsResult `json:"tasks"`
}

type ListDiscountRequestsResult struct {
	// Request ID
	Id uint64 `json:"id"`

	// Request created date
	CreatedAt time.Time `json:"created_at"`

	// End time of the request
	EndAt time.Time `json:"end_at"`

	// Time to change the decision
	EditedTill time.Time `json:"edited_till"`

	// Request status
	Status string `json:"status"`

	// Customer's name
	CustomerName string `json:"customer_name"`

	// Product identifier in the Ozon system, SKU
	SKU uint64 `json:"sku"`

	// Customer's comment on the request
	UserComment string `json:"user_comment"`

	// Seller's comment on the request
	SellerComment string `json:"seller_comment"`

	// Requested price
	RequestedPrice float64 `json:"requested_price"`

	// Approved price
	ApprovedPrice float64 `json:"approved_price"`

	// Product price before all discounts
	OriginalPrice float64 `json:"original_price"`

	// Discount in rubles
	Discount float64 `json:"discount"`

	// Discount percentage
	DiscountPercent float64 `json:"discount_percent"`

	// Base price at which a product is selling on Ozon, if not eligible for a promotion
	BasePrice float64 `json:"base_price"`

	// The minimum price after auto-application of discounts and promotions
	MinAutoPrice float64 `json:"min_auto_price"`

	// ID of the previous customer request for this product
	PrevTaskId uint64 `json:"prev_task_id"`

	// If product is damaged — true
	IsDamaged bool `json:"is_damaged"`

	// Moderation date: review, approval or decline of the request
	ModeratedAt time.Time `json:"moderated_at"`

	// Discount in rubles approved by the seller. Pass the value 0 if the seller did not approve the request
	ApprovedDiscount float64 `json:"approved_discount"`

	// Discount percentage approved by the seller. Pass the value 0 if the seller did not approve the request
	ApprovedDiscountPercent float64 `json:"approved_discount_percent"`

	// Whether the customer has purchased the product. true if purchased
	IsPurchased bool `json:"is_purchased"`

	// Whether the request was moderated automatically. true if moderation was automatic
	IsAutoModerated bool `json:"is_auto_moderated"`

	// Product identifier in the seller's system
	OfferId string `json:"offer_id"`

	// Email of the user who processed the request
	Email string `json:"email"`

	// Last name of the user who processed the request
	LastName string `json:"last_name"`

	// First name of the user who processed the request
	FirstName string `json:"first_name"`

	// Patronymic of the user who processed the request
	Patronymic string `json:"patronymic"`

	// Approved minimum quantity of products
	ApprovedQuantityMin uint64 `json:"approved_quantity_min"`

	// Approved maximum quantity of products
	ApprovedQuantityMax uint64 `json:"approved_quantity_max"`

	// Requested minimum number of products
	RequestedQuantityMin uint64 `json:"requested_quantity_min"`

	// Requested maximum number of products
	RequestedQuantityMax uint64 `json:"requested_quantity_max"`

	// Requested price with fee
	RequestedPriceWithFee float64 `json:"requested_price_with_fee"`

	// Approved price with fee
	ApprovedPriceWithFee float64 `json:"approved_price_with_fee"`

	// Approved price fee percent
	ApprovedPriceFeePercent float64 `json:"approved_price_fee_percent"`

	// Product name
	Name string `json:"name"`

	// Difference between the buyer's price and the seller's price when the request was created
	ReductionFactor float64 `json:"reduction_factor"`
}

// /v1/actions/discounts-task/list 已弃用（未来停用），切换至 /v2/actions/discounts-task/list。
// Method for getting a list of products that customers want to buy with discount
func (c Promotions) ListDiscountRequests(ctx context.Context, params *ListDiscountRequestsParams) (*ListDiscountRequestsResponse, error) {
	url := "/v2/actions/discounts-task/list"

	resp := &ListDiscountRequestsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type DiscountRequestParams struct {
	// List of discount requests
	Tasks []DiscountRequestTask `json:"tasks"`
}

type DiscountRequestTask struct {
	// Request ID
	Id uint64 `json:"id"`

	// Approved price
	ApprovedPrice float64 `json:"approved_price"`

	// Seller's comment on the request
	SellerComment string `json:"seller_comment"`

	// Approved minimum quantity of products
	ApprovedQuantityMin uint64 `json:"approved_quantity_min"`

	// Approved maximum quantity of products
	ApprovedQuantityMax uint64 `json:"approved_quantity_max"`
}

type DiscountRequestResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result DiscountRequestResult `json:"result"`
}

type DiscountRequestResult struct {
	// Errors when creating a request
	FailDetails []DiscountRequestResultFailDetail `json:"fail_details"`

	// The number of requests with a successful status change
	SuccessCount int32 `json:"success_count"`

	// The number of requests that failed to change their status
	FailCount int32 `json:"fail_count"`
}

type DiscountRequestResultFailDetail struct {
	// Request ID
	TaskId uint64 `json:"task_id"`

	// Error message
	ErrorForUser string `json:"error_for_user"`
}

// You can approve applications in statuses:
//   - NEW — new
//   - SEEN — viewed
func (c Promotions) ApproveDiscountRequest(ctx context.Context, params *DiscountRequestParams) (*DiscountRequestResponse, error) {
	url := "/v1/actions/discounts-task/approve"

	resp := &DiscountRequestResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// You can decline applications in statuses:
//   - NEW—new
//   - SEEN—viewed
func (c Promotions) DeclineDiscountRequest(ctx context.Context, params *DiscountRequestParams) (*DiscountRequestResponse, error) {
	url := "/v1/actions/discounts-task/decline"

	resp := &DiscountRequestResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
