package code

const (
	CommonConsulCfgError        = "common.gate_consul_access_get_fail"
	CommonConsulSrvCurlAllError = "common.gate_consul_srv_curl_all_fail"
	CommonRequiredError         = "common.params:%s is required"
	CommonConsulCfgCurlAllError = "common.gate_consul_cfg_curl_all_fail"
)

var commonMap = map[string]int{
	CommonConsulCfgError:        1,
	CommonConsulSrvCurlAllError: 2,
	CommonRequiredError:         3,
	CommonConsulCfgCurlAllError: 4,
}
