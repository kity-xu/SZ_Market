// 现金流量表
package io_finchina

type Cashflow struct {
	FinChinaCashflow
}

//------------------------------------------------------------------------------
func NewCashflow() *Cashflow {
	return &Cashflow{}
}

func (this *Cashflow) GetList(compocode string, listdate string, report_data_type int, per_page int, page int) ([]Cashflow, error) {
	return NewFinChinaCashflow().getCashflowList(compocode, listdate, report_data_type, per_page, page)
}

//------------------------------------------------------------------------------
type FinChinaCashflow struct {
	TQ_FIN_PROCFSTATEMENTNEW
}

func NewFinChinaCashflow() *FinChinaCashflow {
	return &FinChinaCashflow{}
}

func (this *FinChinaCashflow) getCashflowList(compcode string, listdate string, report_data_type int, per_page int, page int) ([]Cashflow, error) {
	var (
		slidb []TQ_FIN_PROCFSTATEMENTNEW
		len1  int
		err   error
	)
	sli := make([]Cashflow, 0, per_page)

	slidb, err = NewTQ_FIN_PROCFSTATEMENTNEW().GetList(compcode, listdate, report_data_type, per_page, page)
	if err != nil {
		return nil, err
	}
	if len1 = len(slidb); 0 == len1 {
		return sli, nil
	}

	for _, v := range slidb {
		one := Cashflow{
			FinChinaCashflow: FinChinaCashflow{
				TQ_FIN_PROCFSTATEMENTNEW: v,
			},
		}
		sli = append(sli, one)
	}

	return sli, nil
}
