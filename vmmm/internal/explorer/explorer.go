package explorer

import (
	"context"
	"net/url"

	"github.com/pscompsci/vmmm/internal/vm"

	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
)

type Repository interface {
	GetVMs(ctx context.Context) (*[]vm.VirtualMachine, error)
}

type Explorer struct {
	explorerRepository Repository
}

// TODO: Clean this up significantly. It works fine, but is too complicated. Simplify and split
func (e *Explorer) GetVMListFromHost(ctx context.Context, url string, insecure bool) (*[]vm.VirtualMachine, error) {
	var vms []vm.VirtualMachine

	u, err := stringToURL(url)
	if err != nil {
		return nil, err
	}

	// Override username and/or password as required
	processURL(u)

	// Share govc's session cache
	s := &cache.Session{
		URL:      u,
		Insecure: insecure,
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
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{}, &movms)
	if err != nil {
		return &vms, err
	}

	var dss []mo.Datastore
	err = v.Retrieve(ctx, []string{"Datastore"}, []string{"summary"}, &dss)
	if err != nil {
		return &vms, err
	}

	// Print summary per vm (see also: govc/vm/info.go)
	for _, movm := range movms {
		ds := getDatastore(dss, movm.Summary.Config.Name)
		vm := vm.VirtualMachine{
			Name:            movm.Summary.Config.Name,
			Parent:          movm.Parent.Value,
			Network:         movm.Network[len(movm.Network)-1].Value,
			OperatingSystem: movm.Guest.GuestFullName,
			IPAddress:       movm.Guest.IpAddress,
			CPU:             movm.Config.Hardware.NumCPU,
			Memory:          movm.Config.Hardware.MemoryMB,
			DiskType:        ds.Summary.Type,
			DiskCapacity:    int32(ds.Summary.Capacity),
			DiskFreeSpace:   int32(ds.Summary.FreeSpace),
			State:           movm.Guest.GuestState,
			OverallStatus:   string(movm.Summary.OverallStatus),
		}
		vms = append(vms, vm)
	}

	return &vms, nil
}

func getDatastore(dss []mo.Datastore, name string) mo.Datastore {
	for _, ds := range dss {
		if ds.Summary.Name == name {
			return ds
		}
	}
	return mo.Datastore{}
}

func stringToURL(url string) (*url.URL, error) {
	return soap.ParseURL(url)
}

func processURL(u *url.URL) {
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

func (e *Explorer) GetVMListFromDB(ctx context.Context) (*[]vm.VirtualMachine, error) {
	return e.explorerRepository.GetVMs(ctx)
}
