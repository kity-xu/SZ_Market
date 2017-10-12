package servers

import (
	"encoding/xml"
	"io/ioutil"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"haina.com/market/hqinit/config"
	"haina.com/market/hqinit/models/fcmysql"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

type IndexComponentXML struct {
	Model `db:"-"`
}

func NewIndexComponentXML() *IndexComponentXML {
	return &IndexComponentXML{
		Model: Model{},
	}
}

type StockInfo struct {
	XMLName xml.Name `xml:"stockinfo"`
	NSID    string   `xml:"nsid,attr"`
}
type IndexInfo struct {
	XMLName    xml.Name `xml:"indexinfo"`
	NSID       string   `xml:"nsid,attr"`
	Value      string   `xml:"name,attr"`
	StockInfos []StockInfo
}
type IndComponent struct {
	XMLName xml.Name `xml:"indexcomponent"`
	Indexs  []IndexInfo
}

func (this *IndexComponentXML) CreateIndexComponentXML(cfg *config.AppConfig) {

	logging.Info("create indexcomponent xml begin ...")

	// 查询所有沪深指数
	indl, err := fcmysql.NewFcSecuNameTab().GetExponentList()
	if err != nil {
		logging.Debug("select index error %v", err)
	}
	var ict IndComponent
	// 根据查询到的指数信息 SECODE 查询指数成分股
	for _, indlit := range indl {
		// 根据SECODE查询指数成分股
		stkl, err := fcmysql.NewTQ_IX_COMP().GetIndexStockL(indlit.SECODE.String)
		if err != nil {
			logging.Debug("select index -> stockinfo error %v", err)
		}
		var idf IndexInfo
		// 遍历指数下成分股
		if len(stkl) <= 0 {
			continue
		}
		for _, stki := range stkl {
			var sif StockInfo
			//			if indlit.EXCHANGE.String == "001002" {
			//				sif.NSID = "100" + stki.SAMPLECODE.String
			//			}
			//			if indlit.EXCHANGE.String == "001003" {
			//				sif.NSID = "200" + stki.SAMPLECODE.String
			//			}

			if len(stki.SAMPLECODE.String) > 0 {
				chara := stki.SAMPLECODE.String[:1]
				if chara == "6" {
					sif.NSID = "100" + stki.SAMPLECODE.String
				} else if chara == "3" || chara == "0" {
					sif.NSID = "200" + stki.SAMPLECODE.String
				} else if chara == "9" {
					continue
				}
			}

			idf.StockInfos = append(idf.StockInfos, sif)
		}
		if indlit.EXCHANGE.String == "001002" {
			idf.NSID = "100" + indlit.SYMBOL.String
		}
		if indlit.EXCHANGE.String == "001003" {
			idf.NSID = "200" + indlit.SYMBOL.String
		}
		idf.Value = indlit.SESNAME.String
		ict.Indexs = append(ict.Indexs, idf)
	}

	data, err := xml.Marshal(&ict)
	if err != nil {
		logging.Info(" xml Marshal error %v", err)
	}
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, data...)

	ioutil.WriteFile(cfg.File.IndexComponentPath, xmlOutPutData, os.ModeAppend)
	logging.Info("create indexcomponent xml end ...")

}
