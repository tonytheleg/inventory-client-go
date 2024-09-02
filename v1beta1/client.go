package v1beta1

import (
	"github.com/authzed/grpcutil"
	kesselrel "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta1/relationships"
	kessel "github.com/project-kessel/inventory-api/api/kessel/inventory/v1beta1/resources"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Inventory interface {
}

type inventoryClient struct {
	K8sClusterService               kessel.KesselK8SClusterServiceClient
	PolicyRelationshipServiceClient kesselrel.KesselPolicyRelationshipServiceClient
	PolicyServiceClient             kessel.KesselK8SPolicyServiceClient
	RhelHostServiceClient           kessel.KesselRhelHostServiceClient
	gRPCConn                        *grpc.ClientConn
	tokenClient                     *TokenClient
}

var _ Inventory = &inventoryClient{}

func New(config *Config) (*inventoryClient, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.EmptyDialOption{})
	var tokencli *TokenClient
	if config.enableOIDCAuth {
		tokencli = NewTokenClient(config)
	}

	if config.insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		tlsConfig, err := grpcutil.WithSystemCerts(grpcutil.VerifyCA)
		if err != nil {
			return nil, err
		}
		opts = append(opts, tlsConfig)
	}

	conn, err := grpc.NewClient(
		config.url,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return &inventoryClient{
		K8sClusterService:               kessel.NewKesselK8SClusterServiceClient(conn),
		PolicyRelationshipServiceClient: kesselrel.NewKesselPolicyRelationshipServiceClient(conn),
		PolicyServiceClient:             kessel.NewKesselK8SPolicyServiceClient(conn),
		RhelHostServiceClient:           kessel.NewKesselRhelHostServiceClient(conn),
		gRPCConn:                        conn,
		tokenClient:                     tokencli,
	}, err
}

func (a inventoryClient) GetTokenCallOption() ([]grpc.CallOption, error) {
	var opts []grpc.CallOption
	opts = append(opts, grpc.EmptyCallOption{})
	token, err := a.tokenClient.GetToken()
	if err != nil {
		return nil, err
	}
	if a.tokenClient.Insecure {
		opts = append(opts, WithInsecureBearerToken(token.AccessToken))
	} else {
		opts = append(opts, WithBearerToken(token.AccessToken))
	}

	return opts, nil
}
