package hedera

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const testClientJson string = `{
    "network": {
        "35.237.200.180:50211": "0.0.3",
        "35.186.191.247:50211": "0.0.4",
        "35.192.2.25:50211": "0.0.5",
        "35.199.161.108:50211": "0.0.6",
        "35.203.82.240:50211": "0.0.7",
        "35.236.5.219:50211": "0.0.8",
        "35.197.192.225:50211": "0.0.9",
        "35.242.233.154:50211": "0.0.10",
        "35.240.118.96:50211": "0.0.11",
        "35.204.86.32:50211": "0.0.12"
    }
}`

const testClientJsonWithOperator string = `{
    "network": {
        "35.237.200.180:50211": "0.0.3",
        "35.186.191.247:50211": "0.0.4",
        "35.192.2.25:50211": "0.0.5",
        "35.199.161.108:50211": "0.0.6",
        "35.203.82.240:50211": "0.0.7",
        "35.236.5.219:50211": "0.0.8",
        "35.197.192.225:50211": "0.0.9",
        "35.242.233.154:50211": "0.0.10",
        "35.240.118.96:50211": "0.0.11",
        "35.204.86.32:50211": "0.0.12"
    },
    "operator": {
        "accountId": "0.0.3",
        "privateKey": "302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10"
    }
}`

func TestClientFromJSON(t *testing.T) {
	client, err := ClientFromJSON([]byte(testClientJson))
	assert.NoError(t, err)

	assert.NotNil(t, client)
	assert.Equal(t, 10, len(client.networkNodeIds))
	assert.Nil(t, client.operator)
}

func TestClientFromJSONWithOperator(t *testing.T) {
	client, err := ClientFromJSON([]byte(testClientJsonWithOperator))
	assert.NoError(t, err)

	assert.NotNil(t, client)

	testOperatorKey, err := Ed25519PrivateKeyFromString("302e020100300506032b657004220420db484b828e64b2d8f12ce3c0a0e93a0b8cce7af1bb8f39c97732394482538e10")
	assert.NoError(t, err)

	assert.Equal(t, 10, len(client.networkNodeIds))
	assert.NotNil(t, client.operator)
	assert.Equal(t, testOperatorKey.keyData, client.operator.privateKey.keyData)
	assert.Equal(t, AccountID{Account: 3}, client.operator.accountID)
}

func TestClientMaxQueryFee(t *testing.T) {
	client, err := clientForTest()
	assert.NoError(t, err)

	client.SetMaxQueryPayment(ZeroHbar)

	query := NewAccountBalanceQuery().
		SetAccountID(AccountID{Account:3})

	_, err = query.Execute(client)

	// Having a client with a max query fee of 0 tinybar and no makually set value on the query should return an error,
	// not panic or nil.
	assert.Error(t, err)

	_, err = query.
		SetMaxQueryPayment(NewHbar(1)).
		Execute(client)

	// Manually setting the max query payment should enable the query to work.
	assert.NoError(t, err)
}
