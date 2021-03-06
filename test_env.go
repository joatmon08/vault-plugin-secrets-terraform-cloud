package tfsecrets

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
)

const (
	envVarRunAccTests           = "VAULT_ACC"
	envVarTerraformToken        = "TF_TOKEN"
	envVarTerraformOrganization = "TF_ORGANIZATION"
	envVarTerraformTeamID       = "TF_TEAM_ID"
)

var runAcceptanceTests = os.Getenv(envVarRunAccTests) == "1"

type testEnv struct {
	Token        string
	Organization string
	TeamID       string

	Backend logical.Backend
	Context context.Context
	Storage logical.Storage

	MostRecentSecret *logical.Secret
}

func (e *testEnv) AddConfig(t *testing.T) {
	req := &logical.Request{
		Operation: logical.CreateOperation,
		Path:      "config",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"token": e.Token,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
	assert.Nil(t, resp)
}

func (e *testEnv) AddOrgTokenRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-org-token",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization": e.Organization,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
}

func (e *testEnv) ReadOrgToken(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/test-org-token",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Data["token"])

	e.MostRecentSecret = resp.Secret
}

func (e *testEnv) RenewOrgToken(t *testing.T) {
	req := &logical.Request{
		Operation: logical.RenewOperation,
		Storage:   e.Storage,
		Secret:    e.MostRecentSecret,
		Data: map[string]interface{}{
			"lease_id": "foo",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
	assert.NotNil(t, resp)
	assert.Equal(t, e.MostRecentSecret, resp.Secret)
}

func (e *testEnv) RevokeOrgToken(t *testing.T) {
	req := &logical.Request{
		Operation: logical.RevokeOperation,
		Storage:   e.Storage,
		Secret:    e.MostRecentSecret,
		Data: map[string]interface{}{
			"lease_id": "foo",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
	assert.Nil(t, resp)
}

func (e *testEnv) AddTeamTokenRole(t *testing.T) {
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/test-team-token",
		Storage:   e.Storage,
		Data: map[string]interface{}{
			"organization": e.Organization,
			"team_id":      e.TeamID,
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
}

func (e *testEnv) ReadTeamToken(t *testing.T) {
	req := &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "creds/test-team-token",
		Storage:   e.Storage,
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.Data["token"])

	e.MostRecentSecret = resp.Secret
}

func (e *testEnv) RenewTeamToken(t *testing.T) {
	req := &logical.Request{
		Operation: logical.RenewOperation,
		Storage:   e.Storage,
		Secret:    e.MostRecentSecret,
		Data: map[string]interface{}{
			"lease_id": "foo",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
	assert.NotNil(t, resp)
	assert.Equal(t, e.MostRecentSecret, resp.Secret)
}

func (e *testEnv) RevokeTeamToken(t *testing.T) {
	req := &logical.Request{
		Operation: logical.RevokeOperation,
		Storage:   e.Storage,
		Secret:    e.MostRecentSecret,
		Data: map[string]interface{}{
			"lease_id": "foo",
		},
	}
	resp, err := e.Backend.HandleRequest(e.Context, req)
	assert.False(t, (err != nil || (resp != nil && resp.IsError())), fmt.Sprintf("bad: resp: %#v\nerr:%v", resp, err))
	assert.Nil(t, resp)
}
