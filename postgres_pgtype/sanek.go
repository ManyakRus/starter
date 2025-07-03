package postgres_pgtype

import (
	"reflect"
)

// getInterfaceName - возвращает имя типа интерфейса
func getInterfaceName(v interface{}) string {
	return reflect.TypeOf(v).String()
}

//func toBigInt(n *pgtype.Numeric) (*big.Int, error) {
//	if n.Exp == 0 {
//		return n.Int, nil
//	}
//
//	num := &big.Int{}
//	num.Set(n.Int)
//	if n.Exp > 0 {
//		mul := &big.Int{}
//		mul.Exp(big10, big.NewInt(int64(n.Exp)), nil)
//		num.Mul(num, mul)
//		return num, nil
//	}
//
//	div := &big.Int{}
//	div.Exp(big10, big.NewInt(int64(-n.Exp)), nil)
//	remainder := &big.Int{}
//	num.DivMod(num, div, remainder)
//	if remainder.Cmp(big0) != 0 {
//		return nil, fmt.Errorf("cannot convert %v to integer", n)
//	}
//	return num, nil
//}
//
//// cardinality returns the number of elements in an array of dimensions size.
//func cardinality(dimensions []pgtype.ArrayDimension) int {
//	if len(dimensions) == 0 {
//		return 0
//	}
//
//	elementCount := int(dimensions[0].Length)
//	for _, d := range dimensions[1:] {
//		elementCount *= int(d.Length)
//	}
//
//	return elementCount
//}
