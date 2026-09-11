package ozon

import (
	"net/http"
	"time"

	ozonCore "github.com/tangwenru/api-ozon-seller"
)

const (
	DefaultAPIBaseUrl = "https://api-seller.ozon.ru"
)

type ClientOptions struct {
	client ozonCore.HttpClient

	baseUri string

	apiKey   string
	clientId string
}

type Client struct {
	client *ozonCore.Client

	analytics     *Analytics
	fbo           *FBO
	fbs           *FBS
	finance       *Finance
	products      *Products
	promotions    *Promotions
	rating        *Rating
	warehouses    *Warehouses
	returns       *Returns
	reports       *Reports
	cancellations *Cancellations
	categories    *Categories
	polygons      *Polygons
	invoices      *Invoices
	brands        *Brands
	chats         *Chats
	certificates  *Certificates
	strategies    *Strategies
	barcodes      *Barcodes
	passes        *Passes
	clusters      *Clusters
	quants        *Quants
	reviews       *Reviews
}

func (c Client) Analytics() *Analytics {
	return c.analytics
}

func (c Client) FBO() *FBO {
	return c.fbo
}

func (c Client) FBS() *FBS {
	return c.fbs
}

func (c Client) Finance() *Finance {
	return c.finance
}

func (c Client) Products() *Products {
	return c.products
}

func (c Client) Promotions() *Promotions {
	return c.promotions
}

func (c Client) Rating() *Rating {
	return c.rating
}

func (c Client) Warehouses() *Warehouses {
	return c.warehouses
}

func (c Client) Returns() *Returns {
	return c.returns
}

func (c Client) Reports() *Reports {
	return c.reports
}

func (c Client) Cancellations() *Cancellations {
	return c.cancellations
}

func (c Client) Categories() *Categories {
	return c.categories
}

func (c Client) Polygons() *Polygons {
	return c.polygons
}

func (c Client) Invoices() *Invoices {
	return c.invoices
}

func (c Client) Brands() *Brands {
	return c.brands
}

func (c Client) Chats() *Chats {
	return c.chats
}

func (c Client) Certificates() *Certificates {
	return c.certificates
}

func (c Client) Strategies() *Strategies {
	return c.strategies
}

func (c Client) Barcodes() *Barcodes {
	return c.barcodes
}

func (c Client) Passes() *Passes {
	return c.passes
}

func (c Client) Clusters() *Clusters {
	return c.clusters
}

func (c Client) Quants() *Quants {
	return c.quants
}

func (c Client) Reviews() *Reviews {
	return c.reviews
}

type ClientOption func(c *ClientOptions)

func WithHttpClient(httpClient ozonCore.HttpClient) ClientOption {
	return func(c *ClientOptions) {
		c.client = httpClient
	}
}

func WithURI(uri string) ClientOption {
	return func(c *ClientOptions) {
		c.baseUri = uri
	}
}

func WithClientId(clientId string) ClientOption {
	return func(c *ClientOptions) {
		c.clientId = clientId
	}
}

func WithAPIKey(apiKey string) ClientOption {
	return func(c *ClientOptions) {
		c.apiKey = apiKey
	}
}

// defaultHttpClient 带超时的默认 HTTP 客户端。
// 修复：http.DefaultClient 无超时，Ozon 请求挂起会永久阻塞调用方。
var defaultHttpClient = &http.Client{
	Timeout: 120 * time.Second,
}

func NewClient(opts ...ClientOption) *Client {
	// default values
	options := &ClientOptions{
		client:  defaultHttpClient,
		baseUri: DefaultAPIBaseUrl,
	}

	for _, opt := range opts {
		opt(options)
	}

	coreClient := ozonCore.NewClient(options.client, options.baseUri, map[string]string{
		"Client-Id": options.clientId,
		"Api-Key":   options.apiKey,
	})

	return &Client{
		client:        coreClient,
		analytics:     &Analytics{client: coreClient},
		fbo:           &FBO{client: coreClient},
		fbs:           &FBS{client: coreClient},
		finance:       &Finance{client: coreClient},
		products:      &Products{client: coreClient},
		promotions:    &Promotions{client: coreClient},
		rating:        &Rating{client: coreClient},
		warehouses:    &Warehouses{client: coreClient},
		returns:       &Returns{client: coreClient},
		reports:       &Reports{client: coreClient},
		cancellations: &Cancellations{client: coreClient},
		categories:    &Categories{client: coreClient},
		polygons:      &Polygons{client: coreClient},
		invoices:      &Invoices{client: coreClient},
		brands:        &Brands{client: coreClient},
		chats:         &Chats{client: coreClient},
		certificates:  &Certificates{client: coreClient},
		strategies:    &Strategies{client: coreClient},
		barcodes:      &Barcodes{client: coreClient},
		passes:        &Passes{client: coreClient},
		clusters:      &Clusters{client: coreClient},
		quants:        &Quants{client: coreClient},
		reviews:       &Reviews{client: coreClient},
	}
}

func NewMockClient(handler http.HandlerFunc) *Client {
	coreClient := ozonCore.NewMockClient(handler)

	return &Client{
		client:        coreClient,
		analytics:     &Analytics{client: coreClient},
		fbo:           &FBO{client: coreClient},
		fbs:           &FBS{client: coreClient},
		finance:       &Finance{client: coreClient},
		products:      &Products{client: coreClient},
		promotions:    &Promotions{client: coreClient},
		rating:        &Rating{client: coreClient},
		warehouses:    &Warehouses{client: coreClient},
		returns:       &Returns{client: coreClient},
		reports:       &Reports{client: coreClient},
		cancellations: &Cancellations{client: coreClient},
		categories:    &Categories{client: coreClient},
		polygons:      &Polygons{client: coreClient},
		invoices:      &Invoices{client: coreClient},
		brands:        &Brands{client: coreClient},
		chats:         &Chats{client: coreClient},
		certificates:  &Certificates{client: coreClient},
		strategies:    &Strategies{client: coreClient},
		barcodes:      &Barcodes{client: coreClient},
		passes:        &Passes{client: coreClient},
		clusters:      &Clusters{client: coreClient},
		quants:        &Quants{client: coreClient},
		reviews:       &Reviews{client: coreClient},
	}
}
