package authentication

import (
	"github.com/kholidss/xyz-skilltest/internal/repositories"
	"github.com/kholidss/xyz-skilltest/pkg/cdn"
	"github.com/kholidss/xyz-skilltest/pkg/config"
)

type authenticationService struct {
	cfg        *config.Config
	repoUser   repositories.UserRepository
	repoBucket repositories.BucketRepository
	repoLimit  repositories.LimitRepository
	cdn        cdn.CDN
}

func NewSvcAuthentication(
	cfg *config.Config,
	repoUser repositories.UserRepository,
	repoBucket repositories.BucketRepository,
	repoLimit repositories.LimitRepository,
	cdn cdn.CDN,
) AuthenticationService {
	return &authenticationService{
		cfg:        cfg,
		repoUser:   repoUser,
		repoBucket: repoBucket,
		repoLimit:  repoLimit,
		cdn:        cdn,
	}
}
