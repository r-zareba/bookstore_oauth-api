package cassandra_db

import (
	"github.com/gocql/gocql"
	"github.com/r-zareba/bookstore_oauth-api/src/clients/cassandra"
	"github.com/r-zareba/bookstore_oauth-api/src/domain/model/access_token"
	"github.com/r-zareba/bookstore_utils-go/rest_errors"
)

const (
	getAccessTokenQuery    = "SELECT access_token, user_id, client_id, expires FROM access_tokens WHERE access_token=?;"
	createAccessTokenQuery = "INSERT INTO access_tokens(access_token, user_id, client_id, expires) VALUES(?, ?, ?, ?);"
	updateExpiresInQuery   = "UPDATE access_tokens SET expires=? WHERE access_token=?;"
)

type CassandraRepository interface {
	GetTokenById(string) (*access_token.AccessToken, *rest_errors.RestError)
	CreateToken(access_token.AccessToken) *rest_errors.RestError
	UpdateExpiresIn(access_token.AccessToken) *rest_errors.RestError
}

func NewCassandraRepository() CassandraRepository {
	return &cassandraRepository{}
}

type cassandraRepository struct {
}

func (r *cassandraRepository) GetTokenById(id string) (*access_token.AccessToken, *rest_errors.RestError) {
	session := cassandra.GetSession()
	var result access_token.AccessToken
	queryErr := session.Query(getAccessTokenQuery, id).Scan(&result.Token, &result.UserId, &result.ClientId, &result.ExpiresIn)
	if queryErr != nil {
		if queryErr == gocql.ErrNotFound {
			return nil, rest_errors.NotFoundError("No access token with given id")
		}
		return nil, rest_errors.InternalServerError(queryErr.Error(), queryErr)
	}
	return &result, nil
}

func (r *cassandraRepository) CreateToken(token access_token.AccessToken) *rest_errors.RestError {
	session := cassandra.GetSession()
	queryErr := session.Query(createAccessTokenQuery, token.Token, token.UserId, token.ClientId, token.ExpiresIn).Exec()
	if queryErr != nil {
		return rest_errors.InternalServerError(queryErr.Error(), queryErr)
	}
	return nil
}

func (r *cassandraRepository) UpdateExpiresIn(token access_token.AccessToken) *rest_errors.RestError {
	session := cassandra.GetSession()
	queryErr := session.Query(updateExpiresInQuery, token.ExpiresIn, token.Token).Exec()
	if queryErr != nil {
		return rest_errors.InternalServerError(queryErr.Error(), queryErr)
	}
	return nil
}
