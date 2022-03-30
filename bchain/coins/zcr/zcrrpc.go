package zcr

import (
	"encoding/json"
	"github.com/golang/glog"
	"github.com/mosqueiro/blockbook/bchain"
	"github.com/mosqueiro/blockbook/bchain/coins/btc"
	"github.com/juju/errors"
)

type ZcrRPC struct {
	*btc.BitcoinRPC
}

func NewZcrRPC(config json.RawMessage, pushHandler func(bchain.NotificationType)) (bchain.BlockChain, error) {
	b, err := btc.NewBitcoinRPC(config, pushHandler)
	if err != nil {
		return nil, err
	}
	s := &ZcrRPC{
		b.(*btc.BitcoinRPC),
	}
	s.RPCMarshaler = btc.JSONMarshalerV1{}
	s.ChainConfig.SupportsEstimateSmartFee = false
	return s, nil
}

func (b *ZcrRPC) Initialize() error {
	ci, err := b.GetChainInfo()
	if err != nil {
		return err
	}
	chainName := ci.Chain
	params := GetChainParams(chainName)
	b.Parser = NewZcrParser(params, b.ChainConfig)
	b.Testnet = false
	b.Network = "livenet"
	glog.Info("rpc: block chain ", params.Name)
	return nil
}

func (b *ZcrRPC) GetBlock(hash string, height uint32) (*bchain.Block, error) {
	var err error
	if hash == "" && height > 0 {
		hash, err = b.GetBlockHash(height)
		if err != nil {
			return nil, err
		}
	}
	glog.V(1).Info("rpc: getblock (verbosity=1) ", hash)
	res := btc.ResGetBlockThin{}
	req := btc.CmdGetBlock{Method: "getblock"}
	req.Params.BlockHash = hash
	req.Params.Verbosity = 1
	err = b.Call(&req, &res)
	if err != nil {
		return nil, errors.Annotatef(err, "hash %v", hash)
	}
	if res.Error != nil {
		return nil, errors.Annotatef(res.Error, "hash %v", hash)
	}
	txs := make([]bchain.Tx, 0, len(res.Result.Txids))
	for _, txid := range res.Result.Txids {
		tx, err := b.GetTransaction(txid)
		if err != nil {
			if err == bchain.ErrTxNotFound {
				glog.Errorf("rpc: getblock: skipping transanction in block %s due error: %s", hash, err)
				continue
			}
			return nil, err
		}
		txs = append(txs, *tx)
	}
	block := &bchain.Block{
		BlockHeader: res.Result.BlockHeader,
		Txs:         txs,
	}
	return block, nil
}

func (b *ZcrRPC) GetTransactionForMempool(txid string) (*bchain.Tx, error) {
	return b.GetTransaction(txid)
}