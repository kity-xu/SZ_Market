package servers

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"haina.com/market/hqinit/config"
	"haina.com/market/hqinit/models/fcmysql"
	tool "haina.com/market/hqtools/printfiletools/util"
	"haina.com/share/logging"
	. "haina.com/share/models"
)

type StockBlockXML struct {
	Model `db:"-"`
}

func NewStockBlockXML() *StockBlockXML {
	return &StockBlockXML{
		Model: Model{},
	}
}

type Femoral struct {
	XMLName xml.Name  `xml:"Plate"`
	Femor   []Servers `xml:"Servers"`
}

type Servers struct {
	XMLName   xml.Name `xml:"servers"`
	BOARDCODE string   `xml:"BOARDCODE,attr"` // 一级板块代码
	Name      string   `xml:"Name,attr"`      // 一级板块名称
	Svs       []server `xml:"server"`
}

type server struct {
	Keycode string    `xml:"Keycode"` // 二级板块代码
	Keyname string    `xml:"Keyname"` // 二级板块名称
	NPreCPx string    `xml:"NPreCPx"` // 板块昨收价
	SerInfo []SerINfo `xml:"SerInfo"`
}

// 最底层详细信息
type SerINfo struct {
	Nsid string `xml:"nsid"`
	Name string `xml:"name"`
}

// 板块K线结构
type BlockStruct struct {
	NSID     int32
	NTime    int32
	NPreCPx  int32
	NOpenPx  int32
	NHighPx  int32
	NLowPx   int32
	NLastPx  int32 // K线最新价 收盘价
	LlVolume int64
	LlValue  int64
	NAvgPx   uint32
}

var BLOCKSIZE int

func (this *StockBlockXML) CreateStockblockXML(cfg *config.AppConfig) {
	logging.Info("stockblock xml begin ...")

	// 查询一级大板块信息
	boar1j, err := fcmysql.NewTQ_COMP_BOARDMAP().GetBoardmapList()
	if err != nil {
		logging.Debug("mysql 1j", err)
	}

	var fem Femoral
	var ser []Servers
	for _, boar1ji := range boar1j {

		var v Servers
		v.BOARDCODE = boar1ji.BOARDCODE.String
		v.Name = boar1ji.BOARDNAME.String
		// 查询二级板块信息
		boar2j, err := fcmysql.NewTQ_COMP_BOARDMAP().GetBoardmap2List(boar1ji.BOARDCODE.String)
		if err != nil {
			logging.Debug("mysql 2j", err)
		}

		var servle []server

		for _, boar2ji := range boar2j {

			var serv server

			serv.Keycode = "81" + strings.Replace(string(boar2ji.KEYCODE.String), "CN", "", -1)
			serv.Keyname = boar2ji.KEYNAME.String

			//---------------------------------------------------------处理板块昨收价
			upath := cfg.File.BlockKLinePath + serv.Keycode

			// 打开文件
			file, err := tool.OpenFile(upath)
			if err != nil {
				serv.NPreCPx = "10000000"
			} else {

				var bkline BlockStruct
				BLOCKSIZE = binary.Size(&bkline)
				var bklines []BlockStruct
				for {
					des := make([]byte, BLOCKSIZE)

					num, err := tool.ReadFiles(file, des)
					if err != nil {
						if err == io.EOF { //读到了文件末尾

							break
						}
						logging.Error("Read file error...%v", err.Error())
						return
					}

					if num < BLOCKSIZE && 0 < num {
						logging.Error("Stock struct size error... or hqtools write file error")
						return
					}
					buffer := bytes.NewBuffer(des)
					binary.Read(buffer, binary.LittleEndian, &bkline)
					bklines = append(bklines, bkline)
				}
				if len(bklines) > 0 {
					zpre := strconv.Itoa(int(bklines[len(bklines)-1].NLastPx))
					serv.NPreCPx = zpre
				} else {
					serv.NPreCPx = "10000000"
				}

				file.Close()
			}

			//----------------------------------------------------------

			// 根据KeyCode查询ComCODE
			boarinfo, err := fcmysql.NewTQ_COMP_BOARDMAP().GetBoardmapInfoList(boar2ji.KEYCODE.String)
			if err != nil {
				logging.Info("二级%v", err)
			}
			var secstr = ""
			for _, item := range boarinfo {
				secstr += "'" + item.SECODE.String + "',"
			}
			secstr = secstr[0 : len(secstr)-1]
			// 根据ComCode查询证券信息
			stcInfo, err := fcmysql.NewFcSecuNameTab().GetComCodeList(secstr)
			var serl []SerINfo
			for _, sitem := range stcInfo {
				var si SerINfo
				if sitem.EXCHANGE.String == "001002" {
					si.Nsid = "100" + sitem.SYMBOL.String
				}
				if sitem.EXCHANGE.String == "001003" {
					si.Nsid = "200" + sitem.SYMBOL.String
				}
				si.Name = sitem.SENAME.String
				serl = append(serl, si)
			}

			serv.SerInfo = serl
			servle = append(servle, serv)
		}
		v.Svs = servle
		ser = append(ser, v)
	}
	fem.Femor = ser
	output, err := xml.MarshalIndent(fem, "  ", "    ")

	if err != nil {
		logging.Info("error: %v\n", err)
	}

	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, output...)

	ioutil.WriteFile(cfg.File.Securitiesplate, xmlOutPutData, os.ModeAppend) // 服务器用
	//ioutil.WriteFile(cfg.File.Securitiesplate, xmlOutPutData, os.ModeAppend)
	logging.Info("stockblock xml end ...")
}
