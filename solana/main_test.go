package solana

import (
	"github.com/olegfomenko/solana-go-sdk/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckTxConfirmed(t *testing.T) {
	solana := Solana{
		client.NewClient("https://api.devnet.solana.com"),
	}

	ok, err := solana.checkTxConfirmed("4cBfXctdwn6ybPsnYztveLxjkzvwr6nr3dRULKmDUGHhtURCyGmkrY7rmJbm2bdFyLfiRi98uKh426AAr3RTK4vg")
	if err != nil {
		panic(err)
	}

	if ok {
		assert.Equal(t, false, ok)
	}

	ok, err = solana.checkTxConfirmed("4fP7cU5RxdE4MQaryXf3LWzfR19kkDd78LF2zq6saSf9ZGzt59CTuQA5B6RtpBtk87diZLc2LXSwEta5d5bCxXKt")
	if err != nil {
		panic(err)
	}

	if ok {
		assert.Equal(t, true, ok)
	}
}
