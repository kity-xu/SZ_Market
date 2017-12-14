// 利润表
package io_finchina

type Profits struct {
	FinChinaProfits
}

//------------------------------------------------------------------------------

func NewProfits() *Profits {
	return &Profits{}
}

func (this *Profits) GetList(scode string, market string, report_data_type int, per_page int, page int) ([]Profits, error) {
	return NewFinChinaProfits().getProfitsList(scode, market, report_data_type, per_page, page)
}

//------------------------------------------------------------------------------
type FinChinaProfits struct {
	TQ_FIN_PROINCSTATEMENTNEW
}

func NewFinChinaProfits() *FinChinaProfits {
	return &FinChinaProfits{}
}

func (this *FinChinaProfits) getProfitsList(scode string, market string, report_data_type int, per_page int, page int) ([]Profits, error) {
	var (
		slidb []TQ_FIN_PROINCSTATEMENTNEW
		len1  int
		err   error
	)
	sli := make([]Profits, 0, per_page)

	slidb, err = NewTQ_FIN_PROINCSTATEMENTNEW().GetList(scode, market, report_data_type, per_page, page)
	if err != nil {
		return nil, err
	}
	if len1 = len(slidb); 0 == len1 {
		return sli, nil
	}

	for _, v := range slidb {
		one := Profits{
			FinChinaProfits: FinChinaProfits{
				TQ_FIN_PROINCSTATEMENTNEW: v,
			},
		}
		sli = append(sli, one)
	}

	return sli, nil
}
