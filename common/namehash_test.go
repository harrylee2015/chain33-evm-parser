package common

import (
	vm "github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"
)

func Test_ParseChinese(t *testing.T){
	input := "苏州.eth"
	hash:=NameHash(input)
	t.Log(hash.String())
	input = "李海磊"
	t.Log(NameHash(input).String())
	t.Log(crypto.Keccak256Hash([]byte(input)).String())

	input ="李海磊.苏州.eth"
	hash=NameHash(input)
	t.Log(hash.String())

}

func Test_labes(t *testing.T) {
	//input := "苏州.eth"
	//hash1:=NameHash(input)
	//hash3:="0xb87c6b0dfe56a8f5a41d21b898570673b713f127ef93da333a2044f90bb8a3e8"
	//
	//t.Log(crypto.Keccak256Hash(hash1.Bytes(),crypto.Keccak256Hash([]byte("李海磊")).Bytes()).String())
	//t.Log(hash3)

	hash := common.Hash{}
	t.Log(hash.String())
	t.Log(NameHash("eth").String())
	t.Log(crypto.Keccak256Hash([]byte("eth")).String())
	t.Log(crypto.Keccak256Hash(hash.Bytes(),crypto.Keccak256Hash([]byte("eth")).Bytes()).String())
	t.Log(crypto.Keccak256Hash([]byte("wallet")).String())
	t.Log(NameHash("wallet.eth").String())
	t.Log("------------------------------------")
	t.Log(crypto.Keccak256Hash([]byte("yuan")).String())
	t.Log(NameHash("yuan").String())
	//t.Log("------------------------------------")
	//t.Log(NameHash("16525216.yuan").String())
	//t.Log(NameHash("litian1.yuan").String())
	//t.Log(NameHash("1723231.yuan").String())
	//t.Log(NameHash("bty").String())
	//t.Log(crypto.Keccak256Hash([]byte("bty")).String())
	//t.Log(crypto.Keccak256Hash(hash.Bytes(),crypto.Keccak256Hash([]byte("bty")).Bytes()).String())
	//t.Log(crypto.Keccak256Hash([]byte("wallet")).String())
	//t.Log(NameHash("wallet.bty").String())
	//t.Log(crypto.Keccak256Hash([]byte("harry")).String())
	//t.Log(NameHash("harry.wallet.bty").String())
	//t.Log("------------------------------------")
	//t.Log(NameHash("bty.litian.eth").String())
	//t.Log(NameHash("BTY.litian.eth").String())
	t.Log("------------------------------------")
	//
	//t.Log(NameHash("litian.eth").String())

	t.Log(crypto.Keccak256Hash([]byte("test")).String())
	t.Log(NameHash("test.bty").String())
	t.Log("------------------------------------")
	t.Log(crypto.Keccak256Hash([]byte("test1")).String())
	t.Log(NameHash("test1.bty").String())
	t.Log("------------------------------------")
	t.Log(crypto.Keccak256Hash([]byte("test2")).String())
	t.Log(NameHash("test2.bty").String())
	t.Log(NameHash("test2.bty").Big())
	t.Log(NameHash("wallet.bty").String())
	t.Log(NameHash("emoemoemoemo.bty").String())
	t.Log(NameHash("emoemoemoemo.bty").Big())
	t.Log(NameHash("test499.bty").String())
	t.Log(crypto.Keccak256Hash([]byte("addr")).String())
	t.Log(NameHash("addr").String())
	t.Log(NameHash("test102.bty").String())
	//number :=crypto.Keccak256Hash([]byte("wallet")).Big()
	//t.Log(crypto.Keccak256Hash([]byte("wallet")).Big())
	//t.Log(number)
	//t.Log(number.Uint64())
	//t.Log(common.BigToHash(number))
	t.Log("------------------------------------")

	//事件参数哈希计算
	t.Log(common.BigToHash(big.NewInt(123)).Hex())
	//反向解析
	hexStr :=common.BigToHash(big.NewInt(123)).Hex()
	t.Log(common.HexToHash(hexStr).Big())
	//14KEKbYtKKQm4wMthSK9J4La4nAiidGozt
	address:=vm.StringToAddress("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
	t.Log(common.BytesToHash(vm.FromHex(address.ToHash160().Hex())).Hex())
	//反向解析
	hexStr=common.BytesToHash(vm.FromHex(address.ToHash160().Hex())).Hex()
	addr:=vm.HexToAddress(hexStr)
	t.Log(addr.ToAddress())
	//0x245Afbf176934ccDD7ca291a8DddaA13c8184822
    t.Log(vm.HexToAddress("0x245Afbf176934ccDD7ca291a8DddaA13c8184822").ToAddress())
	//事件哈希计算
	t.Log(crypto.Keccak256Hash([]byte("BlindBoxMinted(uint256,address)")).String())
	t.Log(crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)")).String())


	//t.Log(common.BytesToHash().Hex())
	t.Log(crypto.Keccak256Hash([]byte("bty")).String())
	t.Log(crypto.Keccak256Hash([]byte("eth")).String())
	t.Log(crypto.Keccak256Hash([]byte("wallet2015")).String())
	t.Log(crypto.Keccak256Hash([]byte("wallet2015")).Big())
	//t.Log(NameHash("1725.eth"))

	//t.Log(crypto.Keccak256Hash([]byte("alice")).String())
	//t.Log(NameHash("alice.wallet.bty").String())
	//
	//
	t.Log("------------------------------------")
	t.Log(crypto.Keccak256Hash([]byte("yuan")).String())
	t.Log(NameHash("yuan").String())
	//t.Log(crypto.Keccak256Hash(hash.Bytes(),crypto.Keccak256Hash([]byte("eth")).Bytes()).String())
	//t.Log(crypto.Keccak256Hash([]byte("test")).String())
	//t.Log(common.Address{}.Bytes())

	//t.Log(NameHash("test.eth").String())
	//t.Log(crypto.Keccak256Hash([]byte("alice")).String())
	//t.Log(NameHash("alice.test.eth").String())
	//t.Log("------------------------------------")
	//t.Log(NameHash("yuan1.eth").String())
	//t.Log(crypto.Keccak256Hash([]byte("bty")).String())
	//t.Log(NameHash("bty.yuan1.eth").String())

	//t.Log(NameHash("wallet.eth").String())
	////t.Log(crypto.Keccak256Hash(NameHash("eth").Bytes(),crypto.Keccak256Hash([]byte("wallet")).Bytes()).String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("lihailei")).String())
	//t.Log(NameHash("lihailei.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("lei")).String())
	//t.Log(NameHash("lei.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("alice")).String())
	//t.Log(NameHash("alice.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("bob")).String())
	//t.Log(NameHash("bob.eth").String())
	//
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("a")).String())
	//t.Log(NameHash("a.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("b")).String())
	//t.Log(NameHash("b.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("c")).String())
	//t.Log(NameHash("c.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("d")).String())
	//t.Log(NameHash("d.eth").String())
	//
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("a")).String())
	//t.Log(NameHash("e.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("b")).String())
	//t.Log(NameHash("f.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("c")).String())
	//t.Log(NameHash("g.eth").String())
	//t.Log("------------------------------------")
	//t.Log(crypto.Keccak256Hash([]byte("d")).String())
	//t.Log(NameHash("h.eth").String())
}
