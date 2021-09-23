package utils

import (
	"crypto/md5"
	"encoding/hex"
)

const (
	DateHyphenYMD  = "2006-01-02"
	DateHyphenYM   = "2006-01"
	DateSlashYMD   = "2006/01/02"
	DateSlashYM    = "2006/01"
	DateChineseYMD = "2006年01月02日"
	DateChineseYM  = "2006年01月"
	DatePointYMD   = "2006.01.02"
	DatePointYM    = "2006.01"
	DateRawYMD     = "20060102"
	DateRawYM      = "200601"
	DateRawMD      = "0102"
	DateRawYear    = "2006"
	DateRawMonth   = "01"
	DateRawDay     = "02"

	TimeHyphenHMS  = "2006-01-02 15:04:05"
	TimeHyphenHM   = "2006-01-02 15:04"
	TimeChineseHMS = "2006年01月02日 15时04分05秒"
	TimeChineseHM  = "2006年01月02日 15时04分"
	TimeRawYMDHMS  = "20060102 150405"
	TimeRawALL     = "20060102150405"
	TimeRawHMS     = "150405"
	TimonHMS       = "15:04:05"
	TimeColonHM    = "15:04"
	TimeColonMS    = "04:05"
	TimeRawHour    = "15"
	TimeRawMinute  = "04"
	TimeRawSecond  = "05"

	//秒数
	DurOneSecond = 1
	DurOneMinute = DurOneSecond * 60
	DurOneHour   = DurOneMinute * 60
	DurOneDay    = DurOneHour * 24
)

func Md5Sum(data []byte) string {
	return hex.EncodeToString(byte16ToBytes(md5.Sum(data)))
}

func byte16ToBytes(in [16]byte) []byte {
	tmp := make([]byte, 16)
	for _, value := range in {
		tmp = append(tmp, value)
	}
	return tmp[16:]
}
