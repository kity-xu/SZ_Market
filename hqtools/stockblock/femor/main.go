package main

import (
	"ProtocolBuffer/projects/hqinit/go/protocol"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	"github.com/golang/protobuf/proto"
	"haina.com/market/hqinit/controllers"
	"haina.com/market/hqtools/stockblock/fcmysql"
	"haina.com/share/lib"
	"haina.com/share/logging"
)

type Femoral struct {
	XMLName xml.Name  `xml:"Plate"`
	Femor   []Servers `xml:"Servers"`
}

type Servers struct {
	XMLName   xml.Name `xml:"servers"`
	BOARDCODE string   `xml:"BOARDCODE,attr"`
	Name      string   `xml:"Name,attr"`
	Svs       []server `xml:"server"`
}

type server struct {
	Keycode string    `xml:"Keycode"`
	Keyname string    `xml:"Keyname"`
	SerInfo []SerINfo `xml:"SerInfo"`
}

// 最底层详细信息
type SerINfo struct {
	Nsid string `xml:"nsid"`
	Name string `xml:"name"`
}

func main() {
	logging.Info("begin ...")
	c, errr := redis.Dial("tcp", "47.94.16.69:61380")
	c.Send("AUTH", "8dc40c2c4598ae5a")
	if errr != nil {
		logging.Info("redis conn error %v", errr)
	}
	//lib.CheckDir("/opt/hqfile/block/")
	lib.CheckDir("E:/hqfile/")
	//_, err := controllers.OpenFile("/opt/hqfile/block/securitiesplate.xml")
	_, err := controllers.OpenFile("E:/hqfile/securitiesplate.xml")

	if err != nil {
		logging.Info("创建文件失败！")
	}
	// 服务器用
	//conn, err := dbr.Open("mysql", "finchina:finchina@tcp(172.16.1.60:3306)/finchina?charset=utf8", nil)
	conn, err := dbr.Open("mysql", "finchina:finchina@tcp(114.55.105.11:3306)/finchina?charset=utf8", nil)
	if err != nil {
		logging.Debug("mysql onn", err)
	}
	sess := conn.NewSession(nil)

	// 查询一级大板块信息
	boar1j, err := new(fcmysql.TQ_COMP_BOARDMAP).GetBoardmapList(sess)
	if err != nil {
		logging.Debug("mysql 1j", err)
	}

	var fem Femoral
	var ser []Servers
	//	var kil []protocol.BlackInfoTable
	for _, boar1ji := range boar1j {

		var v Servers
		v.BOARDCODE = boar1ji.BOARDCODE.String
		v.Name = boar1ji.BOARDNAME.String
		// 查询二级板块信息
		boar2j, err := new(fcmysql.TQ_COMP_BOARDMAP).GetBoardmap2List(sess, boar1ji.BOARDCODE.String)
		if err != nil {
			logging.Debug("mysql 2j", err)
		}

		var servle []server

		for _, boar2ji := range boar2j {
			logging.Info("计=====%v", boar1ji.BOARDCODE)
			var bi protocol.StockBlock
			bi.BoardCode = boar1ji.BOARDCODE.String
			bi.BoardName = boar1ji.BOARDNAME.String
			bi.KeyCode = boar2ji.KEYCODE.String
			bi.KeyName = boar2ji.KEYNAME.String

			//			nsid, err := strconv.Atoi(boar2ji.KEYCODE.String)
			//			if err != nil {
			//				logging.Info("证券id类型转换 error:%v", err)
			//			}
			//			bi.Snid = int32(nsid)
			//			bi.SName = sitem.SENAME.String
			//转PB

			datapb, err := proto.Marshal(&bi)
			if err != nil {
				logging.Error("Encode protocbuf of week Line error...%v", err.Error())
				return
			}
			// 所有板块存redis
			key := fmt.Sprintf("hq:init:bk:set:%v", boar2ji.KEYCODE.String)
			c.Do("SET", key, datapb)

			// 所有地区板块
			if boar1ji.BOARDCODE.String == "1102" {
				keyd := fmt.Sprintf("hq:init:bk:district:%v", boar2ji.KEYCODE.String)
				c.Do("SET", keyd, datapb)
			}
			// 所有概念板块
			if boar1ji.BOARDCODE.String == "1105" {
				keyc := fmt.Sprintf("hq:init:bk:concept:%v", boar2ji.KEYCODE.String)
				c.Do("SET", keyc, datapb)
			}
			// 所有行业板块
			if boar1ji.BOARDCODE.String == "1109" {
				keyv := fmt.Sprintf("hq:init:bk:vocation:%v", boar2ji.KEYCODE.String)
				c.Do("SET", keyv, datapb)
			}

			var serv server
			serv.Keycode = boar2ji.KEYCODE.String
			serv.Keyname = boar2ji.KEYNAME.String

			//			if boar2ji.BOARDCODE == boar1ji.BOARDCODE {

			// 根据KeyCode查询ComCODE
			boarinfo, err := new(fcmysql.TQ_COMP_BOARDMAP).GetBoardmapInfoList(sess, boar2ji.KEYCODE.String)
			if err != nil {
				logging.Info("二级%v", err)
			}
			var secstr = ""
			for _, item := range boarinfo {
				secstr += "'" + item.SECODE.String + "',"
			}
			secstr = secstr[0 : len(secstr)-1]
			// 根据ComCode查询证券信息
			stcInfo, err := new(fcmysql.TQ_OA_STCODE).GetComCodeList(sess, secstr)
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

				var bi protocol.StockBlock

				bi.BoardCode = boar1ji.BOARDCODE.String
				bi.BoardName = boar1ji.BOARDNAME.String
				bi.KeyCode = boar2ji.KEYCODE.String
				bi.KeyName = boar2ji.KEYNAME.String

				nsid, err := strconv.Atoi(si.Nsid)
				if err != nil {
					logging.Info("证券id类型转换 error:%v", err)
				}
				bi.Snid = int32(nsid)
				bi.SName = sitem.SENAME.String

				//转PB
				data, err := proto.Marshal(&bi)
				if err != nil {
					logging.Error("Encode protocbuf of week Line error...%v", err.Error())
					return
				}

				//key := fmt.Sprintf("hq:init:bk:%s:%s:%s", boar1ji.BOARDCODE.String, boar2ji.KEYCODE.String, si.Nsid)
				//c.Do("SET", key, data)

				//kil = append(kil, bi)
			}
			//				serv.SerInfo = serl
			//v.Svs = append(v.Svs, serv)
			serv.SerInfo = serl
			servle = append(servle, serv)
			//			}
		}
		v.Svs = servle
		ser = append(ser, v)
	}
	fem.Femor = ser
	output, err := xml.MarshalIndent(fem, "  ", "    ")

	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// 写入redis
	//	client.Set("hq:black:info", kil)
	//加入XML头
	headerBytes := []byte(xml.Header)
	//拼接XML头和实际XML内容
	xmlOutPutData := append(headerBytes, output...)

	ioutil.WriteFile("/opt/hqfile/block/securitiesplate.xml", xmlOutPutData, os.ModeAppend) // 服务器用
	//	ioutil.WriteFile("E:/hqfile/securitiesplate.xml", xmlOutPutData, os.ModeAppend)
	logging.Info("end ...")
}
