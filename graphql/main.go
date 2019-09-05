package graphql

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/machinebox/graphql"
	"github.com/spf13/viper"
)

func HasValidToken() bool {
	return getGraphQLToken() != ""
}

func getGraphQLToken() string {
	lagoon := viper.GetString("current")
	return viper.GetString("lagoons." + lagoon + ".token")
}

// GraphQLClient returns a new GraphQL client.
func GraphQLClient() *graphql.Client {
	lagoon := viper.GetString("current")
	return graphql.NewClient(viper.GetString("lagoons." + lagoon + ".graphql"))
}

// GraphQLRequest performs a GraphQL request.
func GraphQLRequest(q string, resp interface{}) error {
	client := GraphQLClient()
	req := graphql.NewRequest(q)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", getGraphQLToken()))
	ctx := context.Background()
	return client.Run(ctx, req, &resp)
}

// VerifyTokenExpiry verfies if the current token is valid or not
func VerifyTokenExpiry() bool {
	if HasValidToken() {
		var p jwt.Parser
		token, _, err := p.ParseUnverified(getGraphQLToken(), &jwt.StandardClaims{})
		if err = token.Claims.Valid(); err != nil {
			//handle invalid token
			return false
		}
		return true
	}
	return false
}
