package ozon

import (
	"time"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

// ===== /v1/carriage/create =====

type FbsCarriageCreateParams struct {
	// 表示需要创建包含可追溯商品的发运, 一般为 false
	AllBLRTraceable  bool   `json:"all_blr_traceable,omitempty"`
	DeliveryMethodId int64  `json:"delivery_method_id"`
	DepartureDate    string `json:"departure_date,omitempty"` // 格式: 2006-01-02T15:04:05Z
}

type FbsCarriageCreateResponse struct {
	ozonCore.CommonResponse

	// 发运标识符（OpenAPI：顶层字段，无 result 包裹）
	CarriageId int64 `json:"carriage_id"`
}

// ===== /v1/carriage/approve =====

type FbsCarriageApproveParams struct {
	CarriageId int64 `json:"carriage_id"`
	// 货位数量（大包数）。已开通信任验收并按货位发运时必传，否则跳过。
	ContainersCount int32 `json:"containers_count"`
}

type FbsCarriageApproveResponse struct {
	ozonCore.CommonResponse
}

// ===== /v1/carriage/get =====

type FbsCarriageGetParams struct {
	CarriageId int64 `json:"carriage_id"`
}

type FbsCarriageGetResponse struct {
	ozonCore.CommonResponse

	Result FbsCarriageGetResult `json:"result"`
}

type FbsCarriageGetResult struct {
	ActType                    string                              `json:"act_type"` // 交接单类型（FBS 卖家）
	ArrivalPassIds             []string                            `json:"arrival_pass_ids"`
	AvailableActions           []string                            `json:"available_actions"`
	CancelAvailability         FbsCarriageCancelAvailabilityResult `json:"cancel_availability"`
	CarriageId                 int64                               `json:"carriage_id"`
	CompanyId                  int64                               `json:"company_id"`
	ContainersCount            int32                               `json:"containers_count"` // 货位数量（大包数）
	CreatedAt                  time.Time                           `json:"created_at"`
	DeliveryMethodId           int64                               `json:"delivery_method_id"`
	DepartureDate              string                              `json:"departure_date"` // 发运日期
	FirstMileType              string                              `json:"first_mile_type"`
	HasPostingsForNextCarriage bool                                `json:"has_postings_for_next_carriage"`
	IntegrationType            string                              `json:"integration_type"`
	IsContainerLabelPrinted    bool                                `json:"is_container_label_printed"` // 大包（货位）标签是否已打印
	IsPartial                  bool                                `json:"is_partial"`
	PartialNum                 int64                               `json:"partial_num"`
	RetryCount                 int32                               `json:"retry_count"`
	Status                     string                              `json:"status"` // new / formed / received / closed / sent / cancelled
	TPLProviderId              int64                               `json:"tpl_provider_id"`
	UpdatedAt                  time.Time                           `json:"updated_at"`
	WarehouseId                int64                               `json:"warehouse_id"`
}

type FbsCarriageCancelAvailabilityResult struct {
	IsCancelAvailable bool   `json:"is_cancel_available"`
	Reason            string `json:"reason"`
}

// ===== /v1/carriage/cancel =====

type FbsCarriageCancelParams struct {
	CarriageId int64 `json:"carriage_id"`
}

type FbsCarriageCancelResponse struct {
	ozonCore.CommonResponse

	Error          string `json:"error"`
	CarriageStatus string `json:"carriage_status"`
}

// ===== /v1/carriage/set-postings =====

type FbsCarriageSetPostingsParams struct {
	CarriageId     int64    `json:"carriage_id"`
	PostingNumbers []string `json:"posting_numbers"` // 最新货件列表（全量覆盖）
}

type FbsCarriageSetPostingsResponse struct {
	ozonCore.CommonResponse

	Result []FbsCarriageSetPostingsResult `json:"result"`
}

type FbsCarriageSetPostingsResult struct {
	Error         string `json:"error"`
	PostingNumber string `json:"posting_number"`
	Result        bool   `json:"result"`
}

// ===== /v2/posting/fbs/act/get-postings =====

type FbsActGetPostingsParams struct {
	Id int64 `json:"id"` // 单据标识符（= /v1/carriage/create 返回的 carriage_id）
}

type FbsActGetPostingsResponse struct {
	ozonCore.CommonResponse

	Result []FbsActGetPostingsResult `json:"result"`
}

type FbsActGetPostingsResult struct {
	Id            int64                            `json:"id"` // 单据标识符
	MultiBoxQty   int32                            `json:"multi_box_qty"`
	PostingNumber string                           `json:"posting_number"`
	Status        string                           `json:"status"`
	SellerError   string                           `json:"seller_error"`
	UpdatedAt     time.Time                        `json:"updated_at"`
	CreatedAt     time.Time                        `json:"created_at"`
	Products      []FbsActGetPostingsResultProduct `json:"products"`
}

type FbsActGetPostingsResultProduct struct {
	Name     string `json:"name"`
	OfferId  string `json:"offer_id"`
	Price    string `json:"price"`
	Quantity int64  `json:"quantity"`
	SKU      int64  `json:"sku"`
}

// ===== /v2/carriage/delivery/list（v1 /v1/posting/carriage-available/list 已于 2026-03-20 关闭） =====

type CarriageDeliveryListV2Params struct {
	Cursor string                       `json:"cursor,omitempty"`
	Filter CarriageDeliveryListV2Filter `json:"filter,omitempty"`
	Limit  int64                        `json:"limit"`
}

type CarriageDeliveryListV2Filter struct {
	// 配送方式 ID。FBS 仓库用 warehouse_id 值；缺省 = 全部配送方式。
	DeliveryMethodId int64 `json:"delivery_method_id,omitempty"`
	// 发运日期，格式 YYYY-MM-DD；缺省 = 当天。
	DepartureDate string `json:"departure_date,omitempty"`
}

type CarriageDeliveryListV2Response struct {
	ozonCore.CommonResponse

	Cursor  string                                 `json:"cursor"`
	HasNext bool                                   `json:"has_next"`
	Methods []CarriageDeliveryListV2DeliveryMethod `json:"methods"`
}

type CarriageDeliveryListV2DeliveryMethod struct {
	CarriagePostingsCount    int32                      `json:"carriage_postings_count"` // 所有发运中的货件总数
	Carriages                []CarriageDeliveryCarriage `json:"carriages"`
	CutIn                    string                     `json:"cut_in"`
	CutoffAt                 string                     `json:"cutoff_at"`
	DeliveryMethodId         int64                      `json:"delivery_method_id"`
	DeliveryMethodName       string                     `json:"delivery_method_name"`
	DeliveryMethodStatus     string                     `json:"delivery_method_status"`
	DepartureDate            string                     `json:"departure_date"`
	DropoffAddress           string                     `json:"dropoff_address"`
	FirstMileType            string                     `json:"first_mile_type"`
	HasEntrustedAcceptance   bool                       `json:"has_entrusted_acceptance"` // 信任验收
	IntegrationType          string                     `json:"integration_type"`
	IsOptionalCarriage       bool                       `json:"is_optional_carriage"`
	IsPresort                bool                       `json:"is_presort"`
	IsRfbs                   bool                       `json:"is_rfbs"`
	MandatoryPackagedCount   int32                      `json:"mandatory_packaged_count"` // 已打包（须发运）件数
	MandatoryPostingsCount   int32                      `json:"mandatory_postings_count"` // 须发运件数
	OptionalPackagedCount    int32                      `json:"optional_packaged_count"`
	RecommendedTimeLocal     string                     `json:"recommended_time_local"`
	RecommendedTimeUTCOffset int32                      `json:"recommended_time_utc_offset_in_minutes"`
	TimeslotFrom             string                     `json:"timeslot_from"`
	TimeslotTo               string                     `json:"timeslot_to"`
	TPLProviderIconURL       string                     `json:"tpl_provider_icon_url"`
	TPLProviderName          string                     `json:"tpl_provider_name"`
	WarehouseCity            string                     `json:"warehouse_city"`
	WarehouseId              int64                      `json:"warehouse_id"`
	WarehouseName            string                     `json:"warehouse_name"`
}

type CarriageDeliveryCarriage struct {
	AllBLRTraceable  bool              `json:"all_blr_traceable"`
	AvailableActions []string          `json:"available_actions"`
	CarriageVolume   float64           `json:"carriage_volume"`
	Id               int64             `json:"id"` // 发运 ID；0 = 可创建（尚未创建）
	PickupFee        CarriagePickupFee `json:"pickup_fee"`
	PostingsCount    int32             `json:"postings_count"`
	QuantumCount     int32             `json:"quantum_count"`
	Status           string            `json:"status"` // new / formed / received / closed / sent / cancelled
}

type CarriagePickupFee struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

// ===== /v1/carriage/pass/create =====

type FbsCarriagePassCreateParams struct {
	CarriageId    int64                              `json:"carriage_id"`
	ArrivalPasses []FbsCarriagePassCreateArrivalPass `json:"arrival_passes"`
}

type FbsCarriagePassCreateArrivalPass struct {
	ArrivalTime         string `json:"arrival_time,omitempty"` // 到达时间 HH:MM（仓库时区）
	DriverName          string `json:"driver_name"`
	DriverPhone         string `json:"driver_phone"`
	VehicleLicensePlate string `json:"vehicle_license_plate"`
	VehicleModel        string `json:"vehicle_model"`
	WithReturns         bool   `json:"with_returns,omitempty"`
}

type FbsCarriagePassCreateResponse struct {
	ozonCore.CommonResponse

	// 实际返回数字数组（文档标注 string，兼容数值）
	ArrivalPassIds []int64 `json:"arrival_pass_ids"`
}
