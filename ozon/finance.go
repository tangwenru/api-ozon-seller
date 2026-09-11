package ozon

import (
	"context"
	"net/http"
	"time"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

type Finance struct {
	client *ozonCore.Client
}

type ReportOnSoldProductsParams struct {
	// Month
	Month int32 `json:"month"`

	// Year
	Year int32 `json:"year"`
}

type ReportOnSoldProductsResponse struct {
	ozonCore.CommonResponse

	// Query result
	Result ReportonSoldProductsResult `json:"result"`
}

type ReportonSoldProductsResult struct {
	// Report title page
	Header ReportOnSoldProductsResultHeader `json:"header"`

	// Report table
	Rows []ReportOnSoldProductsResultRow `json:"rows"`
}

type ReportOnSoldProductsResultHeader struct {
	// Report ID
	Id string `json:"number"`

	// Report generation date
	DocDate string `json:"doc_date"`

	// Date of the offer agreement
	ContractDate string `json:"contract_date"`

	// Offer agreement number
	ContractNum string `json:"contract_number"`

	// Currency of your prices
	CurrencySysName string `json:"currency_sys_name"`

	// Amount to accrue
	DocAmount float64 `json:"doc_amount"`

	// Amount to accrue with VAT
	VATAmount float64 `json:"vat_amount"`

	// Payer's TIN
	PayerINN string `json:"payer_inn"`

	// Payer's Tax Registration Reason Code (KPP)
	PayerKPP string `json:"payer_kpp"`

	// Payer's name
	PayerName string `json:"payer_name"`

	// Recipient's TIN
	RecipientINN string `json:"receiver_inn"`

	// Recipient's Tax Registration Reason Code (KPP)
	RecipientKPP string `json:"receiver_kpp"`

	// Recipient's name
	RecipientName string `json:"receiver_name"`

	// Period start in the report
	StartDate string `json:"start_date"`

	// Period end in the report
	StopDate string `json:"stop_date"`
}

type ReportOnSoldProductsResultRow struct {
	// Row number
	RowNumber int32 `json:"rowNumber"`

	// Product Information
	Item ReturnOnSoldProduct `json:"item"`

	// Commission including the quantity of products, discounts and extra charges.
	// Ozon compensates it for the returned products
	ReturnCommission ReturnCommission `json:"return_commission"`

	// Percentage of sales commission by category
	CommissionRatio float64 `json:"commission_ratio"`

	// Delivery fee
	DeliveryCommission ReturnCommission `json:"delivery_commission"`

	// Seller's discounted price
	SellerPricePerInstance float64 `json:"seller_price_per_instance"`
}

type ReturnOnSoldProduct struct {
	// Product name
	ProductName string `json:"name"`

	// Product barcode
	Barcode string `json:"barcode"`

	// Product identifier in the seller's system
	OfferId string `json:"offer_id"`

	SKU int64 `json:"sku"`
}

type ReturnCommission struct {
	// Amount
	Amount float64 `json:"amount"`

	// Points for discounts
	Bonus float64 `json:"bonus"`

	// Commission for sold products, including discounts and extra charges
	Commission float64 `json:"commission"`

	// Additional payment at the expense of Ozon
	Compensation float64 `json:"compensation"`

	// Price per item
	PricePerInstance float64 `json:"price_per_instance"`

	// Product quantity
	Quantity int32 `json:"quantity"`

	// Ozon referral fee
	StandardFee float64 `json:"standard_fee"`

	// Payouts on partner loyalty mechanics: green prices
	BankCoinvestment float64 `json:"bank_coinvestment"`

	// Payouts on partner loyalty mechanics: stars
	Stars float64 `json:"stars"`

	// Total accrual
	Total float64 `json:"total"`
}

// Returns information on products sold and returned within a month. Canceled or non-purchased products are not included.
//
// Report is returned no later than the 5th day of the next month
func (c Finance) ReportOnSoldProducts(ctx context.Context, params *ReportOnSoldProductsParams) (*ReportOnSoldProductsResponse, error) {
	url := "/v2/finance/realization"

	resp := &ReportOnSoldProductsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetTotalTransactionsSumParams struct {
	// Filter by date
	Date GetTotalTransactionsSumDate `json:"date"`

	// Shipment number
	PostingNumber string `json:"posting_number"`

	// Transaction type:
	//
	//   - all — all,
	//   - orders — orders,
	//   - returns — returns and cancellations,
	//   - services — service fees,
	//   - compensation — compensation,
	//   - transferDelivery — delivery cost,
	//   - other — other
	TransactionType string `json:"transaction_type"`
}

type GetTotalTransactionsSumDate struct {
	// Period start.
	//
	// Format: YYYY-MM-DDTHH:mm:ss.sssZ.
	// Example: 2019-11-25T10:43:06.51
	From time.Time `json:"from"`

	// Period end.
	//
	// Format: YYYY-MM-DDTHH:mm:ss.sssZ.
	// Example: 2019-11-25T10:43:06.51
	To time.Time `json:"to"`
}

type GetTotalTransactionsSumResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result GetTotalTransactionsSumResult `json:"result"`
}

type GetTotalTransactionsSumResult struct {
	// Total cost of products and returns for specified period
	AccrualsForSale float64 `json:"accruals_for_sale"`

	// Compensations
	CompensationAmount float64 `json:"compensation_amount"`

	// Charges for delivery and returns when working under rFBS scheme
	MoneyTransfer float64 `json:"money_transfer"`

	// Other accurals
	OthersAmount float64 `json:"others_amount"`

	// Cost of shipment processing, orders packaging, pipeline and last mile services, and delivery cost before the new commissions and rates applied from February 1, 2021.
	//
	// Pipeline is delivery of products from one cluster to another.
	//
	// Last mile is products delivery to the pick-up point, parcle terminal, or by courier
	ProcessingAndDelivery float64 `json:"processing_and_delivery"`

	// Cost of reverse pipeline, returned, canceled and unredeemed orders processing, and return cost before the new commissions and rates applied from February 1, 2021.
	//
	// Pipeline is delivery of products from one cluster to another.
	//
	// Last mile is products delivery to the pick-up point, parcle terminal, or by courier
	RefundsAndCancellations float64 `json:"refunds_and_cancellations"`

	// The commission withheld when the product was sold and refunded when the product was returned
	SaleCommission float64 `json:"sale_commission"`

	// The additional services cost that are not directly related to deliveries and returns.
	// For example, promotion or product placement
	ServicesAmount float64 `json:"services_amount"`
}

// Deprecated: /v3/finance/transaction/totals 已于 2026-09-08 停用。
// 官方要求改用 /v1/finance/accrual/postings、/v1/finance/accrual/types、/v1/finance/accrual/by-day。
// Returns total sums for transactions for specified period
func (c Finance) GetTotalTransactionsSum(ctx context.Context, params *GetTotalTransactionsSumParams) (*GetTotalTransactionsSumResponse, error) {
	url := "/v3/finance/transaction/totals"

	resp := &GetTotalTransactionsSumResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ListTransactionsParams struct {
	// Filter
	Filter ListTransactionsFilter `json:"filter"`

	// Number of the page returned in the request
	Page int64 `json:"page"`

	// Number of items on the page
	PageSize int64 `json:"page_size"`
}

type ListTransactionsFilter struct {
	// Filter by date
	Date ListTransactionsFilterDate `json:"date"`

	// Operation type
	OperationType []string `json:"operation_type"`

	// Shipment number
	PostingNumber string `json:"posting_number"`

	// Transaction type
	TransactionType string `json:"transaction_type"`
}

type ListTransactionsFilterDate struct {
	// Period start.
	//
	// Format: YYYY-MM-DDTHH:mm:ss.sssZ.
	// Example: 2019-11-25T10:43:06.51
	From time.Time `json:"from"`

	// Period end.
	//
	// Format: YYYY-MM-DDTHH:mm:ss.sssZ.
	// Example: 2019-11-25T10:43:06.51
	To time.Time `json:"to"`
}

type ListTransactionsResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result ListTransactionsResult `json:"result"`
}

type ListTransactionsResult struct {
	// Transactions infromation
	Operations []ListTransactionsResultOperation `json:"operations"`

	// Number of pages. If 0, there are no more pages
	PageCount int64 `json:"page_count"`

	// Number of transactions on all pages. If 0, there are no more transactions
	RowCount int64 `json:"row_count"`
}

type ListTransactionsResultOperation struct {
	// Cost of the products with seller's discounts applied
	AccrualsForSale float64 `json:"accruals_for_sale"`

	// Total transaction sum
	Amount float64 `json:"amount"`

	// Delivery cost for charges by rates that were in effect until February 1, 2021, and for charges for bulky products
	DeliveryCharge float64 `json:"delivery_charge"`

	// Product information
	Items []ListTransactionsResultOperationItem `json:"items"`

	// Operation date
	OperationDate string `json:"operation_date"`

	// Operation identifier
	OperationId int64 `json:"operation_id"`

	// Operation type
	OperationType string `json:"operation_type"`

	// Operation type name
	OperationTypeName string `json:"operation_type_name"`

	// Shipment information
	Posting ListTransactionsResultOperationPosting `json:"posting"`

	// Returns and cancellation cost for charges by rates that were in effect until February 1, 2021, and for charges for bulky products
	ReturnDeliveryCharge float64 `json:"return_delivery_charge"`

	// Sales commission or sales commission refund
	SaleCommission float64 `json:"sale_commission"`

	// Additional services
	Services []ListTransactionsResultOperationService `json:"services"`

	// Transaction type
	Type string `json:"type"`
}

type ListTransactionsResultOperationItem struct {
	// Product name
	Name string `json:"name"`

	// Product identifier in the Ozon system, SKU
	SKU int64 `json:"sku"`
}

type ListTransactionsResultOperationPosting struct {
	// Delivery scheme
	DeliverySchema string `json:"delivery_schema"`

	// Date the product was accepted for processing
	OrderDate string `json:"order_date"`

	// Shipment number
	PostingNumber string `json:"posting_number"`

	// Warehouse identifier
	WarehouseId int64 `json:"warehouse_id"`
}

type ListTransactionsResultOperationService struct {
	// Service name
	Name TransactionOperationService `json:"name"`

	// Price
	Price float64 `json:"price"`
}

// Deprecated: /v3/finance/transaction/list 已于 2026-09-08 停用。
// 官方要求改用 /v1/finance/accrual/postings、/v1/finance/accrual/types、/v1/finance/accrual/by-day。
// Returns detailed information on all accruals. The maximum period for which you can get information in one request is 1 month.
//
// If you don't specify the posting_number in request, the response contains all shipments for the specified period or shipments of a certain type
func (c Finance) ListTransactions(ctx context.Context, params *ListTransactionsParams) (*ListTransactionsResponse, error) {
	url := "/v3/finance/transaction/list"

	resp := &ListTransactionsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetReportParams struct {
	// Time period in the YYYY-MM format
	Date string `json:"date"`

	// Response language
	Language string `json:"language"`
}

type ReportResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result ReportResult `json:"result"`
}

type ReportResult struct {
	// Unique report identifier
	Code string `json:"code"`
}

// Use the method to get mutual settlements report.
// Matches the Finance → Documents → Analytical reports → Mutual
// settlements report section in your personal account.
func (c Finance) MutualSettlements(ctx context.Context, params *GetReportParams) (*ReportResponse, error) {
	url := "/v1/finance/mutual-settlement"

	resp := &ReportResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// Use the method to get sales to legal entities report.
// Matches the Finance → Documents → Legal
// entities sales register section in your personal account.
func (c Finance) SalesToLegalEntities(ctx context.Context, params *GetReportParams) (*ReportResponse, error) {
	url := "/v1/finance/mutual-settlement"

	resp := &ReportResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// ===== /v1/finance/accrual/* 系列（/v3/finance/transaction/* 的官方替代品） =====

type AccrualAmount struct {
	// 金额（字符串，防止精度丢失）
	Amount string `json:"amount"`

	// 币种
	Currency string `json:"currency"`
}

// 按货件统计的应计项目
type AccrualPostingsParams struct {
	// 货件编号，[ 1 .. 200 ]
	PostingNumbers []string `json:"posting_numbers"`
}

type AccrualPostingsResponse struct {
	ozonCore.CommonResponse

	// 按货件统计的应计项目列表
	PostingAccruals []AccrualPostingsItem `json:"posting_accruals"`
}

type AccrualPostingsItem struct {
	// 应计项目列表
	Accruals []AccrualPostingsItemAccrual `json:"accruals"`

	// 货件编号
	PostingNumber string `json:"posting_number"`
}

type AccrualPostingsItemAccrual struct {
	// 应计日期
	AccrualDate string `json:"accrual_date"`

	// 应计金额
	Accrued AccrualAmount `json:"accrued"`

	// 数量
	Quantity int64 `json:"quantity"`

	// 卖家价格
	SellerPrice AccrualAmount `json:"seller_price"`

	// SKU
	SKU int64 `json:"sku"`

	// 应计项目类型标识符，参考 AccrualTypes 的 id
	TypeId int64 `json:"type_id"`
}

// 获取按货件统计的应计项目。
func (c Finance) AccrualPostings(ctx context.Context, params *AccrualPostingsParams) (*AccrualPostingsResponse, error) {
	url := "/v1/finance/accrual/postings"

	resp := &AccrualPostingsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// 获取应计项目参考信息。
func (c Finance) AccrualTypes(ctx context.Context) (*AccrualTypesResponse, error) {
	url := "/v1/finance/accrual/types"

	resp := &AccrualTypesResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, nil, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type AccrualTypesResponse struct {
	ozonCore.CommonResponse

	// 应计项目相关信息
	AccrualTypes []AccrualTypesItem `json:"accrual_types"`
}

type AccrualTypesItem struct {
	// 应计项目说明
	Description string `json:"description"`

	// 应计项目标识符
	Id int32 `json:"id"`

	// 应计项目名称
	Name string `json:"name"`
}

// 获取某日应计项目
type AccrualByDayParams struct {
	// 应计日期。YYYY-MM-DD。最早可查询日期为 2022-01-01。
	// 如果指定 last_id，请传递上一个请求中的 date 值。
	Date string `json:"date"`

	// 页面中最后一个值的标识符。首次请求请留空。
	// 标识符有效期为 15 分钟。
	LastId string `json:"last_id"`
}

type AccrualByDayResponse struct {
	ozonCore.CommonResponse

	// 货件的应计项目列表
	Accruals []AccrualByDayItem `json:"accruals"`

	// 页面中最后一个值的标识符（有效期 15 分钟）
	LastId string `json:"last_id"`
}

type AccrualByDayItem struct {
	// 应计项目类别
	AccruedCategory string `json:"accrued_category"`

	// 集装箱相关费用
	ContainerFees AccrualByDayFees `json:"container_fees"`

	// 应计日期
	Date string `json:"date"`

	// 商品相关费用
	ItemFees AccrualByDayItemFees `json:"item_fees"`

	// 非商品相关费用
	NonItemFee AccrualByDayFee `json:"non_item_fee"`

	// 货件信息
	Posting AccrualByDayPosting `json:"posting"`

	// 总金额
	TotalAmount AccrualAmount `json:"total_amount"`

	// 应计项目标识符
	AccrualId int64 `json:"accrual_id"`

	// 单位编号
	UnitNumber string `json:"unit_number"`
}

type AccrualByDayFees struct {
	Fees []AccrualByDayFee `json:"fees"`
}

type AccrualByDayFee struct {
	Accrued AccrualAmount `json:"accrued"`
	TypeId  int64         `json:"type_id"`
}

type AccrualByDayItemFees struct {
	Fees []AccrualByDayItemFee `json:"fees"`
}

type AccrualByDayItemFee struct {
	Fees []AccrualByDayFee `json:"fees"`
	SKU  int64             `json:"sku"`
}

type AccrualByDayPosting struct {
	// 配送方案
	DeliverySchema string `json:"delivery_schema"`

	// 配送速度
	DeliverySpeed int64 `json:"delivery_speed"`

	// 商品信息
	Products []AccrualByDayPostingProduct `json:"products"`
}

type AccrualByDayPostingProduct struct {
	// 佣金信息
	Commission AccrualByDayProductCommission `json:"commission"`

	// 配送费用信息
	Delivery AccrualByDayProductDelivery `json:"delivery"`

	// SKU
	SKU int64 `json:"sku"`
}

type AccrualByDayProductCommission struct {
	Bonus           AccrualAmount `json:"bonus"`
	Coinvestment    AccrualAmount `json:"coinvestment"`
	Commission      AccrualAmount `json:"commission"`
	CommissionRatio string        `json:"commission_ratio"`
	SaleAmount      AccrualAmount `json:"sale_amount"`
	SaleCommission  AccrualAmount `json:"sale_commission"`
	SalePrice       AccrualAmount `json:"sale_price"`
	SellerPrice     AccrualAmount `json:"seller_price"`
}

type AccrualByDayProductDelivery struct {
	Services     []AccrualByDayFee `json:"services"`
	TotalAccrued AccrualAmount     `json:"total_accrued"`
}

// 获取某日应计项目。
// 注意分页：把上一次响应中的 last_id 传回，同时 date 保持同一日期。
func (c Finance) AccrualByDay(ctx context.Context, params *AccrualByDayParams) (*AccrualByDayResponse, error) {
	url := "/v1/finance/accrual/by-day"

	resp := &AccrualByDayResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
