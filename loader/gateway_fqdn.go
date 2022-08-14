package loader

import (
	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

func GatewayFQDNProxyFromZosWorkload(manager deployer.DeploymentManager, nodeID uint32, name string) (workloads.GatewayFQDNProxy, error) {
	wl, err := manager.GetWorkload(nodeID, name)
	if err != nil {
		return workloads.GatewayFQDNProxy{}, errors.Wrapf(err, "couldn't get workload from node %d", nodeID)
	}

	dataI, err := wl.WorkloadData()
	if err != nil {
		return workloads.GatewayFQDNProxy{}, errors.Wrap(err, "failed to get workload data")
	}
	data := dataI.(*zos.GatewayFQDNProxy)

	return workloads.GatewayFQDNProxy{
		NodeId:         nodeID,
		Name:           wl.Name.String(),
		TLSPassthrough: data.TLSPassthrough,
		Backends:       data.Backends,
		FQDN:           data.FQDN,
	}, nil
}
