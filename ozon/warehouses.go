package ozon

import (
	"context"
	"net/http"
	"time"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

type Warehouses struct {
	client *ozonCore.Client
}

type GetListOfWarehousesParams struct {
	Cursor       string  `json:"cursor"`        // 后续数据的选择标志
	Limit        int     `json:"limit"`         //<= 200
	WarehouseIds []int64 `json:"warehouse_ids"` // 仓库识别符
}

type GetListOfWarehousesResponse struct {
	ozonCore.CommonResponse

	GetListOfWarehousesResult // `json:"result"`
}

type GetListOfWarehousesResult struct {
	//// Trusted acceptance attribute. `true` if trusted acceptance is enabled in the warehouse
	//HasEntrustedAcceptance bool `json:"has_entrusted_acceptance"`
	//
	//// Indication that the warehouse works under the rFBS scheme:
	////   - true — the warehouse works under the rFBS scheme;
	////   - false — the warehouse does not work under the rFBS scheme.
	//IsRFBS bool `json:"is_rfbs"`
	//
	//// Warehouse name
	//Name string `json:"name"`
	//
	//// Warehouse identifier
	//WarehouseId int64 `json:"warehouse_id"`
	//
	//// Possibility to print an acceptance certificate in advance. `true` if printing in advance is possible
	//CanPrintActInAdvance bool `json:"can_print_act_in_advance"`
	//
	//// FBS first mile
	//FirstMileType GetListOfWarehousesResultFirstMile `json:"first_mile_type"`
	//
	//// Indication if there is a limit on the minimum number of orders. `true` if there is such a limit
	//HasPostingsLimit bool `json:"has_postings_limit"`
	//
	//// Indication that the warehouse is not working due to quarantine
	//IsKarantin bool `json:"is_karantin"`
	//
	//// Indication that the warehouse accepts bulky products
	//IsKGT bool `json:"is_kgt"`
	//
	//// true if the warehouse handles economy products
	//IsEconomy bool `json:"is_economy"`
	//
	//// Indication that warehouse schedule can be changed
	//IsTimetableEditable bool `json:"is_timetable_editable"`
	//
	//// Minimum limit value: the number of orders that can be brought in one shipment
	//MinPostingsLimit int32 `json:"min_postings_limit"`
	//
	//// Limit value. -1 if there is no limit
	//PostingsLimit int32 `json:"postings_limit"`
	//
	//// Number of warehouse working days
	//MinWorkingDays int64 `json:"min_working_days"`
	//
	//// Warehouse status
	//Status string `json:"status"`
	//
	//// Warehouse working days
	//WorkingDays []WorkingDay `json:"working_days"`

	Cursor     string `json:"cursor"`
	Warehouses []struct {
		AddressInfo struct {
			Address   string  `json:"address"`
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Utc       string  `json:"utc"`
		} `json:"address_info"`
		CarriageLabelType string    `json:"carriage_label_type"`
		CourierComment    string    `json:"courier_comment"`
		CourierPhones     []string  `json:"courier_phones"`
		CreatedAt         time.Time `json:"created_at"`
		FirstMile         struct {
			Type                string `json:"type"`
			DropoffPointId      string `json:"dropoff_point_id"`
			TimeslotFrom        string `json:"timeslot_from"`
			TimeslotId          int64  `json:"timeslot_id"`
			TimeslotTo          string `json:"timeslot_to"`
			FirstMileIsChanging bool   `json:"first_mile_is_changing"`
		} `json:"first_mile"`
		HasEntrustedAcceptance bool   `json:"has_entrusted_acceptance"`
		HasPostingsLimit       bool   `json:"has_postings_limit"`
		IsAutoAssembly         bool   `json:"is_auto_assembly"`
		IsKgt                  bool   `json:"is_kgt"`
		IsRfbs                 bool   `json:"is_rfbs"`
		IsWaybillEnabled       bool   `json:"is_waybill_enabled"`
		MinPostingsLimit       int    `json:"min_postings_limit"`
		IsComfort              bool   `json:"is_comfort"`
		IsExpress              bool   `json:"is_express"`
		WarehouseType          string `json:"warehouse_type"`
		CutInTime              int64  `json:"cut_in_time"`
		Name                   string `json:"name"`
		Phone                  string `json:"phone"`
		PostingsLimit          int    `json:"postings_limit"`
		SlaCutIn               int    `json:"sla_cut_in"`
		Status                 string `json:"status"`
		Timetable              struct {
			TimetableFrom time.Time `json:"timetable_from"`
			TimetableTo   time.Time `json:"timetable_to"`
			WorkingHours  []struct {
				TimeFrom time.Time `json:"time_from"`
				TimeTo   time.Time `json:"time_to"`
			} `json:"working_hours"`
		} `json:"timetable"`
		UpdatedAt    time.Time `json:"updated_at"`
		WarehouseId  int64     `json:"warehouse_id"`
		WithItemList bool      `json:"with_item_list"`
		WorkingDays  []string  `json:"working_days"`
	} `json:"warehouses"`
	HasNext bool `json:"has_next"`
}

type GetListOfWarehousesResultFirstMile struct {
	// DropOff point identifier
	DropoffPointId string `json:"dropoff_point_id"`

	// DropOff timeslot identifier
	DropoffTimeslotId int64 `json:"dropoff_timeslot_id"`

	// Indication that the warehouse settings are being updated
	FirstMileIsChanging bool `json:"first_mile_is_changing"`

	// First mile type:
	//
	// Enum: "DropOff" "Pickup"
	//   - DropOff
	//   - Pickup
	FirstMileType string `json:"first_mile_type"`
}

// /v1/warehouse/list 已于 2026-04-07 关闭，切换至 /v2/warehouse/list。
// Method returns the list of FBS and rFBS warehouses.
// To get the list of FBO warehouses, use the /v1/warehouse/fbo/list method.
func (c Warehouses) GetListOfWarehouses(ctx context.Context) (*GetListOfWarehousesResponse, error) {
	url := "/v2/warehouse/list"

	resp := &GetListOfWarehousesResponse{}

	// v2 要求 limit 必填（<= 200）
	query := &GetListOfWarehousesParams{
		Limit: 200,
	}

	response, err := c.client.Request(ctx, http.MethodPost, url, query, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

func (c Warehouses) GetListOfWarehousesV2(
	ctx context.Context,
	query *GetListOfWarehousesParams,
) (*GetListOfWarehousesResponse, error) {
	url := "/v2/warehouse/list"

	resp := &GetListOfWarehousesResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, query, resp, nil)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("resp 2:", resp)

	response.CopyCommonResponse(&resp.CommonResponse)
	//fmt.Printf("resp 3:", resp)

	return resp, nil
}

type GetListOfDeliveryMethodsParams struct {
	// Cursor for selecting the next batch of data
	Cursor string `json:"cursor,omitempty"`

	// Search filter for delivery methods
	Filter *GetListOfDeliveryMethodsFilter `json:"filter,omitempty"`

	// Number of items in a response. [ 1 .. 100 ]
	Limit int64 `json:"limit,omitempty"`

	// Sorting direction: ASC — ascending, DESC — descending
	SortDir Order `json:"sort_dir,omitempty"`
}

type GetListOfDeliveryMethodsFilter struct {
	// Delivery method identifiers
	DeliveryMethodIds []int64 `json:"delivery_method_ids,omitempty"`

	// Delivery service identifiers
	ProviderIds []int64 `json:"provider_ids,omitempty"`

	// Delivery method statuses:
	//   - NEW—created
	//   - EDITED—being edited
	//   - ACTIVE—active
	//   - DISABLED—inactive
	Status []string `json:"status,omitempty"`

	// Warehouse identifiers
	WarehouseIds []int64 `json:"warehouse_ids,omitempty"`
}

type GetListOfDeliveryMethodsResponse struct {
	ozonCore.CommonResponse

	// Cursor for selecting the next batch of data
	Cursor string `json:"cursor"`

	// Indication that only part of delivery methods was returned in the response:
	//   - true — make a request with the received cursor value for getting the rest of delivery methods;
	//   - false — all delivery methods were returned
	HasNext bool `json:"has_next"`

	// Delivery methods
	DeliveryMethods []GetListOfDeliveryMethodsResult `json:"delivery_methods"`
}

type GetListOfDeliveryMethodsResult struct {
	// Date and time of delivery method creation
	CreatedAt time.Time `json:"created_at"`

	// Time before an order must be packaged
	Cutoff string `json:"cutoff"`

	// Delivery method identifier
	Id int64 `json:"id"`

	// Indication that the delivery method belongs to Ozon Express
	IsExpress bool `json:"is_express"`

	// Delivery method name
	Name string `json:"name"`

	// Delivery service identifier
	ProviderId int64 `json:"provider_id"`

	// Minimum time to package an order in minutes according to warehouse settings
	SLACutIn int64 `json:"sla_cut_in"`

	// Delivery method status:
	//   - NEW—created,
	//   - EDITED—being edited,
	//   - ACTIVE—active,
	//   - DISABLED—inactive
	Status string `json:"status"`

	// Order delivery service identifier
	TemplateId int64 `json:"template_id"`

	// Drop-off point details
	TPLDropoffPoint *GetListOfDeliveryMethodsTPLDropoffPoint `json:"tpl_dropoff_point,omitempty"`

	// Type of integration with the delivery service
	TPLIntegrationType string `json:"tpl_integration_type"`

	// Date and time when the delivery method was last updated
	UpdatedAt time.Time `json:"updated_at"`

	// Warehouse identifier
	WarehouseId int64 `json:"warehouse_id"`
}

type GetListOfDeliveryMethodsTPLDropoffPoint struct {
	// Drop-off point address
	Address string `json:"address"`

	// Drop-off point coordinates
	AddressCoordinates *GetListOfDeliveryMethodsAddressCoordinates `json:"address_coordinates,omitempty"`

	// Drop-off point code
	Code string `json:"code"`

	// Drop-off point name
	Name string `json:"name"`
}

type GetListOfDeliveryMethodsAddressCoordinates struct {
	// Latitude
	Latitude float64 `json:"latitude"`

	// Longitude
	Longitude float64 `json:"longitude"`
}

// /v1/delivery-method/list 已于 2026-04-07 关闭，切换至 /v2/delivery-method/list。
// This methods allows you to get list of all delivery methods that can be applied for this warehouse
func (c Warehouses) GetListOfDeliveryMethods(ctx context.Context, params *GetListOfDeliveryMethodsParams) (*GetListOfDeliveryMethodsResponse, error) {
	url := "/v2/delivery-method/list"

	resp := &GetListOfDeliveryMethodsResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type ListForShippingParams struct {
	// Supply type
	FilterBySupplyType []string `json:"filter_by_supply_type"`

	// Search by warehouse name. To search for pick-up points, specify the full name
	Search string `json:"search"`
}

type ListForShippingResponse struct {
	ozonCore.CommonResponse

	// Warehouse search result
	Search []ListForShippingSearch `json:"search"`
}

type ListForShippingSearch struct {
	// Warehouse address
	Address string `json:"address"`

	// Warehouse coordinates
	Coordinates Coordinates `json:"coordinates"`

	// Warehouse name
	Name string `json:"name"`

	// Identifier of the warehouse, pick-up point, or sorting center
	WarehouseId int64 `json:"warehouse_id"`

	// Type of warehouse, pick-up point, or sorting center
	WarehouseType string `json:"warehouse_type"`
}

type Coordinates struct {
	// Latitude
	Latitude float64 `json:"latitude"`

	// Longitude
	Longitude float64 `json:"longitude"`
}

// Use the method to find sorting centres, pick-up points, and drop-off points available for cross-docking and direct supplies.
//
// You can view the addresses of all points on the map and in a table in the Knowledge Base.
func (c Warehouses) ListForShipping(ctx context.Context, params *ListForShippingParams) (*ListForShippingResponse, error) {
	url := "/v1/warehouse/fbo/list"

	resp := &ListForShippingResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
