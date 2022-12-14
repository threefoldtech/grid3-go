package workloads

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/threefoldtech/grid3-go/deployer"
	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

var ErrDuplicateName = errors.New("node names are not unique")

type K8sNodeData struct {
	Name      string
	Node      uint32
	DiskSize  int
	PublicIP  bool
	PublicIP6 bool
	Planetary bool
	IP        string
	Cpu       int
	Memory    int
}

type K8sCluster struct {
	Master      *K8sNodeData
	Workers     []K8sNodeData
	Token       string
	SSHKey      string
	NetworkName string
}

func (k *K8sCluster) Stage(
	ctx context.Context,
	manager deployer.DeploymentManager,
) error {

	err := k.validateNames()
	if err != nil {
		return err
	}

	workloads := map[uint32][]gridtypes.Workload{}

	workloads[k.Master.Node] = append(workloads[k.Master.Node], k.Master.GenerateK8sWorkload(manager, k, false)...)
	for _, worker := range k.Workers {
		workloads[worker.Node] = append(workloads[worker.Node], worker.GenerateK8sWorkload(manager, k, true)...)
	}

	err = manager.SetWorkloads(workloads)
	if err != nil {
		return err
	}

	return nil
}

func (k *K8sCluster) validateNames() error {
	names := map[string]bool{}
	names[k.Master.Name] = true
	for _, worker := range k.Workers {
		if _, ok := names[worker.Name]; ok {
			return errors.Wrapf(ErrDuplicateName, "name %s is duplicated", worker.Name)
		}
		names[worker.Name] = true
	}
	return nil
}

func flistChecksumURL(url string) string {
	return fmt.Sprintf("%s.md5", url)
}

func getFlistChecksum(url string) (string, error) {
	response, err := http.Get(flistChecksumURL(url))
	if err != nil {
		return "", err
	}
	hash, err := ioutil.ReadAll(response.Body)
	return strings.TrimSpace(string(hash)), err
}

func (k *K8sNodeData) GenerateK8sWorkload(manager deployer.DeploymentManager, deployer *K8sCluster, worker bool) []gridtypes.Workload {
	diskName := fmt.Sprintf("%sdisk", k.Name)
	workloads := make([]gridtypes.Workload, 0)
	diskWorkload := gridtypes.Workload{
		Name:        gridtypes.Name(diskName),
		Version:     0,
		Type:        zos.ZMountType,
		Description: "",
		Data: gridtypes.MustMarshal(zos.ZMount{
			Size: gridtypes.Unit(k.DiskSize) * gridtypes.Gigabyte,
		}),
	}
	workloads = append(workloads, diskWorkload)
	publicIPName := ""
	if k.PublicIP || k.PublicIP6 {
		publicIPName = fmt.Sprintf("%sip", k.Name)
		workloads = append(workloads, constructPublicIPWorkload(publicIPName, k.PublicIP, k.PublicIP6))
	}
	envVars := map[string]string{
		"SSH_KEY":           deployer.SSHKey,
		"K3S_TOKEN":         deployer.Token,
		"K3S_DATA_DIR":      "/mydisk",
		"K3S_FLANNEL_IFACE": "eth0",
		"K3S_NODE_NAME":     k.Name,
		"K3S_URL":           "",
	}
	if worker {
		// K3S_URL marks where to find the master node
		envVars["K3S_URL"] = fmt.Sprintf("%d:%s", deployer.Master.Node, deployer.Master.Name)
	}
	log.Printf("env k3s: %+v", envVars)
	workload := gridtypes.Workload{
		Version: 0,
		Name:    gridtypes.Name(k.Name),
		Type:    zos.ZMachineType,
		Data: gridtypes.MustMarshal(zos.ZMachine{
			FList: "https://hub.grid.tf/tf-official-apps/threefoldtech-k3s-latest.flist",
			Network: zos.MachineNetwork{
				Interfaces: []zos.MachineInterface{
					{
						Network: gridtypes.Name(deployer.NetworkName),
						IP:      net.ParseIP(k.IP),
					},
				},
				PublicIP:  gridtypes.Name(publicIPName),
				Planetary: k.Planetary,
			},
			ComputeCapacity: zos.MachineCapacity{
				CPU:    uint8(k.Cpu),
				Memory: gridtypes.Unit(uint(k.Memory)) * gridtypes.Megabyte,
			},
			Entrypoint: "/sbin/zinit init",
			Mounts: []zos.MachineMount{
				{Name: gridtypes.Name(diskName), Mountpoint: "/mydisk"},
			},
			Env: envVars,
		}),
	}
	workloads = append(workloads, workload)
	return workloads
}
