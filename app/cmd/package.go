package cmd

import (
	"diviner/common/csp"
	"fmt"

	"github.com/hyperledger/fabric/bccsp"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	pbm "diviner/protos/member"
	pbs "diviner/protos/service"
)

var (
	priv   bccsp.Key
	member *pbm.Member
	conn   *grpc.ClientConn
	client pbs.DivinerSerivceClient

	eventId  = ""
	title    = ""
	outcomes []string

	marketId = ""
	number   = 0.0
	isFund   = false
)

func Init() {
	var err error
	ski := viper.GetString("ski")
	priv, err = csp.ImportPrivFromSKI(ski)
	if err != nil {
		panic(fmt.Sprintf("import private key from ski error: %v\n", err))
	}

	member, err = pbm.NewMember(priv, 0.0)
	if err != nil {
		panic(fmt.Sprintf("create member structure error: %v", err))
	}

	conn, err = grpc.Dial(viper.GetString("host"), grpc.WithInsecure())
	if err != nil {
		panic(fmt.Sprintf("dial grpc server error: %v\n", err))
	}

	client = pbs.NewDivinerSerivceClient(conn)
}

func CloseConnection() {
	if conn != nil {
		conn.Close()
	}
}
