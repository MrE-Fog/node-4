/*
 * Copyright (C) 2020 The "MysteriumNetwork/node" Authors.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package pingpong

import (
	"math/big"
	"time"

	"github.com/mysteriumnetwork/node/datasize"
	"github.com/mysteriumnetwork/node/market"
	"github.com/rs/zerolog/log"
)

// CalculatePaymentAmount calculates the required payment amount.
func CalculatePaymentAmount(timePassed time.Duration, bytesTransferred DataTransferred, price market.Price) *big.Int {
	if price.IsFree() {
		return big.NewInt(0)
	}

	timeComponent := big.NewFloat(0)
	if price.PerHour.Cmp(big.NewInt(0)) > 0 {
		timeQuote := timePassed.Seconds() / time.Hour.Seconds()
		timeComponent = new(big.Float).Mul(new(big.Float).SetInt(price.PerHour), big.NewFloat(timeQuote))
	}

	dataComponent := big.NewFloat(0)
	if price.PerGiB.Cmp(big.NewInt(0)) > 0 {
		dataQuote := float64(bytesTransferred.sum()) / float64(datasize.GiB.Bytes())
		dataComponent = new(big.Float).Mul(new(big.Float).SetInt(price.PerGiB), big.NewFloat(dataQuote))
	}

	tc, _ := timeComponent.Int(nil)
	bc, _ := dataComponent.Int(nil)

	total := new(big.Int).Add(tc, bc)
	log.Debug().Msgf("Calculated price %v. Time component: %v, data component: %v ", total, timeComponent, dataComponent)
	return total
}
