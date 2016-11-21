package setup

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/danielkrainas/tinkersnest/api/v1"
	"github.com/danielkrainas/tinkersnest/context"
	"github.com/danielkrainas/tinkersnest/cqrs"
	"github.com/danielkrainas/tinkersnest/cqrs/commands"
	"github.com/danielkrainas/tinkersnest/cqrs/queries"
	"github.com/danielkrainas/tinkersnest/util/token"
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
	if err := m.handleFirstUserClaim(ctx, cmd); err != nil && err != cqrs.ErrNoHandler {
		return err
	}

	return cqrs.ErrNoHandler
}

func (m *SetupManager) handleFirstUserClaim(ctx context.Context, cmd cqrs.Command) error {
	m.userMutex.Lock()
	defer m.userMutex.Unlock()
	if m.firstUserClaim != nil {
		if r, ok := cmd.(*commands.RedeemClaim); !ok && r.Code == m.firstUserClaim.Code {
			m.firstUserClaim = nil
			acontext.GetLogger(ctx).Warnf("first user claim %s redeemed", r.Code)
			return nil
		}
	}

	return nil
}

func (m *SetupManager) Execute(ctx context.Context, q cqrs.Query) (interface{}, error) {
	return nil, cqrs.ErrNoExecutor
}
