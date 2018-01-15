// 资产负债表
package io_finchina

type Liabilities struct {
	FinChinaLiabilities
}

//------------------------------------------------------------------------------
func NewLiabilities() *Liabilities {
	return &Liabilities{}
}

func (this *Liabilities) GetList(sid int, report_type int, per_page int, page int) ([]Liabilities, error) {
	return NewFinChinaLiabilities().getLiabilitiesList(sid, report_type, per_page, page)
}

//------------------------------------------------------------------------------

type FinChinaLiabilities struct {
	TQ_FIN_PROBALSHEETNEW
}

func NewFinChinaLiabilities() *FinChinaLiabilities {
	return &FinChinaLiabilities{}
}

func (this *FinChinaLiabilities) getLiabilitiesList(sid int, report_data_type int, per_page int, page int) ([]Liabilities, error) {
	var (
		slidb []TQ_FIN_PROBALSHEETNEW
		len1  int
		err   error
	)
	sli := make([]Liabilities, 0, per_page)

	slidb, err = NewTQ_FIN_PROBALSHEETNEW().GetList(sid, report_data_type, per_page, page)
	if err != nil {
		return nil, err
	}
	if len1 = len(slidb); 0 == len1 {
		return sli, nil
	}

	for _, v := range slidb {
		one := Liabilities{
			FinChinaLiabilities: FinChinaLiabilities{
				TQ_FIN_PROBALSHEETNEW: v,
			},
		}
		sli = append(sli, one)
	}

	return sli, nil
}
