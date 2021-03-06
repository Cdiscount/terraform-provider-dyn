package api

// DSFSResponse is used for holding the data returned by a call to
// "https://api.dynect.net/REST/DSF/" with 'detail: Y'.
type AllDSFDetailedResponse struct {
	ResponseBlock
	Data []DSFService `json:"data"`
}

// DSFResponse is used for holding the data returned by a call to
// "https://api.dynect.net/REST/DSF/SERVICE_ID".
type DSFResponse struct {
	ResponseBlock
	Data DSFService `json:"data"`
}

// Type DSFService is used as a nested struct, which holds the data for a
// DSF Service returned by a call to "https://api.dynect.net/REST/DSF/SERVICE_ID".
type DSFService struct {
	ID            string       `json:"service_id"`
	Label         string       `json:"label"`
	Active        string       `json:"active"`
	TTL           SInt         `json:"ttl"`
	PendingChange string       `json:"pending_change"`
	Notifiers     []Notifier   `json:"notifiers"`
	Nodes         []DSFNode    `json:"nodes"`
	Rulesets      []DSFRuleset `json:"rulesets"`
}

type DSFServiceRequest struct {
	PublishBlock
	Label string `json:"label"`
	TTL   SInt   `json:"ttl"`
}

type DSFResponsePoolRef struct {
	ID string `json:"dsf_response_pool_id"`
}
type DSFRulesetRequest struct {
	PublishBlock
	Label        string                `json:"label"`
	CriteriaType string                `json:"criteria_type"`
	ResponsePool *[]DSFResponsePoolRef `json:"response_pools"`
}
type DSFRuleset struct {
	ID            string            `json:"dsf_ruleset_id"`
	Label         string            `json:"label"`
	CriteriaType  string            `json:"criteria_type"`
	Criteria      interface{}       `json:"criteria"`
	Ordering      string            `json:"ordering"`
	Eligible      string            `json:"eligible"`
	PendingChange string            `json:"pending_change"`
	ResponsePools []DSFResponsePool `json:"response_pools"`
}
type DSFRulesetResponse struct {
	ResponseBlock
	Data DSFRuleset `json:"data"`
}

type DSFResponsePoolResponse struct {
	ResponseBlock
	Data DSFResponsePool `json:"data"`
}
type DSFResponsePoolRequest struct {
	PublishBlock
	Label      string `json:"label"`
	Automation string `json:"automation",omit_empty`
}
type DSFResponsePool struct {
	ID            string              `json:"dsf_response_pool_id"`
	Label         string              `json:"label"`
	Automation    string              `json:"automation"`
	CoreSetCount  string              `json:"core_set_count"`
	Eligible      string              `json:"eligible"`
	PendingChange string              `json:"pending_change"`
	RsChains      []DSFRecordSetChain `json:"rs_chains"`
	Rulesets      []DSFRuleset        `json:"rulesets"`
	Status        string              `json:"status"`
	LastMonitored string              `json:"last_monitored"`
	Notifier      string              `json:"notifier"`
}

type DSFRecordSetChain struct {
	ID                string         `json:"dsf_record_set_failover_chain_id"`
	Status            string         `json:"status"`
	Core              string         `json:"core"`
	Label             string         `json:"label"`
	DSFResponsePoolID string         `json:"dsf_response_pool_id"`
	DSFServiceID      string         `json:"service_id"`
	PendingChange     string         `json:"pending_change"`
	DSFRecordSets     []DSFRecordSet `json:"record_sets"`
}

type DSFRecordSet struct {
	Status        string      `json:"status"`
	Eligible      SBool       `json:"eligible"`
	ID            string      `json:"dsf_record_set_id"`
	MonitorID     string      `json:"dsf_monitor_id"`
	Label         string      `json:"label"`
	TroubleCount  SInt        `json:"trouble_count"`
	Records       []DSFRecord `json:"records"`
	FailCount     SInt        `json:"fail_count"`
	TorpidityMax  string      `json:"torpidity_max"`
	TTLDerived    string      `json:"ttl_derived"`
	LastMonitored string      `json:"last_monitored"`
	TTL           SInt        `json:"ttl"`
	ServiceID     string      `json:"service_id"`
	ServeCount    SInt        `json:"serve_count"`
	Automation    string      `json:"automation"`
	PendingChange string      `json:"pending_change"`
	RDataClass    string      `json:"rdata_class"`
}

