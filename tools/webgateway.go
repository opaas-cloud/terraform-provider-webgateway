package tools

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type WebGateway struct {
	Token string
	Host  string
}

type WebApi struct {
	Domain    types.String `tfsdk:"domain"`
	IpAddress types.String `tfsdk:"ip_address"`
}

type WebApiModel struct {
	Domain    string `json:"domain"`
	IpAddress string `json:"ip_address"`
}
