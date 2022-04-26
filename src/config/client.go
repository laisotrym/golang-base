package config

import (
    "flag"
    "fmt"

    "github.com/pkg/errors"
    "github.com/spf13/viper"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/testdata"
)

var (
    tls    = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
    caFile = flag.String("ca_file", "", "The file containing the CA root cert file")
    // serverAddr         = flag.String("server_addr", "localhost:10000", "The server address in the format of host:port")
    serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

type ClientAlias string

const (
    ClientRuleDefault      ClientAlias = "default"
    ClientRuleBlackList                = "rule_blacklist"
    ClientRuleInfo                     = "rule_info"
    ClientRuleMerchant                 = "rule_merchant"
    ClientRuleBank                     = "rule_bank"
    ClientRuleCustomer                 = "rule_customer"
    ClientCalculateBenefit             = "calculate_benefit"
    ClientTransaction                  = "transaction"
)

type ClientServiceAddr struct {
    Addr     string `json:"addr" yaml:"addr" validate:"required,hostname"`
    Port   int64  `json:"port" yaml:"port" validate:"required,numeric"`
    ApiKey string `json:"api_key" yaml:"apiKey" validate:"required"`
}

func (c ClientServiceAddr) String() string {
    return fmt.Sprintf("%s", c.DSN())
}

func (c ClientServiceAddr) DSN() string {
    return fmt.Sprintf("%s:%d", c.Addr, c.Port)
}

func ConnectClient(serverAddr string) (*grpc.ClientConn, error) {
    var opts []grpc.DialOption
    if *tls {
        if *caFile == "" {
            *caFile = testdata.Path("ca.pem")
        }

        if transportCredentials, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride); err != nil {
            return nil, errors.New(fmt.Sprintf("Failed to create TLS credentials %v", err))
        } else {
            opts = append(opts, grpc.WithTransportCredentials(transportCredentials))
        }
    } else {
        opts = append(opts, grpc.WithInsecure())
    }

    opts = append(opts, grpc.WithBlock())

    if conn, err := grpc.Dial(serverAddr, opts...); err != nil {
        return nil, errors.New(fmt.Sprintf("Fail to dial: %v", err))
    } else {
        return conn, nil
    }
}

func DefaultClientConfig() map[ClientAlias]ClientServiceAddr {
    return map[ClientAlias]ClientServiceAddr{
        ClientRuleDefault: {
            // Addr: "localhost",
            Port:   10443,
            ApiKey: "NoSecret",
        },
    }
}

func InitClientConfig(vip *viper.Viper) map[ClientAlias]ClientServiceAddr {
    cnf := DefaultClientConfig()
    if vip != nil {
        if err := vip.Unmarshal(&cnf); err != nil {
            fmt.Printf("unable to decode into config struct, %v", err)
        }
    }
    return cnf
}
