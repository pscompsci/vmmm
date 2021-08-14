package explorer

import (
	"context"
	"net/url"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
)

type VirtualMachine struct {
	ID            int    `db:"vm_id" json:"id,omitempty"`
	Name          string `db:"name"  json:"name"`
	Parent        string `db:"parent" json:"parent"`
	OverallStatus string `db:"overall_status" json:"overallStatus"`
}

type Repository interface {
	GetVMs(ctx context.Context) (*[]VirtualMachine, error)
}

type Explorer struct {
	explorerRepository Repository
}

// TODO: Clean this up significantly. It works fine, but is too complicated. Simplify and split
func (e *Explorer) GetVMListFromHost(ctx context.Context, url string) (*[]VirtualMachine, error) {
	var vms []VirtualMachine

	u, err := soap.ParseURL(url)
	if err != nil {
		return nil, err
	}

	// Override username and/or password as required
	processOverride(u)

	// Share govc's session cache
	s := &cache.Session{
		URL:      u,
		Insecure: true,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	m := view.NewManager(c)

	v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		return &vms, err
	}

	defer v.Destroy(ctx)

	var movms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary", "parent"}, &movms)
	if err != nil {
		return &vms, err
	}

	// Print summary per vm (see also: govc/vm/info.go)
	for _, movm := range movms {
		vm := VirtualMachine{
			Name:          movm.Summary.Config.Name,
			Parent:        movm.Parent.Value,
			OverallStatus: string(movm.Summary.OverallStatus),
		}
		vms = append(vms, vm)
	}

	return &vms, nil
}

func processOverride(u *url.URL) {
	var envUsername string
	var envPassword string

	// Override username if provided
	if envUsername != "" {
		var password string
		var ok bool

		if u.User != nil {
			password, ok = u.User.Password()
		}

		if ok {
			u.User = url.UserPassword(envUsername, password)
		} else {
			u.User = url.User(envUsername)
		}
	}

	// Override password if provided
	if envPassword != "" {
		var username string

		if u.User != nil {
			username = u.User.Username()
		}

		u.User = url.UserPassword(username, envPassword)
	}
}

func (e *Explorer) GetVMListFromDB(ctx context.Context) (*[]VirtualMachine, error) {
	return e.explorerRepository.GetVMs(ctx)
}
