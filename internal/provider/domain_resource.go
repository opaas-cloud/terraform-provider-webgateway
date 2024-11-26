package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"net/http"
	"terraform-provider-webgateway/tools"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &domainResource{}
	_ resource.ResourceWithConfigure = &domainResource{}
)

// NewRepoResource is a helper function to simplify the provider implementation.
func NewDomainResource() resource.Resource {
	return &domainResource{}
}

// repoResource is the resource implementation.
type domainResource struct {
	client *tools.WebGateway
}

// Configure adds the provider configured client to the resource.
func (r *domainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tools.WebGateway)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *domainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Schema defines the schema for the resource.
func (r *domainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Optional: true,
			},
			"ip_address": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

// Create a new resource.
func (r *domainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tools.WebApi
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var convert = tools.WebApiModel{
		Domain:    plan.Domain.ValueString(),
		IpAddress: plan.IpAddress.ValueString(),
	}

	out, err := json.Marshal(convert)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	request, err := http.NewRequest("POST", r.client.Host+"/create", bytes.NewBuffer(out))
	request.Header.Add("Authorization", "Bearer "+r.client.Token)
	request.Header.Add("Content-Type", "application/json")

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	defer response.Body.Close()

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *domainResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *domainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state tools.WebApi
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *domainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tools.WebApi
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	var convert = tools.WebApiModel{
		Domain:    state.Domain.ValueString(),
		IpAddress: state.IpAddress.ValueString(),
	}

	out, err := json.Marshal(convert)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	request, _ := http.NewRequest("DELETE", r.client.Host+"/delete", bytes.NewBuffer(out))
	request.Header.Add("Authorization", "Bearer "+r.client.Token)
	request.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	defer response.Body.Close()

	if resp.Diagnostics.HasError() {
		return
	}
}
