package internal

import (
	"fmt"
	"net/http"

	"github.com/TencentBlueKing/bk-apigateway-sdks/bkapi-client-core/define"
)

// BkApiResponseDetail implements the define.BkApiResponseDetail interface.
type BkApiResponseDetail struct {
	requestId    string
	errorCode    string
	errorMessage string
}

// Cause always return define.ErrBkApiRequest
func (detail *BkApiResponseDetail) Cause() error {
	return define.ErrBkApiRequest
}

// Error renders the error message.
func (detail *BkApiResponseDetail) Error() string {
	return fmt.Sprintf(
		"requestId: %s, errorCode: %s, errorMessage: %s",
		detail.requestId, detail.errorCode, detail.errorMessage,
	)
}

// RequestId returns the request id.
func (detail *BkApiResponseDetail) RequestId() string {
	return detail.requestId
}

// ErrorCode returns the error code.
func (detail *BkApiResponseDetail) ErrorCode() string {
	return detail.errorCode
}

// ErrorMessage returns the error message.
func (detail *BkApiResponseDetail) ErrorMessage() string {
	return detail.errorMessage
}

// GetError returns the error when errorCode is not empty.
func (detail *BkApiResponseDetail) GetError() error {
	if detail.errorCode == "" {
		return nil
	}

	return detail
}

// Map will return a map[string]interface{} that contains the non-empty details.
func (detail *BkApiResponseDetail) Map() map[string]interface{} {
	result := make(map[string]interface{}, 3)

	if detail.requestId != "" {
		result["bkapi_request_id"] = detail.requestId
	}
	if detail.errorCode != "" {
		result["bkapi_error_code"] = detail.errorCode
	}
	if detail.errorMessage != "" {
		result["bkapi_error_message"] = detail.errorMessage
	}

	return result

}

// NewBkApiResponseDetail creates a new BkApiResponseDetail.
func NewBkApiResponseDetail(requestId string, errorCode string, errorMessage string) *BkApiResponseDetail {
	return &BkApiResponseDetail{
		requestId:    requestId,
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}

// NewBkApiResponseDetailFromResponse creates a new BkApiResponseDetail from the response.
func NewBkApiResponseDetailFromResponse(response *http.Response) *BkApiResponseDetail {
	header := response.Header

	return NewBkApiResponseDetail(
		header.Get("X-Bkapi-Request-Id"),
		header.Get("X-Bkapi-Error-Code"),
		header.Get("X-Bkapi-Error-Message"),
	)
}
