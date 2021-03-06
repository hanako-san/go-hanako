package repo

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type DB struct {
	data map[string]*Entry
}

type Entry struct {
	Name string
	Observes []Observation
}

type Observation struct {
	Date time.Time
	Count int
}

var kRequestMap = map[string]string{
	"Kanto": "__EVENTTARGET=&__EVENTARGUMENT=&__LASTFOCUS=&__VIEWSTATE=%2FwEPDwUKMjA4NDcxNzAyMw9kFgICAQ9kFhgCBQ8QDxYCHgtfIURhdGFCb3VuZGdkEBUBBDIwMjAVAQQyMDIwFCsDAWdkZAIHDxAPFgIfAGdkEBUMATEBMgEzATQBNQE2ATcBOAE5AjEwAjExAjEyFQwBMQEyATMBNAE1ATYBNwE4ATkCMTACMTECMTIUKwMMZ2dnZ2dnZ2dnZ2dnZGQCCQ8QDxYCHwBnZBAVHwExATIBMwE0ATUBNgE3ATgBOQIxMAIxMQIxMgIxMwIxNAIxNQIxNgIxNwIxOAIxOQIyMAIyMQIyMgIyMwIyNAIyNQIyNgIyNwIyOAIyOQIzMAIzMRUfATEBMgEzATQBNQE2ATcBOAE5AjEwAjExAjEyAjEzAjE0AjE1AjE2AjE3AjE4AjE5AjIwAjIxAjIyAjIzAjI0AjI1AjI2AjI3AjI4AjI5AjMwAjMxFCsDH2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dkZAILDxAPFgIfAGdkEBUYATEBMgEzATQBNQE2ATcBOAE5AjEwAjExAjEyAjEzAjE0AjE1AjE2AjE3AjE4AjE5AjIwAjIxAjIyAjIzAjI0FRgBMQEyATMBNAE1ATYBNwE4ATkCMTACMTECMTICMTMCMTQCMTUCMTYCMTcCMTgCMTkCMjACMjECMjICMjMCMjQUKwMYZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZGQCDQ8QDxYCHwBnZBAVAQQyMDIwFQEEMjAyMBQrAwFnZGQCDw8QDxYCHwBnZBAVDAExATIBMwE0ATUBNgE3ATgBOQIxMAIxMQIxMhUMATEBMgEzATQBNQE2ATcBOAE5AjEwAjExAjEyFCsDDGdnZ2dnZ2dnZ2dnZ2RkAhEPEA8WAh8AZ2QQFR8BMQEyATMBNAE1ATYBNwE4ATkCMTACMTECMTICMTMCMTQCMTUCMTYCMTcCMTgCMTkCMjACMjECMjICMjMCMjQCMjUCMjYCMjcCMjgCMjkCMzACMzEVHwExATIBMwE0ATUBNgE3ATgBOQIxMAIxMQIxMgIxMwIxNAIxNQIxNgIxNwIxOAIxOQIyMAIyMQIyMgIyMwIyNAIyNQIyNgIyNwIyOAIyOQIzMAIzMRQrAx9nZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZGQCEw8QDxYCHwBnZBAVGAExATIBMwE0ATUBNgE3ATgBOQIxMAIxMQIxMgIxMwIxNAIxNQIxNgIxNwIxOAIxOQIyMAIyMQIyMgIyMwIyNBUYATEBMgEzATQBNQE2ATcBOAE5AjEwAjExAjEyAjEzAjE0AjE1AjE2AjE3AjE4AjE5AjIwAjIxAjIyAjIzAjI0FCsDGGdnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZ2RkAhUPEA8WBh4ORGF0YVZhbHVlRmllbGQFBVZhbHVlHg1EYXRhVGV4dEZpZWxkBQRUZXh0HwBnZBAVBgzmnbHljJflnLDln58M6Zai5p2x5Zyw5Z%2BfDOS4remDqOWcsOWfnwzplqLopb%2FlnLDln58V5Lit5Zu944O75Zub5Zu95Zyw5Z%2BfDOS5neW3nuWcsOWfnxUGAjAyAjAzAjA1AjA2AjA3AjA4FCsDBmdnZ2dnZxYBAgFkAhsPEA8WBh8BBQVWYWx1ZR8CBQRUZXh0HwBnZBAVFEw8Zm9udCBjb2xvcj0iYmx1ZSI%2BKOiMqOWfjuecjCnmsLTmiLjnn7Plt53kuIDoiKznkrDlooPlpKfmsJfmuKzlrprlsYA8L2ZvbnQ%2BTDxmb250IGNvbG9yPSJibHVlIj4o6Iyo5Z%2BO55yMKeeLrOeri%2BihjOaUv%2BazleS6uuWbveeri%2BeSsOWig%2BeglOeptuaJgDwvZm9udD46PGZvbnQgY29sb3I9ImJsdWUiPijojKjln47nnIwp5pel56uL5biC5raI6Ziy5pys6YOoPC9mb250PjUo5qCD5pyo55yMKeWuh%2BmDveWuruW4guS4reWkrueUn%2Ba2r%2BWtpue%2FkuOCu%2BODs%2BOCv%2BODvCAo5qCD5pyo55yMKeagg%2BacqOecjOmCo%2BmgiOW6geiIjiYo5qCD5pyo55yMKeaXpeWFieW4guW9ueaJgOesrO%2B8lOW6geiIjkM8Zm9udCBjb2xvcj0iYmx1ZSI%2BKOe%2BpOmmrOecjCnnvqTppqznnIzooZvnlJ%2FnkrDlooPnoJTnqbbmiYA8L2ZvbnQ%2BQDxmb250IGNvbG9yPSJibHVlIj4o576k6aas55yMKemkqOael%2BS%2FneWBpeemj%2BelieS6i%2BWLmeaJgDwvZm9udD4gKOWfvOeOieecjCnjgZXjgYTjgZ%2Fjgb7luILlvbnmiYAmKOWfvOeOieecjCnnhorosLfluILkv53lgaXjgrvjg7Pjgr%2Fjg7waKOWfvOeOieecjCnpo6%2Fog73luILlvbnmiYAxPGZvbnQgY29sb3I9ImJsdWUiPijljYPokYnnnIwp5p2x6YKm5aSn5a2mPC9mb250PkY8Zm9udCBjb2xvcj0iYmx1ZSI%2BKOWNg%2BiRieecjCnljYPokYnnnIznkrDlooPnoJTnqbbjgrvjg7Pjgr%2Fjg7w8L2ZvbnQ%2BTzxmb250IGNvbG9yPSJibHVlIj4o5Y2D6JGJ55yMKeWNsOaXm%2BWBpeW6t%2Bemj%2BelieOCu%2BODs%2BOCv%2BODvOaIkOeUsOaUr%2BaJgDwvZm9udD49PGZvbnQgY29sb3I9ImJsdWUiPijljYPokYnnnIwp5ZCb5rSl5biC57Og55Sw5ris5a6a5bGAPC9mb250Piko5p2x5Lqs6YO9KeadseS6rOmDveWkmuaRqeWwj%2BW5s%2BS%2FneWBpeaJgCko5p2x5Lqs6YO9KeaWsOWuv%2BWMuuW9ueaJgOesrOS6jOWIhuW6geiIjkM8Zm9udCBjb2xvcj0iYmx1ZSI%2BKOelnuWliOW3neecjCnnpZ7lpYjlt53nnIzluoHkuozliIbluoHoiI48L2ZvbnQ%2BVTxmb250IGNvbG9yPSJibHVlIj4o56We5aWI5bed55yMKeW3neW0jueUn%2BWRveenkeWtpuODu%2BeSsOWig%2BeglOeptuOCu%2BODs%2BOCv%2BODvDwvZm9udD5MPGZvbnQgY29sb3I9ImJsdWUiPijnpZ7lpYjlt53nnIwp56We5aWI5bed55yM55Kw5aKD56eR5a2m44K744Oz44K%2F44O8PC9mb250PhUUCDUwODEwMTAwCDUwODEwMjAwCDUwODIwMTAwCDUwOTEwMTAwCDUwOTIwMTAwCDUwOTIwMjAwCDUxMDEwMTAwCDUxMDIwMTAwCDUxMTEwMjAwCDUxMTIwNDAwCDUxMTIwMzAwCDUxMjEwMTAwCDUxMjEwMjAwCDUxMjIwMTAwCDUxMjIwMjAwCDUxMzEwMjAwCDUxMzIwMTAwCDUxNDEwMTAwCDUxNDEwMjAwCDUxNDIwMTAwFCsDFGdnZ2dnZ2dnZ2dnZ2dnZ2dnZ2dnZGQCHQ8PZBYCHgdvbmNsaWNrBRlyZXR1cm4gQ2hhbmdlUGFnZShGb3JtMSk7ZAIfDw9kFgIfAwUZcmV0dXJuIENoYW5nZVBhZ2UoRm9ybTEpO2QYAQUeX19Db250cm9sc1JlcXVpcmVQb3N0QmFja0tleV9fFhUFEUNoZWNrQm94TXN0TGlzdCQwBRFDaGVja0JveE1zdExpc3QkMQURQ2hlY2tCb3hNc3RMaXN0JDIFEUNoZWNrQm94TXN0TGlzdCQzBRFDaGVja0JveE1zdExpc3QkNAURQ2hlY2tCb3hNc3RMaXN0JDUFEUNoZWNrQm94TXN0TGlzdCQ2BRFDaGVja0JveE1zdExpc3QkNwURQ2hlY2tCb3hNc3RMaXN0JDgFEUNoZWNrQm94TXN0TGlzdCQ5BRJDaGVja0JveE1zdExpc3QkMTAFEkNoZWNrQm94TXN0TGlzdCQxMQUSQ2hlY2tCb3hNc3RMaXN0JDEyBRJDaGVja0JveE1zdExpc3QkMTMFEkNoZWNrQm94TXN0TGlzdCQxNAUSQ2hlY2tCb3hNc3RMaXN0JDE1BRJDaGVja0JveE1zdExpc3QkMTYFEkNoZWNrQm94TXN0TGlzdCQxNwUSQ2hlY2tCb3hNc3RMaXN0JDE4BRJDaGVja0JveE1zdExpc3QkMTkFEkNoZWNrQm94TXN0TGlzdCQxOX6KPmKi0aq8Zz3EEGJjBDeqaAtl&__VIEWSTATEGENERATOR=DE1917C9&__EVENTVALIDATION=%2FwEWqQECipaWuwsC0pnr3g0CvfDYuwkC7rOKoA8C77OKoA8C7LOKoA8C7bOKoA8C6rOKoA8C67OKoA8C6LOKoA8C%2BbOKoA8C9rOKoA8C7rPKow8C7rPGow8C7rPCow8C%2BfykvAwC%2BPykvAwC%2B%2FykvAwC%2BvykvAwC%2FfykvAwC%2FPykvAwC%2F%2FykvAwC7vykvAwC4fykvAwC%2BfzkvwwC%2BfzovwwC%2BfzsvwwC%2BfzQvwwC%2BfzUvwwC%2BfzYvwwC%2BfzcvwwC%2BfzAvwwC%2BfyEvAwC%2BfyIvAwC%2BPzkvwwC%2BPzovwwC%2BPzsvwwC%2BPzQvwwC%2BPzUvwwC%2BPzYvwwC%2BPzcvwwC%2BPzAvwwC%2BPyEvAwC%2BPyIvAwC%2B%2FzkvwwC%2B%2FzovwwC2dfPjg0C2NfPjg0C29fPjg0C2tfPjg0C3dfPjg0C3NfPjg0C39fPjg0CztfPjg0CwdfPjg0C2dePjQ0C2deDjQ0C2deHjQ0C2de7jQ0C2de%2FjQ0C2dezjQ0C2de3jQ0C2derjQ0C2dfvjg0C2dfjjg0C2NePjQ0C2NeDjQ0C2NeHjQ0C2Ne7jQ0C2Ne%2FjQ0Cg774pgMCq%2F66zwsCqv66zwsCqf66zwsCqP66zwsCr%2F66zwsCrv66zwsCrf66zwsCvP66zwsCs%2F66zwsCq%2F76zAsCq%2F72zAsCq%2F7yzAsC%2BeK61AkC%2BOK61AkC%2B%2BK61AkC%2BuK61AkC%2FeK61AkC%2FOK61AkC%2F%2BK61AkC7uK61AkC4eK61AkC%2BeL61wkC%2BeL21wkC%2BeLy1wkC%2BeLO1wkC%2BeLK1wkC%2BeLG1wkC%2BeLC1wkC%2BeLe1wkC%2BeKa1AkC%2BeKW1AkC%2BOL61wkC%2BOL21wkC%2BOLy1wkC%2BOLO1wkC%2BOLK1wkC%2BOLG1wkC%2BOLC1wkC%2BOLe1wkC%2BOKa1AkC%2BOKW1AkC%2B%2BL61wkC%2B%2BL21wkC8YeS3gcC8IeS3gcC84eS3gcC8oeS3gcC9YeS3gcC9IeS3gcC94eS3gcC5oeS3gcC6YeS3gcC8YfS3QcC8Yfe3QcC8Yfa3QcC8Yfm3QcC8Yfi3QcC8Yfu3QcC8Yfq3QcC8Yf23QcC8Yey3gcC8Ye%2B3gcC8IfS3QcC8Ife3QcC8Ifa3QcC8Ifm3QcC8Ifi3QcCq6WQZAK7yvKJDAK7ys6JDAK7ysaJDAK7ysKJDAK7yt6JDAK7ypqKDAK2yfrdAQKA6KnRAgKd%2BNSJDAKc%2BNSJDAKb%2BNSJDAKa%2BNSJDAKh%2BNSJDAKg%2BNSJDAKf%2BNSJDAKe%2BNSJDAKl%2BNSJDAKk%2BNSJDAKc%2BJSJDAKc%2BJiJDAKc%2BIyJDAKc%2BJCJDAKc%2BKSJDAKc%2BKiJDAKc%2BJyJDAKc%2BKCJDAKc%2BLSJDAKc%2BLiJDAK0%2F8YSAtjQ4LIPUutxh%2FW2zjCTEqo3l7X5dHkNjl4%3D&StartTime=2020020101&ddlStartYear=2020&ddlStartMonth=2&ddlStartDay=1&ddlStartHour=1&ddlEndYear=2020&ddlEndMonth=2&ddlEndDay=24&ddlEndHour=16&ddlArea=03&CheckBoxMstList%240=on&CheckBoxMstList%241=on&CheckBoxMstList%242=on&CheckBoxMstList%243=on&CheckBoxMstList%244=on&CheckBoxMstList%245=on&CheckBoxMstList%246=on&CheckBoxMstList%247=on&CheckBoxMstList%248=on&CheckBoxMstList%249=on&CheckBoxMstList%2410=on&CheckBoxMstList%2411=on&CheckBoxMstList%2412=on&CheckBoxMstList%2413=on&CheckBoxMstList%2414=on&CheckBoxMstList%2415=on&CheckBoxMstList%2416=on&CheckBoxMstList%2417=on&CheckBoxMstList%2418=on&CheckBoxMstList%2419=on&download=%E3%83%80%E3%82%A6%E3%83%B3%E3%83%AD%E3%83%BC%E3%83%89",
}

