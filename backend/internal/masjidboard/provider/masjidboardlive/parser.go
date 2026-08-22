package masjidboardlive

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/X-Calibre/MasjidPi/backend/internal/masjidboard/model"
)

const (
	rowJumuah       = 1
	rowClock        = 2
	rowSalah        = 3
	rowDisplayTimes = 4
	rowAstronomical = 5
	rowMasjid       = 6
)

// Parse normalises the verified core portion of a MasjidBoard Live response.
// The 29-row upstream structure remains confined to this provider package.
func Parse(rows []json.RawMessage, boardID string, now time.Time) (model.Board, error) {
	if boardID == "" {
		return model.Board{}, fmt.Errorf("masjidboardlive: board ID is required")
	}
	if len(rows) < 29 {
		return model.Board{}, fmt.Errorf("masjidboardlive: expected 29 rows, got %d", len(rows))
	}

	masjid, err := parseMasjidRow(rows[rowMasjid])
	if err != nil { return model.Board{}, err }
	clock, err := parseClockRow(rows[rowClock])
	if err != nil { return model.Board{}, err }
	if clock.Timezone == "" { return model.Board{}, fmt.Errorf("masjidboardlive: clock row has no timezone") }

	loc, err := parseTimezone(clock.Timezone, masjid.OffsetMS)
	if err != nil { return model.Board{}, err }
	localNow := now.In(loc)

	prayers, err := parseSalahRow(rows[rowSalah])
	if err != nil { return model.Board{}, err }
	jumuah, err := parseJumuahRow(rows[rowJumuah])
	if err != nil { return model.Board{}, err }
	if jumuah != nil { prayers.Jumuah = []model.JumuahService{*jumuah} }

	astronomical, err := parseAstronomicalRow(rows[rowAstronomical])
	if err != nil { return model.Board{}, err }
	astronomical, err = applyDisplayTimes(rows[rowDisplayTimes], astronomical)
	if err != nil { return model.Board{}, err }

	return model.Board{
		Identity: model.BoardIdentity{ID: clock.MasjidID, Name: masjid.Name1, AlternateName: masjid.Name2, TimeZone: clock.Timezone},
		DateContext: model.DateContext{GregorianDate: dateOnly(localNow, loc), IslamicDate: stringValueFromRow(rows[rowClock], 5)},
		PrayerTimes: prayers,
		Astronomical: astronomical,
	}, nil
}

type masjidRow struct { Name1, Name2, URL string; OffsetMS int64 }
func parseMasjidRow(raw json.RawMessage) (masjidRow, error) {
	values, err := rowValues(raw); if err != nil { return masjidRow{}, fmt.Errorf("masjidboardlive: parse row 6: %w", err) }
	if len(values) < 5 { return masjidRow{}, fmt.Errorf("masjidboardlive: row 6 has %d fields, need at least 5", len(values)) }
	offset, err := strconv.ParseInt(stringValue(values, 4), 10, 64); if err != nil { return masjidRow{}, fmt.Errorf("masjidboardlive: invalid row 6 timezone offset %q: %w", stringValue(values, 4), err) }
	return masjidRow{Name1:stringValue(values,0), Name2:stringValue(values,1), URL:stringValue(values,2), OffsetMS:offset}, nil
}

type clockRow struct { MasjidID, Location, Timezone string }
func parseClockRow(raw json.RawMessage) (clockRow, error) {
	values, err := rowValues(raw); if err != nil { return clockRow{}, fmt.Errorf("masjidboardlive: parse row 2: %w", err) }
	if len(values) < 16 { return clockRow{}, fmt.Errorf("masjidboardlive: row 2 has %d fields, need at least 16", len(values)) }
	return clockRow{MasjidID:stringValue(values,12), Location:stringValue(values,14), Timezone:stringValue(values,15)}, nil
}

