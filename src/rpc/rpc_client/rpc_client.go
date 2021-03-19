package rpc_client

import (
	"net/rpc"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	"github.com/SkycoinPro/skywire-services-util/src/rpc/authorization"
)

// FetchRightsFromRemoteServices - Fetching authorization rights
func FetchRightsFromRemoteServices(username string) []authorization.Right {
	var response authorization.GetResponse

	// Fetch rights from whitelist system
	rights, err := fetchRightsFromRemoteServices(username, viper.GetString("rpc.whitelist.protocol"), viper.GetString("rpc.whitelist.address"))
	if err != nil {
		log.Error("Whitelist authorization access rights fetch error: ", err)
	} else {
		log.Debug("Whitelist authorization rights fetched successfully for ", username)
		response.Rights = append(response.Rights, rights...)
	}

	// Fetch rights from chb system
	rights, err = fetchRightsFromRemoteServices(username, viper.GetString("rpc.chb.protocol"), viper.GetString("rpc.chb.address"))
	if err != nil {
		log.Error("Coinhour bank authorization access rights fetch error: ", err)
	} else {
		log.Debug("Coinhour bank authorization rights fetched successfully for ", username)
		response.Rights = append(response.Rights, rights...)
	}

	return response.Rights
}

func fetchRightsFromRemoteServices(username, protocol, address string) ([]authorization.Right, error) {
	args := &authorization.GetRequest{Username: username}
	var reply authorization.GetResponse

	client, err := rpc.DialHTTP(protocol, address)
	if err != nil {
		log.Error("dialing:", err)
	} else {
		err = client.Call("Handler.GetUserAuthorization", args, &reply)
	}

	return reply.Rights, err
}