func parse(raw io.Reader) (*DB, error) {
	db := &DB{
		data: make(map[string]*Entry),
	}
	var err error
	r := csv.NewReader(raw)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	/*
		ダウンロードデータ(Data.csv)の説明

		列番号,項目,備考
		1,測定局コード,
		2,アメダス測定局コード,
		3,年月日,
		4,時,
		5,測定局名,
		6,測定局種別,1：都市部、2：山間部、0：区分なし
		7,都道府県コード,01～47
		8,都道府県名,
		9,市区町村コード,5桁
		10,市区町村名,
		11,花粉飛散数[個/m3],
		12,風向,0：静穏、1：北北東、2：北東、3：東北東、4：東、5：東南東、6：南東、7：南南東、8：南、9：南南西、10：南西、11：西南西、12：西、13：西北西、14：北西、15：北北西、16：北
		13,風速[m/s],
		14,気温[℃],
		15,降水量[mm],
		16,レーダー降水量[mm],
	*/
	loc, _ := time.LoadLocation("Asia/Tokyo")
	for _, rec := range records {
		var hour int
		var obs Observation
		if obs.Date, err = time.ParseInLocation("20060102 15", fmt.Sprintf("%s 00", rec[3-1]), loc); err != nil {
			return nil,err
		}
		if hour, err = strconv.Atoi(rec[4-1]); err != nil {
			return nil, err
		}
		obs.Date = obs.Date.Add(time.Duration(hour) * time.Hour)
		if obs.Count, err = strconv.Atoi(rec[11-1]); err != nil {
			return nil, err
		}
		ent := db.data[rec[5-1]]
		if ent == nil {
			ent = &Entry{}
			ent.Name = rec[5-1]
			db.data[ent.Name] = ent
		}
		ent.Observes = append(ent.Observes, obs)
	}
	return db,nil
}

func FetchFromFile(fname string) (*DB, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil,err
	}
	defer f.Close()
	return parse(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
}
func FetchFromInternet(name string) (*DB, error) {
	body := kRequestMap[name]
	req, err := http.NewRequest("POST", "http://kafun.taiki.go.jp/DownLoad1.aspx", strings.NewReader(body))
	if err != nil {
		return nil,err
	}
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:73.0) Gecko/20100101 Firefox/73.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "ja,en;q=0.5")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Referer", "http://kafun.taiki.go.jp/DownLoad1.aspx")
	ret,err := http.DefaultClient.Do(req)
	if err != nil {
		return nil,err
	}
	return parse(transform.NewReader(ret.Body, japanese.ShiftJIS.NewDecoder()))
}