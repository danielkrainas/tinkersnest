package setup

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/danielkrainas/gobag/context"
	"github.com/danielkrainas/gobag/decouple/cqrs"
	"github.com/danielkrainas/gobag/util/token"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/commands"
	"github.com/danielkrainas/tinkersnest/queries"
)

type SetupManager struct {
	firstUserClaim *v1.Claim
	userMutex      sync.Mutex
}

func (m *SetupManager) Bootstrap(ctx context.Context) error {
	countRaw, err := cqrs.DispatchQuery(ctx, &queries.CountUsers{})
	if err != nil {
		return err
	}

	count, ok := countRaw.(int)
	if !ok {
		return fmt.Errorf("couldn't convert raw value (%#+v) to user count", countRaw)
	} else if count > 0 {
		// user exists so we don't need to do anything
		return nil
	}

	claim := m.AddFirstUserClaim(ctx)
	acontext.GetLogger(ctx).Warnf("no users found, use claim %s to create one", claim.Code)
	return nil
}

func (m *SetupManager) AddFirstUserClaim(ctx context.Context) *v1.Claim {
	m.userMutex.Lock()
	defer m.userMutex.Unlock()
	m.firstUserClaim = &v1.Claim{
		Code:         token.Generate(string(v1.UserResource)),
		Created:      time.Now().Unix(),
		ResourceType: v1.UserResource,
	}

	return m.firstUserClaim
}

func (m *SetupManager) Handle(ctx context.Context, cmd cqrs.Command) error {
	switch ct := cmd.(type) {
	case *commands.RedeemClaim:
		return m.handleFirstUserClaim(ctx, ct)
	}

	return cqrs.ErrNoHandler
}

func (m *SetupManager) handleFirstUserClaim(ctx context.Context, cmd *commands.RedeemClaim) error {
	m.userMutex.Lock()
	defer m.userMutex.Unlock()
	if m.firstUserClaim != nil {
		if cmd.Code == m.firstUserClaim.Code {
			m.firstUserClaim = nil
			acontext.GetLogger(ctx).Warnf("first user claim %s redeemed", cmd.Code)
			return nil
		}
	}

	return cqrs.ErrNoHandler
}

func (m *SetupManager) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	switch qt := q.(type) {
	case *queries.FindClaim:
		return m.executeFindClaim(ctx, qt)
	}

	return nil, cqrs.ErrNoExecutor
}

func (m *SetupManager) executeFindClaim(ctx context.Context, q *queries.FindClaim) (*v1.Claim, error) {
	m.userMutex.Lock()
	defer m.userMutex.Unlock()
	if m.firstUserClaim != nil && m.firstUserClaim.Code == q.Code {
		return m.firstUserClaim, nil
	}

	return nil, cqrs.ErrNoExecutor
}
