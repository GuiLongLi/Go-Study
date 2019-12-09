package localSpike

type LocalSpike struct {
	LocalInStock int64 //本地库存
	LocalSalesVolume int64 //本地销量
}

//本地扣库存，
//销量+1
//返回 销量是否小于等于库存的 bool 值
func (spike *LocalSpike) LocalDeductionStock() bool{
	spike.LocalSalesVolume = spike.LocalSalesVolume + 1
	return spike.LocalSalesVolume <= spike.LocalInStock
}