package tests

import "common.queueb.org/tests/http"

// Shortcuts for http sub-module.
var (
	NewResponse       = http.NewResponse
	NewResponseString = http.NewResponseString
	NewStackedRouter  = http.NewStackedRouter
	NewHTTPServer     = http.NewHTTPServer

	// options

	ROptAccepted           = http.ROptAccepted
	ROptBadGateway         = http.ROptBadGateway
	ROptBadRequest         = http.ROptBadRequest
	ROptCreated            = http.ROptCreated
	ROptForbidden          = http.ROptForbidden
	ROptGatewayTimeout     = http.ROptGatewayTimeout
	ROptNoContent          = http.ROptNoContent
	ROptNotFound           = http.ROptNotFound
	ROptNotImplemented     = http.ROptNotImplemented
	ROptServerError        = http.ROptServerError
	ROptServiceUnavailable = http.ROptServiceUnavailable
	ROptStatusOK           = http.ROptStatusOK
	ROptUnauthorized       = http.ROptUnauthorized
	ROptUnknownError       = http.ROptUnknownError
	ROptWithHeaders        = http.ROptWithHeaders
	ROptWithOK             = http.ROptWithOK
)
