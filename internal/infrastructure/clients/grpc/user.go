package clients

import (
	"context"

	"github.com/mc-lovin-132/auth/internal/domain"

	userspb "github.com/mc-lovin-132/users/pb"
)

type UserClient struct {
	client userspb.UserServiceClient
}

func New(client userspb.UserServiceClient) *UserClient {
	return &UserClient{client: client}
}

func (c *UserClient) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	response, err := c.client.Get(ctx, &userspb.GetRequest{
		Selector: &userspb.GetRequest_Email{
			Email: email,
		},
	})
	if err != nil {
		// mapping
		return nil, err
	}
	return &domain.User{
		ID:       int(response.User.Id),
		Name:     response.User.Name,
		Email:    response.User.Email,
		Password: response.User.Password,
	}, nil
}
