package delivery

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"

	uuid "github.com/satori/go.uuid"
	sess "no_homomorphism/configs/proto/session"
	"no_homomorphism/sessions/internal"
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

func (uc *SessionDelivery) Create(ctx context.Context, in *sess.Session) (*sess.SessionID, error) {
	sid, err := uc.UseCase.Create(in.Login, uc.ExpireTime)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &sess.SessionID{ID: sid.String()}, nil
}

func (uc *SessionDelivery) Delete(ctx context.Context, in *sess.SessionID) (*sess.Nothing, error) {
	sid, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse uuid from string")
	}
	if err := uc.UseCase.Delete(sid); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &sess.Nothing{Dummy: true}, nil
}

func (uc *SessionDelivery) Check(ctx context.Context, in *sess.SessionID) (*sess.Session, error) {
	sid, err := uuid.FromString(in.ID)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "can't parse uuid from string")
	}
	login, err := uc.UseCase.Check(sid)
	if err != nil {
		return nil, status.Error(codes.Aborted, err.Error())
	}
	return &sess.Session{Login: login}, nil
}
