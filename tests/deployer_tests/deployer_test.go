package deployertests

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/threefoldtech/grid3-go/deployer"
	client "github.com/threefoldtech/grid3-go/node"
	mock "github.com/threefoldtech/grid3-go/tests/mocks"
	"github.com/threefoldtech/grid3-go/workloads"
	"github.com/threefoldtech/substrate-client"

	"github.com/threefoldtech/zos/pkg/gridtypes"
	"github.com/threefoldtech/zos/pkg/gridtypes/zos"
)

const Words = "secret add bag cluster deposit beach illness letter crouch position rain arctic"
const twinID = 214

var identity, _ = substrate.NewIdentityFromEd25519Phrase(Words)

func deployment1(identity substrate.Identity, TLSPassthrough bool, version uint32) gridtypes.Deployment {
	dl := workloads.NewDeployment(twinID)
	dl.Version = version

	gw := workloads.GatewayNameProxy{
		Name:           "name",
		TLSPassthrough: TLSPassthrough,
		Backends:       []zos.Backend{"http://1.1.1.1"},
	}
	workload, err := gw.GenerateWorkloadFromGName(gw)
	if err != nil {
		panic(err)
	}
	dl.Workloads = append(dl.Workloads, workload)
	dl.Workloads[0].Version = version
	err = dl.Sign(twinID, identity)
	if err != nil {
		panic(err)
	}

	return dl
}

func deployment2(identity substrate.Identity) gridtypes.Deployment {
	dl := workloads.NewDeployment(uint32(twinID))
	gw := workloads.GatewayFQDNProxy{
		Name:     "fqdn",
		FQDN:     "a.b.com",
		Backends: []zos.Backend{"http://1.1.1.1"},
	}

	workload, err := gw.GenerateWorkloadFromFQDN(gw)
	if err != nil {
		panic(err)
	}
	dl.Workloads = append(dl.Workloads, workload)
	err = dl.Sign(twinID, identity)
	if err != nil {
		panic(err)
	}
	return dl
}

func hash(dl *gridtypes.Deployment) string {
	hash, err := dl.ChallengeHash()
	if err != nil {
		panic(err)
	}
	hashHex := hex.EncodeToString(hash)
	return hashHex
}

type EmptyValidator struct{}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	gridClient := mock.NewMockClient(ctrl)
	cl := mock.NewRMBMockClient(ctrl)
	sub := mock.NewMockSubstrateExt(ctrl)
	ncPool := mock.NewMockNodeClientCollection(ctrl)
	deployer := deployer.NewDeployer(
		identity,
		214,
		gridClient,
		ncPool,
		true,
	)

	dl1, dl2 := deployment1(identity, true, 0), deployment2(identity)
	newDls := map[uint32]gridtypes.Deployment{
		10: dl1,
		20: dl2,
	}

	dl1.ContractID = 100
	dl2.ContractID = 200

	sub.EXPECT().
		CreateNodeContract(
			identity,
			uint32(10),
			nil,
			hash(&dl1),
			uint32(0),
			true,
		).Return(uint64(100), nil)
	sub.EXPECT().
		CreateNodeContract(
			identity,
			uint32(20),
			nil,
			hash(&dl2),
			uint32(0),
			true,
		).Return(uint64(200), nil)
	ncPool.EXPECT().
		GetNodeClient(sub, uint32(10)).
		Return(client.NewNodeClient(13, cl), nil)
	ncPool.EXPECT().
		GetNodeClient(sub, uint32(20)).
		Return(client.NewNodeClient(23, cl), nil)
	cl.EXPECT().
		Call(gomock.Any(), uint32(13), "zos.deployment.deploy", dl1, gomock.Any()).
		DoAndReturn(func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
			dl1.Workloads[0].Result.State = gridtypes.StateOk
			dl1.Workloads[0].Result.Data, _ = json.Marshal(zos.GatewayProxyResult{})
			return nil
		})
	cl.EXPECT().
		Call(gomock.Any(), uint32(23), "zos.deployment.deploy", dl2, gomock.Any()).
		DoAndReturn(func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
			dl2.Workloads[0].Result.State = gridtypes.StateOk
			dl2.Workloads[0].Result.Data, _ = json.Marshal(zos.GatewayFQDNResult{})
			return nil
		})
	cl.EXPECT().
		Call(gomock.Any(), uint32(13), "zos.deployment.get", gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
			var res *gridtypes.Deployment = result.(*gridtypes.Deployment)
			*res = dl1
			return nil
		})
	cl.EXPECT().
		Call(gomock.Any(), uint32(23), "zos.deployment.get", gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, twin uint32, fn string, data, result interface{}) error {
			var res *gridtypes.Deployment = result.(*gridtypes.Deployment)
			*res = dl2
			return nil
		})
	deployer.(*DeployerImpl).validator = &EmptyValidator{}
	contracts, err := deployer.Deploy(context.Background(), sub, nil, newDls)
	assert.NoError(t, err)
	assert.Equal(t, contracts, map[uint32]uint64{10: 100, 20: 200})
}

// func TestUpdate(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	gridClient := mock.NewMockClient(ctrl)
// 	// cl := mock.
// }