type DSFRecord struct {
	Status         string   `json:"status"`
	Endpoints      []string `json:"endpoints"`
	RDataClass     string   `json:"rdata_class"`
	Weight         int      `json:"weight"`
	Eligible       SBool    `json:"eligible"`
	ID             string   `json:"dsf_record_id"`
	DSFRecordSetID string   `json:"dsf_record_set_id"`
	//RData           interface{} `json:"rdata"`
	EndpointUpCount int    `json:"endpoint_up_count"`
	Label           string `json:"label"`
	MasterLine      string `json:"master_line"`
	Torpidity       int    `json:"torpidity"`
	LastMonitored   int    `json:"last_monitored"`
	TTL             string `json:"ttl"`
	DSFServiceID    string `json:"service_id"`
	PendingChange   string `json:"pending_change"`
	Automation      string `json:"automation"`
	ReponseTime     int    `json:"response_time"`
	Publish         string `json:"publish",omit_empty`
}

type DSFNodeRequest struct {
	PublishBlock
	Node []DSFNode `json:"nodes"`
}

type DSFNodeResponse struct {
	ResponseBlock
	Data []DSFNode `json:"data"`
}

type DSFNode struct {
	Zone string `json:"zone"`
	FQDN string `json:"fqdn"`
}

type Notifier struct {
	ID         int    `json:"notifier_id"`
	Label      string `json:"label"`
	Recipients string `json:"recipients"`
	Active     string `json:"active"`
}

type DSFRsfcRequest struct {
	PublishBlock
	Label string `json:"label"`
}
type DSFRsfcResponse struct {
	ResponseBlock
	Data DSFRecordSetChain `json:"data"`
}
type DSFRecordSetRequest struct {
	PublishBlock
	Label          string  `json:"label"`
	RDataClass     string  `json:"rdata_class,omitempty"`
	TTL            SInt    `json:"ttl,omitempty"`
	Automation     string  `json:"automation,omitempty"`
	ServeCount     SInt    `json:"serve_count,omitempty"`
	FailCount      SInt    `json:"fail_count,omitempty"`
	TroubleCount   SInt    `json:"trouble_count,omitempty"`
	Eligible       *SBool  `json:"eligible,omitempty"`
	MonitorID      *string `json:"dsf_monitor_id"`
	DSFRsfc        string  `json:"dsf_record_set_failover_chain_id,omitempty"`
	ResponsePoolId string  `json:"dsf_response_pool_id,omitempty"`
}

type DSFRecordSetResponse struct {
	ResponseBlock
	Data DSFRecordSet `json:"data"`
}

type DSFRecordRequest struct {
	PublishBlock
	Label      string `json:"label"`
	Weight     int    `json:"weight,omitempty"`
	Automation string `json:"automation,omitempty"`
	Eligible   *SBool `json:"eligible,omitempty"`
	MasterLine string `json:"master_line,omitempty"`
}

type DSFRecordResponse struct {
	ResponseBlock
	Data DSFRecord `json:"data"`
}

type DSFMonitorResponse struct {
	ResponseBlock
	Data DSFMonitor `json:"data"`
}

type DSFMonitor struct {
	ID            string             `json:"dsf_monitor_id,omitempty"`
	Label         string             `json:"label"`
	Protocol      string             `json:"protocol"`
	Active        YNBool             `json:"active,omitempty"`
	ResponseCount SInt               `json:"response_count"`
	ProbeInterval SInt               `json:"probe_interval"`
	Retries       SInt               `json:"retries"`
	Options       *DSFMonitorOptions `json:"options,omitempty"`
}

type DSFMonitorOptions struct {
	Timeout  SInt   `json:"timeout,omitempty"`
	Port     SInt   `json:"port,omitempty"`
	Path     string `json:"path,omitempty"`
	Host     string `json:"host,omitempty"`
	Header   string `json:"header,omitempty"`
	Expected string `json:"expected,omitempty"`
}
