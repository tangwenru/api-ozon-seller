package ozon

import (
	"context"
	"net/http"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

type Categories struct {
	client *ozonCore.Client
}

type GetProductTreeParams struct {
	// Response language
	Language Language `json:"language,omitempty"`
}

type GetProductTreeResponse struct {
	ozonCore.CommonResponse

	// Categories list
	Result []GetProductTreeResult `json:"result"`
}

type GetProductTreeResult struct {
	// Category identifier
	DescriptionCategoryId int64 `json:"description_category_id"`

	// Category name
	CategoryName string `json:"category_name"`

	// `true`, if you can't create products in the category. `false`, if you can
	Disabled bool `json:"disabled"`

	// Product type identifier
	TypeId int64 `json:"type_id"`

	// Product type name
	TypeName string `json:"type_name"`

	// Subcategory tree
	Children []GetProductTreeResult `json:"children"`
}

// Returns product categories in the tree view.
//
// New products can be created in the last level categories only.
// This means that you need to match these particular categories with the categories of your site.
// We don't create new categories by user request.
func (c *Categories) Tree(ctx context.Context, params *GetProductTreeParams) (*GetProductTreeResponse, error) {
	url := "/v1/description-category/tree"

	resp := &GetProductTreeResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetCategoryAttributesParams struct {
	// Category identifier
	DescriptionCategoryId int64 `json:"description_category_id"`

	// Response language
	Language Language `json:"language,omitempty"`

	// Product type identifier
	TypeId int64 `json:"type_id"`
}

type GetCategoryAttributesResponse struct {
	ozonCore.CommonResponse

	// Method result
	Result []GetCategoryAttributesResult `json:"result"`
}

type GetCategoryAttributesResult struct {
	// Indication that the dictionary attribute values depend on the category:
	//
	// 	true — the attribute has its own set of values for each category.
	// 	false — the attribute has the same set of values for all categories
	CategoryDependent bool `json:"category_dependent"`

	// Characteristic description
	Description string `json:"description"`

	// Directory identifier
	DictionaryId int64 `json:"dictionary_id"`

	// Characteristics group identifier
	GroupId int64 `json:"group_id"`

	// Characteristics group name
	GroupName string `json:"group_name"`

	// Number of document generation task
	Id int64 `json:"id"`

	// Indicates that the attribute is aspect. An aspect attribute is a characteristic that distinguishes products of the same model.
	//
	// For example, clothes or shoes of the same model may have different colors and sizes. That is, color and size are aspect attributes.
	//
	// Values description:
	//
	// 	- `true`—the attribute is aspect and can't be changed after the products are delivered to the warehouse or sold from the seller's warehouse.
	// 	- `false`—the attribute is not aspect and can be changed at any time
	IsAspect bool `json:"is_aspect"`

	// Indicates that the characteristic is a set of values:
	//
	// 	- `true`—the characteristic is a set of values,
	// 	- `false`—the characteristic consists of a single value
	IsCollection bool `json:"is_collection"`

	// Indicates that the characteristic is mandatory:
	//
	// 	- `true`—a mandatory characteristic,
	// 	- `false`—an optional characteristic
	IsRequired bool `json:"is_required"`

	// Name
	Name string `json:"name"`

	// Characteristic type
	Type string `json:"type"`

	// Complex attribute identifier
	AttributeComplexId int64 `json:"attribute_complex_id"`

	// Maximum number of values for attribute
	MaxValueCount int64 `json:"max_value_count"`

	// Indicates whether the complex attribute is a set of values
	ComplexIsCollection bool `json:"complex_is_collection"`
}

// Getting characteristics for specified product category and type.
//
// If the dictionary_id value is 0, there is no directory.
// If the value is different, there are directories.
// Get them using the `/v1/description-category/attribute/values` method.
func (c *Categories) Attributes(ctx context.Context, params *GetCategoryAttributesParams) (*GetCategoryAttributesResponse, error) {
	url := "/v1/description-category/attribute"

	resp := &GetCategoryAttributesResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetAttributeDictionaryParams struct {
	// Characteristics identifier
	AttributeId int64 `json:"attribute_id"`

	// Category identifier
	DescriptionCategoryId int64 `json:"description_category_id"`

	// Response language
	Language Language `json:"language,omitempty"`

	// Identifier of the directory to start the response with.
	// If `last_value_id` is 10, the response will contain directories starting from the 11th
	LastValueId int64 `json:"last_value_id"`

	// Number of values in the response:
	//
	// 	- maximum—5000,
	// 	- minimum—1.
	Limit int64 `json:"limit,omitempty"`

	// Product type identifier
	TypeId int64 `json:"type_id"`
}

type GetAttributeDictionaryResponse struct {
	ozonCore.CommonResponse

	// Indication that only part of characteristic values was returned in the response:
	//
	// 	- true—make a request with a new last_value_id parameter value for getting the rest of characteristic values;
	// 	- false—all characteristic values were returned
	HasNext bool `json:"has_next"`

	// Characteristic values
	Result []GetAttributeDictionaryResult `json:"result"`
}

type GetAttributeDictionaryResult struct {
	// Characteristic value identifier
	Id int64 `json:"id"`

	// Additional description
	Info string `json:"info"`

	// Image link
	Picture string `json:"picture"`

	// Product characteristic value
	Value string `json:"value"`
}

// Returns characteristics value directory.
//
// To check if an attribute has a nested directory,
// use the `/v1/description-category/attribute` method.
func (c *Categories) AttributesDictionary(ctx context.Context, params *GetAttributeDictionaryParams) (*GetAttributeDictionaryResponse, error) {
	url := "/v1/description-category/attribute/values"

	resp := &GetAttributeDictionaryResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type SearchAttributeDictionaryParams struct {
	// Characteristics identifier
	AttributeId int64 `json:"attribute_id"`

	// Category identifier
	DescriptionCategoryId int64 `json:"description_category_id"`

	// The value to be searched for
	//  - minimum—2 characters
	Value string `json:"value"`

	// Number of values in the response:
	//
	// 	- maximum—100,
	// 	- minimum—1.
	Limit int64 `json:"limit,omitempty"`

	// Product type identifier
	TypeId int64 `json:"type_id"`
}

type SearchAttributeDictionaryResponse struct {
	ozonCore.CommonResponse

	// Characteristic values
	Result []SearchAttributeDictionaryResult `json:"result"`
}

type SearchAttributeDictionaryResult struct {
	// Characteristic value identifier
	Id int64 `json:"id"`

	// Additional description
	Info string `json:"info"`

	// Image link
	Picture string `json:"picture"`

	// Product characteristic value
	Value string `json:"value"`
}

// Returns found characteristics value directory.
//
// To check if an attribute has a nested directory,
// use the `/v1/description-category/attribute` method.
func (c *Categories) SearchAttributesDictionary(ctx context.Context, params *SearchAttributeDictionaryParams) (*SearchAttributeDictionaryResponse, error) {
	url := "/v1/description-category/attribute/values/search"

	resp := &SearchAttributeDictionaryResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

// 依赖特征（Beta，2026-09-24 加入）：
// 某些特征只有先选了父特征的值才有意义，例如 品牌 → 车型。
// 1）DependentAttributes 返回「父特征 → 子特征」的 id 对；
// 2）DependentAttributeValues 返回某个父子特征对的「父值 → 允许的子值」映射。

type GetDependentAttributesParams struct {
	// Category identifier
	DescriptionCategoryId int64 `json:"description_category_id"`

	// Product type identifier
	TypeId int64 `json:"type_id,omitempty"`
}

type GetDependentAttributesResponse struct {
	ozonCore.CommonResponse

	// Parent attribute → child attribute pairs
	Result []GetDependentAttributesResult `json:"result"`
}

type GetDependentAttributesResult struct {
	// Child attribute identifier
	ChildAttributeId int64 `json:"child_attribute_id"`

	// Parent attribute identifier
	// parent_attribute_id = 8229 means the parent value is the product type id (type_id)
	ParentAttributeId int64 `json:"parent_attribute_id"`
}

// Returns pairs of parent and child attributes. Use DependentAttributeValues
// to get the child values allowed for each parent value.
func (c *Categories) DependentAttributes(ctx context.Context, params *GetDependentAttributesParams) (*GetDependentAttributesResponse, error) {
	url := "/v1/description-category/dependent-attributes"

	resp := &GetDependentAttributesResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}

type GetDependentAttributeValuesParams struct {
	// Child attribute identifier
	ChildAttributeId int64 `json:"child_attribute_id"`

	// Cursor for the next page; empty string starts from the beginning
	Cursor string `json:"cursor,omitempty"`

	// Category identifier
	DescriptionCategoryId int64 `json:"description_category_id,omitempty"`

	// Number of child values in the response, 1..1000 (default 100).
	// OZON counts children, not parent rows; a single parent can span pages.
	Limit int64 `json:"limit,omitempty"`

	// Parent attribute identifier
	ParentAttributeId int64 `json:"parent_attribute_id"`

	// Product type identifier
	TypeId int64 `json:"type_id,omitempty"`
}

type GetDependentAttributeValuesResponse struct {
	ozonCore.CommonResponse

	// Cursor for the next page; empty when there are no more values
	Cursor string `json:"cursor"`

	// Parent value → allowed child values
	Result []GetDependentAttributeValuesResult `json:"result"`
}

type GetDependentAttributeValuesResult struct {
	// Allowed child values for this parent value
	Children []GetDependentAttributeValueChild `json:"children"`

	// Parent attribute value
	ParentValue string `json:"parent_value"`

	// Parent attribute value identifier
	ParentValueId int64 `json:"parent_value_id"`
}

type GetDependentAttributeValueChild struct {
	// Child attribute value identifier
	ChildValueId int64 `json:"child_value_id"`

	// Child attribute value
	ChildValue string `json:"child_value"`
}

// Returns child values allowed for each parent value of a dependent attribute pair.
// Page through the result with Cursor until it is empty.
func (c *Categories) DependentAttributeValues(ctx context.Context, params *GetDependentAttributeValuesParams) (*GetDependentAttributeValuesResponse, error) {
	url := "/v1/description-category/dependent-attributes/values"

	resp := &GetDependentAttributeValuesResponse{}

	response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
	if err != nil {
		return nil, err
	}
	response.CopyCommonResponse(&resp.CommonResponse)

	return resp, nil
}
