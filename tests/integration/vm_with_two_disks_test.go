package integration

import (
	"context"
	"net"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/grid3-go/loader"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/zos/pkg/gridtypes"
)

func TestVMWithTwoDisk(t *testing.T) {
	manager, apiClient := setup()
	publicKey := os.Getenv("PUBLICKEY")
	network := workloads.TargetNetwork{
		Name:        "testingNetwork123",
		Description: "network for testing",
		Nodes:       []uint32{14},
		IPRange: gridtypes.NewIPNet(net.IPNet{
			IP:   net.IPv4(10, 1, 0, 0),
			Mask: net.CIDRMask(16, 32),
		}),
		AddWGAccess: false,
	}
	disk1 := workloads.Disk{
		Name: "diskTest1",
		Size: 1,
	}
	disk2 := workloads.Disk{
		Name: "diskTest2",
		Size: 2,
	}
	vm := workloads.VM{
		Name:       "vm",
		Flist:      "https://hub.grid.tf/tf-official-apps/threefoldtech-ubuntu-20.04.flist",
		Cpu:        2,
		Planetary:  true,
		Memory:     1024,
		RootfsSize: 20 * 1024,
		Entrypoint: "/init.sh",
		EnvVars: map[string]string{
			"SSH_KEY":  publicKey,
			"TEST_VAR": "this value for test",
		},
		Mounts: []workloads.Mount{
			{DiskName: "diskTest1", MountPoint: "/disk1"},
			{DiskName: "diskTest2", MountPoint: "/disk2"},
		},
		IP:          "10.1.0.2",
		NetworkName: "testingNetwork123",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	_, err := network.Stage(ctx, apiClient)
	assert.NoError(t, err)
	err = disk1.Stage(manager, 14)
	assert.NoError(t, err)
	err = disk2.Stage(manager, 14)
	assert.NoError(t, err)

	err = manager.Commit(ctx)
	assert.NoError(t, err)
	defer manager.CancelAll()

	err = vm.Stage(manager, 14)
	assert.NoError(t, err)
	err = manager.Commit(ctx)
	assert.NoError(t, err)

	result, err := loader.LoadVmFromGrid(manager, 14, "vm")
	assert.NoError(t, err)

	yggIP := result.YggIP

	// check d1 & d2 sizes

	res, err := RemoteRun("root", yggIP, "df /disk1/ | tail -1 | awk '{print $2}' | tr -d '\\n'")
	assert.NoError(t, err)
	assert.Equal(t, res, strconv.Itoa(1024*1024))
	res, err = RemoteRun("root", yggIP, "df /disk2/ | tail -1 | awk '{print $2}' | tr -d '\\n'")
	assert.NoError(t, err)
	assert.Equal(t, res, strconv.Itoa(2*1024*1024))

	// create file -> d1, check file size, move file -> d2, check file size

	_, err = RemoteRun("root", yggIP, "dd if=/dev/vda bs=1M count=512 of=/disk1/test.txt")
	assert.NoError(t, err)

	res, err = RemoteRun("root", yggIP, "du /disk1/test.txt | head -n1 | awk '{print $1;}' | tr -d -c 0-9")
	assert.NoError(t, err)
	assert.Equal(t, res, strconv.Itoa(512*1024))

	_, err = RemoteRun("root", yggIP, "mv /disk1/test.txt /disk2/")
	assert.NoError(t, err)

	res, err = RemoteRun("root", yggIP, "du /disk2/test.txt | head -n1 | awk '{print $1;}' | tr -d -c 0-9")
	assert.NoError(t, err)
	assert.Equal(t, res, strconv.Itoa(512*1024))

	// create file -> d2, check file size, copy file -> d1, check file size

	_, err = RemoteRun("root", yggIP, "dd if=/dev/vdb bs=1M count=512 of=/disk2/test.txt")
	assert.NoError(t, err)

	res, err = RemoteRun("root", yggIP, "du /disk2/test.txt | head -n1 | awk '{print $1;}' | tr -d -c 0-9")
	assert.NoError(t, err)
	assert.Equal(t, res, strconv.Itoa(512*1024))

	_, err = RemoteRun("root", yggIP, "cp /disk2/test.txt /disk1/")
	assert.NoError(t, err)

	res, err = RemoteRun("root", yggIP, "du /disk1/test.txt | head -n1 | awk '{print $1;}' | tr -d -c 0-9")
	assert.NoError(t, err)
	assert.Equal(t, res, strconv.Itoa(512*1024))

	// copy same file -> d1 (not enough space)

	_, err = RemoteRun("root", yggIP, "cp /disk2/test.txt /disk1/test2.txt")
	assert.Error(t, err)

}
