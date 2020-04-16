package delivery

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	session "github.com/2020_1_no_homomorphism/no_homo_sessions/internal"
	"time"

	uuid "github.com/satori/go.uuid"
)

type SessionDelivery struct {
	UseCase    session.UseCase
	ExpireTime time.Duration
}

func NewSessionDelivery(useCase session.UseCase, expire time.Duration) *SessionDelivery {
	return &SessionDelivery{
		UseCase:    useCase,
		ExpireTime: expire,
	}
}

func (uc *SessionDelivery) Create(ctx context.Context, in *session.Session) (*session.SessionID, error) {
	sid, err := uc.UseCase.Create(in.Login, uc.ExpireTime)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &session.SessionID{ID: sid.String()}, nil
}

func (uc *SessionDelivery) Delete(ctx context.Context, in *session.SessionID) (*session.Nothing, error) {
	sid, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse uuid from string")
	}
	if err := uc.UseCase.Delete(sid); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &session.Nothing{Dummy: true}, nil
}

func (uc *SessionDelivery) Check(ctx context.Context, in *session.SessionID) (*session.Session, error) {
	sid, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse uuid from string")
	}
	login, err := uc.UseCase.Check(sid)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &session.Session{Login: login}, nil
}
