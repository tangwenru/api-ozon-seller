package ozon

//主要解决 ogl 仓库的 fbs
//根据文档，此接口用于创建首次商品发货，所有轨迹为 awaiting_deliver 的订单将自动加入该发货中【3†L3-L5】。departure_date 默认为当天【3†L12】
//func (c FBS) CarriageCreate(
//  ctx context.Context,
//  params *FbsCarriageCreateParams,
//) (*FbsCarriageCreateResponse, error) {
//  url := "/v1/carriage/create"
//
//  resp := &FbsCarriageCreateResponse{}
//
//  response, err := c.client.Request(ctx, http.MethodPost, url, params, resp, nil)
//  if err != nil {
//    return nil, err
//  }
//  response.CopyCommonResponse(&resp.CommonResponse)
//
//  return resp, nil
//}
