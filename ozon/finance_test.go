package ozon

import (
	"context"
	"net/http"
	"testing"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

func TestReportOnSoldProducts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode   int
		headers      map[string]string
		params       *ReportOnSoldProductsParams
		response     string
		errorMessage string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&ReportOnSoldProductsParams{
				Month: 9,
				Year:  2022,
			},
			`{
				"result": {
				  "header": {
					"contract_date": "string",
					"contract_number": "string",
					"currency_sys_name": "string",
					"doc_amount": 0,
					"doc_date": "string",
					"number": "string",
					"payer_inn": "string",
					"payer_kpp": "string",
					"payer_name": "string",
					"receiver_inn": "string",
					"receiver_kpp": "string",
					"receiver_name": "string",
					"start_date": "string",
					"stop_date": "string",
					"vat_amount": 0
				  },
				  "rows": [
					{
					  "commission_ratio": 0,
					  "delivery_commission": {
						"amount": 0,
						"bonus": 0,
						"commission": 0,
						"compensation": 0,
						"price_per_instance": 0,
						"quantity": 0,
						"standard_fee": 0,
						"bank_coinvestment": 0,
						"stars": 0,
						"total": 0
					  },
					  "item": {
						"barcode": "string",
						"name": "string",
						"offer_id": "string",
						"sku": 0
					  },
					  "return_commission": {
						"amount": 0,
						"bonus": 0,
						"commission": 0,
						"compensation": 0,
						"price_per_instance": 0,
						"quantity": 0,
						"standard_fee": 0,
						"bank_coinvestment": 0,
						"stars": 0,
						"total": 0
					  },
					  "rowNumber": 0,
					  "seller_price_per_instance": 0
					}
				  ]
				}
			}`,
			"",
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			&ReportOnSoldProductsParams{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
			"Client-Id and Api-Key headers are required",
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Finance().ReportOnSoldProducts(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &ReportOnSoldProductsResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.Message != test.errorMessage {
				t.Errorf("got wrong error message: got: %s, expected: %s", resp.Message, test.errorMessage)
			}
		}
	}
}

func TestGetTotalTransactionsSum(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode   int
		headers      map[string]string
		params       *GetTotalTransactionsSumParams
		response     string
		errorMessage string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&GetTotalTransactionsSumParams{
				Date: GetTotalTransactionsSumDate{
					From: ozonCore.TimeFromString(t, "2006-01-02T15:04:05Z", "2021-11-01T00:00:00.000Z"),
					To:   ozonCore.TimeFromString(t, "2006-01-02T15:04:05Z", "2021-11-02T00:00:00.000Z"),
				},
				TransactionType: "ALL",
			},
			`{
				"result": {
				  "accruals_for_sale": 96647.58,
				  "sale_commission": -11456.65,
				  "processing_and_delivery": -24405.68,
				  "refunds_and_cancellations": -330,
				  "services_amount": -1307.57,
				  "compensation_amount": 0,
				  "money_transfer": 0,
				  "others_amount": 113.05
				}
			}`,
			"",
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			&GetTotalTransactionsSumParams{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
			"Client-Id and Api-Key headers are required",
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Finance().GetTotalTransactionsSum(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &GetTotalTransactionsSumResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.Message != test.errorMessage {
				t.Errorf("got wrong error message: got: %s, expected: %s", resp.Message, test.errorMessage)
			}
		}
	}
}

func TestListTransactions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode   int
		headers      map[string]string
		params       *ListTransactionsParams
		response     string
		errorMessage string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&ListTransactionsParams{
				Filter: ListTransactionsFilter{
					Date: ListTransactionsFilterDate{
						From: ozonCore.TimeFromString(t, "2006-01-02T15:04:05Z", "2021-11-01T00:00:00.000Z"),
						To:   ozonCore.TimeFromString(t, "2006-01-02T15:04:05Z", "2021-11-02T00:00:00.000Z"),
					},
					TransactionType: "ALL",
				},
				Page:     1,
				PageSize: 1000,
			},
			`{
				"result": {
				  "operations": [
					{
					  "operation_id": 11401182187840,
					  "operation_type": "MarketplaceMarketingActionCostOperation",
					  "operation_date": "2021-11-01 00:00:00",
					  "operation_type_name": "Услуги продвижения товаров",
					  "delivery_charge": 0,
					  "return_delivery_charge": 0,
					  "accruals_for_sale": 0,
					  "sale_commission": 0,
					  "amount": -6.46,
					  "type": "services",
					  "posting": {
						"delivery_schema": "",
						"order_date": "",
						"posting_number": "",
						"warehouse_id": 0
					  },
					  "items": [],
					  "services": []
					}
				  ],
				  "page_count": 1,
				  "row_count": 355
				}
			}`,
			"",
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			&ListTransactionsParams{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
			"Client-Id and Api-Key headers are required",
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Finance().ListTransactions(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &ListTransactionsResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.Message != test.errorMessage {
				t.Errorf("got wrong error message: got: %s, expected: %s", resp.Message, test.errorMessage)
			}
		}
	}
}

