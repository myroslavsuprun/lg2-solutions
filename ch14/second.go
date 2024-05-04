package main

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

func secondTask() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	ctx, cancelCause := context.WithCancelCause(ctx)

	var total = big.NewInt(0)
	var i int
loop:
	for ; ; i++ {
		select {
		case <-ctx.Done():
			break loop

		default:
			n, err := randNum()
			if err != nil {
				cancelCause(fmt.Errorf("in for loop: %w", err))
				break loop
			}

			if n.Cmp(big.NewInt(1234)) == 0 {
				cancelCause(errors.New("number reached"))
				break loop
			}

			total.Add(n, total)
		}
	}

	fmt.Printf("total: %v, i: %v, err: %v\n", total.String(), i, context.Cause(ctx))

}

func randNum() (n *big.Int, err error) {
	return rand.Int(rand.Reader, big.NewInt(100_000_000))
}
