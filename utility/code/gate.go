package code

const (
	GateSearchRegServiceFail     = "gate.search_regservice_fail"
	GateSearchRegServiceLenError = "gate.search_regservice_len_error"
	GateSearchServiceFail        = "gate.search_service_fail"
	GateSearchMethodFail         = "gate.search_method_fail"
	GatePayloadParamsError       = "gate.payload_params_error"
	GateRpcTimeout               = "gate.rpc_timeout.%s|%s|%s"
	GateLimiterError             = "gate.limiter_error"
)

var gateMap = map[string]int{
	GateSearchRegServiceFail:     1,
	GateSearchRegServiceLenError: 2,
	GateSearchServiceFail:        3,
	GateSearchMethodFail:         4,
	GatePayloadParamsError:       5,
	GateRpcTimeout:               6,
	GateLimiterError:             7,
}
