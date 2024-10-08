package semrush

import (
	"fmt"
	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	"net/http"
	"net/url"
)

const (
	apiName          string = "Semrush"
	apiUrlAnalytics  string = "https://api.semrush.com/analytics/v1"
	apiUrlManagement string = "https://api.semrush.com/management/v1"
	apiUrlReports    string = "https://api.semrush.com/reports/v1"
)

type Service struct {
	apiKey      string
	httpService *go_http.Service
}

type ServiceConfig struct {
	ApiKey string
}

func NewService(serviceConfig *ServiceConfig) (*Service, *errortools.Error) {
	if serviceConfig == nil {
		return nil, errortools.ErrorMessage("ServiceConfig must not be a nil pointer")
	}

	if serviceConfig.ApiKey == "" {
		return nil, errortools.ErrorMessage("Service ApiKey not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiKey:      serviceConfig.ApiKey,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add authentication header
	header := http.Header{}
	header.Set("Authorization", service.apiKey)
	(*requestConfig).NonDefaultHeaders = &header

	if requestConfig.Parameters == nil {
		requestConfig.Parameters = &url.Values{}
	}
	requestConfig.Parameters.Set("key", service.apiKey)

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if errorResponse.Message != "" {
		e.SetMessage(errorResponse.Message)
	}

	return request, response, e
}

func (service *Service) urlAnalytics(path string) string {
	return fmt.Sprintf("%s/%s", apiUrlAnalytics, path)
}

func (service *Service) urlManagement(path string) string {
	return fmt.Sprintf("%s/%s", apiUrlManagement, path)
}

func (service *Service) urlReports(path string) string {
	return fmt.Sprintf("%s/%s", apiUrlReports, path)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.apiKey
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}
