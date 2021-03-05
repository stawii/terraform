package exec

import (
	"context"
	// "crypto/tls"
	// "fmt"
	// "net/http"
	// "net/url"
	// "time"

	// "github.com/hashicorp/go-cleanhttp"
	// "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform/backend"
	"github.com/hashicorp/terraform/internal/legacy/helper/schema"
	"github.com/hashicorp/terraform/states/remote"
	"github.com/hashicorp/terraform/states/statemgr"
)

func New() backend.Backend {
	s := &schema.Backend{
		Schema: map[string]*schema.Schema{
			"load_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_HTTP_ADDRESS", nil),
				Description: "The command to execute when reading state",
			},
			"save_command": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_HTTP_UPDATE_METHOD", nil),
				Description: "The command to execute when updating state",
			},
			"lock_command": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_HTTP_LOCK_ADDRESS", nil),
				Description: "The command to execute to lock the state",
			},
			"unlock_command": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_HTTP_UNLOCK_ADDRESS", nil),
				Description: "The command to execute to unlock the state",
			},
		},
	}

	b := &Backend{Backend: s}
	b.Backend.ConfigureFunc = b.configure
	return b
}

type Backend struct {
	*schema.Backend
	client *execClient
}

func (b *Backend) configure(ctx context.Context) error {
	data := schema.FromContextBackendConfig(ctx)
	b.client = &execClient{
		loadCommand : data.Get("load_command").(string),
		saveCommand : data.Get("save_command").(string),
		lockCommand : data.Get("lock_command").(string),
		unlockCommand : data.Get("unlock_command").(string),
	}
	return nil
}


func (b *Backend) StateMgr(name string) (statemgr.Full, error) {
	if name != backend.DefaultStateName {
		return nil, backend.ErrWorkspacesNotSupported
	}
	return &remote.State{Client: b.client}, nil
}

func (b *Backend) Workspaces() ([]string, error) {
	return nil, backend.ErrWorkspacesNotSupported
}

func (b *Backend) DeleteWorkspace(string) error {
	return backend.ErrWorkspacesNotSupported
}
