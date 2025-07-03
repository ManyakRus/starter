package postgres_pgtype

//
//import (
//	"math/big"
//)
//
//// PostgreSQL internal numeric storage uses 16-bit "digits" with base of 10,000
//const nbase = 10000
//
//const (
//	pgNumericNaN     = 0x00000000c0000000
//	pgNumericNaNSign = 0xc000
//
//	pgNumericPosInf     = 0x00000000d0000000
//	pgNumericPosInfSign = 0xd000
//
//	pgNumericNegInf     = 0x00000000f0000000
//	pgNumericNegInfSign = 0xf000
//)
//
//var big0 *big.Int = big.NewInt(0)
//var big1 *big.Int = big.NewInt(1)
//var big10 *big.Int = big.NewInt(10)
//var big100 *big.Int = big.NewInt(100)
//var big1000 *big.Int = big.NewInt(1000)
//
//var bigNBase *big.Int = big.NewInt(nbase)
//var bigNBaseX2 *big.Int = big.NewInt(nbase * nbase)
//var bigNBaseX3 *big.Int = big.NewInt(nbase * nbase * nbase)
//var bigNBaseX4 *big.Int = big.NewInt(nbase * nbase * nbase * nbase)