func TestMutualSettlements(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode   int
		headers      map[string]string
		params       *GetReportParams
		response     string
		errorMessage string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&GetReportParams{
				Date:     "2024-08",
				Language: "DEFAULT",
			},
			`{
				"result": {
				  "code": "string"
				}
			}`,
			"",
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			&GetReportParams{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
			"Client-Id and Api-Key headers are required",
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Finance().MutualSettlements(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &ReportResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.Message != test.errorMessage {
				t.Errorf("got wrong error message: got: %s, expected: %s", resp.Message, test.errorMessage)
			}
		}
	}
}

func TestSalesToLegalEntities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode   int
		headers      map[string]string
		params       *GetReportParams
		response     string
		errorMessage string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&GetReportParams{
				Date:     "2024-08",
				Language: "DEFAULT",
			},
			`{
				"result": {
				  "code": "string"
				}
			}`,
			"",
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			&GetReportParams{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
			"Client-Id and Api-Key headers are required",
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Finance().SalesToLegalEntities(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &ReportResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.Message != test.errorMessage {
				t.Errorf("got wrong error message: got: %s, expected: %s", resp.Message, test.errorMessage)
			}
		}
	}
}

func TestGetRealizationPosting(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode   int
		headers      map[string]string
		params       *RealizationPostingParams
		response     string
		errorMessage string
	}{
		// Test Ok：人民币销售月报，逐货件明细。
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&RealizationPostingParams{
				Month: 10,
				Year:  2026,
			},
			`{
				"header": {
					"doc_date": "2026-11-05",
					"contract_date": "2025-01-01",
					"contract_number": "39-808132/26",
					"currency_sys_name": "CNY",
					"number": "2026-10",
					"payer_inn": "7704217370",
					"payer_kpp": "997750001",
					"payer_name": "OOO Ozon",
					"receiver_inn": "91110000MA01",
					"receiver_kpp": "",
					"receiver_name": "Seller Ltd",
					"start_date": "2026-10-01",
					"stop_date": "2026-10-31"
				},
				"rows": [
					{
						"row_number": 1,
						"commission_ratio": 18,
						"seller_price_per_instance": 100,
						"item": {"barcode": "123", "name": "商品", "offer_id": "SKU-1", "sku": 111},
						"order": {"posting_number": "12345-0001-1", "created_date": "2026-10-02"},
						"delivery_commission": {
							"amount": 100,
							"bonus": 0,
							"commission": 0,
							"compensation": 0,
							"price_per_instance": 100,
							"quantity": 1,
							"standard_fee": 18,
							"bank_coinvestment": 2,
							"stars": 0,
							"pick_up_point_coinvestment": 1,
							"total": 85
						},
						"return_commission": {
							"amount": 0,
							"bonus": 0,
							"commission": 0,
							"compensation": 0,
							"price_per_instance": 0,
							"quantity": 0,
							"standard_fee": 0,
							"bank_coinvestment": 0,
							"stars": 0,
							"pick_up_point_coinvestment": 0,
							"total": 0
						},
						"legal_entity_document": {"number": "", "sale_date": ""}
					}
				]
			}`,
			"",
		},
		// Test No report：当月报告尚未生成，Ozon 返回 404。
		{
			http.StatusNotFound,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&RealizationPostingParams{Month: 11, Year: 2026},
			`{
				"code": 5,
				"message": "Report not found"
			}`,
			"Report not found",
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Finance().GetRealizationPosting(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &RealizationPostingResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode != http.StatusOK {
			if resp.Message != test.errorMessage {
				t.Errorf("got wrong error message: got: %s, expected: %s", resp.Message, test.errorMessage)
			}
		}
	}
}
