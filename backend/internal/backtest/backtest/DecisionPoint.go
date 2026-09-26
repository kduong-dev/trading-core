package backtest

import (
	"time"

	"github.com/kduong-dev/trading-core/backend/internal/tradingstrategy"
)

type DecisionPoint struct {
	At       time.Time
	Price    float64
	Action   tradingstrategy.Action
	Quantity float64
	Reason   string
}