var gmtOffsetRE = regexp.MustCompile(`^GMT([+-])(\d{2})(?::?(\d{2}))?$`)
func parseGMTOffset(value string) (int64,error) { match:=gmtOffsetRE.FindStringSubmatch(strings.TrimSpace(value)); if match==nil{return 0,fmt.Errorf("masjidboardlive: unsupported timezone %q",value)}; hours,_:=strconv.Atoi(match[2]); minutes:=0;if match[3]!=""{minutes,_=strconv.Atoi(match[3])}; seconds:=hours*3600+minutes*60;if match[1]=="-"{seconds=-seconds};return int64(seconds)*1000,nil }
func parseTimezone(label string, offsetMS int64)(*time.Location,error){label=strings.TrimSpace(label);if label==""{return nil,fmt.Errorf("masjidboardlive: timezone is empty")};if offsetMS==0{parsed,err:=parseGMTOffset(label);if err!=nil{return nil,err};offsetMS=parsed};return time.FixedZone(label,int(offsetMS/1000)),nil}

func parseSalahRow(raw json.RawMessage)(model.PrayerTimes,error){values,err:=rowValues(raw);if err!=nil{return model.PrayerTimes{},fmt.Errorf("masjidboardlive: parse row 3: %w",err)};if len(values)<10{return model.PrayerTimes{},fmt.Errorf("masjidboardlive: row 3 has %d fields, need at least 10",len(values))};fields:=make([]*model.ClockTime,10);for i:=range fields{value:=stringValue(values,i);if isAbsent(value){continue};parsed,err:=parseClockTime(value);if err!=nil{return model.PrayerTimes{},fmt.Errorf("masjidboardlive: row 3 column %d value %q: %w",i,value,err)};fields[i]=parsed};prayers:=model.PrayerTimes{Fajr:model.PrayerTime{Adhan:fields[0],Jamaah:fields[1]},Dhuhr:model.PrayerTime{Adhan:fields[2],Jamaah:fields[3]},Asr:model.PrayerTime{Adhan:fields[4],Jamaah:fields[5]},Maghrib:model.PrayerTime{Adhan:fields[6],Jamaah:fields[7]},Esha:model.PrayerTime{Adhan:fields[8],Jamaah:fields[9]}};checks:=[]struct{name string;value model.PrayerTime}{{"Fajr",prayers.Fajr},{"Dhuhr",prayers.Dhuhr},{"Asr",prayers.Asr},{"Maghrib",prayers.Maghrib},{"Esha",prayers.Esha}};for _,check:=range checks{if check.value.Adhan==nil&&check.value.Jamaah==nil{return model.PrayerTimes{},fmt.Errorf("masjidboardlive: missing core prayer time for %s",check.name)}};return prayers,nil}

func parseJumuahRow(raw json.RawMessage)(*model.JumuahService,error){values,err:=rowValues(raw);if err!=nil{return nil,fmt.Errorf("masjidboardlive: parse row 1: %w",err)};if len(values)<12{return nil,fmt.Errorf("masjidboardlive: row 1 has %d fields, need at least 12",len(values))};hasData:=false;for i:=0;i<12;i++{if !isAbsent(stringValue(values,i)){hasData=true;break}};if !hasData{return nil,nil};codes:=strings.Split(stringValue(values,11),",");events:=make([]model.JumuahEvent,0,3);for i:=0;i<3;i++{heading:=stringValue(values,i*2);timeValue:=stringValue(values,i*2+1);var parsed *model.ClockTime;if !isAbsent(timeValue){parsed,err=parseClockTime(timeValue);if err!=nil{return nil,fmt.Errorf("masjidboardlive: row 1 Jumuah time %d value %q: %w",i+1,timeValue,err)}};code:="";if i<len(codes){code=strings.TrimSpace(codes[i])};if !isAbsent(heading)||parsed!=nil||code!=""{events=append(events,model.JumuahEvent{Code:code,Heading:heading,Time:parsed})}};parseOptional:=func(index int)(*model.ClockTime,error){value:=stringValue(values,index);if isAbsent(value){return nil,nil};parsed,parseErr:=parseClockTime(value);if parseErr!=nil{return nil,fmt.Errorf("masjidboardlive: row 1 column %d value %q: %w",index,value,parseErr)};return parsed,nil};adhan,err:=parseOptional(7);if err!=nil{return nil,err};jamaah,err:=parseOptional(8);if err!=nil{return nil,err};alternateAdhan,err:=parseOptional(9);if err!=nil{return nil,err};alternateJamaah,err:=parseOptional(10);if err!=nil{return nil,err};return &model.JumuahService{Adhan:adhan,Jamaah:jamaah,AlternateAdhan:alternateAdhan,AlternateJamaah:alternateJamaah,Khateeb:stringValue(values,6),Events:events},nil}

