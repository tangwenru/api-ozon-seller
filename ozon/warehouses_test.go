package ozon

import (
	"context"
	"net/http"
	"testing"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

func TestGetListOfWarehouses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode int
		headers    map[string]string
		response   string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			`{
				"cursor": "",
				"warehouses": [
					{
						"warehouse_id": 1020000177886000,
						"name": "This is a test",
						"is_rfbs": false,
						"has_entrusted_acceptance": false,
						"first_mile": {
							"type": "DROP_OFF",
							"dropoff_point_id": "",
							"timeslot_from": "09:00",
							"timeslot_id": 0,
							"timeslot_to": "18:00",
							"first_mile_is_changing": false
						},
						"is_kgt": false,
						"has_postings_limit": false,
						"postings_limit": -1,
						"working_days": [
							"MONDAY",
							"TUESDAY",
							"WEDNESDAY",
							"THURSDAY",
							"FRIDAY",
							"SATURDAY",
							"SUNDAY"
						],
						"min_postings_limit": 10,
						"status": "disabled"
					}
				],
				"has_next": false
			}`,
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Warehouses().GetListOfWarehouses(ctx)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &GetListOfWarehousesResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode == http.StatusOK {
			if len(resp.Warehouses) > 0 {
				if resp.Warehouses[0].WarehouseId == 0 {
					t.Errorf("Warehouse id cannot be 0")
				}
				if resp.Warehouses[0].Name == "" {
					t.Errorf("Name cannot be empty")
				}
			}
		}
	}
}

func TestGetListOfDeliveryMethods(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode int
		headers    map[string]string
		params     *GetListOfDeliveryMethodsParams
		response   string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&GetListOfDeliveryMethodsParams{
				Filter: &GetListOfDeliveryMethodsFilter{
					WarehouseIds: []int64{15588127982000},
				},
				Limit:   100,
				SortDir: Ascending,
			},
			`{
				"cursor": "",
				"delivery_methods": [
				  {
					"id": 15588127982000,
					"name": "Ozon Логистика курьеру, Есипово",
					"status": "ACTIVE",
					"cutoff": "13:00",
					"provider_id": 24,
					"template_id": 0,
					"warehouse_id": 15588127982000,
					"created_at": "2019-04-04T15:22:31.048202Z",
					"updated_at": "2021-08-15T10:21:44.854209Z",
					"sla_cut_in": 1440
				  }
				],
				"has_next": false
			}`,
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			&GetListOfDeliveryMethodsParams{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Warehouses().GetListOfDeliveryMethods(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &GetListOfDeliveryMethodsResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}

		if resp.StatusCode == http.StatusOK {
			if len(resp.DeliveryMethods) > 0 {
				if resp.DeliveryMethods[0].Id == 0 {
					t.Errorf("Id cannot be 0")
				}
				if resp.DeliveryMethods[0].Name == "" {
					t.Errorf("Name cannot be empty")
				}
				if resp.DeliveryMethods[0].Status == "" {
					t.Errorf("Status cannot be empty")
				}
				if resp.DeliveryMethods[0].WarehouseId == 0 {
					t.Errorf("Warehouse id cannot be 0")
				}
			}
		}
	}
}

func TestListForShipping(t *testing.T) {
	t.Parallel()

	tests := []struct {
		statusCode int
		headers    map[string]string
		params     *ListForShippingParams
		response   string
	}{
		// Test Ok
		{
			http.StatusOK,
			map[string]string{"Client-Id": "my-client-id", "Api-Key": "my-api-key"},
			&ListForShippingParams{
				FilterBySupplyType: []string{"CREATE_TYPE_UNKNOWN"},
				Search:             "string",
			},
			`{
				"search": [
				  {
					"address": "string",
					"coordinates": {
					  "latitude": 0,
					  "longitude": 0
					},
					"name": "string",
					"warehouse_id": 0,
					"warehouse_type": "WAREHOUSE_TYPE_UNKNOWN"
				  }
				]
			}`,
		},
		// Test No Client-Id or Api-Key
		{
			http.StatusUnauthorized,
			map[string]string{},
			&ListForShippingParams{},
			`{
				"code": 16,
				"message": "Client-Id and Api-Key headers are required"
			}`,
		},
	}

	for _, test := range tests {
		c := NewMockClient(ozonCore.NewMockHttpHandler(test.statusCode, test.response, test.headers))

		ctx, _ := context.WithTimeout(context.Background(), testTimeout)
		resp, err := c.Warehouses().ListForShipping(ctx, test.params)
		if err != nil {
			t.Error(err)
			continue
		}

		compareJsonResponse(t, test.response, &ListForShippingResponse{})

		if resp.StatusCode != test.statusCode {
			t.Errorf("got wrong status code: got: %d, expected: %d", resp.StatusCode, test.statusCode)
		}
	}
}
