package zcr

import (
	"github.com/mosqueiro/blockbook/bchain"
	"github.com/mosqueiro/blockbook/bchain/coins/btc"

	"github.com/martinboehm/btcd/wire"
	"github.com/martinboehm/btcutil/chaincfg"
)

const (
	MainnetMagic wire.BitcoinNet = 0xCC645C66
)

var (
	MainNetParams chaincfg.Params
)

func init() {
	MainNetParams = chaincfg.MainNetParams
	MainNetParams.Net = MainnetMagic

	MainNetParams.PubKeyHashAddrID = []byte{142}
	MainNetParams.ScriptHashAddrID = []byte{145}
}

type ZcrParser struct {
	*btc.BitcoinParser
	baseparser *bchain.BaseParser
}

func NewZcrParser(params *chaincfg.Params, c *btc.Configuration) *ZcrParser {
	return &ZcrParser{BitcoinParser: btc.NewBitcoinParser(params, c), baseparser: &bchain.BaseParser{}}
}

func GetChainParams(chain string) *chaincfg.Params {
	if !chaincfg.IsRegistered(&MainNetParams) {
		err := chaincfg.Register(&MainNetParams)
		if err != nil {
			panic(err)
		}
	}
	return &MainNetParams
}

func (p *ZcrParser) PackTx(tx *bchain.Tx, height uint32, blockTime int64) ([]byte, error) {
	return p.baseparser.PackTx(tx, height, blockTime)
}

func (p *ZcrParser) UnpackTx(buf []byte) (*bchain.Tx, uint32, error) {
	return p.baseparser.UnpackTx(buf)
}