func parseAstronomicalRow(raw json.RawMessage)(*model.AstronomicalTimes,error){values,err:=rowValues(raw);if err!=nil{return nil,fmt.Errorf("masjidboardlive: parse row 5: %w",err)};if len(values)<9{return nil,fmt.Errorf("masjidboardlive: row 5 has %d fields, need at least 9",len(values))};parsed:=make([]*model.ClockTime,9);any:=false;for i:=range parsed{value:=stringValue(values,i);if isAbsent(value){continue};valueTime,err:=parseClockTime(value);if err!=nil{return nil,fmt.Errorf("masjidboardlive: row 5 column %d value %q: %w",i,value,err)};parsed[i]=valueTime;any=true};if !any{return nil,nil};return &model.AstronomicalTimes{Suhur:parsed[0],FajrStart:parsed[1],Sunrise:parsed[2],Ishraaq:parsed[3],Duha:parsed[4],AsrShafii:parsed[5],AsrHanafi:parsed[6],Sunset:parsed[7],EshaStart:parsed[8]},nil}

// applyDisplayTimes adds the Istiwa/Zawaal values that MasjidBoard Live stores
// separately from its main astronomical row.
func applyDisplayTimes(raw json.RawMessage, astronomical *model.AstronomicalTimes)(*model.AstronomicalTimes,error){values,err:=rowValues(raw);if err!=nil{return nil,fmt.Errorf("masjidboardlive: parse row 4: %w",err)};if len(values)<5{return astronomical,nil};parse:=func(index int)(*model.ClockTime,error){value:=stringValue(values,index);if isAbsent(value){return nil,nil};parsed,err:=parseClockTime(value);if err!=nil{return nil,fmt.Errorf("masjidboardlive: row 4 column %d value %q: %w",index,value,err)};return parsed,nil};caution,err:=parse(0);if err!=nil{return nil,err};zawaal,err:=parse(1);if err!=nil{return nil,err};istiwa,err:=parse(4);if err!=nil{return nil,err};if astronomical==nil&&(caution!=nil||zawaal!=nil||istiwa!=nil){astronomical=&model.AstronomicalTimes{}};if astronomical!=nil{astronomical.IstiwaCaution=caution;astronomical.Istiwa=istiwa;astronomical.ZawaalEnd=zawaal};return astronomical,nil}

func rowValues(raw json.RawMessage)([]json.RawMessage,error){var values []json.RawMessage;if err:=json.Unmarshal(raw,&values);err!=nil{return nil,err};return values,nil}
func stringValue(values []json.RawMessage,index int)string{if index<0||index>=len(values){return ""};var value string;if err:=json.Unmarshal(values[index],&value);err!=nil{return ""};return strings.TrimSpace(value)}
func stringValueFromRow(raw json.RawMessage,index int)string{values,err:=rowValues(raw);if err!=nil{return ""};return stringValue(values,index)}
func isAbsent(value string)bool{switch strings.TrimSpace(value){case "","-","–","—","~~~~","FALSE","false","Hide","hide","#N/A":return true;default:return false}}
func parseClockTime(value string)(*model.ClockTime,error){value=strings.TrimSpace(value);layouts:=[]string{"15:04","15:04:05","3:04 PM","3:04PM","3:04 pm","3:04pm"};for _,layout:=range layouts{parsed,err:=time.Parse(layout,value);if err==nil{return &model.ClockTime{Hour:parsed.Hour(),Minute:parsed.Minute()},nil}};return nil,fmt.Errorf("unsupported time format")}
func dateOnly(t time.Time,loc *time.Location)time.Time{local:=t.In(loc);return time.Date(local.Year(),local.Month(),local.Day(),0,0,0,0,loc)}
