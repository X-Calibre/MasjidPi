
//------------------------------------------------------------------------------------------------
// Declaration of all global variables
var nowDate                     = Date.now();
var currentDate                 = new Date();
var currentDay                  = currentDate.getDay();
var currentHours                = currentDate.getHours();
var currentMinutes              = currentDate.getMinutes();
var currentSeconds              = currentDate.getSeconds();
var fullYear                    = currentDate.getFullYear();
var fullMonth                   = currentDate.getMonth();
var fullDate                    = currentDate.getDate();
var oneDayMilli                 = parseInt(86400000);
var timeInSeconds               = (currentDate.getHours() * 60 * 60) + (currentDate.getMinutes() * 60) + currentDate.getSeconds();
var fajrNextDate                = 0;
var fajrNextTime                = 0;
var fajrNextMilli               = 0;
var asrNextDate                 = 0;
var asrNextTime                 = 0;
var asrNextMilli                = 0;
var eshaNextDate                = 0;
var eshaNextTime                = 0;
var eshaNextMilli               = 0;
var nextChangeDisplay           = false;
var jumuahHead1                 = 0;
var jumuahTime1                 = 0;
var jumuahHead2                 = 0;
var jumuahTime2                 = 0;
var jumuahHead3                 = 0;
var jumuahTime3                 = 0;
var jumuahKhateeb               = 0;
var jumuahHeadingsArray         = 0
var moonBirth                   = 0;
var moonAge                     = 0;
var bestVisibil                 = 0;
var moonRise                    = 0;
var moonSet                     = 0;
//var moonVisibilFor              = 0;
var moonAzimuth0                = 0;
var moonAzimuth1                = 0;
var moonAltitude0               = 0;
var moonAltitude1               = 0;
var moonVisibilScript           = 0;
var moonBirthScript             = 0;
var masjidImage                 = 0;
var fajrAthan                   = 1;
var fajrJamaah                  = 0;
var dhuhrAthan                  = 0;
var dhuhrJamaah                 = 0;
var asrAthan                    = 0;
var asrJamaah                   = 0;
var maghribAthan                = 0;
var maghribJamaah               = 0;
var eshaAthan                   = 0;
var eshaJamaah                  = 0;
var istiwaCaution               = 0;
var istiwa                      = 0;
var zawaalEnd                   = 0;
var theme                       = 0;
var displayLanguage             = 0;
var sehriEnds                   = 0;
var fajrStarts                  = 0;
var sunrise                     = 0;
var ishraaq                     = 0;
var duha                        = 0;
var asrShafi                    = 0;
var asrHanafi                   = 0;
var sunset                      = "00:00";
var iftar                       = 0;
var eshaStarts                  = 0;
var sunsetI                     = 0;
var eshaStartsI                 = 0;
var suhurEndsI                  = 0;
var fajrStartsI                 = 0;
var sunriseI                    = 0;
var ishraaqI                    = 0;
var duhaI                       = 0;
var istiwaCI                    = 0;
var istiwaI                     = 0;
var zuhrStartsI                 = 0;
var asrShafiI                   = 0;
var asrHanafiI                  = 0;
var maghribAthanI               = 0;
var maghribJamaahI              = 0;
var eshaAthanI                  = 0;
var eshaJamaahI                 = 0;
var fajrAthanI                  = 0;
var fajrJamaahI                 = 0;
var dhuhrAthanI                 = 0;
var dhuhrJamaahI                = 0;
var asrAthanI                   = 0;
var asrJamaahI                  = 0;
var masjidName1I                = 0;
var masjidName2I                = 0;
var sundayDhuhr                 = 0;
var salaahLanguage              = 0;
var sundayzuhr                  = 0;
var moonBirthMilli              = 0;
var moonVisibilMilli            = 0;
var sunsetMilli                 = 0;
var masjidName1                 = 0;
var masjidName2                 = 0;
var masjidUrl                   = 0;
var timeZone                    = 0;
var timeZoneMilli               = 0;
var islamicYear                 = 0;
var islamicMonthNum             = 0;
var islamicDateNum              = 0;
var islamicDateAdjust           = 0;
var forceDate30                 = 0;
var masjidname1TopPad           = 0;
var masjidname2TopPad           = 0;
var bigClockTopPad              = 0;
var masjidNameFont              = 0;
var ayahSurah                   = 0;
var AyahNo                      = 0;
var ayah                        = 0;
var ayahTrueFalse               = 0;
var hadith                      = 0;
var hadithRef                   = 0;
var hadithTrueFalse             = 0;
var sunnahHeading               = 0;
var sunnah                      = 0;
var sunnahRef                   = 0;
var sunnahRefCheck              = 0;
var sunnahTrueFalse             = 0;
var taleemTime1                 = 0;
var taleemDate1                 = 0;
var taleemResi1                 = 0;
var taleemAddy1                 = 0;
var taleemTime2                 = 0;
var taleemDate2                 = 0;
var taleemResi2                 = 0;
var taleemAddy2                 = 0;
var taleemTrueFalse1            = 0;
var taleemTrueFalse2            = 0;
var announceHead2               = 0;
var announce2                   = 0;
var announceTrueFalse2          = 0;
var announceHead3               = 0;
var announce3                   = 0;
var announceTrueFalse3          = 0;
var announceHead4               = 0;
var announce4                   = 0;
var announceTrueFalse4          = 0;
var announceHead5               = 0;
var announce5                   = 0;
var announceTrueFalse5          = 0;
var announceHead6               = 0;
var announce6                   = 0;
var announceTrueFalse6          = 0;
var announceHead7               = 0;
var announce7                   = 0;
var announceTrueFalse7          = 0;
var announceHead8               = 0;
var announce8                   = 0;
var announceTrueFalse8          = 0;
var announceHead9               = 0;
var announce9                   = 0;
var announceTrueFalse9          = 0;
var announceHead10              = 0;
var announce10                  = 0;
var announceTrueFalse10         = 0;
var funeralName                 = 0;
var funeralRelation             = 0;
var funeralAddy                 = 0;
var funeralPickup               = 0;
var funeralCemetery             = 0;
var funeralSalaahVenue          = 0;
var funeralSalaahTime           = 0;
var funeralTrueFalse            = 0;
var nikahNameOne                = 0;
var nikahGroomRelation          = 0;
var nikahRelationOne            = 0;
var nikahRelationTwo            = 0;
var nikahNameTwo                = 0;
var nikahDate                   = 0;
var nikahTime                   = 0;
var nikahPopUp                  = 0;
var nikahTrueFalse              = 0;
var masjidTaleem                = 0;
var gashtDayIn                  = 0;
var gashtDayOut                 = 0;
var gashtOut1                   = 0;
var gashtOut2                   = 0;
//var gashtOut3                   = 0;
//var gashtOut4                   = 0;
var gashtTrueFalse              = 'Hide';
var threeDaysHeader             = 0;
var threeDays1                  = 0;
var threeDays1Date              = 0;
var threeDays2                  = 0;
var threeDays2Date              = 0;
//var threeDays3                  = 0;
//var threeDays3Date              = 0;
var dawahTrueFalse              = 'Hide';
var commBroadHead               = 0;
var commBroadText               = 0;
var commBroadTrueFalse          = 0;
var commPosterImage1            = 0;
var commPosterTrueFalse1        = 0;
var commPosterImage2            = 0;
var commPosterTrueFalse2        = 0;
var commPosterImage3            = 0;
var commPosterTrueFalse3        = 0;
var commPosterImage4            = 0;
var commPosterTrueFalse4        = 0;
var commPosterImage5            = 0;
var commPosterTrueFalse5        = 0;
var commPosterImage6            = 0;
var commPosterTrueFalse6        = 0;
var commPosterImage7            = 0;
var commPosterTrueFalse7        = 0;
var commPosterImage8            = 0;
var commPosterTrueFalse8        = 0;
var commPosterImage9            = 0;
var commPosterTrueFalse9        = 0;
var commPosterImage10           = 0;
var commPosterTrueFalse10       = 0;
var eidDate                     = 0;
var eidVenue                    = 0;
var eidAddress                  = 0;
var eidLecture                  = 0;
var eidSalaah                   = 0;
var eidgahTrueFalse             = 0;
var posterTopPadding            = 0;
var posterHeight                = 0;
var posterWidth                 = 0;
var posterImage1                = 0;
var posterImage2                = 0;
var posterImage3                = 0;
var posterImage4                = 0;
var posterImage5                = 0;
var posterImage6                = 0;
var posterImage7                = 0;
var posterImage8                = 0;
var posterImage9                = 0;
var posterImage10               = 0;
var posterTrueFalse1            = 0;
var posterTrueFalse2            = 0;
var posterTrueFalse3            = 0;
var posterTrueFalse4            = 0;
var posterTrueFalse5            = 0;
var posterTrueFalse6            = 0;
var posterTrueFalse7            = 0;
var posterTrueFalse8            = 0;
var posterTrueFalse9            = 0;
var posterTrueFalse10           = 0;
var bigPosterImage0             = 0;
var bigPosterTrueFalse0         = 0;
var bigPosterImage1             = 0;
var bigPosterTrueFalse1         = 0;
var bigPosterImage2             = 0;
var bigPosterTrueFalse2         = 0;
var bigPosterImage3             = 0;
var bigPosterTrueFalse3         = 0;
var bigPosterImage4             = 0;
var bigPosterTrueFalse4         = 0;
var bigPosterImage5             = 0;
var bigPosterTrueFalse5         = 0;
var bigPosterImage6             = 0;
var bigPosterTrueFalse6         = 0;
var bigPosterImage7             = 0;
var bigPosterTrueFalse7         = 0;
var bigPosterImage8             = 0;
var bigPosterTrueFalse8         = 0;
var bigPosterImage9             = 0;
var bigPosterTrueFalse9         = 0;
var bigPosterDuration0          = 0;
var bigPosterDuration1          = 0;
var bigPosterDuration2          = 0;
var bigPosterDuration3          = 0;
var bigPosterDuration4          = 0;
var bigPosterDuration5          = 0;
var bigPosterDuration6          = 0;
var bigPosterDuration7          = 0;
var bigPosterDuration8          = 0;
var bigPosterDuration9          = 0;
var bigPosterSecondScreen       = false;
var bigPosterSecondScreenActive = false;
var oldTotalPosterDuration      = 120;
//var bigPosterInterval           = 300000;
var sheetRefreshRate            = 15000; //120000;
var slideDuration               = 30;
var bankHeader                  = 0;
var bankName                    = 0;
//var bankBranch                  = 0;
var bankCode                    = 0;
var accountName                 = 0;
var accountNum                  = 0;
var bankTrueFalse               = 0;
var bankAUSBSB                  = 0;
var islamicDateEng              = 0;
var islamicDateAra              = 0;
var sick1                       = 0;
var sick2                       = 0;
var sick3                       = 0;
var sick4                       = 0;
var sick5                       = 0;
var sick6                       = 0;
var sick7                       = 0;
var sick8                       = 0;
var sick9                       = 0;
var sick10                      = 0;
var sickTrueFalse               = 0;
var highlightSalaah             = 0;
var masjidNameOnTop             = 0;

var cBigPosterImage0             = 0;
var cBigPosterTrueFalse0         = 0;
var cBigPosterImage1             = 0;
var cBigPosterTrueFalse1         = 0;
var cBigPosterImage2             = 0;
var cBigPosterTrueFalse2         = 0;
var cBigPosterImage3             = 0;
var cBigPosterTrueFalse3         = 0;
var cBigPosterImage4             = 0;
var cBigPosterTrueFalse4         = 0;
var cBigPosterImage5             = 0;
var cBigPosterTrueFalse5         = 0;
var cBigPosterImage6             = 0;
var cBigPosterTrueFalse6         = 0;
var cBigPosterImage7             = 0;
var cBigPosterTrueFalse7         = 0;
var cBigPosterImage8             = 0;
var cBigPosterTrueFalse8         = 0;
var cBigPosterImage9             = 0;
var cBigPosterTrueFalse9         = 0;
var cBigPosterDuration0          = 0;
var cBigPosterDuration1          = 0;
var cBigPosterDuration2          = 0;
var cBigPosterDuration3          = 0;
var cBigPosterDuration4          = 0;
var cBigPosterDuration5          = 0;
var cBigPosterDuration6          = 0;
var cBigPosterDuration7          = 0;
var cBigPosterDuration8          = 0;
var cBigPosterDuration9          = 0;

var pageHardReload              = false;
var boardLayout                 = 0;
var salaah12HourFormat          = false;
var islamicTime                 = false;
var islamicLiveTime             = false
var islamicTimeTicker           = false;
var liveStreamingServer         = null;
var liveStreamingURL            = null;
var puAlternatingTime           = false;
var islamicTimeInterval         = 1000;
var westernTimeInterval         = 1000;
var rightToLeft;
var defaultRightToLeft;
var islamicTimesArabic;
var puIslamicTimes;
var islamicTimeMaghribFirst;
var masjidNameLanguage;
var islamicTime12Hour;
var tickerSpeed;
var isBigPopUpShowing           = false
var highlightOwwal;
var disableBoard;
var isBoardDisabled             = false;
var nextSalaahFlash             = true;

// ~~~~~~~~~ Variables used in JS ~~~~~~~~ //
const urlParams        = new URLSearchParams(window.location.search);
var mobile = true;
(urlParams.get('mid')) ? mobile = true: mobile = false;

const debug            = urlParams.get('mblbug');
var popUpsDuplicated   = false;
var oldInfo            = [];
let posterSetInterval  = 450000
var posters            = {
    iftaarPoster      : { id: "iftaarPoster"      , state: 0 },
    thursDuroodPoster : { id: "thursDuroodPoster" , state: 0 },
    jumuahSunnah      : { id: "jumuahSunnah"      , state: 0 },
    jumuahDurood0     : { id: "jumuahDurood0"     , state: 0 },
    jumuahDurood1     : { id: "jumuahDurood1"     , state: 0 },
    jumuahDurood2     : { id: "jumuahDurood2"     , state: 0 },
    jumuahDurood3     : { id: "jumuahDurood3"     , state: 0 },
    jumuahDurood4     : { id: "jumuahDurood4"     , state: 0 },
    jumuahDurood5     : { id: "jumuahDurood5"     , state: 0 },
    jumuahDurood6     : { id: "jumuahDurood6"     , state: 0 },
    jumuahDurood7     : { id: "jumuahDurood7"     , state: 0 },
    jumuahDurood8     : { id: "jumuahDurood8"     , state: 0 },
    jumuahDurood9     : { id: "jumuahDurood9"     , state: 0 },
    ayaamPoster       : { id: "ayaamPoster"       , state: 0 },
    islamicMonthPoster: { id: "islamicMonthPoster", state: 0 },
    muharramPoster1   : { id: "muharramPoster1"   , state: 0 },
    muharramPoster2   : { id: "muharramPoster2"   , state: 0 },
    muharramPoster3   : { id: "muharramPoster3"   , state: 0 },
    ashuraSpendPoster : { id: "ashuraSpendPoster" , state: 0 },
    rajabPoster       : { id: "rajabPoster"       , state: 0 },
    shabaanPoster1    : { id: "shabaanPoster1"    , state: 0 },
    shabaanPoster2    : { id: "shabaanPoster2"    , state: 0 },
    ramadanApproachPoster: { id: "ramadanApproachPoster", state: 0 },
    ramadanPoster1    : { id: "ramadanPoster1"    , state: 0 },
    ramadanPoster2    : { id: "ramadanPoster2"    , state: 0 },
    ramadanPoster3    : { id: "ramadanPoster3"    , state: 0 },
    laylatulQadrPoster: { id: "laylatulQadrPoster", state: 0 },
    eidDuaPoster      : { id: "eidDuaPoster"      , state: 0 },
    hajjPoster1       : { id: "hajjPoster1"       , state: 0 },
    hajjPoster2       : { id: "hajjPoster2"       , state: 0 },
    hajjPoster3       : { id: "hajjPoster3"       , state: 0 },
    hajjPoster4       : { id: "hajjPoster4"       , state: 0 },
    hajjPoster5       : { id: "hajjPoster5"       , state: 0 },
    hajjPoster6       : { id: "hajjPoster6"       , state: 0 },
    hajjPoster7       : { id: "hajjPoster7"       , state: 0 },
    hajjPoster8       : { id: "hajjPoster8"       , state: 0 },
    hajjPoster9       : { id: "hajjPoster9"       , state: 0 },
    hajjPoster10      : { id: "hajjPoster10"      , state: 0 },
    tashreeqPoster    : { id: "tashreeqPoster"    , state: 0 },
    bigPosterImage0   : { id: ""                  , state: 0 },
    bigPosterImage1   : { id: ""                  , state: 0 },
    bigPosterImage2   : { id: ""                  , state: 0 },
    bigPosterImage3   : { id: ""                  , state: 0 },
    bigPosterImage4   : { id: ""                  , state: 0 },
    bigPosterImage5   : { id: ""                  , state: 0 },
    bigPosterImage6   : { id: ""                  , state: 0 },
    bigPosterImage7   : { id: ""                  , state: 0 },
    bigPosterImage8   : { id: ""                  , state: 0 },
    bigPosterImage9   : { id: ""                  , state: 0 },
    cBigPosterImage0  : { id: ""                  , state: 0 },
    cBigPosterImage1  : { id: ""                  , state: 0 },
    cBigPosterImage2  : { id: ""                  , state: 0 },
    cBigPosterImage3  : { id: ""                  , state: 0 },
    cBigPosterImage4  : { id: ""                  , state: 0 },
    cBigPosterImage5  : { id: ""                  , state: 0 },
    cBigPosterImage6  : { id: ""                  , state: 0 },
    cBigPosterImage7  : { id: ""                  , state: 0 },
    cBigPosterImage8  : { id: ""                  , state: 0 },
    cBigPosterImage9  : { id: ""                  , state: 0 },
    moonVisibility    : { id: "moonVisibility"    , state: 0 },
    hideAll           : { id: "hideAll"           , state: 0 }
}
var posterHide         = {
    "funeral"   : false,
    "countdown" : false,
    "zawaal"    : false,
    "moon"      : false,
    "daa"       : false,
    "cell"      : false,
    "nikaah"    : false,
    "tashreeq"  : false,
    "iftaar"    : false,
    "blackout"  : false
};
var masjidNameCounterG = 0
let chast              = false;
let colors             = {  //text     //bg
    primary     : ["#084298","#cfe2ff"],
    secondary   : ["#41464b","#e2e3e5"],
    success     : ["#0f5132","#d1e7dd"],
    danger      : ["#842029","#f8d7da"],
    warning     : ["#664d03","#fff3cd"],
    info        : ["#055160","#cff4fc"] // time
}
let logStyles          = {
    time        : `color:${colors.info[1]};background-color: ${colors.info[0]};padding: 3px 10px; border-radius: 5px;`,
    primary     : `color:${colors.primary[1]};`,
    secondary   : `color:${colors.secondary[1]};`,
    success     : `color:${colors.success[1]};`,
    warning     : `color:${colors.warning[1]};`,
    danger      : `color:${colors.danger[1]};`,
    green       : `color:#29871c; font-weight: bold; border: thin solid #0f0; padding: 3px 10px; background-color: white;`,
    red         : `color:#a1151a; font-weight: bold; border: thin solid #f00; padding: 3px 10px; background-color: white;`,
    successbg   : `background-color: ${colors.success[0]};   padding: 3px 10px;`,
    primarybg   : `background-color: ${colors.primary[0]};   padding: 3px 10px;`,
    secondarybg : `background-color: ${colors.secondary[0]}; padding: 3px 10px;`,
    warningbg   : `background-color: ${colors.warning[0]};   padding: 3px 10px;`,
    dangerbg    : `background-color: ${colors.danger[0]};    padding: 3px 10px;`,
    info        : `color:#1E74FD; font-weight: bold;`,
    timer       : "color: #0da2bf; font-style: italic",
    poster      : `color:#cff4fc ;background-color: #0f5132; padding: 3px 10px; border-radius: 5px;`,
    function    : `color:#cff4fc ;background-color: #b5781d; padding: 3px 10px; border-radius: 5px;`,
    api         : `color:#cff4fc ;background-color: #467fcf; padding: 3px 10px; border-radius: 5px;`,
    error       : `color:#f8d7da ;background-color: #842029; padding: 3px 10px; border-radius: 5px;`,
    general     : `color:#141619 ;background-color: #d3d3d4; padding: 3px 10px; border-radius: 5px;`,
    warn        : `color:#fff3cd ;background-color: #664d03; padding: 3px 10px; border-radius: 5px;`,
    //color:${colors.success[1]};background-color: ${colors.success[0]};padding: 3px;`
}
let postersToHide      = []
// ~~~~~~~~~~ Interval Functions ~~~~~~~~~ //
let currentTimeI;
let duaAfterAdhanI;
let cellPhoneNoticeI;
let moonPopUpI;
let timesBoxI;
let salaahRowHighlightI;
let tashreeqPopUpI;
let fullDayPostersI;

// ~~~~~~~~~~ Translation Variables ~~~~~~~~~ //
let islamicDateDay   = 0;
let islamicDateDate  = 0;
let islamicDateMonth = 0;
let islamicDateYear  = 0;


//------------------------------------------------------------------------------------------------



// !====================================================== //
// !=============== Starting MBL Functions =============== //
// !====================================================== //
//  Function starts MBL, its async so code runs after the other is complete  //
async function startMBL() {
    _debug("startup","بسم الله الرحمن الرحيم","","",false,false)
    if (debug) console.groupCollapsed("%c----- MBL Start Up -----", logStyles.info)
    try{
        const checkDataResp = await checkData(fajrAthan);
            loadingScreen(checkDataResp, true, false, 0);
        const translateResp = await translateDOM(['loadingText','pu_fajr', 'pu_dhuhr', 'pu_asr', 'pu_maghrib', 'pu_esha', 'puh_sehriEnds', 'puh_fajrStarts', 'puh_sunrise', 'puh_ishraaq', 'puh_duha', 'puh_istiwa', 'puh_istiwaH', 'puh_asrHanafiH', 'puh_asrShafiH', 'puh_sundayDhuhr', 'puh_asrShafi', 'puh_asrHanafi', 'puh_sunset', 'puh_eshaStarts', 
                                                  'zawaalPopUpText', 'nextChangeHeading', 'nsch_fajr', 'nsch_asr', 'nsch_esha', 
                                                  'jumuahHeading', 'jumuahHead1', 'jumuahHead2', 'jumuahHead3', 'jumuahKhateeb', 
                                                  'adhanColumn', 'jamaahColumn', 'ayahSurah', 'hadithHeading', 'sickSlideHeading',
                                                  'commBroadHead', 'announceHead2', 'announceHead3', 'announceHead4', 'announceHead5', 'announceHead6', 'announceHead7', 'announceHead8', 'announceHead9', 'announceHead10', 
                                                  'fajrSuhurH', 'fajrStartsH', 'sunriseH', 'ishraaqH', 'duhaH', 'istiwaH', 'zawaalEndH', 'sundayzuhrH', 'asrShafiH', 'asrHanafiH', 'sunsetH', 'eshaStartsH', 
                                                  'duaAfterAdhanH', 'duaAfterAdhanEnglish', 'cellphoneText', 'newMoonH', 'moonBirthH', 'moonAgeH', 'moonSetH', 'moonVisibilH', 'moonAzimuthH', 'moonAltitudeH', 'moonDirectionHeading', 
                                                  'ayahSurah', 'AyahNo', 'ayah', 'hadithHeading', 'hadith', 'hadithRef', 'sunnahHeading', 'sunnah', 'sunnahRef'])
            loadingScreen(translateResp, true, false, 20);
        const dataResp = await getAPI(boardId)
            loadingScreen(dataResp, true, false, 40);
        const startFuncs = await runStartFunctions();
            loadingScreen(startFuncs, true, false, 60);
            const intervalFuncs = await intervalFunctions();
        loadingScreen(intervalFuncs, true, false, 80);
            const loadingDoneI = await loadingDone()
        loadingScreen(loadingDoneI, false, false, 100);
    }
    catch (err) {
      loadingScreen(err + "<br><span class='text-muted'>Contact Masjid Board Live: mbl@masjidboardlive.com</span>",true,true,100)
    }
    if (debug) console.groupEnd("MBL Start Up")
}
async function startMBL2() {
    _debug("startup","بسم الله الرحمن الرحيم","","",false,false)
    if (debug) console.groupCollapsed("%c----- MBL Second Screen Start Up -----", logStyles.info)
    try{
        const checkDataResp = await checkData(fajrAthan);
            loadingScreen(checkDataResp, true, false, 0);
        const dataResp = await getAPI(boardId)
            loadingScreen(dataResp, true, false, 40);
        const loadingDoneI = await loadingDone()
            loadingScreen(loadingDoneI, false, false, 100);
        fullDayPosters();
        detectSleepI = setInterval(detectSleep,5000)
        detectSleepI;
    }
    catch (err) {
      loadingScreen(err + "<br><span class='text-muted'>Contact Masjid Board Live: mbl@masjidboardlive.com</span>",true,true,100)
    }
    if (debug) console.groupEnd("MBL Second Screen Start Up")
}
//  Function controls loading screen status & to hide it  //
function loadingScreen(message, state, error, percent) {
    if (state == true && error == false) {
        $("#progBarPercent").animate({ "width": percent + "%" }, { duration: 10 })
        //$("#scaleable-wrapper").fadeOut(1)
        if (debug) {
            $("#loadingAlert").removeClass("hidden")
            $("#loadingAlert").removeClass("alert-danger").addClass("alert-success")
            $("#alertText").html(message)
            $.toast({
                title: message,
                subtitle: currentDate.getHours() + ":" + currentDate.getMinutes() + ":" + currentDate.getSeconds(),
                content: '',
                type: 'success',
                delay: 3000,
                dismissible: false,
            });
        }
        scaleBoard()
    } else if (state == true && error == true) {
        $("#loadingAlert").removeClass("alert-success").addClass("alert-danger")
        $("#progBarPercent").removeClass("bg-success").addClass("bg-danger")
        $("#progBarPercent").animate({ "width": percent + "%" }, { duration: 10 })
        $("#alertText").html(message)
    } else if (state == false) {
        $("#scaleable-wrapper").removeClass("hiddenNd")
        $("#progBarPercent").animate({ "width": percent + "%" }, { duration: 10 })
        $("#loading").fadeOut("slow")
        $("#toast-container").fadeOut(100)
        scaleBoard()
    }
    _debug("Function",message,"Executed Successfully",logStyles.function,false,false)
}

// ~~ Function checks if any data is 0's ~ //
function checkData(data) {
    return new Promise((resolve, reject) => {
        (data !== 0) ? setTimeout(resolve("Data updated"), 15000): reject("<br> MasjidBoard Live initialising...")
    })
}

// ~~~~ Function translates DOM items ~~~~ //
function translateDOM(id) {
    return new Promise((resolve, reject) => {
        if (lang == null || lang == 0) {
            resolve('Page did not translate. check language code')
        } else {
            for (i = 0; i < id.length; i++) {
                document.getElementById(id[i]).innerHTML = translations[id[i]][lang]
            }
            resolve(`Page Translated To ${lang}`)
        }
    })
}

//  Function to return translated JS items  //
function translate(id, lcode) {
    function checkIfArray(info) {
        if (translations[id][(lcode) ? lcode : lang][0] == '[') { return JSON.parse(info) } else { return info }
    }

    //returns eng as a default
    if (lang == null || lang == 0) {
        return checkIfArray(translations[id].en)
    }
    //if no translation available default to eng
    else if (translations[id][(lcode) ? lcode : lang] == null || translations[id][(lcode) ? lcode : lang] == undefined) {
        return checkIfArray(translations[id].en)
    } else {
        return checkIfArray(translations[id][(lcode) ? lcode : lang])
    }
}

//  Function to move item rtl  //
function rtl() {
    if (translations.lang_dir[lang] == "ltr") {
        return
    } else if (translations.lang_dir[lang] == "rtl"){
        let ids = ["ayah","hadith","sunnah","fajr","dhuhr","asr","maghrib","esha"]
        for(i = 0; i < ids.length; i++){
            console.log(ids[i])
            document.getElementById(ids[i]).style.float = "right";
            document.getElementById(ids[i]).style.textAlign = "right";
        }
    } else {
        return
    }
}

//  Function runs all functions that only run on startup once  //
function runStartFunctions() {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            currentTime();
            todaysDate();
            salaahRowHighlight();
            owwalHighlight();
            countdownCalc();
            //cellPhoneNoticeCalc();
            timesBoxCalc();
            //tashreeqPopUp();
            clocknamePositions();
            setTimeout(() => { tickerAnimate() }, 1000);
            moonPhaseText();
            moonPopUp();
            runOnChange();
            //sliderMain();
            //duaAfterAdhanCalc();
            salaahLanguageCheck();
            rtl();
            jumuah();
            function iftaarPopUp() {
                let maginsecs = minutestoseconds(maghribAthan) - 10
                let timeout = (timeInSeconds < maginsecs) ? (maginsecs - timeInSeconds) * 1000 : (timeInSeconds + maginsecs) * 1000;
                let timeoutLog = new Date(timeout).toISOString().substr(11, 8)
                _debug("Iftaar Poster","Poster showing in",timeoutLog,logStyles.poster)
                setTimeout(() => {
                    _debug("Iftaar Poster","Interval Starting",timeoutLog,logStyles.poster)
                    var iftarPopUpJI = setInterval(iftarPopUpJ, 1000); //clearInterval(iftarPopUpJI);
                    iftarPopUpJI;
                    setTimeout(() => {
                        clearInterval(iftarPopUpJI)
                        _debug("Iftaar Poster","Interval Cleared",timeoutLog,logStyles.poster)
                    }, 120000)
                }, timeout)
            }
            iftaarPopUp();
            turnOffBoard(disableBoard)
            resolve('1 Time Run Functions Started')
        }, 2000)
    })
}

// ~~~~ Strats all interval functions ~~~~ //
function intervalFunctions() {
    return new Promise((resolve, reject) => {
        currentTimeI = setInterval(currentTime, 500);
        //fullDayPostersI = setInterval(fullDayPosters, parseInt(posterSetInterval, 10) + 1000)
        if(!bigPosterSecondScreen && !posterDisabled) fullDayPostersI = setInterval(fullDayPosters, 30000) //1min
        //detectSleepI = setInterval(detectSleep,5000)
        //detectSleepI;
        resolve('Interval Functions Started')
    })
}

//  Removes loading screen & creates delay  //
function loadingDone() {
    return new Promise((resolve, reject) => {
        setTimeout(() => {
            resolve('Loading Completed')
        }, 3000)
    })
}

//  Checks if any data has changed, if chgange all true/false functions will run  //
function checkForChange(data) {
    if (JSON.stringify(oldInfo) != JSON.stringify(data)) {
        _debug("API","Data Downloaded","Info Changed, running functions",logStyles.api,false,false)
        if (debug) {
            console.groupCollapsed("%c----- New importArray Info -----", logStyles.info)
            console.table(data)
            console.groupEnd("----- New importArray Info -----")
        }
        if (debug) document.getElementById("apiCheck").innerHTML = "🔄"
        if(!bigPosterSecondScreenActive) runOnChange();
    } else {
        _debug("API","Data Downloaded","No Change Detected",logStyles.api,true,false)
        return;
    }
}

// ~~ Functions to run when data changes ~ //
function runOnChange() {
    resetTextFit();
    writeTimesToDOM()
    populateSlider()
    //fajrSuhurCheck()
    islamicDateDOM()
    funeralBoard()
    nikahPopUpFunction()
    themeChange();
    pageReload();
    clocknamePositions();
    nextSalaahChangeFlash();
}





// !====================================================== //
// !============= Get Info from API Functions ============ //
// !====================================================== //
// ~~~~~~~~~~ Get info from API ~~~~~~~~~~ //
function getAPI(url) {
    return new Promise((resolve, reject) => {
        if (debug) document.getElementById("apiCheck").innerHTML = "⌛"
        fetch(`https://api.masjidboardlive.com/mblapi?id=${url}`)
            .then(response => response.json())
            .then(json => {
                if (debug) document.getElementById("apiCheck").innerHTML = "✅"
                handleResults(json)
                checkForChange(json)
                oldInfo = json
            })
            .then(setTimeout(() => {
                getAPI(boardId)
            }, sheetRefreshRate))
            .then(setTimeout(()=>{if (debug) document.getElementById("apiCheck").innerHTML = ""},3000))
            .then(resolve("Data Downloaded"))
            .catch((err) => {_debug("API","MBL API GET Failed",err,logStyles.error,false); if (debug) document.getElementById("apiCheck").innerHTML = "❌" });
            //setTimeout(()=>{if (debug) document.getElementById("apiCheck").innerHTML = ""},10000)
    })
}

// ~~~~ Send all info to global VAR's ~~~~ //
function handleResults(spreadsheetArray) {
    return new Promise((resolve, reject) => {
        nextChangeDisplay               = stringToBoolean(spreadsheetArray[0][9]);
        fajrNextDate                    = (nextChangeDisplay) ? translateShortMonth(checkForError(spreadsheetArray[0][10])) : translateShortMonth(checkForError(spreadsheetArray[0][0]));
        fajrNextTime                    = (nextChangeDisplay) ? checkForError(spreadsheetArray[0][11]) : checkForError(spreadsheetArray[0][1]);
        fajrNextMilli                   = (nextChangeDisplay) ? checkForError(spreadsheetArray[0][12]) : checkForError(spreadsheetArray[0][2]);
        asrNextDate                     = (nextChangeDisplay) ? translateShortMonth(checkForError(spreadsheetArray[0][13])) : translateShortMonth(checkForError(spreadsheetArray[0][3]));
        asrNextTime                     = (nextChangeDisplay) ? checkForError(spreadsheetArray[0][14]) : checkForError(spreadsheetArray[0][4]);
        asrNextMilli                    = (nextChangeDisplay) ? checkForError(spreadsheetArray[0][15]) : checkForError(spreadsheetArray[0][5]);
        eshaNextDate                    = (nextChangeDisplay) ? translateShortMonth(checkForError(spreadsheetArray[0][16])) : translateShortMonth(checkForError(spreadsheetArray[0][6]));
        eshaNextTime                    = (nextChangeDisplay) ? checkForError(spreadsheetArray[0][17]) : checkForError(spreadsheetArray[0][7]);
        eshaNextMilli                   = (nextChangeDisplay) ? checkForError(spreadsheetArray[0][18]) : checkForError(spreadsheetArray[0][8]);
                  
        jumuahHead1                     = spreadsheetArray[1][0];
        jumuahTime1                     = spreadsheetArray[1][1];
        jumuahHead2                     = spreadsheetArray[1][2];
        jumuahTime2                     = spreadsheetArray[1][3];
        jumuahHead3                     = spreadsheetArray[1][4];
        jumuahTime3                     = spreadsheetArray[1][5];
        jumuahKhateeb                   = spreadsheetArray[1][6];
        jumuahAdhan                     = spreadsheetArray[1][7];
        jumuahJamaah                    = spreadsheetArray[1][8];
        jumuahAdhanI                    = spreadsheetArray[1][9];
        jumuahJamaahI                   = spreadsheetArray[1][10];
        jumuahHeadingsArray             = spreadsheetArray[1][11]

        moonBirth                       = spreadsheetArray[2][0];
        moonSet0                        = spreadsheetArray[2][1];
        moonAge0                        = spreadsheetArray[2][2];
        //moonVisibilFor0                 = spreadsheetArray[2][3];
        moonAzimuth0                    = spreadsheetArray[2][3];
        moonAltitude0                   = spreadsheetArray[2][4];
        bestVisibil                     = spreadsheetArray[2][5];
        moonSet1                        = spreadsheetArray[2][6];
        moonAge1                        = spreadsheetArray[2][7];
        //moonVisibilFor1                = spreadsheetArray[2][9];
        moonAzimuth1                    = spreadsheetArray[2][8];
        moonAltitude1                   = spreadsheetArray[2][9];
        moonBirthScript                 = spreadsheetArray[2][10];
        moonVisibilScript               = spreadsheetArray[2][11];
        masjidImage                     = spreadsheetArray[2][12];

        fajrAthan                       = checkForError(spreadsheetArray[3][0]);
        fajrJamaah                      = checkForError(spreadsheetArray[3][1]);
        dhuhrAthan                      = checkForError(spreadsheetArray[3][2]);
        dhuhrJamaah                     = checkForError(spreadsheetArray[3][3]);
        asrAthan                        = checkForError(spreadsheetArray[3][4]);
        asrJamaah                       = checkForError(spreadsheetArray[3][5]);
        iftar                           = checkForError(spreadsheetArray[3][6]);
        maghribAthan                    = checkForError(spreadsheetArray[3][6]);
        maghribJamaah                   = checkForError(spreadsheetArray[3][7]);
        eshaAthan                       = checkForError(spreadsheetArray[3][8]);
        eshaJamaah                      = checkForError(spreadsheetArray[3][9]);
        sundayDhuhr                     = checkForError(spreadsheetArray[3][10]);
        salaahLanguage                  = spreadsheetArray[3][11];
        sundayzuhr                      = spreadsheetArray[3][12];
        highlightSalaah                 = stringToBoolean(spreadsheetArray[3][13]);
        highlightOwwal                  = stringToBoolean(spreadsheetArray[3][14]);
        liveStreamingServer             = spreadsheetArray[3][15];
        liveStreamingURL                = spreadsheetArray[3][16];

        istiwaCaution                   = checkForError(spreadsheetArray[4][0]);
        zawaalEnd                       = checkForError(spreadsheetArray[4][1]);
        theme                           = spreadsheetArray[4][2];
        displayLanguage                 = spreadsheetArray[4][3];
        istiwa                          = spreadsheetArray[4][4];
        //boardLayout                     = spreadsheetArray[4][5];
        islamicTime                     = stringToBoolean(spreadsheetArray[4][6]);
        salaah12HourFormat              = stringToBoolean(spreadsheetArray[4][7]);
        pageHardReload                  = stringToBoolean(spreadsheetArray[4][8]);

        islamicTimeInterval             = spreadsheetArray[4][9] * 1000;
        westernTimeInterval             = spreadsheetArray[4][10] * 1000;
        puAlternatingTime               = stringToBoolean(spreadsheetArray[4][11]);
        islamicOwwalTimes               = stringToBoolean(spreadsheetArray[4][12]);
        rightToLeft                     = stringToBoolean(spreadsheetArray[4][13]);
        islamicTimesArabic              = stringToBoolean(spreadsheetArray[4][14]);
        islamicTimeMaghribFirst         = stringToBoolean(spreadsheetArray[4][15]);

        sehriEnds                       = checkForError(spreadsheetArray[5][0]);
        fajrStarts                      = checkForError(spreadsheetArray[5][1]);
        sunrise                         = checkForError(spreadsheetArray[5][2]);
        ishraaq                         = checkForError(spreadsheetArray[5][3]);
        duha                            = checkForError(spreadsheetArray[5][4]);
        asrShafi                        = checkForError(spreadsheetArray[5][5]);
        asrHanafi                       = checkForError(spreadsheetArray[5][6]);
        sunset                          = checkForError(spreadsheetArray[5][7]);
        eshaStarts                      = checkForError(spreadsheetArray[5][8]);
        secondLanguage                  = spreadsheetArray[5][9]
        mblPosterHide                   = (spreadsheetArray[5][10] != "-" || spreadsheetArray[5][10] != null) ? spreadsheetArray[5][10].split(",") : []
        mblPosterTime                   = spreadsheetArray[5][11]
        if(spreadsheetArray[5][12] != undefined) disableBoard = spreadsheetArray[5][12] 

        masjidName1                     = spreadsheetArray[6][0];
        masjidName2                     = spreadsheetArray[6][1];
        masjidUrl                       = spreadsheetArray[6][2];
        timeZone                        = spreadsheetArray[6][3];
        timeZoneMilli                   = spreadsheetArray[6][4];
        //islamicYear                     = spreadsheetArray[6][5];
        //islamicMonthNum                 = spreadsheetArray[6][6];
        //islamicDateNum                  = spreadsheetArray[6][7];
        islamicDateAdjust               = spreadsheetArray[6][8];
        forceDate30                     = spreadsheetArray[6][9];
        masjidname1TopPad               = spreadsheetArray[6][10];
        puIslamicTimes                  = stringToBoolean(spreadsheetArray[6][11]);
        bigClockTopPad                  = spreadsheetArray[6][12];
        clockNameFormat                 = spreadsheetArray[6][13];
        islamicTime12Hour               = stringToBoolean(spreadsheetArray[6][14]);
        defaultRightToLeft              = stringToBoolean(spreadsheetArray[6][15]);
 
        slideDuration                   = Number(spreadsheetArray[7][3]);
        ayahTrueFalse                   = spreadsheetArray[7][4]
 
        hadithTrueFalse                 = spreadsheetArray[8][2];
         
        sunnahTrueFalse                 = spreadsheetArray[9][3];
        sunnahRefCheck                  = spreadsheetArray[9][4];
 
        taleemTime1                     = spreadsheetArray[10][0];
        taleemDate1                     = spreadsheetArray[10][1];
        taleemResi1                     = spreadsheetArray[10][2];
        taleemAddy1                     = spreadsheetArray[10][3];
        taleemTrueFalse1                = spreadsheetArray[10][4];
        taleemTime2                     = spreadsheetArray[10][5];
        taleemDate2                     = spreadsheetArray[10][6];
        taleemResi2                     = spreadsheetArray[10][7];
        taleemAddy2                     = spreadsheetArray[10][8];
        taleemTrueFalse2                = spreadsheetArray[10][9];
 
        announceHead2                   = spreadsheetArray[11][3];
        announce2                       = spreadsheetArray[11][4];
        announceTrueFalse2              = spreadsheetArray[11][5];
 
        announceHead3                   = spreadsheetArray[11][6];
        announce3                       = spreadsheetArray[11][7];
        announceTrueFalse3              = spreadsheetArray[11][8];
 
        announceHead4                   = spreadsheetArray[11][9];
        announce4                       = spreadsheetArray[11][10];
        announceTrueFalse4              = spreadsheetArray[11][11];
 
        announceHead5                   = spreadsheetArray[11][12];
        announce5                       = spreadsheetArray[11][13];
        announceTrueFalse5              = spreadsheetArray[11][14];
 
        announceHead6                   = spreadsheetArray[12][0];
        announce6                       = spreadsheetArray[12][1];
        announceTrueFalse6              = spreadsheetArray[12][2];
 
        announceHead7                   = spreadsheetArray[12][3];
        announce7                       = spreadsheetArray[12][4];
        announceTrueFalse7              = spreadsheetArray[12][5];
 
        announceHead8                   = spreadsheetArray[12][6];
        announce8                       = spreadsheetArray[12][7];
        announceTrueFalse8              = spreadsheetArray[12][8];
 
        announceHead9                   = spreadsheetArray[12][9];
        announce9                       = spreadsheetArray[12][10];
        announceTrueFalse9              = spreadsheetArray[12][11];
 
        announceHead10                  = spreadsheetArray[12][12];
        announce10                      = spreadsheetArray[12][13];
        announceTrueFalse10             = spreadsheetArray[12][14];
 
        nikahNameOne                    = spreadsheetArray[13][0];
        nikahGroomRelation              = spreadsheetArray[13][1];
        nikahRelationOne                = spreadsheetArray[13][2];
        nikahRelationTwo                = spreadsheetArray[13][3];
        nikahNameTwo                    = spreadsheetArray[13][4];
        nikahDate                       = spreadsheetArray[13][5];
        nikahTime                       = spreadsheetArray[13][6];
        nikahPopUp                      = spreadsheetArray[13][7];
        nikahTrueFalse                  = spreadsheetArray[13][8];
        nikahPopUpTrueFalse             = stringToBoolean(spreadsheetArray[13][9]);
        nikahBride                      = spreadsheetArray[13][10];
 
        funeralName                     = spreadsheetArray[14][0];
        funeralRelation                 = spreadsheetArray[14][1];
        funeralAddy                     = spreadsheetArray[14][2];
        funeralPickup                   = spreadsheetArray[14][3];
        funeralCemetery                 = spreadsheetArray[14][4];
        funeralSalaahVenue              = spreadsheetArray[14][5];
        funeralSalaahTime               = spreadsheetArray[14][6];
        funeralTrueFalse                = spreadsheetArray[14][7];
 
        masjidTaleem                    = spreadsheetArray[15][0];
        gashtDayIn                      = spreadsheetArray[15][1];
        gashtDayOut                     = spreadsheetArray[15][2];
        gashtOut1                       = spreadsheetArray[15][3];
        gashtOut2                       = spreadsheetArray[15][4];
        //gashtOut3                       = spreadsheetArray[15][5];
        //gashtOut4                       = spreadsheetArray[15][6];
        gashtTrueFalse                  = spreadsheetArray[15][5];
        threeDaysHeader                 = spreadsheetArray[15][6];
        threeDays1                      = spreadsheetArray[15][7];
        threeDays1Date                  = spreadsheetArray[15][8];
        threeDays2                      = spreadsheetArray[15][9];
        threeDays2Date                  = spreadsheetArray[15][10];
        //threeDays3                      = spreadsheetArray[15][13];
        //threeDays3Date                  = spreadsheetArray[15][14];
        dawahTrueFalse                  = spreadsheetArray[15][11];
 
        //commBroadHead                   = spreadsheetArray[16][0];
        //commBroadText                   = spreadsheetArray[16][1];
        //commBroadTrueFalse              = spreadsheetArray[16][2];
        commPosterImage1                = spreadsheetArray[16][0];
        commPosterTrueFalse1            = spreadsheetArray[16][1];
        commPosterImage2                = spreadsheetArray[16][2];
        commPosterTrueFalse2            = spreadsheetArray[16][3];
        commPosterImage3                = spreadsheetArray[16][4];
        commPosterTrueFalse3            = spreadsheetArray[16][5];
        commPosterImage4                = spreadsheetArray[16][6];
        commPosterTrueFalse4            = spreadsheetArray[16][7];
        commPosterImage5                = spreadsheetArray[16][8];
        commPosterTrueFalse5            = spreadsheetArray[16][9];
        commPosterImage6                = spreadsheetArray[16][10];
        commPosterTrueFalse6            = spreadsheetArray[16][11];
        commPosterImage7                = spreadsheetArray[16][12];
        commPosterTrueFalse7            = spreadsheetArray[16][13];
        commPosterImage8                = spreadsheetArray[16][14];
        commPosterTrueFalse8            = spreadsheetArray[16][15];
        commPosterImage9                = spreadsheetArray[16][16];
        commPosterTrueFalse9            = spreadsheetArray[16][17];
        commPosterImage10               = spreadsheetArray[16][18];
        commPosterTrueFalse10           = spreadsheetArray[16][19];
 
        eidDate                         = spreadsheetArray[17][0];
        eidVenue                        = spreadsheetArray[17][1];
        eidAddress                      = spreadsheetArray[17][2];
        eidLecture                      = spreadsheetArray[17][3];
        eidSalaah                       = spreadsheetArray[17][4];
        eidgahTrueFalse                 = spreadsheetArray[17][5];
 
        posterImage1                    = spreadsheetArray[18][0];
        posterTopPadding                = spreadsheetArray[18][1];
        posterHeight                    = spreadsheetArray[18][2];
        posterWidth                     = spreadsheetArray[18][3];
        posterTrueFalse1                = spreadsheetArray[18][4];
        posterImage2                    = spreadsheetArray[18][5];
        posterTrueFalse2                = spreadsheetArray[18][6];
        posterImage3                    = spreadsheetArray[18][7];
        posterTrueFalse3                = spreadsheetArray[18][8];
        posterImage4                    = spreadsheetArray[18][9];
        posterTrueFalse4                = spreadsheetArray[18][10];
        posterImage5                    = spreadsheetArray[18][11];
        posterTrueFalse5                = spreadsheetArray[18][12];
        posterImage6                    = spreadsheetArray[18][13];
        posterTrueFalse6                = spreadsheetArray[18][14];
        posterImage7                    = spreadsheetArray[18][15];
        posterTrueFalse7                = spreadsheetArray[18][16];
        posterImage8                    = spreadsheetArray[18][17];
        posterTrueFalse8                = spreadsheetArray[18][18];
        posterImage9                    = spreadsheetArray[18][19];
        posterTrueFalse9                = spreadsheetArray[18][20];
        posterImage10                   = spreadsheetArray[18][21];
        posterTrueFalse10               = spreadsheetArray[18][22];
 
        bigPosterImage0                 = spreadsheetArray[19][0];
        bigPosterTrueFalse0             = spreadsheetArray[19][1];
        oldTotalPosterDuration          = Number(spreadsheetArray[19][2]);
        bigPosterImage1                 = spreadsheetArray[19][3];
        bigPosterTrueFalse1             = spreadsheetArray[19][4];
        bigPosterImage2                 = spreadsheetArray[19][5];
        bigPosterTrueFalse2             = spreadsheetArray[19][6];
        bigPosterImage3                 = spreadsheetArray[19][7];
        bigPosterTrueFalse3             = spreadsheetArray[19][8];
        bigPosterImage4                 = spreadsheetArray[19][9];
        bigPosterTrueFalse4             = spreadsheetArray[19][10];
        bigPosterImage5                 = spreadsheetArray[19][11];
        bigPosterTrueFalse5             = spreadsheetArray[19][12];
        bigPosterImage6                 = spreadsheetArray[19][13];
        bigPosterTrueFalse6             = spreadsheetArray[19][14];
        bigPosterImage7                 = spreadsheetArray[19][15];
        bigPosterTrueFalse7             = spreadsheetArray[19][16];
        bigPosterImage8                 = spreadsheetArray[19][17];
        bigPosterTrueFalse8             = spreadsheetArray[19][18];
        bigPosterImage9                 = spreadsheetArray[19][19];
        bigPosterTrueFalse9             = spreadsheetArray[19][20];
        bigPosterDuration0              = spreadsheetArray[19][21];
        bigPosterDuration1              = spreadsheetArray[19][22];
        bigPosterDuration2              = spreadsheetArray[19][23];
        bigPosterDuration3              = spreadsheetArray[19][24];
        bigPosterDuration4              = spreadsheetArray[19][25];
        bigPosterDuration5              = spreadsheetArray[19][26];
        bigPosterDuration6              = spreadsheetArray[19][27];
        bigPosterDuration7              = spreadsheetArray[19][28];
        bigPosterDuration8              = spreadsheetArray[19][29];
        bigPosterDuration9              = spreadsheetArray[19][30];
        bigPosterSecondScreen           = stringToBoolean(spreadsheetArray[19][31]);
        //bigPosterInterval               = Number(spreadsheetArray[19][11]);
 
        sheetRefreshRate                = Number(spreadsheetArray[20][0]);
        bankHeader                      = spreadsheetArray[20][1];
        bankName                        = spreadsheetArray[20][2];
        accountName                     = spreadsheetArray[20][3];
        bankCode                        = spreadsheetArray[20][4];
        accountNum                      = spreadsheetArray[20][5];
        bankTrueFalse                   = spreadsheetArray[20][6];
        bankAUSBSB                      = spreadsheetArray[20][7];
 
        sick1                           = spreadsheetArray[21][0];
        sick2                           = spreadsheetArray[21][1];
        sick3                           = spreadsheetArray[21][2];
        sick4                           = spreadsheetArray[21][3];
        sick5                           = spreadsheetArray[21][4];
        sick6                           = spreadsheetArray[21][5];
        sick7                           = spreadsheetArray[21][6];
        sick8                           = spreadsheetArray[21][7];
        sick9                           = spreadsheetArray[21][8];
        sick10                          = spreadsheetArray[21][9];
        sickTrueFalse                   = spreadsheetArray[21][10];
 
        sunsetI                         = checkForError(spreadsheetArray[22][0]);
        eshaStartsI                     = checkForError(spreadsheetArray[22][1]);
        suhurEndsI                      = checkForError(spreadsheetArray[22][2]);
        fajrStartsI                     = checkForError(spreadsheetArray[22][3]);
        sunriseI                        = checkForError(spreadsheetArray[22][4]);
        ishraaqI                        = checkForError(spreadsheetArray[22][5]);
        duhaI                           = checkForError(spreadsheetArray[22][6]);
        istiwaCI                        = checkForError(spreadsheetArray[22][7]);
        istiwaI                         = checkForError(spreadsheetArray[22][8]);
        zuhrStartsI                     = checkForError(spreadsheetArray[22][9]);
        asrShafiI                       = checkForError(spreadsheetArray[22][10]);
        asrHanafiI                      = checkForError(spreadsheetArray[22][11]);
        tickerSpeed                     = Number(spreadsheetArray[22][12]) * 1000;
        //islamicTimeTicker               = stringToBoolean(spreadsheetArray[22][13]);
        tsehriEnds                       = checkForError(spreadsheetArray[22][14]);
        tfajrStarts                      = checkForError(spreadsheetArray[22][15]);
        tsunrise                         = checkForError(spreadsheetArray[22][16]);
        tishraaq                         = checkForError(spreadsheetArray[22][17]);
        tduha                            = checkForError(spreadsheetArray[22][18]);
        tisiwaCaution                    = checkForError(spreadsheetArray[22][19]);
        tistiwa                          = checkForError(spreadsheetArray[22][20]);
        tzuhrStarts                      = checkForError(spreadsheetArray[22][21]);
        tasrShafi                        = checkForError(spreadsheetArray[22][22]);
        tasrHanafi                       = checkForError(spreadsheetArray[22][23]);
        tsunset                          = checkForError(spreadsheetArray[22][24]);
        teshaStarts                      = checkForError(spreadsheetArray[22][25]);

 
        maghribAthanI                   = checkForError(spreadsheetArray[23][0]);
        maghribJamaahI                  = checkForError(spreadsheetArray[23][1]);
        eshaAthanI                      = checkForError(spreadsheetArray[23][2]);
        eshaJamaahI                     = checkForError(spreadsheetArray[23][3]);
        fajrAthanI                      = checkForError(spreadsheetArray[23][4]);
        fajrJamaahI                     = checkForError(spreadsheetArray[23][5]);
        dhuhrAthanI                     = checkForError(spreadsheetArray[23][6]);
        dhuhrJamaahI                    = checkForError(spreadsheetArray[23][7]);
        asrAthanI                       = checkForError(spreadsheetArray[23][8]);
        asrJamaahI                      = checkForError(spreadsheetArray[23][9]);
        masjidName1I                    = spreadsheetArray[23][10];
        masjidName2I                    = spreadsheetArray[23][11];
        masjidNameLanguage              = stringToBoolean(spreadsheetArray[23][12]);
        masjidNameOnTop                 = stringToBoolean(spreadsheetArray[23][13]);

        cBigPosterImage0                = spreadsheetArray[25][0];
        cBigPosterTrueFalse0            = spreadsheetArray[25][1];
        cBigPosterImage1                = spreadsheetArray[25][2];
        cBigPosterTrueFalse1            = spreadsheetArray[25][3];
        cBigPosterImage2                = spreadsheetArray[25][4];
        cBigPosterTrueFalse2            = spreadsheetArray[25][5];
        cBigPosterImage3                = spreadsheetArray[25][6];
        cBigPosterTrueFalse3            = spreadsheetArray[25][7];
        cBigPosterImage4                = spreadsheetArray[25][8];
        cBigPosterTrueFalse4            = spreadsheetArray[25][9];
        cBigPosterImage5                = spreadsheetArray[25][10];
        cBigPosterTrueFalse5            = spreadsheetArray[25][11];
        cBigPosterImage6                = spreadsheetArray[25][12];
        cBigPosterTrueFalse6            = spreadsheetArray[25][13];
        cBigPosterImage7                = spreadsheetArray[25][14];
        cBigPosterTrueFalse7            = spreadsheetArray[25][15];
        cBigPosterImage8                = spreadsheetArray[25][16];
        cBigPosterTrueFalse8            = spreadsheetArray[25][17];
        cBigPosterImage9                = spreadsheetArray[25][18];
        cBigPosterDuration0             = spreadsheetArray[25][19];
        cBigPosterDuration1             = spreadsheetArray[25][20];
        cBigPosterDuration2             = spreadsheetArray[25][21];
        cBigPosterDuration3             = spreadsheetArray[25][22];
        cBigPosterDuration4             = spreadsheetArray[25][23];
        cBigPosterDuration5             = spreadsheetArray[25][24];
        cBigPosterDuration6             = spreadsheetArray[25][25];
        cBigPosterDuration7             = spreadsheetArray[25][26];
        cBigPosterDuration8             = spreadsheetArray[25][27];
        cBigPosterDuration9             = spreadsheetArray[25][28];

        resolve("Write To DOM Complete")
    })
}

// ~~ Write all times to HTML front end ~~ //
function writeTimesToDOM() {
    document.getElementById("fajrNextDate").innerHTML = fajrNextDate;
    document.getElementById("asrNextDate").innerHTML = asrNextDate;
    document.getElementById("eshaNextDate").innerHTML = eshaNextDate;
    document.getElementById("fajrNextTime").innerHTML = (salaah12HourFormat) ? time24to12(fajrNextTime) : fajrNextTime
    document.getElementById("asrNextTime").innerHTML = (salaah12HourFormat) ? time24to12(asrNextTime) : asrNextTime
    document.getElementById("eshaNextTime").innerHTML = (salaah12HourFormat) ? time24to12(eshaNextTime) : eshaNextTime
    document.getElementById("masjidUrl").innerHTML = masjidUrl;
    document.getElementById("announceHead2").innerHTML = announceHead2;
    document.getElementById("announce2").innerHTML = announce2;
    document.getElementById("announceHead3").innerHTML = announceHead3;
    document.getElementById("announce3").innerHTML = announce3;
    document.getElementById("announceHead4").innerHTML = announceHead4;
    document.getElementById("announce4").innerHTML = announce4;
    document.getElementById("announceHead5").innerHTML = announceHead5;
    document.getElementById("announce5").innerHTML = announce5;
    document.getElementById("announceHead6").innerHTML = announceHead6;
    document.getElementById("announce6").innerHTML = announce6;
    document.getElementById("announceHead7").innerHTML = announceHead7;
    document.getElementById("announce7").innerHTML = announce7;
    document.getElementById("announceHead8").innerHTML = announceHead8;
    document.getElementById("announce8").innerHTML = announce8;
    document.getElementById("announceHead9").innerHTML = announceHead9;
    document.getElementById("announce9").innerHTML = announce9;
    document.getElementById("announceHead10").innerHTML = announceHead10;
    document.getElementById("announce10").innerHTML = announce10;
    document.getElementById("commBroadHead").innerHTML = commBroadHead;
    document.getElementById("commBroadText").innerHTML = commBroadText;
    setTimeout(() => { mblTextFit(fitTextGroups) }, 3000)
}





// !====================================================== //
// !=============== MBL Run Once Functions =============== //
// !====================================================== //
// ~~~~~~~ Islamic Date Calculation ~~~~~~ //
function islamicDateDOM() {
    function gmod(n, m) {
        return ((n % m) + m) % m;
    }
    function kuwaiticalendar(adjust) {
        var today = new Date();
        if (adjust) {
            adjustmili = 1000 * 60 * 60 * 24 * adjust;
            todaymili = today.getTime() + adjustmili;
            today = new Date(todaymili);
        }

        var sunst = sunset;
        //var a = tnow.split(":");
        //var timeinsecs = (+a[0]) * 60 * 60 + (+a[1]) * 60 + (+a[2]);
        var b = sunst.split(":");
        var sunss = (+b[0]) * 60 * 60 + (+b[1]) * 60;
        var suns2 = sunss + 185;

        let midn = 86400 - timeInSeconds;
        let timeout = (timeInSeconds < sunss) ? ((sunss + 186) - timeInSeconds) * 1000 : timeInSeconds - midn + sunss * 1000;

        _debug("Islamic Date","Date changing in",new Date(timeout).toISOString().substr(11, 8),logStyles.function,false,false)

        setTimeout(() => {
            islamicDateDOM()
        }, timeout)


        if (timeInSeconds > suns2) { day = today.getDate() + 1; } else { day = today.getDate(); }

        month = today.getMonth();
        year = today.getFullYear();
        m = month + 1;
        y = year;
        if (m < 3) {
            y -= 1;
            m += 12;
        }

        a = Math.floor(y / 100.);
        b = 2 - a + Math.floor(a / 4.);
        if (y < 1583) b = 0;
        if (y == 1582) {
            if (m > 10) b = -10;
            if (m == 10) {
                b = 0;
                if (day > 4) b = -10;
            }
        }
        jd = Math.floor(365.25 * (y + 4716)) + Math.floor(30.6001 * (m + 1)) + day + b - 1524;
        b = 0;
        if (jd > 2299160) {
            a = Math.floor((jd - 1867216.25) / 36524.25);
            b = 1 + a - Math.floor(a / 4.);
        }
        bb = jd + b + 1524;
        cc = Math.floor((bb - 122.1) / 365.25);
        dd = Math.floor(365.25 * cc);
        ee = Math.floor((bb - dd) / 30.6001);
        day = (bb - dd) - Math.floor(30.6001 * ee);
        month = ee - 1;
        if (ee > 13) {
            cc += 1;
            month = ee - 13;
        }
        year = cc - 4716;
        if (adjust) {
            wd = gmod(jd + 1 - adjust, 7) + 1;
        } else {
            wd = gmod(jd + 1, 7) + 1;
        }
        iyear = 10631. / 30.;
        epochastro = 1948084;
        epochcivil = 1948085;
        shift1 = 8.01 / 60.;
        z = jd - epochastro;
        cyc = Math.floor(z / 10631.);
        z = z - 10631 * cyc;
        j = Math.floor((z - shift1) / iyear);
        iy = 30 * cyc + j;
        z = z - Math.floor(j * iyear + shift1);
        im = Math.floor((z + 28.5001) / 29.5);
        if (im == 13) im = 12;
        id = z - Math.floor(29.5001 * im - 29);

        var myRes = new Array(10);
        myRes[0] = day;                           //calculated day (CE)
        myRes[1] = month - 1;                     //calculated month (CE)
        myRes[2] = year;                          //calculated year (CE)
        myRes[3] = jd - 1;                        //julian day number
        myRes[4] = wd - 1;                        //weekday number
        myRes[5] = id;                            //islamic date
        myRes[6] = im - 1;                        //islamic month
        myRes[7] = iy;                            //islamic year
        myRes[8] = 30;                            //force prev islamic date
        myRes[9] = im - 2;                        //force prev islamic month
        if(myRes[9] < 0 && forceDate30 == "N"){   //Check if month is Muharram & roll back to prev year
            myRes[9] = 11
            myRes[7] = iy -1
        }
        return myRes;
    }
    
    function writeIslamicDateEng(adjustment) {
        var wdNames = translate("islamicDayEng");
        var iMonthNames = translate("islamicMonthEng");
        var iDate = kuwaiticalendar(adjustment);

        islamicDateNum  = iDate[5];
        islamicMonthNum = iDate[6] + 1;
        islamicDateYear = iDate[7];

        var outputIslamicDate = wdNames[iDate[4]] + ", " + iDate[5] + " " + iMonthNames[iDate[6]] + " " + iDate[7];
        return outputIslamicDate;
    }

    function writeIslamicDateAra(adjustment) {
        var wdNames = new Array("الأحَد", "الإثْنَين", "الثُلَاثَاء", "الأرْبِعَاء", "الخَمِيْس", "الجُمُعَة", "السَبْت");
        var iMonthNames = new Array("مُحَرَّمْ", "صَفَر", "رَبِيْعُ الْأَوَّل", "رَبِيْعُ الثَّانِي", "جُمَادَى الْاُوْلىٰ", "جُمَادَى الْاُخْرٰى", "رَجَب", "شَعْبان", "رَمَضَان", "شَوَّال", "ذُو الْقَعْدَة", "ذُو الْحِجَّة");
        var iDateArabic = new Array("٠", "١", "٢", "٣", "٤", "٥", "٦", "٧", "٨", "٩", "١٠", "١١", "١٢", "١٣", "١٤", "١٥", "١٦", "١٧", "١٨", "١٩", "٢٠", "٢١", "٢٢", "٢٣", "٢٤", "٢٥", "٢٦", "٢٧", "٢٨", "٢٩", "٣٠")
        var iDate = kuwaiticalendar(adjustment);

        islamicDateNum  = iDate[5];
        islamicMonthNum = iDate[6] + 1;
        islamicDateYear = iDate[7];

        var outputIslamicDate = wdNames[iDate[4]] + " ، " + iDateArabic[iDate[5]] + " " + iMonthNames[iDate[6]] + " " + convertArabicNumber(iDate[7]);
        return outputIslamicDate;
    }

    function forceIslamicDate30Eng(adjustment) {
        var wdNames = translate("islamicDayEng");
        var iMonthNames = translate("islamicMonthEng");
        var iDate = kuwaiticalendar(0);

        islamicDateNum  = 30;
        islamicMonthNum = iDate[9];
        islamicDateYear = iDate[7];

        var forceIslamicDate = wdNames[iDate[4]] + ", " + islamicDateNum + " " + iMonthNames[iDate[9]] + " " + iDate[7];
        return forceIslamicDate;
    }

    function forceIslamicDate30Ara(adjustment) {
        var wdNames = new Array("الأحَد", "الإثْنَين", "الثُلَاثَاء", "الأرْبِعَاء", "الخَمِيْس", "الجُمُعَة", "السَبْت");
        var iMonthNames = new Array("مُحَرَّمْ", "صَفَر", "رَبِيْعُ الْأَوَّل", "رَبِيْعُ الثَّانِي", "جُمَادَى الْاُوْلىٰ", "جُمَادَى الْاُخْرٰى", "رَجَب", "شَعْبان", "رَمَضَان", "شَوَّال", "ذُو الْقَعْدَة", "ذُو الْحِجَّة");
        var iDateArabic = new Array("٠", "١", "٢", "٣", "٤", "٥", "٦", "٧", "٨", "٩", "١٠", "١١", "١٢", "١٣", "١٤", "١٥", "١٦", "١٧", "١٨", "١٩", "٢٠", "٢١", "٢٢", "٢٣", "٢٤", "٢٥", "٢٦", "٢٧", "٢٨", "٢٩", "٣٠")
        var iDate = kuwaiticalendar(0);

        islamicDateNum  = 30;
        islamicMonthNum = iDate[9];
        islamicDateYear = iDate[7];

        var forceIslamicDate = wdNames[iDate[4]] + " ، " + iDateArabic[islamicDateNum] + " " + iMonthNames[iDate[9]] + " " + convertArabicNumber(iDate[7]);
        return forceIslamicDate;
    }

    let Gsheetadj = parseInt(islamicDateAdjust);
    if (forceDate30 === "N") {
        islamicDateAra = forceIslamicDate30Ara(Gsheetadj);
        islamicDateEng = forceIslamicDate30Eng(Gsheetadj);
    } else {
        islamicDateAra = writeIslamicDateAra(Gsheetadj);
        islamicDateEng = writeIslamicDateEng(Gsheetadj);
    }

    document.getElementById("islamicDate").innerHTML = islamicDateEng;

    /*let iDate     = kuwaiticalendar(parseInt(islamicDateAdjust));
    if (forceDate30 === "N") {
        islamicDateDay   = iDate[4];
        islamicDateDate  = iDate[5];
        islamicDateMonth = iDate[6];
        islamicDateYear  = iDate[7];
    } else {
        islamicDateDay   = iDate[4];
        islamicDateDate  = iDate[8];
        islamicDateMonth = iDate[9];
        islamicDateYear  = iDate[7];
    }*/
}

// ~~~~~~~~~~~~~ Todays Date ~~~~~~~~~~~~~ //
function todaysDate() {
    var weekdayNames     = translate('weekdayNames');
    var monthNames       = translate('monthNames');
    var weekdayShortname = translate('weekdayShortname')
    var monthShortname   = translate('monthShortname')
    
    document.getElementById("todaysDate").innerHTML = (weekdayNames[currentDay] + ", " + fullDate + " " + monthNames[fullMonth] + " " + fullYear);
    document.getElementById("popUpDate").innerHTML  = "<span>"+(weekdayShortname[currentDay] + ", " + fullDate + " " + monthShortname[fullMonth])+"</span>"
}

// ~~ Masjid Name & Big Clock positions ~~ //
function clocknamePositions() {
    if (isBigPopUpShowing) return;
    if (clockNameFormat == 1) {
        $("#masjidNameTable").css({ "top": '207px', });
        if(masjidNameCounterG != 1){
            $("#tableTime"  ).css({ "top": '125px', })
        }
        $("#adhanColumn").css({ "color": 'var(--adhanHeading1)'})
        $("#msttop").css({ "color": 'var(--jamaahHeading1)'})
        $("#jamaahColumn").css({ "color": 'var(--jamaahHeading1)'})
        //document.getElementById("#jamaahColumn").style.color = 'var(--jamaahHeading1)'
    } else {
        $("#masjidNameTable").css({ "top": masjidname1TopPad + 'px'});
        $("#tableTime").css({ "top": bigClockTopPad + 'px'})
        $("#adhanColumn").css({ "color": 'var(--adhanHeading0)'})
        $("#msttop").css({ "color": 'var(--jamaahHeading0)'})
        $("#jamaahColumn").css({ "color": 'var(--jamaahHeading0)'})
        //document.getElementById("#jamaahColumn").style.color = 'var(--jamaahHeading0)'
    }
}

//  Next salaah change highlights 1 day before  //
function nextSalaahChangeFlash() {
    if(!nextSalaahFlash) return;

    var ftd = (parseInt(fajrNextMilli) - parseInt(timeZoneMilli)) - oneDayMilli;
    var atd = (parseInt(asrNextMilli) - parseInt(timeZoneMilli)) - oneDayMilli;
    var etd = (parseInt(eshaNextMilli) - parseInt(timeZoneMilli)) - oneDayMilli;
    var fod = (parseInt(fajrNextMilli) - parseInt(timeZoneMilli)) + oneDayMilli;
    var aod = (parseInt(asrNextMilli) - parseInt(timeZoneMilli)) + oneDayMilli;
    var eod = (parseInt(eshaNextMilli) - parseInt(timeZoneMilli)) + oneDayMilli;


    //Fajr
    if (ftd < nowDate && nowDate < fod) {
        document.getElementById("fcc").classList.add("flashnextfajr");
        document.getElementById("fccs").classList.add("flashnextfajr");
        document.getElementById("fccss").classList.add("flashnextfajr");
        document.getElementById("nsch_fajr").classList.add("flashnextfajr");
        document.getElementById("fajrNextDate").classList.add("flashnextfajr");
        document.getElementById("fajrNextTime").classList.add("flashnextfajr");
    }
    //Asr
    if (atd < nowDate && nowDate < aod) {
        document.getElementById("acc").classList.add("flashnextasr");
        document.getElementById("accs").classList.add("flashnextasr");
        document.getElementById("accss").classList.add("flashnextasr");
        document.getElementById("nsch_asr").classList.add("flashnextasr");
        document.getElementById("asrNextDate").classList.add("flashnextasr");
        document.getElementById("asrNextTime").classList.add("flashnextasr");
    }
    //Esha
    if (etd < nowDate && nowDate < eod) {
        document.getElementById("ecc").classList.add("flashnextesha");
        document.getElementById("eccs").classList.add("flashnextesha");
        document.getElementById("eccss").classList.add("flashnextesha");
        document.getElementById("nsch_esha").classList.add("flashnextesha");
        document.getElementById("eshaNextDate").classList.add("flashnextesha");
        document.getElementById("eshaNextTime").classList.add("flashnextesha");
    }
}

// ~~~~~~~~~ Starting the ticker ~~~~~~~~~ //
let tickerCounter = 1
function tickerAnimate() {

    let hasTimePassed = (time) => (minutestoseconds(time) < timeInSeconds)

    let sehriEndsTime       = (hasTimePassed(sehriEnds))   ? ((salaah12HourFormat) ? time24to12(tsehriEnds)    : tsehriEnds )   : ((salaah12HourFormat) ? time24to12(sehriEnds)   : sehriEnds)
    let fajrStartsTime      = (hasTimePassed(fajrStarts))  ? ((salaah12HourFormat) ? time24to12(tfajrStarts)   : tfajrStarts)   : ((salaah12HourFormat) ? time24to12(fajrStarts)  : fajrStarts)
    let sunriseTime         = (hasTimePassed(sunrise))     ? ((salaah12HourFormat) ? time24to12(tsunrise)      : tsunrise)      : ((salaah12HourFormat) ? time24to12(sunrise)     : sunrise)
    let ishraaqTime         = (hasTimePassed(ishraaq))     ? ((salaah12HourFormat) ? time24to12(tishraaq)      : tishraaq)      : ((salaah12HourFormat) ? time24to12(ishraaq)     : ishraaq)
    let duhaTime            = (hasTimePassed(duha))        ? ((salaah12HourFormat) ? time24to12(tduha)         : tduha)         : ((salaah12HourFormat) ? time24to12(duha)        : duha)
    let istiwaTime          = (hasTimePassed(istiwa))      ? ((salaah12HourFormat) ? time24to12(tisiwaCaution) : tisiwaCaution) : ((salaah12HourFormat) ? time24to12(istiwa)      : istiwa)
    let zawaalEndTime       = (hasTimePassed(zawaalEnd))   ? ((salaah12HourFormat) ? time24to12(tistiwa)       : tistiwa)       : ((salaah12HourFormat) ? time24to12(zawaalEnd)   : zawaalEnd)
    let sundayDhuhrTime     = (hasTimePassed(sundayDhuhr)) ? ((salaah12HourFormat) ? time24to12(tzuhrStarts)   : tzuhrStarts)   : ((salaah12HourFormat) ? time24to12(sundayDhuhr) : sundayDhuhr)
    let asrShafiTime        = (hasTimePassed(asrShafi))    ? ((salaah12HourFormat) ? time24to12(tasrShafi)     : tasrShafi)     : ((salaah12HourFormat) ? time24to12(asrShafi)    : asrShafi)
    let asrHanafiTime       = (hasTimePassed(asrHanafi))   ? ((salaah12HourFormat) ? time24to12(tasrHanafi)    : tasrHanafi)    : ((salaah12HourFormat) ? time24to12(asrHanafi)   : asrHanafi)
    let sunsetTime          = (hasTimePassed(sunset))      ? ((salaah12HourFormat) ? time24to12(tsunset)       : tsunset)       : ((salaah12HourFormat) ? time24to12(sunset)      : sunset)
    let eshaStartsTime      = (hasTimePassed(eshaStarts))  ? ((salaah12HourFormat) ? time24to12(teshaStarts)   : teshaStarts)   : ((salaah12HourFormat) ? time24to12(eshaStarts)  : eshaStarts)



    let timesInfo = {
        fajrSuhur: {
            name: ["", translate("fajrSuhurH")],
            time: [convertArabicNumber(0, "ticker"), sehriEndsTime]
        },
        fajrStarts: {
            name: ["الصبح الصادق", translate("fajrStartsH")],
            time: [convertArabicNumber(fajrStartsI, "ticker"), fajrStartsTime]
        },
        sunrise: {
            name: ["طلوع الشمس", translate("sunriseH")],
            time: [convertArabicNumber(sunriseI, "ticker"), sunriseTime]
        },
        ishraaq: {
            name: ["الإشراق", translate("ishraaqH")],
            time: [convertArabicNumber(ishraaqI, "ticker"), ishraaqTime]
        },
        duha: {
            name: ["الضحى", translate("duhaH")],
            time: [convertArabicNumber(duhaI, "ticker"), duhaTime]
        },
        chast: {
            name: ["چاست", translate("chast")],
            time: [0,0]
        },
        istiwa: {
            name: ["إستواء الشمس", translate("istiwaH")],
            time: [convertArabicNumber(istiwaI, "ticker"), istiwaTime]
        },
        zawaalEnd: {
            name: ["", translate("zawaalEndH")],
            time: [convertArabicNumber(0, "ticker"), zawaalEndTime]
        },
        sundayzuhr: {
            name: ["الظهر (كل يوم)", translate("sundayzuhrH")],
            time: [convertArabicNumber(0, "ticker"), sundayDhuhrTime]
        },
        asrShafi: {
            name: ["(الشافعي) العصر", translate("asrShafiH")],
            time: [convertArabicNumber(asrShafiI, "ticker"), asrShafiTime]
        },
        asrHanafi: {
            name: ["(الحنفي) العصر", translate("asrHanafiH")],
            time: [convertArabicNumber(asrHanafiI, "ticker"), asrHanafiTime]
        },
        sunset: {
            name: ["غروب الشمس", translate("sunsetH")],
            time: [convertArabicNumber(sunsetI, "ticker"), sunsetTime]
        },
        eshaStarts: {
            name: ["العشاء", translate("eshaStartsH")],
            time: [convertArabicNumber(eshaStartsI, "ticker"), eshaStartsTime]
        }
    }
    let container = document.getElementById("outer");
    let box = document.getElementById("tick");
    let engCSS = { "font-family": "'Open Sans'   , serif", "font-weight": "700", "font-size": "35px" };
    let araCSS = { "font-family": "'Scheherazade', serif", "font-weight": "400", "font-size": "60px" };

    let engTicker = `   <li><h1 class="tickerOdd" id="fajrSuhurH">${timesInfo.fajrSuhur.name[tickerCounter]}&nbsp;<span id="fajrSuhur"></span>&nbsp;&nbsp;<span id="sehriEnds">${timesInfo.fajrSuhur.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="fajrStartsH">${timesInfo.fajrStarts.name[tickerCounter]}&nbsp;&nbsp;<span id="fajrStarts">${timesInfo.fajrStarts.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="sunriseH">${timesInfo.sunrise.name[tickerCounter]}&nbsp;&nbsp;<span id="sunrise">${timesInfo.sunrise.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="ishraaqH">${timesInfo.ishraaq.name[tickerCounter]}&nbsp;&nbsp;<span id="ishraaq">${timesInfo.ishraaq.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="duhaH">${timesInfo.duha.name[tickerCounter]}&nbsp;${(chast)? timesInfo.chast.name[tickerCounter]+'&nbsp;' : ''}&nbsp;<span id="duha">${timesInfo.duha.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="istiwaH">${timesInfo.istiwa.name[tickerCounter]}&nbsp;&nbsp;<span id="istiwa">${timesInfo.istiwa.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="zawaalEndH">${timesInfo.zawaalEnd.name[tickerCounter]}&nbsp;&nbsp;<span id="zawaalEnd">${timesInfo.zawaalEnd.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="sundayzuhrH">${timesInfo.sundayzuhr.name[tickerCounter]}${(sundayzuhr=="") ? "" : "&nbsp;"}<span id="sundayzuhr">${translate(sundayzuhr.replace(/\s/g, "").replace(/\(/g, "").replace(/\)/g, "").replace(/&/g, "").toLowerCase())}</span>&nbsp;<span id="sundayDhuhr">&nbsp;${timesInfo.sundayzuhr.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="asrShafiH">${timesInfo.asrShafi.name[tickerCounter]}&nbsp;&nbsp;<span id="asrShafi">${timesInfo.asrShafi.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="asrHanafiH">${timesInfo.asrHanafi.name[tickerCounter]}&nbsp;&nbsp;<span id="asrHanafi">${timesInfo.asrHanafi.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="sunsetH">${timesInfo.sunset.name[tickerCounter]}&nbsp;&nbsp;<span id="sunset">${timesInfo.sunset.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="eshaStartsH">${timesInfo.eshaStarts.name[tickerCounter]}&nbsp;&nbsp;<span id="eshaStarts">${timesInfo.eshaStarts.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>`;

    let araTicker = `   <li><h1 class="tickerOdd" id="asrHanafiH">${timesInfo.asrHanafi.name[tickerCounter]}&nbsp;&nbsp;<span id="asrHanafi">${timesInfo.asrHanafi.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="asrShafiH">${timesInfo.asrShafi.name[tickerCounter]}&nbsp;&nbsp;<span id="asrShafi">${timesInfo.asrShafi.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="istiwaH">${timesInfo.istiwa.name[tickerCounter]}&nbsp;&nbsp;<span id="istiwa">${timesInfo.istiwa.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="duhaH">${timesInfo.duha.name[tickerCounter]}&nbsp;&nbsp;<span id="duha">${timesInfo.duha.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="ishraaqH">${timesInfo.ishraaq.name[tickerCounter]}&nbsp;&nbsp;<span id="ishraaq">${timesInfo.ishraaq.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerOdd" id="sunriseH">${timesInfo.sunrise.name[tickerCounter]}&nbsp;&nbsp;<span id="sunrise">${timesInfo.sunrise.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="fajrStartsH">${timesInfo.fajrStarts.name[tickerCounter]}&nbsp;<span id="fajrSuhur"></span> <span id="fajrStarts">${timesInfo.fajrStarts.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1><span class="hidden" id="fajrSuhur"></span></li>
                        <li><h1 class="tickerOdd" id="eshaStartsH">${timesInfo.eshaStarts.name[tickerCounter]}&nbsp;&nbsp;<span id="eshaStarts">${timesInfo.eshaStarts.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>
                        <li><h1 class="tickerEven" id="sunsetH">${timesInfo.sunset.name[tickerCounter]}&nbsp;&nbsp;<span id="sunset">${timesInfo.sunset.time[tickerCounter]}</span>&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;</h1></li>`;
    let mainTicker = [araTicker, engTicker]

    box.innerHTML = mainTicker[tickerCounter]
    let width = box.scrollWidth;

    function startAnimation() {
        let duration = tickerSpeed
        console.log("Ticker Duration: "+duration)
        setTimeout(() => {
            tickerCounter++;
            if (!islamicTime || !islamicTimeTicker) tickerCounter = 1;
            (tickerCounter >= 2) ? tickerCounter = 0: tickerCounter = tickerCounter
            tickerAnimate()
        }, duration-50)
        if (tickerCounter == 0) {
            fajrSuhurCheck(3)
            $("#outer").css({ "top": "976px", "height": "100px" })
            $("#tick").find('*').css(araCSS)
            width = box.scrollWidth;
            //console.log((width/1920)*100)
            box.animate(
                [
                    { transform: `translateX(-${(width/1920)*100}%)` },
                    { transform: `translateX(100%)` }
                ], {
                    duration: duration,
                    iterations: 1,
                    easing: "linear"
                }
            )
        } else {
            fajrSuhurCheck(4)
            $("#outer").css({ "top": "995px", "height": "85px" })
            $("#tick").find('*').css(engCSS)
            width = box.scrollWidth;
            //console.log((width/1920)*100)
            box.animate(
                [
                    { transform: `translateX(100%)` },
                    { transform: `translateX(-${(width/1920)*100}%)` }
                ], {
                    duration: duration,
                    iterations: 1,
                    easing: "linear"
                }
            )
        }
    }
    startAnimation()
}

// ~~~~~~~~~~~~ MoonPhase Text ~~~~~~~~~~~ //
function moonPhaseText() {
    let currentDate = new Date();
    let day = currentDate.getDate();
    let month = currentDate.getMonth() + 1;
    let year = currentDate.getFullYear();
    let fullDate = day + "-" + month + "-" + year;
    
    let moonBirthScriptDate = new Date(moonBirthScript);
    currentDate.setHours(0, 0, 0, 0); // Reset today's time to 00:00:00
    let diffInTime = moonBirthScriptDate.getTime() - currentDate.getTime();
    let diffInDays = diffInTime / (1000 * 3600 * 24); // Convert time difference from milliseconds to days
    let isNewMoonMoreThan5Days = diffInDays > 10;

    function getMoonPhase(day, month, year) {
        let c = 0;
        let e = 0;
        let jd = 0;
        let b = 0;

        day = parseInt(day, 10);
        month = parseInt(month, 10);
        year = parseInt(year, 10);

        if (month < 3) {
            year--;
            month += 12;
        }
        ++month;
        c = 365.25 * year;
        e = 30.6 * month;
        jd = c + e + day - 694039.09; //jd is total days elapsed
        jd /= 29.5305882; //divide by the moon cycle
        b = parseInt(jd); //int(jd) -> b, take integer part of jd
        jd -= b; //subtract integer part to leave fractional part of original jd
        b = (jd * 100).toFixed(4);
        
        if (b > 0 && b <= 0.75) { //0
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/27.png";
            return translate('waningCrescent')
        } else if (b > 0.75 && b <= 3.3863) { //0.5
            if(isNewMoonMoreThan5Days){
                document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/3.png";
                return translate('waxingCrescent')
            } else{
                document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/3.png";
                return translate('newMoon') + moonBirth
            }
        } else if (b > 3.3863 && b < 6.7726) { //1
            if(isNewMoonMoreThan5Days){
                document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/3.png";
                return translate('waxingCrescent')
            } else{
                document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/3.png";
                return translate('newMoon') + moonBirth
            }
        } else if (b >= 6.7726 && b < 10.1590) { //2
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/3.png";
            return translate('waxingCrescent')
        } else if (b >= 10.1590 && b < 13.5453) { //3
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/3.png";
            return translate('waxingCrescent')
        } else if (b >= 13.5453 && b < 16.9316) { //4
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/4.png";
            return translate('waxingCrescent')
        } else if (b >= 16.9316 && b < 20.3179) { //5
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/5.png";
            return translate('waxingCrescent')
        } else if (b >= 20.3179 && b < 23.7042) { //6
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/6.png";
            return translate('waxingCrescent')
        } else if (b >= 23.7042 && b < 25) { //7
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/7.png";
            return translate('firstQuarter')
        } else if (b >= 25 && b < 27.0906) { //7
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/7.png";
            return translate('waxingGibbous')
        } else if (b >= 27.0906 && b < 30.4769) { //8
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/7.png";
            return translate('waxingGibbous')
        } else if (b >= 30.4769 && b < 33.8632) { //9
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/9.png";
            return translate('waxingGibbous')
        } else if (b >= 33.8632 && b < 37.2495) { //10
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/9.png";
            return translate('waxingGibbous')
        } else if (b >= 37.2495 && b < 40.6358) { //11
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/10.png";
            return translate('waxingGibbous')
        } else if (b >= 40.6358 && b < 44.0222) { //12
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/10.png";
            return translate('waxingGibbous')
        } else if (b >= 44.0222 && b < 47.4085) { //13
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/13.png";
            return translate('waxingGibbous')
        } else if (b >= 47.4085 && b < 50) { //14-50.7948
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/14.png";
            return translate('waxingGibbous')
        } else if (b >= 50 && b < 54.1811) { //15
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/14.png";
            return translate('fullMoon')
        } else if (b >= 54.1811 && b < 57.5674) { //16
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/16.png";
            return translate('waningGibbous')
        } else if (b >= 57.5674 && b < 60.9537) { //17
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/17.png";
            return translate('waningGibbous')
        } else if (b >= 60.9537 && b < 64.3401) { //18
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/18.png";
            return translate('waningGibbous')
        } else if (b >= 64.3401 && b < 67.7264) { //19
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/19.png";
            return translate('waningGibbous')
        } else if (b >= 67.7264 && b < 71.1127) { //20
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/20.png";
            return translate('waningGibbous')
        } else if (b >= 71.1127 && b < 74.4990) { //21
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/21.png";
            return translate('waningGibbous')
        } else if (b >= 74.4990 && b < 77.8853) { //22
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/21.png";
            return translate('lastQuarter')
        } else if (b >= 77.8853 && b < 81.2717) { //23
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/22.png";
            return translate('waningCrescent')
        } else if (b >= 81.2717 && b < 84.6580) { //24
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/23.png";
            return translate('waningCrescent')
        } else if (b >= 84.6580 && b < 88.0443) { //25
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/23.png";
            return translate('waningCrescent')
        } else if (b >= 88.0443 && b < 91.4306) { //26
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/24.png";
            return translate('waningCrescent')
        } else if (b >= 91.4306 && b < 94.8169) { //27
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/24.png";
            return translate('waningCrescent')
        } else if (b >= 94.8169 && b < 98.2033) { //28
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/26.png";
            return translate('waningCrescent')
        } else if (b >= 98.2033 && b <= 100) { //29
            document.getElementById("moonImage").src = "https://premium.masjidboardlive.com/v2/images/moon_images/27.png";
            return translate('waningCrescent')
        }

    };
    document.getElementById("currentMoon").innerHTML = getMoonPhase(day, month, year);
}

// ~~~~~~~~~~~~~ Theme Select ~~~~~~~~~~~~ //
function themeChange() {
    if (window.location.href.indexOf("demo") > -1) return
    document.getElementById('theme').href = `https://premium.masjidboardlive.com/v2/themes/${theme.toLowerCase()}.css`;
}

// ~~~~~~~ Function to reload page ~~~~~~~ //
function pageReload() {
    if (pageHardReload == true) {
        window.location.reload(true);
    } else {
        //var reload = 87120 - timeInSeconds;
        let reload = ((86400 + 600) + (Math.floor(Math.random() * 300) + 1)) - timeInSeconds
        setTimeout(() => {
            console.warn("Reloading Page...")
            window.location.reload(true);
        }, reload * 1000)
        _debug("Page reload","reloading in",new Date(reload * 1000).toISOString().substr(11, 8),logStyles.function,false,false)
    }
    var newNode = document.createElement('div');
    //if (debug) 
}

//  function to check if suhur is the same as fajr starts, if it is it will remove (x mins early)  //
function fajrSuhurCheck(t) {
    var sehriEndsSec = minutestoseconds(sehriEnds);
    var fajrStartsSec = minutestoseconds(fajrStarts);

    if (sehriEndsSec == fajrStartsSec) {
        document.getElementById('pu_fajrSuhur').innerHTML = `(${translate("puh_fajrSuhurFinal")})`
        if (t == 1) return `<p id="countdownInfo" class='countdownInfo' >(${translate("puh_fajrSuhurFinal")})</p>`
        else if (t == 2) return `<p style='font-size: 16px; margin-top: -4px;'>(${translate("puh_fajrSuhurFinal")})</p>`
        else if (t == 3) document.getElementById('fajrSuhur').innerHTML = "";
        else if (t == 4) document.getElementById('fajrSuhur').innerHTML = `(${translate("puh_fajrSuhurFinal")})`;
    } else {
        var timeDiff = (fajrStartsSec - sehriEndsSec) / 60;
        document.getElementById('pu_fajrSuhur').innerHTML = `(${timeDiff} ${translate("puh_fajrSuhurMinutes")})`
        if (t == 1) return `<p id="countdownInfo" class='countdownInfo'> (${timeDiff} ${translate("puh_fajrSuhurMinutes")})</p>`
        else if (t == 2) return `<p style='font-size: 16px; margin-top: -4px;'>(${timeDiff} ${translate("puh_fajrSuhurMinutes")})</p>`
        else if (t == 3) document.getElementById('fajrSuhur').innerHTML = "";
        else if (t == 4) document.getElementById('fajrSuhur').innerHTML = `(${timeDiff} ${translate("fajrSuhurMins")})`;
    }
}

function turnOffBoard(configuration){
    console.log(configuration)
    if(configuration.length == 0 || configuration == null || configuration == "" || configuration == 0 || mobile) return
    
    let timeNow      = currentTime("s")
    let config       = configuration.split(",")
    let cont         = document.getElementById("boardDisable")
    let time         = document.getElementById("currentTime")
    let maghribTime  = (minutestoseconds(maghribJamaah) == "") ? minutestoseconds(maghribAthan) : minutestoseconds(maghribJamaah);
    let jumuahTime   = jumuahTime3 

    if(jumuahHeadingsArray.length > 1 && jumuahHeadingsArray != "#N/A"){
        let newJumuahHeading = jumuahHeadingsArray.split(",")
             if(['2', '8'          ].includes(newJumuahHeading[0]) && timeNow < minutestoseconds(jumuahTime1)) {jumuahTime = jumuahTime1; console.log("1");}
        else if(['4', '9', '2', '8'].includes(newJumuahHeading[1]) && timeNow < minutestoseconds(jumuahTime2)) {jumuahTime = jumuahTime2; console.log("2");}
        else if(['6', '5', '10'    ].includes(newJumuahHeading[2]) && timeNow < minutestoseconds(jumuahTime3)) {jumuahTime = jumuahTime3; console.log("3");}
        else                                                                                                   {jumuahTime = jumuahTime3; console.log("else");}
    }

    //                                        "Fajr"              |                            "Jumuah"               |         "Dhuhr"              |          "Asr"             |       "Maghrib"               |       "Esha"
    let times        = [0           , Number(config[1])           , (currentDay == 5) ? Number(config[5])             : Number(config[3])            , Number(config[7])          , Number(config[9])             , Number(config[11])          , 0     ]
    let salaahs      = [false       , stringToBoolean(config[0])  , (currentDay == 5) ? stringToBoolean(config[4])    : stringToBoolean(config[2])   , stringToBoolean(config[6]) , stringToBoolean(config[8])    , stringToBoolean(config[10]) , false ]
    let stimes       = [0           , minutestoseconds(fajrJamaah), (currentDay == 5) ? minutestoseconds(jumuahTime)  : minutestoseconds(dhuhrJamaah), minutestoseconds(asrJamaah), maghribTime                   , minutestoseconds(eshaJamaah), 86400 ]

    for(let h=0; h<salaahs.length; h++){
        if(!salaahs[h]) continue;
        let jamaahTime1 = stimes[h]
        let jamaahTime2 = stimes[h-1]
        if(timeNow < jamaahTime1 && timeNow > jamaahTime2){
            let calcTime = jamaahTime1 - timeNow;

            //Time of day opacity
            let tod = 95;
                 if(calcTime == stimes[0]) tod = 999; //Fajr
            else if(calcTime == stimes[2]) tod = 999; //Dhuhr/Jumuah
            else if(calcTime == stimes[3]) tod = 999; //Asr
            else if(calcTime == stimes[4]) tod = 999; //Maghrib
            else if(calcTime == stimes[5]) tod = 999; //Esha
            else tod = 95;                            //None

            _debug("Disabling Board","in: ",calcTime+"s - "+secondstotime(calcTime),logStyles.function,false,false)
            setTimeout(()=>{
                _debug("Re-enabling Board","in: ",times[h]+"s - "+secondstotime(times[h]),logStyles.function,false,false)
                isBoardDisabled = true
                hidePoster("blackout","Screen Blackout", true)
                cont.style.display    = "block"
                cont.style.visibility = "visible"
                cont.style.background = `rgba(0,0,0,0.${tod});`
                time.style.color      = "#d4d4d4"
                time.style.textShadow = "0px 0 0 transparent"
                setTimeout(()=>{cont.style.opacity = 1},500)
                setTimeout(()=>{
                    cont.style.opacity    = 0
                    time.style.color      = "var(--bigClock)"
                    time.style.textShadow = "4px 0 0 var(--clockShadow)"
                    setTimeout(()=>{cont.style.display = "none",cont.style.visibility = "hidden"},700)
                    isBoardDisabled = false
                    turnOffBoard(disableBoard)
                },times[h]*1000)
            },calcTime*1000)
        }
    }
}

let fitTextGroups      = [
    {
        ids: ['masjidNameTable'],
        inner: 'p',
        min: 25,
        widthOnly: false,
        originalFS: 36
    },
    {
        ids: ['jumuahAdhaanHeading', 'jumuahSunnanHeading', 'jumuahKhutbahHeading'],
        inner: 'h2',
        min: 27,
        widthOnly: true,
        originalFS: 30
    }, {
        ids: ['currentMoonBox'],
        inner: 'h1',
        min: 25,
        widthOnly: false,
        originalFS: 35
    }, {
        ids: ['jumuahKhateebHeading'],
        inner: 'span',
        min: 10,
        widthOnly: false,
        originalFS: 30
    }
]
let namedFitTextGroups = {
    hadithSlide: [{
        ids: ['hadithC'],
        inner: 'span',
        min: 16,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33
    }],
    ayahSlide: [{
        ids: ['ayahC'],
        inner: 'span',
        min: 16,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33
    }],
    sunnahSlide: [{
        ids: ['sunnahC'],
        inner: 'p',
        min: 16,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33
    }],
    taleemSlide: [{
        ids: ['taleemSlide'],
        inner: 'span',
        min: 16,
        height: 515,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33,
        hasChildren: true
    }],
    nikahSlide: [{
        ids: ['nikahSlide'],
        inner: 'span',
        min: 10,
        height: 575,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33,
        hasChildren: true
    }],
    nikahPopUp: [{
        ids:["nikahPopUp"],
        inner: 'span',
        min: 10,
        height: 550,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33,
        hasChildren: true
    }],
    eidSlide: [{
        ids: ['eidSlide'],
        inner: 'span',
        min: 13,
        widthOnly: false,
        height: 580,
        isTextFitted: false,
        originalFS: 33,
        hasChildren: true
    }],
    bankSlide: [{
        ids: ['bankSlide'],
        inner: 'span',
        min: 13,
        widthOnly: false,
        height: 580,
        isTextFitted: false,
        originalFS: 33,
        hasChildren: true
    }],
    sickSlide: [{
        ids: ['sickSlideC'],
        inner: 'div',
        min: 16,
        height: 518,
        width: 445,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33
    }],
    dawahSlide: [{
        ids: ['dawahSlide'],
        inner: 'span',
        min: 16,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33,
        height: 580,
        hasChildren: true
    }],
    funeralBoard: [{
        ids: ['funeralBoard'],
        inner: 'span',
        min: 5,
        height: 600,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 33,
        hasChildren: true
    }],
    newMoonPopUp: [{
        ids: ['moonBirthHC', 'moonBirthC', 'moonAgeHC', 'moonAgeC', 'moonSetHC', 'moonSetC', 'moonVisibilHC', 'moonVisibilC', 'moonAzimuthHC', 'moonAzimuthC', 'moonAltitudeHC', 'moonAltitudeC',],
        inner: 'h3',
        min: 20,
        widthOnly: false,
        isTextFitted: false,
        originalFS: 30
    }],
    popUpOwwalTimes: [{
        ids: ['puh_sehriEnds', 'puh_fajrStarts','puh_sunrise','puh_ishraaq','puh_duha','puh_istiwa','puh_sundayDhuhr','puh_asrShafi','puh_asrHanafi','puh_sunset','puh_eshaStarts',
              'pu_sehriEnds' , 'pu_fajrStarts' ,'pu_sunrise' ,'pu_ishraaq' ,'pu_duha' ,'pu_istiwa' ,'pu_sundayDhuhr' ,'pu_asrShafi' ,'pu_asrHanafi' ,'pu_sunset' ,'pu_eshaStarts'],
        inner: 'span',
        min: 5,
        widthOnly: true,
        width: 200,
        isTextFitted: false,
        originalFS: 36
    }],
    timesBoxHeader: [{
        ids: ['timesBoxHeadC'],
        inner: 'h1',
        min: 25,
        widthOnly: true,
        isTextFitted: false,
        originalFS: 35
    }],
    dates: [{
        ids: ['todaysDateC', 'islamicDateC'],
        inner: 'p',
        min: 30,
        widthOnly: true,
        isTextFitted: false,
        originalFS: 35
    }],
    popUpDate: [{
        ids: ['popUpDate'],
        inner: 'span',
        min: '20',
        width: 354,
        widthOnly: true,
        isTextFitted: false,
        originalsFS: 43
    }]
}

// ~~~~~~~~~~ Text fit functions ~~~~~~~~~~ //
function textFit(id, inner, min, width, height, widthOnly, originalFS, hasChildren, callback) {

    let box = $("#" + id)
    let text = $("#" + id + ">" + inner)

    //reset to original font size
    text.css({ fontSize: `${originalFS}px` })

    //get initial sizings
    let boxWidth = (width) ? width : box.innerWidth();
    let boxHeight = (height) ? height : box.innerHeight();

    let textWidth = Math.max(text.innerWidth(), text[0].scrollWidth);
    let textHeight = Math.max(text.innerHeight(), text[0].scrollHeight);
    let initTextHeight = 0
    let initTextWidth = 0
    if(hasChildren){
        text.children().each(function(){
            textWidth  = $(this).innerWidth();
            textHeight = textHeight + $(this).innerHeight();
            initTextWidth = textWidth
            initTextHeight = textHeight
        });
    }

    let intFS = parseInt(text.css('font-size'), 10)


    //Width ONLY
    if (widthOnly) {
        //run loop to drop font size till it fits in width
        while (boxWidth < textWidth) {
            //get curr font size
            let CurrFontSize = parseInt(text.css('font-size'), 10)
                //check iif text is too small
            if (CurrFontSize == 1 || CurrFontSize <= min) {
                if(debug && CurrFontSize <= min) console.log(`Reached minimum font size ${min}px, cant text fix`);
                if(debug && CurrFontSize == 1) console.log(`Cant text fit lower than 1px`); 
                break;
            }
                //-1 from font size
            //text.css({ fontSize: `${(CurrFontSize - 1)}px !important` })
            text.get(0).style.setProperty('font-size', `${CurrFontSize - 1}px`, 'important');
                //update box & text widths
            boxWidth = (width) ? width : box.innerWidth();
            if(hasChildren){
                textWidth=0
                text.children().each(function(){
                    textWidth  = textWidth  + $(this).innerWidth();
                });
            }
            else{
                textWidth = Math.max(text.innerWidth(), text[0].scrollWidth);
            }
        }
        if (debug) console.table([
            ["Container Width", boxWidth],
            ["Text Width", (hasChildren) ? initTextWidth : textWidth],
            ["Initial Font Size", intFS]
        ])
        if (debug) console.log("Text Fit Complete (Width Only): " + id + "\n" + "Old Font Size: " + intFS + "   -   " + "New Font Size: " + parseInt(text.css('font-size'), 10))
        if (debug && hasChildren) console.log("Width calculated on child elements")

        return parseInt(text.css('font-size'), 10)
    } else {
        //run loop to drop font size till it fits in width & height
        while (boxWidth < textWidth || boxHeight < textHeight) {
            //get curr font size
            let CurrFontSize = parseInt(text.css('font-size'), 10)
                //check if text is too small
            if (CurrFontSize == 1 || CurrFontSize <= min) {
                if(debug && CurrFontSize <= min) console.log(`Reached minimum font size ${min}px, cant text fix`);
                if(debug && CurrFontSize == 1) console.log(`Cant text fit lower than 1px`); 
                break;
            }
                //-1 from font size
            //text.css({ fontSize: `${(CurrFontSize - 1)}px !important` })
            text.get(0).style.setProperty('font-size', `${CurrFontSize - 1}px`, 'important');
                //update box & text widths & heights
            boxWidth = (width) ? width : box.innerWidth();
            boxHeight = (height) ? height : box.innerHeight();
            if(hasChildren){
                textWidth=0
                textHeight=0
                text.children().each(function(){
                        //Single width dosnt suppoort multiple widths
                    textWidth  = $(this).innerWidth()
                        //Adds all heights toogethe
                    textHeight = textHeight + $(this).innerHeight()
                });
            }
            else{
                textWidth  = Math.max(text.innerWidth(), text[0].scrollWidth);
                textHeight = Math.max(text.innerHeight(), text[0].scrollHeight);
            }

        }
        if (debug) console.table([
            ["Container Width", boxWidth],
            ["Container Height", boxHeight],
            ["Text Width", (hasChildren) ? initTextWidth : textWidth],
            ["Text Height",(hasChildren) ? initTextHeight : textHeight],
            ["Initial Font Size", intFS]
        ])
        if (debug) console.log("Text Fit Complete (Width & Height): " + id + "\n" + "Old Font Size: " + intFS + "   -   " + "New Font Size: " + parseInt(text.css('font-size'), 10))
        if (debug && hasChildren) console.log("Width & Height calculated on child elements")
        return parseInt(text.css('font-size'), 10)
    }
}

function mblTextFit(array, name) {
    for (a = 0; a < array.length; a++) {
        if (debug) console.groupCollapsed("%c***** Text Fit: " + array[a].ids + " *****", logStyles.info);
        let fittedText = [];
        for (b = 0; b < array[a].ids.length; b++) {
            _debug("Text Fit","",array[a].ids[b],logStyles.general,false)
            fittedText.push(textFit(array[a].ids[b],
                array[a].inner,
                array[a].min,
                array[a].width,
                array[a].height,
                array[a].widthOnly,
                array[a].originalFS,
                array[a].hasChildren))
        }
        fittedText.sort(function(a, b) { return a - b });
        if (debug) console.log(fittedText)
        if (debug) console.groupEnd()
        for (c = 0; c < array[a].ids.length; c++) {
            document.querySelector(`#${array[a].ids[c]} > ${array[a].inner}`).style.fontSize = fittedText[0] + "px"
        }
    }
}

function resetTextFit() {
    const keys = Object.keys(namedFitTextGroups);
    keys.forEach((key, index) => {
        namedFitTextGroups[key][0].isTextFitted = false
    });
}

function scaleBoard() {
    var $el = $("#very-specific-design");
    var elHeight = $el.outerHeight();
    var elWidth = $el.outerWidth();
    var $wrapper = $("#scaleable-wrapper");
    var isMobile = /Mobi|Android/i.test(navigator.userAgent);

    if(!isMobile){
        $wrapper.resizable({
            animate: true,
            maxHeight: 1808,
            maxWidth: 1920,
            aspectRatio: 16 / 9,
            resize: doResize
        });
    }

    function doResize(event, ui) {
        var scale, origin;
        scale = Math.min(
            ui.size.width / elWidth,
            ui.size.height / elHeight
        );
        $el.css({
            transform: "translate(-50%, -50%) " + "scale(" + scale + ")"
        });
    }

    var starterData = {
        size: {
            width: $wrapper.width(),
            height: $wrapper.height()
        }
    };
    doResize(null, starterData);

    if (!isMobile) {
        window.addEventListener('resize', function () {
            // Debounce resize event
            clearTimeout(this.resizeTimeout);
            this.resizeTimeout = setTimeout(function () {
                scaleBoard();
            }, 100);
        });
    }
    if(isMobile){
        window.addEventListener('resize', function () {
            $el.removeClass("very-specific-design")
            $el.css({
                transform: "",
                left: "0%",
                top: "0%",
                transformOrigin: "center center",
                width: "1920px",
                height: "1080px",
                padding: "0px",
                textAlign: "center",
                background: "transparent",
                position: "absolute",
                backgroundImage: "url('https://premium.masjidboardlive.com/v2/images/pillars.png')",
                backgroundColor: "var(--background)"
            });
        })
    }
}
scaleBoard();



//Original Function
function scaleBoardOri() {
	var $el      = $("#very-specific-design");
	var elHeight = $el.outerHeight();
	var elWidth  = $el.outerWidth();
	var $wrapper = $("#scaleable-wrapper");

	$wrapper.resizable({
	  resize: doResize
	});

	function doResize(event, ui) {
	    var scale, origin;
	    scale = Math.min(
		ui.size.width / elWidth,    
		ui.size.height / elHeight
	  );
	  $el.css({
		transform: "translate(-50%, -50%) " + "scale(" + scale + ")"
	  });
	}
	var starterData = { 
	  size: {
		width: $wrapper.width(),
		height: $wrapper.height()
	  }
	}
	doResize(null, starterData)
	
}
//(mobile) ? scaleBoard() : window.addEventListener('resize', scaleBoard);
//scaleBoard()

/*function scaleContainer() {
    const container = document.getElementById('scaleable-wrapper');
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight;

    const scaleWidth = viewportWidth / 1920;
    const scaleHeight = viewportHeight / 1080;

    // Use the smaller scale to ensure the entire container is visible
    const scale = Math.min(scaleWidth, scaleHeight);

    //container.style.transform = `scale(${scale})`;
    container.style.width = scaleWidth;
    container.style.height = scaleHeight;
    console.log(scaleWidth)
    console.log(scaleHeight)
}

(mobile) ? scaleContainer() : scaleContainer();window.addEventListener('resize', scaleContainer, 300);
scaleContainer();
let resizeTimeout;
window.addEventListener('resize', () => {
    clearTimeout(resizeTimeout);
    resizeTimeout = setTimeout(scaleContainer, 100);
});*/
/*document.addEventListener('gesturestart', function (e) {
    e.preventDefault();
});
document.addEventListener('gesturechange', function (e) {
e.preventDefault();
});
document.addEventListener('gestureend', function (e) {
e.preventDefault();
});*/

function jumuah(){
    document.getElementById("jumuahTime1"  ).innerHTML = (salaah12HourFormat) ? time24to12(jumuahTime1) : jumuahTime1
    document.getElementById("jumuahTime2"  ).innerHTML = (salaah12HourFormat) ? time24to12(jumuahTime2) : jumuahTime2
    document.getElementById("jumuahTime3"  ).innerHTML = (salaah12HourFormat) ? time24to12(jumuahTime3) : jumuahTime3
    document.getElementById("jumuahKhateeb").innerHTML = jumuahKhateeb;
    
    let headings = [document.getElementById("jumuahAdhaanHeading"),document.getElementById("jumuahSunnanHeading"),document.getElementById("jumuahKhutbahHeading")]
    let times    = [document.getElementById("jumuahAdaanTime"    ),document.getElementById("jumuahSunanTime"    ),document.getElementById("jumuahKhutbahTime"   )]
    
    if(jumuahHeadingsArray == undefined){
        document.getElementById("jumuahHead1").innerHTML = jumuahHead1;
        document.getElementById("jumuahHead2").innerHTML = jumuahHead2;
        document.getElementById("jumuahHead3").innerHTML = jumuahHead3;
        mblTextFit(fitTextGroups)
        return;
    }
    
    if(jumuahHeadingsArray.length > 1 && jumuahHeadingsArray != "#N/A"){
        let newJumuahHeading = jumuahHeadingsArray.split(",")
        let totalCols  = 3;
        //cater for blanks
        for(j=0; j<newJumuahHeading.length; j++){
            if(newJumuahHeading[j] == ""){
                totalCols--
                headings[j].style.display  = "none"
                times   [j].style.display  = "none"
                //headings[j]
                //times   [j].colSpan        = 2
                newJumuahHeading[j] = j
            }
        }
        
        let colgroup = document.querySelector("#jumuahTable colgroup");
        let cols = colgroup.querySelectorAll("col");
        let width = 405 / totalCols;
        for (c = 0; c < cols.length; c++) {
          cols[c].setAttribute("width", width);
        }
        document.getElementById("jumuahKhateebHeading").colSpan = totalCols
        
        document.getElementById("jumuahHead1").innerHTML = translate("jumuahHeading"+ newJumuahHeading[0])
        document.getElementById("jumuahHead2").innerHTML = translate("jumuahHeading"+ newJumuahHeading[1])
        document.getElementById("jumuahHead3").innerHTML = translate("jumuahHeading"+ newJumuahHeading[2])
        mblTextFit(fitTextGroups)
    } else {
        document.getElementById("jumuahHead1").innerHTML = jumuahHead1;
        document.getElementById("jumuahHead2").innerHTML = jumuahHead2;
        document.getElementById("jumuahHead3").innerHTML = jumuahHead3;
        mblTextFit(fitTextGroups)
    }
    
    
}



// !====================================================== //
// !================= Interval Functions ================= //
// !====================================================== //
// ~~~~~~~~ Salaah Countdown Timer  ~~~~~~~ //
function countdownCalcOri() { //Original
    let countDownI
    let count = null;
    let countDownTextFit = false;

    function updateCountdown(time, label, zwalPopUp, hide) {
        _debug("Countdown Timer",label,secondstotime(time),logStyles.function,false)
        count = time;
        countDownI = setInterval(() => { countdown(label, zwalPopUp, hide) }, 1000);
    }

    function countdown(label, zwalPopUp, hide) {
        let timenow = timeInSeconds;
        let countdown = (count > timenow) ? count - timenow : (86400 - timenow) + count;
        let timer = document.getElementById("countdownTimer");
        let name = document.getElementById("countdownTimerName");
        let zawaal = document.getElementById('zawaalPopUp');
        let zawaalPopUpText = document.getElementById("zawaalPopUpText");
        let masjidURL = document.getElementById('masjidUrl');

        timer.innerHTML = secondstotime(countdown);
        name.innerHTML = label;

        //Flash when 5mins remain
        if (countdown < 300 && countdown > 3) {
            if(hide) hidePoster("countdown","Salaah countdown flash", true)
            name.classList.add("flashcountdown");
            timer.classList.remove("shadow");
            timer.classList.add("flashcountdown");
        } else if (countdown < 3) {
            hidePoster("countdown","Salaah countdown flash", false)
        }

        if (!countDownTextFit) {
            //mblTextFit();
            countDownTextFit = true;
        }

        //Zawaal pop up
        if (zwalPopUp) {
            hidePoster("zawaal","Zawaal Flash", true);
            zawaalPopUp.classList.remove("hidden");
            zawaalPopUp.classList.add("show");
            zawaalPopUp.classList.add("flashcountdown");
            zawaalPopUpText.innerHTML = translate('zawaalPopUpText');
            masjidURL.classList.add("hidden");
            timer.innerHTML = "";
            name.innerHTML = "";
        } else {
            hidePoster("zawaal","Zawaal Flash", false);
        }

        //expired if 0
        if (countdown < 2) {
            if (zwalPopUp) {
                zawaalPopUp.classList.add("hidden");
                zawaalPopUp.classList.remove("show");
                zawaalPopUp.classList.remove("flashcountdown");
                masjidURL.classList.remove("hidden");
            }
            timer.innerHTML = translate('cDownExpired');
            clearInterval(countDownI);
            setTimeout(() => {
                countDownTextFit = false;
                name.classList.remove("flashcountdown");
                timer.classList.remove("flashcountdown");
                timer.classList.add("shadow");
                countdownCalc();
                //return
            }, 3000)
        }
    }
    //Variables
    let timenow = timeInSeconds
    let suhr = minutestoseconds(sehriEnds);
    let fjrs = minutestoseconds(fajrStarts);
    let fjra = minutestoseconds(fajrAthan);
    let fjrj = minutestoseconds(fajrJamaah);
    let sunr = minutestoseconds(sunrise);
    let isrq = minutestoseconds(ishraaq);
    let dha = minutestoseconds(duha);
    let istc = minutestoseconds(istiwaCaution);
    let istw = minutestoseconds(istiwa);
    let zwle = minutestoseconds(zawaalEnd);
    let zhra = minutestoseconds(dhuhrAthan);
    let zhrj = minutestoseconds(dhuhrJamaah);
    let asrs = minutestoseconds(asrShafi);
    let asrh = minutestoseconds(asrHanafi);
    let asra = minutestoseconds(asrAthan);
    let asrj = minutestoseconds(asrJamaah);
    let suns = minutestoseconds(sunset);
    let iftr = minutestoseconds(maghribAthan);
    //let mags = minutestoseconds(maghribAthan);
    let estr = minutestoseconds(eshaStarts);
    let esha = minutestoseconds(eshaAthan);
    let eshj = minutestoseconds(eshaJamaah);
    let midn = 86400;

    let jum1 = minutestoseconds(jumuahTime1)
    let jum2 = minutestoseconds(jumuahTime2)
    let jum3 = minutestoseconds(jumuahTime3)
    let jzhe;
    let jzct;

    let suhurFajrLabel = (fjrs == suhr)    ? translate('cDownSuhurFajrIn') + fajrSuhurCheck(1) : translate('cDownSuhurExpiresIn') + fajrSuhurCheck(1);
    let jumuahZuhr     = (currentDay == 5) ? translate('cDownJumuahAdhanIn') : translate('cDownZuhrAdhanIn');
    let jumuahKhutbah  = (currentDay == 5) ? translate('cDownKhutbahIn')     : translate('cDownZuhrJamaahIn');

    //Jumuah Calculations
    /*if(currentDay == 5){
        if(jumuahHeadingsArray.length > 1 && jumuahHeadingsArray != "#N/A"){
            let newJumuahHeading = jumuahHeadingsArray.split(",")
            let times            = [zwle,jum1,jum2,jum3]
            let futureTimes      = times.filter(time => time > timenow);
            let nextTime         = Math.min(...futureTimes);
            let index            = times.indexOf(nextTime);
            if(nextTime != Infinity || nextTime != undefined || newJumuahHeading[(index==0) ? index : -1] != undefined || nextTime == -1){
                console.log(newJumuahHeading)
                console.log(index)
                console.log('jumuahHeading'+newJumuahHeading[(index==0) ? index : -1])
                jzhe = `${translate('jumuahHeading')} ${translate('jumuahHeading'+newJumuahHeading[(index==0) ? index : -1])}`
                jzct = nextTime
            }
        }
        if (timenow > zwle && timenow < jzct) {
            updateCountdown(jzct, jzhe, false, false);
        } else if(timenow > jum3 && timenow < asrs){
            updateCountdown(asrs, translate('cDownZuhrSExpiry')  , false, true )
        }
    }*/

    if (asra < asrh) {
        if      (timenow > fjrs && timenow < fjra) { updateCountdown(fjra, translate('cDownFajrAdhanIn')  , false, false) } 
        else if (timenow > fjra && timenow < fjrj) { updateCountdown(fjrj, translate('cDownFajrJamaahIn') , false, false) } 
        else if (timenow > fjrj && timenow < sunr) { updateCountdown(sunr, translate('cDownFajrExpiry')   , false, true ) } 
        else if (timenow > sunr && timenow < isrq) { updateCountdown(isrq, translate('cDownIshraqIn')     , false, true ) } 
        else if (timenow > isrq && timenow < dha ) { updateCountdown(dha, (chast) ? translate('cDownChastIn') : translate('cDownDuhaIn'), false, true) } 
        else if (timenow > dha && timenow < istc ) { updateCountdown(istw, translate('cDownZawaalIn')     , false, true ) } 
        else if (timenow > istc && timenow < istw) { updateCountdown(istw, translate('cDownIstiwaCaution'), false, true ) } 
        else if (timenow > istw && timenow < zwle) { updateCountdown(zwle, ""                             , true , true ) } 
        else if (timenow > zwle && timenow < zhra) { updateCountdown(zhra, jumuahZuhr                     , false, false) } 
        else if (timenow > zhra && timenow < zhrj) { updateCountdown(zhrj, jumuahKhutbah                  , false, false) } 
        else if (timenow > zhrj && timenow < asrs) { updateCountdown(asrs, translate('cDownZuhrSExpiry')  , false, true ) } 
        else if (timenow > asrs && timenow < asra) { updateCountdown(asra, translate('cDownAsrAdhanIn')   , false, false) } 
        else if (timenow > asra && timenow < asrj) { updateCountdown(asrj, translate('cDownAsrJamaahIn')  , false, false) } 
        else if (timenow > asrj && timenow < asrh) { updateCountdown(asrh, translate('cDownZuhrHExpiry')  , false, true ) } 
        else if (timenow > asrj && timenow < suns) { updateCountdown(suns, translate('cDownSunsetIn')     , false, true ) } 
        else if (timenow > suns && timenow < iftr) { updateCountdown(iftr, translate('cDownMaghribIn')    , false, false) } 
        else if (timenow > suns && timenow < estr) { updateCountdown(estr, translate('cDownMagribExpiry') , false, true ) } 
        else if (timenow > estr && timenow < esha) { updateCountdown(esha, translate('cDownEshaAdhanIn')  , false, false) } 
        else if (timenow > esha && timenow < eshj) { updateCountdown(eshj, translate('cDownEshaJamaahIn') , false, false) } 
        else if (timenow > eshj && timenow > suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow < eshj && timenow < suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow > suhr && timenow < fjrs) { updateCountdown(fjrs, translate('cDownFajrStarts')   , false, true ) }
    }

    else if (asra >= asrh) {
        if      (timenow > fjrs && timenow < fjra) { updateCountdown(fjra, translate('cDownFajrAdhanIn')  , false, false) } 
        else if (timenow > fjra && timenow < fjrj) { updateCountdown(fjrj, translate('cDownFajrJamaahIn') , false, false) } 
        else if (timenow > fjrj && timenow < sunr) { updateCountdown(sunr, translate('cDownFajrExpiry')   , false, true ) } 
        else if (timenow > sunr && timenow < isrq) { updateCountdown(isrq, translate('cDownIshraqIn')     , false, true ) } 
        else if (timenow > isrq && timenow < dha ) { updateCountdown(dha, (chast) ? translate('cDownChastIn') : translate('cDownDuhaIn'), false, true) } 
        else if (timenow > dha && timenow < istc ) { updateCountdown(istw, translate('cDownZawaalIn')     , false, true ) } 
        else if (timenow > istc && timenow < istw) { updateCountdown(istw, translate('cDownIstiwaCaution'), false, true ) } 
        else if (timenow > istw && timenow < zwle) { updateCountdown(zwle, ""                             , true , true ) } 
        else if (timenow > zwle && timenow < zhra) { updateCountdown(zhra, jumuahZuhr                     , false, false) } 
        else if (timenow > zhra && timenow < zhrj) { updateCountdown(zhrj, jumuahKhutbah                  , false, false) } 
        else if (timenow > zhrj && timenow < asrs) { updateCountdown(asrs, translate('cDownZuhrSExpiry')  , false, true ) } 
        else if (timenow > asrs && timenow < asrh) { updateCountdown(asrh, translate('cDownZuhrHExpiry')  , false, true ) } 
        else if (timenow > asrh && timenow < asra) { updateCountdown(asra, translate('cDownAsrAdhanIn')   , false, false) } 
        else if (timenow > asra && timenow < asrj) { updateCountdown(asrj, translate('cDownAsrJamaahIn')  , false, false) } 
        else if (timenow > asrj && timenow < suns) { updateCountdown(suns, translate('cDownSunsetIn')     , false, true ) } 
        else if (timenow > suns && timenow < iftr) { updateCountdown(iftr, translate('cDownMaghribIn')    , false, false) } 
        else if (timenow > suns && timenow < estr) { updateCountdown(estr, translate('cDownMagribExpiry') , false, true ) } 
        else if (timenow > estr && timenow < esha) { updateCountdown(esha, translate('cDownEshaAdhanIn')  , false, false) } 
        else if (timenow > esha && timenow < eshj) { updateCountdown(eshj, translate('cDownEshaJamaahIn') , false, false) } 
        else if (timenow > eshj && timenow > suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow < eshj && timenow < suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow > suhr && timenow < fjrs) { updateCountdown(fjrs, translate('cDownFajrStarts')   , false, true ) }
    }
}

function countdownCalc() { //Mali
    let countDownI
    let count = null;
    let countDownTextFit = false;

    function updateCountdown(time, label, zwalPopUp, hide) {
        _debug("Countdown Timer",label,secondstotime(time),logStyles.function,false)
        count = time;
        clearInterval(countDownI)
        countDownI = setInterval(() => { countdown(label, zwalPopUp, hide) }, 1000);
    }

    function countdown(label, zwalPopUp, hide) {
        let timenow = timeInSeconds;
        let countdown = (count > timenow) ? count - timenow : (86400 - timenow) + count;
        let timer = document.getElementById("countdownTimer");
        let name = document.getElementById("countdownTimerName");
        let zawaal = document.getElementById('zawaalPopUp');
        let zawaalPopUpText = document.getElementById("zawaalPopUpText");
        let masjidURL = document.getElementById('masjidUrl');

        timer.innerHTML = secondstotime(countdown);
        name.innerHTML = label;

        //Flash when 5mins remain
        if (countdown < 300 && countdown > 3) {
            if(hide) hidePoster("countdown","Salaah countdown flash", true)
            name.classList.add("flashcountdown"); 
            timer.classList.remove("shadow");
            timer.classList.add("flashcountdown");
        } else if (countdown < 3) {
            hidePoster("countdown","Salaah countdown flash", false)
        }

        if (!countDownTextFit) {
            //mblTextFit();
            countDownTextFit = true;
        }

        //Zawaal pop up
        if (zwalPopUp) {
            hidePoster("zawaal","Zawaal Flash", true);
            zawaalPopUp.classList.remove("hidden");
            zawaalPopUp.classList.add("show");
            zawaalPopUp.classList.add("flashcountdown");
            zawaalPopUpText.innerHTML = translate('zawaalPopUpText');
            masjidURL.classList.add("hidden");
            timer.innerHTML = "";
            name.innerHTML = "";
        } else {
            hidePoster("zawaal","Zawaal Flash", false);
        }

        //expired if 0
        if (countdown < 2) {
            if (zwalPopUp) {
                zawaalPopUp.classList.add("hidden");
                zawaalPopUp.classList.remove("show");
                zawaalPopUp.classList.remove("flashcountdown");
                masjidURL.classList.remove("hidden");
            }
            timer.innerHTML = translate('cDownExpired');
            clearInterval(countDownI);
            setTimeout(() => {
                countDownTextFit = false;
                name.classList.remove("flashcountdown");
                timer.classList.remove("flashcountdown");
                timer.classList.add("shadow");
                countdownCalc();
                //return
            }, 3000)
        }
    }
    //Variables
    let timenow = timeInSeconds
    let suhr = minutestoseconds(sehriEnds);
    let fjrs = minutestoseconds(fajrStarts);
    let fjra = minutestoseconds(fajrAthan);
    let fjrj = minutestoseconds(fajrJamaah);
    let sunr = minutestoseconds(sunrise);
    let isrq = minutestoseconds(ishraaq);
    let dha = minutestoseconds(duha);
    let istc = minutestoseconds(istiwaCaution);
    let istw = minutestoseconds(istiwa);
    let zwle = minutestoseconds(zawaalEnd);
    let zhra = minutestoseconds(dhuhrAthan);
    let zhrj = minutestoseconds(dhuhrJamaah);
    let asrs = minutestoseconds(asrShafi);
    let asrh = minutestoseconds(asrHanafi);
    let asra = minutestoseconds(asrAthan);
    let asrj = minutestoseconds(asrJamaah);
    let suns = minutestoseconds(sunset);
    let iftr = minutestoseconds(maghribAthan);
    //let mags = minutestoseconds(maghribAthan);
    let estr = minutestoseconds(eshaStarts);
    let esha = minutestoseconds(eshaAthan);
    let eshj = minutestoseconds(eshaJamaah);
    let midn = 86400;

    let jum1 = minutestoseconds(jumuahTime1)
    let jum2 = minutestoseconds(jumuahTime2)
    let jum3 = minutestoseconds(jumuahTime3)
    let jzhe;
    let jzct;

    let suhurFajrLabel = (fjrs == suhr)    ? translate('cDownSuhurFajrIn') + fajrSuhurCheck(1) : translate('cDownSuhurExpiresIn') + fajrSuhurCheck(1);
    let jumuahZuhr     = (currentDay == 5) ? translate('cDownJumuahAdhanIn') : translate('cDownZuhrAdhanIn');
    let jumuahKhutbah  = (currentDay == 5) ? translate('cDownKhutbahIn')     : translate('cDownZuhrJamaahIn');

    if (asra < asrh) { //Shafii timings
        if      (timenow > fjrs && timenow < fjra) { updateCountdown(fjra, translate('cDownFajrAdhanIn')  , false, false) } 
        else if (timenow > fjra && timenow < fjrj) { updateCountdown(fjrj, translate('cDownFajrJamaahIn') , false, false) } 
        else if (timenow > fjrj && timenow < sunr) { updateCountdown(sunr, translate('cDownFajrExpiry')   , false, true ) } 
        else if (timenow > sunr && timenow < isrq) { updateCountdown(isrq, translate('cDownIshraqIn')     , false, true ) } 
        else if (timenow > isrq && timenow < dha ) { updateCountdown(dha, (chast) ? translate('cDownChastIn') : translate('cDownDuhaIn'), false, true) } 
        else if (timenow > dha && timenow < istc ) { updateCountdown(istw, translate('cDownZawaalIn')     , false, true ) } 
        else if (timenow > istc && timenow < istw) { updateCountdown(istw, translate('cDownIstiwaCaution'), false, true ) } 
        else if (timenow > istw && timenow < zwle) { updateCountdown(zwle, ""                             , true , true ) } 
        else if (timenow > zwle && timenow < zhra) { updateCountdown(zhra, translate('cDownZuhrAdhanIn')  , false, false) } 
        else if (timenow > zhra && timenow < zhrj) { updateCountdown(zhrj, translate('cDownZuhrJamaahIn') , false, false) } 
        else if (timenow > zhrj && timenow < asrs) { updateCountdown(asrs, translate('cDownZuhrSExpiry')  , false, true ) } 
        else if (timenow > asrs && timenow < asra) { updateCountdown(asra, translate('cDownAsrAdhanIn')   , false, false) } 
        else if (timenow > asra && timenow < asrj) { updateCountdown(asrj, translate('cDownAsrJamaahIn')  , false, false) } 
        else if (timenow > asrj && timenow < asrh) { updateCountdown(asrh, translate('cDownZuhrHExpiry')  , false, true ) } 
        else if (timenow > asrj && timenow < suns) { updateCountdown(suns, translate('cDownSunsetIn')     , false, true ) } 
        else if (timenow > suns && timenow < iftr) { updateCountdown(iftr, translate('cDownMaghribIn')    , false, false) } 
        else if (timenow > suns && timenow < estr) { updateCountdown(estr, translate('cDownMagribExpiry') , false, true ) } 
        else if (timenow > estr && timenow < esha) { updateCountdown(esha, translate('cDownEshaAdhanIn')  , false, false) } 
        else if (timenow > esha && timenow < eshj) { updateCountdown(eshj, translate('cDownEshaJamaahIn') , false, false) } 
        else if (timenow > eshj && timenow > suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow < eshj && timenow < suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow > suhr && timenow < fjrs) { updateCountdown(fjrs, translate('cDownFajrStarts')   , false, true ) }
    }

    else if (asra >= asrh) {
        if      (timenow > fjrs && timenow < fjra) { updateCountdown(fjra, translate('cDownFajrAdhanIn')  , false, false) } 
        else if (timenow > fjra && timenow < fjrj) { updateCountdown(fjrj, translate('cDownFajrJamaahIn') , false, false) } 
        else if (timenow > fjrj && timenow < sunr) { updateCountdown(sunr, translate('cDownFajrExpiry')   , false, true ) } 
        else if (timenow > sunr && timenow < isrq) { updateCountdown(isrq, translate('cDownIshraqIn')     , false, true ) } 
        else if (timenow > isrq && timenow < dha ) { updateCountdown(dha, (chast) ? translate('cDownChastIn') : translate('cDownDuhaIn'), false, true) } 
        else if (timenow > dha && timenow < istc ) { updateCountdown(istw, translate('cDownZawaalIn')     , false, true ) } 
        else if (timenow > istc && timenow < istw) { updateCountdown(istw, translate('cDownIstiwaCaution'), false, true ) } 
        else if (timenow > istw && timenow < zwle) { updateCountdown(zwle, ""                             , true , true ) } 
        else if (timenow > zwle && timenow < zhra) { updateCountdown(zhra, translate('cDownZuhrAdhanIn')  , false, false) } 
        else if (timenow > zhra && timenow < zhrj) { updateCountdown(zhrj, translate('cDownZuhrJamaahIn') , false, false) } 
        else if (timenow > zhrj && timenow < asrs) { updateCountdown(asrs, translate('cDownZuhrSExpiry')  , false, true ) } 
        else if (timenow > asrs && timenow < asrh) { updateCountdown(asrh, translate('cDownZuhrHExpiry')  , false, true ) } 
        else if (timenow > asrh && timenow < asra) { updateCountdown(asra, translate('cDownAsrAdhanIn')   , false, false) } 
        else if (timenow > asra && timenow < asrj) { updateCountdown(asrj, translate('cDownAsrJamaahIn')  , false, false) } 
        else if (timenow > asrj && timenow < suns) { updateCountdown(suns, translate('cDownSunsetIn')     , false, true ) } 
        else if (timenow > suns && timenow < iftr) { updateCountdown(iftr, translate('cDownMaghribIn')    , false, false) } 
        else if (timenow > suns && timenow < estr) { updateCountdown(estr, translate('cDownMagribExpiry') , false, true ) } 
        else if (timenow > estr && timenow < esha) { updateCountdown(esha, translate('cDownEshaAdhanIn')  , false, false) } 
        else if (timenow > esha && timenow < eshj) { updateCountdown(eshj, translate('cDownEshaJamaahIn') , false, false) } 
        else if (timenow > eshj && timenow > suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow < eshj && timenow < suhr) { updateCountdown(suhr, suhurFajrLabel                 , false, true ) } 
        else if (timenow > suhr && timenow < fjrs) { updateCountdown(fjrs, translate('cDownFajrStarts')   , false, true ) }
    }

    //Jumuah Calculations
    if(currentDay == 5){
        if(jumuahHeadingsArray.length > 1 && jumuahHeadingsArray != "#N/A"){
            let newJumuahHeading = jumuahHeadingsArray.split(",")
            let times            = [zwle,jum1,jum2,jum3]
            let futureTimes      = times.filter(time => time > timenow);
            let nextTime         = Math.min(...futureTimes);
            let index            = times.indexOf(nextTime);

            console.log(newJumuahHeading)
            console.log(times)
            console.log(futureTimes)
            console.log(nextTime)
            console.log(index)

            if(nextTime != Infinity || nextTime != undefined || newJumuahHeading[(index-1<0) ? 0 : index-1] != undefined || nextTime != -1 || (index-1) != -1){
                console.log("In IF")
                jzhe = `${translate('jumuahHeading')} ${translate('jumuahHeading'+newJumuahHeading[(index-1<0) ? 0 : index-1])}`
                jzct = nextTime
            }

            if (timenow > zwle && timenow < jzct && futureTimes.length > 0) {
                updateCountdown(jzct, jzhe, false, false);
                console.log("Jumuah Time Calculated")
            } else if(timenow > jum3 && timenow < asrs){
                updateCountdown(asrs, translate('cDownZuhrSExpiry')  , false, true )
                console.log("Other Time Calculated in Jumuah if else")
            }
        }
    }
}







// ~~~~~~~~~~~ Salaah Language ~~~~~~~~~~~ //
let salaahLanguageCounter = 0
let masjidNameCounter     = 0
let nextSalaahHighlight   = 0
let nextOwwalHighlight    = 0
let OrightToLeft          = false
let idfs;
let salaahLanguageTimout ;
function salaahLanguageCheck() {
    let columnA = ["msttop", "fajr", "dhuhr", "asr", "maghrib", "esha", "msttop", "fajr", "dhuhr", "asr", "maghrib", "esha"]
    let columnB = ["adhanColumn", "fajrAthan", "dhuhrAthan", "asrAthan", "maghribAthan", "eshaAthan", "adhanColumn", "fajrAthan", "dhuhrAthan", "asrAthan", "maghribAthan", "eshaAthan"]
    let columnC = ["jamaahColumn", "fajrJamaah", "dhuhrJamaah", "asrJamaah", "maghribJamaah", "eshaJamaah", "jamaahColumn", "fajrJamaah", "dhuhrJamaah", "asrJamaah", "maghribJamaah", "eshaJamaah"]
    let puColA  = ["puh_sehriEnds","puh_fajrStarts","puh_sunrise","puh_ishraaq","puh_duha","puh_istiwa","puh_sundayDhuhr","puh_asrShafi","puh_asrHanafi","puh_sunset","puh_eshaStarts"];
    let puColB  = ["pu_sehriEnds","pu_fajrStarts","pu_sunrise","pu_ishraaq","pu_duha","pu_istiwa","pu_sundayDhuhr","pu_asrShafi","pu_asrHanafi","pu_sunset","pu_eshaStarts"];
    let date = document.getElementById("islamicDate")
    let masjidName1C = document.getElementById("masjidName1")
    let masjidName2C = document.getElementById("masjidName2")

    let engNames = ["", translate('fajr'),
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? translate('jumuah') : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? translate('jumuah') : translate('zuhr'),
        translate('asr'), translate('maghrib'), translate('esha')
                    ,"", translate('fajr'),
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? translate('jumuah') : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? translate('jumuah') : translate('zuhr'),
        translate('asr'), translate('maghrib'), translate('esha')
    ]

    let araNames = [ "",
                    "الفجر",
    (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? "الجمعة" : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? "الجمعة" : "الظهر",
                    "العصر",
                    "المغرب",
                    "العشاء",
                    "", "المغرب",
                    "العشاء",
                    "الفجر",
    (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? "الجمعة" : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? "الجمعة" : "الظهر",
                    "العصر"
                    ]
    let owwalNames = ["إنتهى السحور"//Suhur Ends
                    ,"الصبح الصادق" //Fajr Starts
                    ,"طلوع الشمس"   //Sunrise
                    ,"الإشراق"       //Ishraaq
                    ,"الضحى"        // Duha
                    ,"إستواء الشمس" //Istiwa
                    ,"إبتداء الظهر" //Zuhr
                    ,"(العصر (ش"    //Asr Shafii
                    ,"(العصر (ح"    //Asr Hanafi
                    ,"غروب الشمس"   //Sunset
                    ,"العشاء"       //Esha
                    ]

    //Add chast to Duha
         if (chast == true && salaahLanguageCounter == 0 && puIslamicTimes)  document.getElementById('puh_chast').innerHTML = translate('puh_chast');
    else if (chast == true && salaahLanguageCounter == 1 && puIslamicTimes)  document.getElementById('puh_chast').innerHTML =                 'چاست';
    else                                                                     document.getElementById('puh_chast').innerHTML = '';

    let masjidNameE    = [masjidName1, masjidName2]
    let masjidNameA    = [masjidName1I, masjidName2I]
    
    let adhaanTimes    = [
        (salaahLanguageCounter == 1) ? "الاذان" : translate("adhanColumn"), 
        fajrAthan,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahAdhan : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahAdhan : dhuhrAthan,
        asrAthan, 
        maghribAthan, 
        eshaAthan,
        (salaahLanguageCounter == 1) ? "الاذان" : translate("adhanColumn"),
        maghribAthan,
        eshaAthan,
        fajrAthan,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahAdhan : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahAdhan : dhuhrAthan,
        asrAthan
    ]
    let jamaahTimes    = [
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"),
        fajrJamaah,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaah : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaah : dhuhrJamaah,
        asrJamaah,
        maghribJamaah,
        eshaJamaah,
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"),
        maghribJamaah,
        eshaJamaah,
        fajrJamaah,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaah : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaah : dhuhrJamaah,
        asrJamaah
    ]
    let puJamaahTimes  = [
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"),
        fajrJamaah,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaah : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaah : dhuhrJamaah,
        asrJamaah,
        maghribAthan,
        eshaJamaah,
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"),
        maghribAthan,
        eshaJamaah,
        fajrJamaah,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaah : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaah : dhuhrJamaah,
        asrJamaah
    ]
    let adhaanTimesI   = [
        (salaahLanguageCounter == 1) ? "الاذان" : translate("adhanColumn"),
        fajrAthanI, 
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahAdhanI : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahAdhanI : dhuhrAthanI,
        asrAthanI,
        maghribAthanI,
        eshaAthanI, 
        (salaahLanguageCounter == 1) ? "الاذان" : translate("adhanColumn"),
        maghribAthanI,
        eshaAthanI,
        fajrAthanI,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahAdhanI : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahAdhanI : dhuhrAthanI,
        asrAthanI
    ]
    let jamaahTimesI   = [
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"),
        fajrJamaahI,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaahI : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaahI : dhuhrJamaahI,
        asrJamaahI,
        maghribJamaahI,
        eshaJamaahI, 
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"),
        maghribJamaahI,
        eshaJamaahI,
        fajrJamaahI,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaahI : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaahI : dhuhrJamaahI,
        asrJamaahI
    ]
    let puJamaahTimesI = [
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"),
        fajrJamaahI,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaahI : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaahI : dhuhrJamaahI,
        asrJamaahI,
        maghribAthanI,
        eshaJamaahI,
        (salaahLanguageCounter == 1) ? "الصلاة" : translate("jamaahColumn"), 
        maghribAthanI,
        eshaJamaahI,
        fajrJamaahI,
        (timeInSeconds > minutestoseconds(sunset) && currentDay == 4) ? jumuahJamaahI : (timeInSeconds < minutestoseconds(sunset) && currentDay == 5) ? jumuahJamaahI : dhuhrJamaahI,
        asrJamaahI
    ]
    
    
    let oTimes         = [sehriEnds,fajrStarts,sunrise,ishraaq,duha,istiwa,zawaalEnd,asrShafi,asrHanafi,sunset,eshaStarts];
    let oTimesI        = [suhurEndsI,fajrStartsI,sunriseI,ishraaqI,duhaI,istiwaI,zuhrStartsI,asrShafiI,asrHanafiI,sunsetI,eshaStartsI];


    let engCss    = { "text-align": "left", "font-family": "'Open Sans'   , serif", "font-weight": "700", "font-size": "40px", "line-height": "1" }
    let araCss    = { "text-align": "center", "font-family": "'Scheherazade', serif", "font-weight": "400", "font-size": "80px", "line-height": "1" }
    let timeCss   = { "text-align": "center", "font-family": "'Open Sans'   , serif", "font-weight": "700", "font-size": "53px", "line-height": "1" }
    let engCssTop = { "text-align": "center", "font-family": "'Open Sans'   , serif", "font-weight": "600", "font-size": "30px", "margin": "0px", "margin": "-10px"}
    let araCssTop = { "text-align": "center", "font-family": "'Scheherazade', serif", "font-weight": "400", "font-size": "55px", "margin": "-10px"}
    let araNumCss = { "text-align": "center", "font-size": "53px", "font-weight": "700", "line-height": "1.5", "margin-top": "0", "margin-bottom": "0", "font-family": "Open Sans"}
    let puSEngCss = { /*"font-size": "50px",*/ "font-weight": "700", "line-height": "1", "margin-top": "0", "margin-bottom": "0", "font-family": "Open Sans"}
    let puSAraCss = { /*"font-size": "100px",*/ "font-weight": "400", "line-height": "1", "margin-top": "0", "margin-bottom": "0", "font-family": "'Scheherazade', serif"}
    let puOEngCss = { /*"font-size": "37px",*/ "font-weight": "700", "line-height": "1", "margin-top": "0", "margin-bottom": "0", "font-family": "Open Sans"}
    let puOAraCss = { /*"font-size": "59px",*/ "font-weight": "400", "line-height": "1", "margin-top": "0", "margin-bottom": "0", "font-family": "'Scheherazade', serif"}
    let topBorder = { "border": "thin solid var(--borderColors)" }
    let noBorder  = { "border": "0px" }

    let lang = [engNames, araNames]
    let masjidNameL = [masjidNameE, masjidNameA]
    let dates = [islamicDateEng, islamicDateAra]

    //Masjd Name
    //if(mobile){
        if (!masjidNameLanguage) {
            masjidName1C.innerHTML = masjidNameE[(masjidNameOnTop) ? 1 : 0];
            masjidName2C.innerHTML = masjidNameE[(masjidNameOnTop) ? 0 : 1];
        }
        if (masjidNameLanguage) {
            masjidName1C.innerHTML = masjidNameL[masjidNameCounter][(masjidNameOnTop) ? 1 : 0];
            masjidName2C.innerHTML = masjidNameL[masjidNameCounter][(masjidNameOnTop) ? 0 : 1];
        }
    //} else {
        //if (!masjidNameLanguage) masjidName1C.innerHTML = masjidNameE[1]
        //if (masjidNameLanguage ) masjidName1C.innerHTML = masjidNameL[masjidNameCounter][1]
    //}

    //Islamic Date
    date.innerHTML = dates[salaahLanguageCounter]
        //Islamic date text fit
    if (!namedFitTextGroups["dates"][0].isTextFitted) {
        setTimeout(() => {
            mblTextFit(namedFitTextGroups["dates"]);
            idfs = $("#islamicDate").css('font-size');
            namedFitTextGroups["dates"][0].isTextFitted = true;
        }, 2500);
    }
    //once off  CSS
    if (salaahLanguageCounter == 0) {
        $("#islamicDate").css({ "font-family": "Open Sans", "font-size": `${idfs}` });
        $("#islamicDateC").css({ "top": "50px" });
        $("#adhanColumn, #jamaahColumn").css({ "text-align": "center", "font-family": "Open Sans", "font-weight": "600", "font-size": "30px" })
        $("#pu_fajrSuhur").removeClass("hiddenNd")
        $("#puh_istiwaH").removeClass("hiddenNd")
        $("#pu_sundayzuhr").removeClass("hiddenNd")
        $("#puh_asrShafiH").removeClass("hiddenNd")
        $("#puh_asrHanafiH").removeClass("hiddenNd")
        $("#masjidName1, #masjidName2").css({ "text-align": "center", "font-family": "Open Sans", "font-weight": "600", "font-size": "36px" })
        $("#masjidName2").css({ "margin-top": "-33px" })
        //$(".pu_salaah_times").css({"font-size":"50px"})
        //$("#popUpPerpetual h1").css({"font-size":"40px"})
        if(clockNameFormat == 1){
            if(masjidNameE[0] == "" /*|| !mobile*/){
                $("#masjidNameTable").css({ "top": "256px" })
            } else {
                $("#masjidNameTable").css({ "top": "228px" })
            }
        }
        else{
            if(masjidNameE[0] == "" /*|| !mobile*/){
                $("#masjidNameTable").css({ "top": "170px" })
            } else {
                $("#masjidNameTable").css({ "top":  masjidname1TopPad+"px" })
            }
        }
        
    } else if (salaahLanguageCounter == 1) {
        $("#islamicDate").css({ "font-family": "'Scheherazade', serif", "font-size": "58px" });
        $("#islamicDateC").css({ "top": "31px" });
        $("#adhanColumn, #jamaahColumn").css({ "text-align": "center", "font-family": "'Scheherazade', serif", "font-weight": "400", "font-size": "40px" })
        if(islamicTimesArabic && puAlternatingTime){
            $(".pu_salaah_times").css({"font-size":"62px"})
            $("#popUpPerpetual h1").css({"font-size":"44px"})
        }
        if(puAlternatingTime){
            $("#pu_fajrSuhur").addClass("hiddenNd")
            $("#puh_istiwaH").addClass("hiddenNd")
            $("#pu_sundayzuhr").addClass("hiddenNd")
            $("#puh_asrShafiH").addClass("hiddenNd")
            $("#puh_asrHanafiH").addClass("hiddenNd")
        }
    }
    if(masjidNameLanguage && masjidNameCounter == 1) {
        $("#masjidName1, #masjidName2").css({ "text-align": "center", "font-family": "'Scheherazade', serif", "font-weight": "400", "font-size": "58px" })
        $("#masjidName2").css({ "margin-top": "-58px" })
        if(clockNameFormat == 1){
            if(masjidNameA[0] == "" /*|| !mobile*/){
                $("#masjidNameTable").css({ "top": "270px" })
            } else {
                $("#masjidNameTable").css({ "top": "207px" })
            }
        }
        else{
            if(masjidNameA[0] == "" /*|| !mobile*/){
                $("#masjidNameTable").css({ "top": "145px" })
            } else {
                $("#masjidNameTable").css({ "top": "109px" })
            }
        }
    }

    //Live time in islamic ?:
    islamicLiveTime = (islamicTime && salaahLanguageCounter == 1) ? true : false

    //Apply css to row 0 
    if (salaahLanguageCounter == 1) {
        $(`#${(rightToLeft)?columnC[0]:columnA[0]}`).parent().css(noBorder)
        $(`#${(rightToLeft)?columnA[0]:columnC[0]}`).parent().css(topBorder)
        $(`#${(rightToLeft)?columnA[0]:columnC[0]}`).css(araCssTop)
        $(`#adhanColumn`).css(araCssTop)
        $(`#${(rightToLeft)?columnC[0]:columnA[0]}`).attr('style', '');//.removeAttr("style")
    } else {
        $(`#${columnA[0]}`).parent().css(noBorder)
        $(`#${columnA[0]}`).attr('style', '')
        $(`#${columnC[0]}`).parent().css(topBorder)
        $(`#${columnC[0]}`).css(engCssTop)
        $(`#adhanColumn`).css(engCssTop)
    }

    //Salaah Highlight checker
    let nextSalaahHighlightC = 1
    switch (nextSalaahHighlight) {
        case 'fajr':
        case 'midnightN':
            nextSalaahHighlightC = (islamicTimeMaghribFirst && salaahLanguageCounter == 1 /*&& islamicTime*/) ? 3 : 1;
            break;
        case 'dhuhr':
            nextSalaahHighlightC = (islamicTimeMaghribFirst && salaahLanguageCounter == 1/* && islamicTime*/) ? 4 : 2;
            break;
        case 'asr':
            nextSalaahHighlightC = (islamicTimeMaghribFirst && salaahLanguageCounter == 1/* && islamicTime*/) ? 5 : 3;
            break;
        case 'maghrib':
            nextSalaahHighlightC = (islamicTimeMaghribFirst && salaahLanguageCounter == 1/* && islamicTime*/) ? 1 : 4;
            break;
        case 'esha':
            nextSalaahHighlightC = (islamicTimeMaghribFirst && salaahLanguageCounter == 1/* && islamicTime*/) ? 2 : 5;
            break;
    }

    //Loop for salaah time table population
    let q  = (islamicTimeMaghribFirst && salaahLanguageCounter == 1) ? 6  : 0
    let qc = (islamicTimeMaghribFirst && salaahLanguageCounter == 1) ? 12 : 6
    for (q; q < qc; q++) {
        //calculations
        let adhaanTimeCurrent    = (islamicTime && salaahLanguageCounter == 1 && q != 0) ? (islamicTime12Hour && q != 0) ? time24to12(adhaanTimesI[q]) : adhaanTimesI[q] : (salaah12HourFormat && q != 0) ? time24to12(adhaanTimes[q]) : adhaanTimes[q];
        let adhaanTimeCurrent2   = (islamicTimesArabic && salaahLanguageCounter == 1 && q!=0 && q!=6) ? convertArabicNumber(adhaanTimeCurrent, "salaah", "") : adhaanTimeCurrent
        let jamaahTimeCurrent    = (islamicTime && salaahLanguageCounter == 1 && q != 0) ? (islamicTime12Hour && q != 0) ? time24to12(jamaahTimesI[q]) : jamaahTimesI[q] : (salaah12HourFormat && q != 0) ? time24to12(jamaahTimes[q]) : jamaahTimes[q];
        let jamaahTimeCurrent2   = (islamicTimesArabic && salaahLanguageCounter == 1 && q!=0 && q!=6) ? convertArabicNumber(jamaahTimeCurrent, "salaah", "") : jamaahTimeCurrent
        let puJamaahTimeCurrent  = (islamicTime && puIslamicTimes && salaahLanguageCounter == 1 && q != 0) ? (islamicTime12Hour && q != 0) ? time24to12(puJamaahTimesI[q]) : puJamaahTimesI[q] : (salaah12HourFormat && q != 0) ? time24to12(puJamaahTimes[q]) : puJamaahTimes[q];
        let puJamaahTimeCurrent2 = (islamicTimesArabic && salaahLanguageCounter == 1 && q!=0 && q!=6) ? convertArabicNumber(puJamaahTimeCurrent, "salaah", "") : puJamaahTimeCurrent

        //Write salaah times in
        document.getElementById((rightToLeft && salaahLanguageCounter == 1 || defaultRightToLeft) ? columnC[q] : columnA[q]).innerHTML = lang[salaahLanguageCounter][q]
        document.getElementById(columnB[q]).innerHTML = adhaanTimeCurrent2
        document.getElementById((rightToLeft && salaahLanguageCounter == 1 || defaultRightToLeft) ? columnA[q] : columnC[q]).innerHTML = jamaahTimeCurrent2

        //Write main Salaah to pop up table
        document.getElementById("pu_" + columnA[(q == 0) ? 1 : q]).innerHTML = (puAlternatingTime) ? lang[salaahLanguageCounter][q] : lang[0][q]
        document.getElementById("pu_" + columnC[(q == 0) ? 1 : q]).innerHTML = (puAlternatingTime) ? puJamaahTimeCurrent2 : puJamaahTimeCurrent

        //Apply CSS for salaah rows
        if (salaahLanguageCounter == 1) {
            $(`#${(rightToLeft)?columnC[(q==0||q==6)?1:q]:columnA[(q==0||q==6)?1:q]}`).css(araCss)
            $(`#${(rightToLeft)?columnA[(q==0||q==6)?1:q]:columnC[(q==0||q==6)?1:q]}`).css(araNumCss)
            //pop up
            if(puAlternatingTime){
                $(`#${"pu_" + columnA[(q == 0) ? 1 : q]}`).css(puSAraCss)
            }
        } else {
            $(`#${columnA[(q==0)?1:q]}`).css(engCss)
            $(`#${columnC[(q==0)?1:q]}`).attr('style', '');
            //pop up
            if(puAlternatingTime){
                $(`#${"pu_" + columnA[(q == 0) ? 1 : q]}`).css(puSEngCss)
            }
        }

        

        //Apply next salaah highlighting
        hs = (q > 6) ? q - 6 : q
        if (hs == nextSalaahHighlightC && highlightSalaah) {
            document.getElementById(columnA[nextSalaahHighlightC]).style.color = "var(--highlightSalaah)";
            document.getElementById(columnB[nextSalaahHighlightC]).style.color = "var(--highlightSalaah)";
            document.getElementById(columnC[nextSalaahHighlightC]).style.color = "var(--highlightSalaah)";
            document.getElementById("pu_" + columnA[(hs==0) ? 1 : hs]).style.color = "var(--highlightSalaah)";
            document.getElementById("pu_" + columnC[(hs==0) ? 1 : hs]).style.color = "var(--highlightSalaah)";
        } else {
            document.getElementById((rightToLeft || defaultRightToLeft) ? columnC[(hs == 0) ? 1 : hs] : columnA[(hs == 0) ? 1 : hs]).style.color = "var(--arabicSalaahNames)";
            document.getElementById(                                                                    columnB[(hs == 0) ? 1 : hs]).style.color = "var(--adhanColumn)";
            document.getElementById((rightToLeft || defaultRightToLeft) ? columnA[(hs == 0) ? 1 : hs] : columnC[(hs == 0) ? 1 : hs]).style.color = "var(--jamaahColumn)";
            document.getElementById("pu_" + columnA[(hs==0) ? 1 : hs]).style.color = "var(--tickerEven)";
            document.getElementById("pu_" + columnC[(hs==0) ? 1 : hs]).style.color = "var(--tickerOdd)";
        }
    }

    //Pop up times & labels
    for(bp=0; bp < puColA.length; bp++) {
        let currentOwwalTime = (islamicTime && puIslamicTimes && salaahLanguageCounter == 1) ? (islamicTime12Hour) ? time24to12(oTimesI[bp]) : oTimesI[bp] : (salaah12HourFormat) ? time24to12(oTimes[bp]) : oTimes[bp]
        document.getElementById((OrightToLeft) ? puColB[bp] : puColA[bp]).innerHTML = `<span>${(salaahLanguageCounter == 1 && puAlternatingTime) ? owwalNames[bp] : translate(puColA[bp])}</span>`
        document.getElementById((OrightToLeft) ? puColA[bp] : puColB[bp]).innerHTML = `<span>${(islamicTimesArabic && salaahLanguageCounter == 1 && puAlternatingTime) ? convertArabicNumber(currentOwwalTime, "salaah", "") : currentOwwalTime}</span>`;

        //Highlight Owwal Times in Pop Up
        nextOwwalHighlight = (nextOwwalHighlight > 10) ? 0 : nextOwwalHighlight;
        if(bp == nextOwwalHighlight && highlightOwwal){
            document.getElementById(puColA[nextOwwalHighlight]).style.color = "var(--highlightSalaah)";
            document.getElementById(puColB[nextOwwalHighlight]).style.color = "var(--highlightSalaah)";
        } else{
            document.getElementById(puColA[bp]).style.color = "var(--tickerEven)";
            document.getElementById(puColB[bp]).style.color = "var(--tickerOdd)";
        }

        //Apply CSS for rows
        if (salaahLanguageCounter == 1 && puAlternatingTime) {
            $(`#${(OrightToLeft) ? puColB[bp] : puColA[bp]}`).css(puOAraCss)
        } else {
            $(`#${(OrightToLeft) ? puColB[bp] : puColA[bp]}`).css(puOEngCss)
        }
    }

    if (checkPoster()) {
        mblTextFit(namedFitTextGroups["popUpOwwalTimes"])
        //namedFitTextGroups["popUpOwwalTimes"][0].isTextFitted = true;
    }
    
    //Set counter according to language
    if (salaahLanguage == "Default") salaahLanguageCounter = 0;
    if (salaahLanguage == "Arabic") salaahLanguageCounter = 1
    if (salaahLanguage == "Alternating") salaahLanguageCounter++
    if (masjidNameLanguage) masjidNameCounter++
    if (!masjidNameLanguage) masjidNameCounter = 0
 
        //Reset counters 
        salaahLanguageCounter =  (salaahLanguageCounter > 1) ? salaahLanguageCounter = 0 : salaahLanguageCounter = salaahLanguageCounter
        masjidNameCounter     =  (masjidNameCounter     > 1) ? masjidNameCounter     = 0 : masjidNameCounter     = masjidNameCounter

    //Setting timeouts
    let salaahLanguageTimouts = [islamicTimeInterval,westernTimeInterval]
                           clearTimeout(salaahLanguageTimout)
    salaahLanguageTimout = setTimeout(salaahLanguageCheck, salaahLanguageTimouts[salaahLanguageCounter]);

}


function boardTranslator(){
    let primaryLang = "en"
    let secondLang = "ur"
    let currentLang = primaryLang;
    
    let columnA = ["msttop", "fajr", "dhuhr", "asr", "maghrib", "esha", "msttop", "fajr", "dhuhr", "asr", "maghrib", "esha"]
    let columnB = ["adhanColumn", "fajrAthan", "dhuhrAthan", "asrAthan", "maghribAthan", "eshaAthan", "adhanColumn", "fajrAthan", "dhuhrAthan", "asrAthan", "maghribAthan", "eshaAthan"]
    let columnC = ["jamaahColumn", "fajrJamaah", "dhuhrJamaah", "asrJamaah", "maghribJamaah", "eshaJamaah", "jamaahColumn", "fajrJamaah", "dhuhrJamaah", "asrJamaah", "maghribJamaah", "eshaJamaah"]
}





// ~~~~~~~~~~~~~~ Times Box ~~~~~~~~~~~~~~ //
function timesBoxCalc() {
    function timesBox(time, endTime, label, zawaal) {
        let head = document.getElementById("timesBoxHead");
        let start = document.getElementById("timesBoxStart");
        let to = document.getElementById("timesBoxTo");
        let end = document.getElementById("timesBoxEnd");

        if (salaah12HourFormat) {
            time = time24to12(time)
            endTime = time24to12(endTime)
        }

        _debug("Times Box","Calculating next time - "+label,secondstotime(time),logStyles.function,false)

        if (zawaal) {
            head.innerHTML = label;
            start.innerHTML = time;
            to.innerHTML = translate('TBZawaalTo');
            end.innerHTML = endTime;
        } else {
            head.innerHTML = label;
            start.innerHTML = "";
            to.innerHTML = time;
            end.innerHTML = "";
        }

        let endTimeSec = minutestoseconds(time);
        // endtime above check why its working with time & not endtime 
        let timeToNext = (endTimeSec > timeInSeconds) ? endTimeSec - timeInSeconds : (86400 - timeInSeconds) + endTimeSec;

        _debug("Times Box","Next calculation - "+label,secondstotime(timeToNext),logStyles.function,false)

        if (!namedFitTextGroups["timesBoxHeader"][0].isTextFitted) {
            mblTextFit(namedFitTextGroups["timesBoxHeader"])
            namedFitTextGroups["timesBoxHeader"][0].isTextFitted = true;
        }


        setTimeout(() => { timesBoxCalc() }, timeToNext * 1000 + 2000)
    }
    namedFitTextGroups["timesBoxHeader"][0].isTextFitted = false;
    if (timeInSeconds > 0 && timeInSeconds < minutestoseconds(sehriEnds)) {
        timesBox(sehriEnds, fajrStarts, (minutestoseconds(fajrStarts) == minutestoseconds(sehriEnds)) ? translate('TBSuhrFajrStarts') + fajrSuhurCheck(2) : translate('TBSuhrEnds') + fajrSuhurCheck(2), false)
    } 
    else if (timeInSeconds > minutestoseconds(sehriEnds) && timeInSeconds < minutestoseconds(fajrStarts) && minutestoseconds(fajrStarts) > minutestoseconds(sehriEnds)) {
        timesBox(fajrStarts, sunrise, translate('TBFajrComm'), false)
    } 
    else if (timeInSeconds > minutestoseconds(fajrStarts) && timeInSeconds < minutestoseconds(sunrise)) {
        timesBox(sunrise, ishraaq, translate('TBSunrise'), false)
    } 
    else if (timeInSeconds > minutestoseconds(sunrise) && timeInSeconds < minutestoseconds(ishraaq)) {
        timesBox(ishraaq, duha, translate('TBIshraq'), false)
    } 
    else if (timeInSeconds > minutestoseconds(ishraaq) && timeInSeconds < minutestoseconds(duha)) {
        timesBox(duha, istiwa, (chast) ? translate('TBChast') : translate('TBDuha'), false)
    } 
    else if (timeInSeconds > minutestoseconds(duha) && timeInSeconds < minutestoseconds(istiwa)) {
        timesBox(istiwa, zawaalEnd, translate('TBZawaal'), false)
    }
    else if (timeInSeconds > minutestoseconds(istiwa) && timeInSeconds < minutestoseconds(zawaalEnd)) {
        timesBox(zawaalEnd, asrShafi, translate('TBZuhr'), false)
    } 
    else if (timeInSeconds > minutestoseconds(zawaalEnd) && timeInSeconds < minutestoseconds(asrShafi)) {
        timesBox(asrShafi, asrHanafi, translate('TBAsrS'), false)
    } 
    else if (timeInSeconds > minutestoseconds(asrShafi) && timeInSeconds < minutestoseconds(asrHanafi)) {
        timesBox(asrHanafi, sunset, translate('TBAsrH'), false)
    } 
    else if (timeInSeconds > minutestoseconds(asrHanafi) && timeInSeconds < minutestoseconds(sunset)) {
        timesBox(sunset, maghribAthan, translate('TBSunset'), false)
    } 
    else if (timeInSeconds > minutestoseconds(sunset) && timeInSeconds < minutestoseconds(maghribAthan)) {
        timesBox(maghribAthan, eshaStarts, translate('TBMaghrib'), false)
    } 
    else if (timeInSeconds > minutestoseconds(maghribAthan) && timeInSeconds < minutestoseconds(eshaStarts)) {
        timesBox(eshaStarts, sehriEnds, translate('TBEsha'), false)
    } 
    else if (timeInSeconds > minutestoseconds(eshaStarts) && timeInSeconds < 86400) {
        timesBox(sehriEnds, "00:00", (minutestoseconds(fajrStarts) == minutestoseconds(sehriEnds)) ? translate('TBSuhrFajrStarts') + fajrSuhurCheck(2) : translate('TBSuhrEnds') + fajrSuhurCheck(2), false)
    }
}

// ~~~~~~ Highlight Salaah Times Row ~~~~~ //
function salaahRowHighlight() {
    //Check if Salaah Highlight is enabled before continuing
    if (!highlightSalaah) { return; }

    //Date from PC, this will automatically have the correct TimeZone unlike pHp
    let dateTime = new Date();
    let timeMil = Date.now();

    //Send dateTime to midnight last & next
    let lastMidnight = dateTime.setHours(0, 0, 0, 0); // last midnight
    let nextMidnight = lastMidnight + 86400000; // next midnight

    //Arrays
    let salaahTimes = [["midnightL", lastMidnight],["fajr", fajrJamaah],["dhuhr", dhuhrJamaah],["asr", asrJamaah],["maghrib", maghribAthan],["esha", eshaJamaah],["midnightN", nextMidnight]];
    let salaahWithTimes = []; //Array with populate only with valid Salaah times

    //Check for any Salaah times that are not valid
    for (x = 0; x < salaahTimes.length; x++) {
        if (salaahTimes[x][1] != "~~~~" || salaahTimes[x][1] != "#REF" || salaahTimes[x][1] != null || salaahTimes[x][1] != 0 || salaahTimes[x][1] != "") {
            salaahWithTimes.push([salaahTimes[x][0], salaahTimes[x][1]]);
        }
    }

    //Loop through Salaah with times to check for next Salaah time
    for (i = 0; i < salaahWithTimes.length - 1; i++) {
        let salaahTimeOne = lastMidnight + minutestoseconds(salaahWithTimes[i][1]) * 1000;
        let salaahTimeTwo = lastMidnight + minutestoseconds(salaahWithTimes[i + 1][1]) * 1000;

        let nexSalaahTime = minutestoseconds(salaahWithTimes[i + 1][1]);
        let timeout = (nexSalaahTime > timeInSeconds) ? (nexSalaahTime - timeInSeconds) * 1000 : 86400 * 1000;

        //Highlight next salaah
        if (timeMil > salaahTimeOne && timeMil < salaahTimeTwo) {
            nextSalaahHighlight = salaahWithTimes[i + 1][0]
            $("#" + salaahWithTimes[i  ][0]           ).css("color", "var(--arabicSalaahNames)");
            $("#" + salaahWithTimes[i  ][0] + "Athan" ).css("color", "var(--adhanColumn)");
            $("#" + salaahWithTimes[i  ][0] + "Jamaah").css("color", "var(--jamaahColumn)");
            $("#" + salaahWithTimes[i+1][0]           ).css("color", "var(--highlightSalaah)"); //Salaah Name
            $("#" + salaahWithTimes[i+1][0] + "Athan" ).css("color", "var(--highlightSalaah)"); //Adhaan Column
            $("#" + salaahWithTimes[i+1][0] + "Jamaah").css("color", "var(--highlightSalaah)"); //Jamaah Column
            setTimeout(salaahRowHighlight, timeout)
            return;
        }
        // this is for esha to midnight, must highlight fajr
        else if (i == salaahWithTimes.length - 2) {
            nextSalaahHighlight = salaahWithTimes[i + 1][0]
            $("#" + salaahWithTimes[i][0]           ).css("color", "var(--arabicSalaahNames)");
            $("#" + salaahWithTimes[i][0] + "Athan" ).css("color", "var(--adhanColumn)");
            $("#" + salaahWithTimes[i][0] + "Jamaah").css("color", "var(--jamaahColumn)");
            $("#" + salaahWithTimes[1][0]           ).css("color", "var(--highlightSalaah)"); //Salaah Name
            $("#" + salaahWithTimes[1][0] + "Athan" ).css("color", "var(--highlightSalaah)"); //Adhaan Column
            $("#" + salaahWithTimes[1][0] + "Jamaah").css("color", "var(--highlightSalaah)"); //Jamaah Column
            setTimeout(salaahRowHighlight, timeout)
            return;
        }
    }
}

// ~~~~~~~~ Highlight Owwal Times ~~~~~~~~ //
function owwalHighlight(){
    //Check if Salaah Highlight is enabled before continuing
    //if (!highlightSalaah) { return; }

    //Date from PC, this will automatically have the correct TimeZone unlike pHp
    let dateTime = new Date();
    let timeMil = Date.now();

    //Send dateTime to midnight last & next
    let lastMidnight = dateTime.setHours(0, 0, 0, 0); //last midnight
    let nextMidnight = lastMidnight + 86400000; // next midnight

    let owwalTimes = [lastMidnight,sehriEnds,fajrStarts,sunrise,ishraaq,duha,istiwa,zawaalEnd,asrShafi,asrHanafi,sunset,eshaStarts,nextMidnight];

    //loop for owwal times
    for(o = 0; o < owwalTimes.length; o++){
        let owwalTimeOne = lastMidnight + minutestoseconds(owwalTimes[o]) * 1000;
        let owwalTimeTwo = lastMidnight + minutestoseconds(owwalTimes[o + 1]) * 1000;

        //Highlight next salaah
        if (timeMil > owwalTimeOne && timeMil < owwalTimeTwo) {
            nextOwwalHighlight = o
            return;
        }
        // this is for esha to midnight, must highlight fajr
        else if (o == owwalTimes.length - 2) {
            nextOwwalHighlight = o;
            return;
        }
    }

    //let nexOwwalTime = minutestoseconds(owwalTimes[o + 1]);
    //let oTimeout = (nexOwwalTime > timeInSeconds) ? (nexOwwalTime - timeInSeconds) * 1000 : 86400 * 1000;

    //console.log("Next Owwal Hlight Value: "+nextOwwalHighlight)
    //console.log("Next Owwal time:"+nexOwwalTime)
    //console.log("Owwal Timeout: "+oTimeout)

    //setTimeout(owwalHighlight, oTimeout)
}

// ~~~~~~~~~~~~ Moon Direction ~~~~~~~~~~~ //
function moonPopUp() {
    var az0 = parseInt(moonAzimuth0);
    var az1 = parseInt(moonAzimuth1);

    var mba = maghribAthan;

    var todaydate = new Date();
    var today = todaydate.setHours(0, 0, 0, 0);
    var date = todaydate.getDate();

    var b = mba.split(":");
    var mga = (+b[0]) * 60 * 60 + (+b[1]) * 60;

    let moonBirthScriptDate = new Date(moonBirthScript);
    let moonVisibilScriptDate = new Date(moonVisibilScript);

    let showTimeOut = (mga + 540) - timeInSeconds;

    if ((moonBirthScriptDate.setHours(0, 0, 0, 0) == today || moonVisibilScriptDate.setHours(0, 0, 0, 0) == today) && timeInSeconds < mga + 2700) {
        document.getElementById("masjidImage").src = 'https://premium.masjidboardlive.com/v2/images/new_moon/' + masjidImage + ".jpg";

        //if (debug) console.log("Moon pop in: " + showTimeOut + "s - Hiding in: " + (mga + 2700) - timeInSeconds + "s")
        _debug("Moon Pop Up","Show/Hide",showTimeOut+"/"+((mga + 2700) - timeInSeconds) + "s",logStyles.function,false)

        setTimeout(() => {
            //hidePoster("moon","Moon pop up", true)
            document.getElementById("moonShow").classList.remove("hidden");
            document.getElementById("moonShow").classList.add("moonShow");
            document.getElementById("moonVisiCell").classList.add("flashmoonvis");
            $('#moonPopUp').removeClass('hidden');
            $('#moonPopUp').addClass('show');
            $('#moonPopUp').css("z-index", "130")

            if (moonBirthScriptDate.getDate() == date) {
                _debug("Moon Pop Up","Pop Up Activated","Birth Date",logStyles.function,false)
                document.getElementById("moonBirth").innerHTML        = `<span>${moonBirth}</span>`;
                document.getElementById("moonAge").innerHTML          = `<span>${moonAge0}</span>`;
                document.getElementById("moonVisibil").innerHTML      = `<span>${bestVisibil}</span>`;
                document.getElementById("moonSet").innerHTML          = `<span>${moonSet0}</span>`;
                //document.getElementById("moonVisibilFor").innerHTML = `<span>${moonVisibilFor0}</span>`;
                document.getElementById("moonAzimuth").innerHTML      = `<span>${moonAzimuth0 + "°"}</span>`;
                document.getElementById("moonAltitude").innerHTML     = `<span>${moonAltitude0 + "°"}</span>`;
                $("#circle").css({ "transform": "rotate(" + az0 + "deg)", });
                $("#circle2").css({ "transform": "rotate(" + ( /*azim-360*/ 0) + "deg)", });
                $("#masjidImage").css({ "transform": "rotate(" + 0 + "deg)", });
            } else if (moonVisibilScriptDate.getDate() == date) {
                _debug("Moon Pop Up","Pop Up Activated","Visibility Date",logStyles.function,false)
                document.getElementById("moonBirth").innerHTML        = `<span>${moonBirth}</span>`;
                document.getElementById("moonAge").innerHTML          = `<span>${moonAge1}</span>`;
                document.getElementById("moonVisibil").innerHTML      = `<span>${bestVisibil}</span>`;
                document.getElementById("moonSet").innerHTML          = `<span>${moonSet1}</span>`;
                //document.getElementById("moonVisibilFor").innerHTML = `<span>${moonVisibilFor1}</span>`;
                document.getElementById("moonAzimuth").innerHTML      = `<span>${moonAzimuth1 + "°"}</span>`;
                document.getElementById("moonAltitude").innerHTML     = `<span>${moonAltitude1 + "°"}</span>`;
                $("#circle").css({ "transform": "rotate(" + az1 + "deg)", });
                $("#circle2").css({ "transform": "rotate(" + ( /*azim-360*/ 0) + "deg)", });
                $("#masjidImage").css({ "transform": "rotate(" + 0 + "deg)", });
            }
            let hideTimeOut = 2700
            setTimeout(() => {
                if (debug) console.log("Moon pop up now hiding")
                _debug("Moon Pop Up","Pop Up Deactivated","Hiding pop up",logStyles.function,false)
                //hidePoster("moon","Moon pop up", false)
                document.getElementById("moonShow").classList.remove("moonShow");
                document.getElementById("moonShow").classList.add("hidden");
                $('#moonPopUp').removeClass('show');
                $('#moonPopUp').addClass('hidden');
                $('#moonPopUp').css("z-index", "-1")
            }, hideTimeOut * 1000)
            mblTextFit(namedFitTextGroups["newMoonPopUp"])
        }, showTimeOut * 1000)
    }
}

// ~~~~~~~~~ Dua After Athan Show ~~~~~~~~ //
function duaAfterAdhanCalc() {
    //Date from PC, this will automatically have the correct TimeZone unlike pHp
    let dateTime = new Date();
    let timeMil = Date.now();

    //Send dateTime to midnight last & next
    let lastMidnight = dateTime.setHours(0, 0, 0, 0); // last midnight
    let nextMidnight = lastMidnight + 86400000; // next midnight

    //Arrays
    let adhanTimes = [
        ["midnightL", lastMidnight],
        ["fajr", minutestoseconds(fajrAthan)],
        ["dhuhr", minutestoseconds(dhuhrAthan)],
        ["asr", minutestoseconds(asrAthan)],
        ["maghrib", minutestoseconds(maghribAthan)],
        ["esha", minutestoseconds(eshaAthan)],
        ["midnightN", nextMidnight]
    ];
    let adhanWithTimes = []; //Array with populate only with valid Adhan times

    //Check for any Adhan times that are not valid
    for (x = 0; x < adhanTimes.length; x++) {
        if (adhanTimes[x][1] != "~~~~" || adhanTimes[x][1] != "#REF" || adhanTimes[x][1] != null || adhanTimes[x][1] != 0 || adhanTimes[x][1] != "" || adhanTimes[x][1] != "FALSE" || adhanTimes[x][1] != "TRUE") {
            adhanWithTimes.push([adhanTimes[x][0], adhanTimes[x][1]]);
        }
    }
    
    //Loop through Adhan with times to check for next Adhan time
    for (i = 0; i < adhanWithTimes.length - 1; i++) {
        let adhanTimeOne = lastMidnight + adhanWithTimes[i][1] * 1000;
        let adhanTimeTwo = lastMidnight + adhanWithTimes[i + 1][1] * 1000;
        
        
        //Pop up Dua After Adhan pop up
        if (timeMil > adhanTimeOne && timeMil < adhanTimeTwo) {
            
            let showTimeout = (adhanTimeTwo) - timeMil;
            let hideTimeout = 180000;
            
            
            _debug("Dua after Adhan",adhanWithTimes[i+1][0],"Showing in: " + Math.round(showTimeout / 1000) + "s",logStyles.poster,false)

            setTimeout(() => {
                _debug("Dua after Adhan",adhanWithTimes[i+1][0],"Hiding in: " + Math.round(hideTimeout / 1000) + "s",logStyles.poster,false)
                document.getElementById("duaAfterAdhan").classList.remove("hidden");
                document.getElementById("duaAfterAdhan").classList.add("duaAfterAdhanShow");
                document.getElementById("duaAfterAdhanPopUp").classList.remove("hidden");
                document.getElementById("duaAfterAdhanPopUp").classList.add("show");
                document.getElementById("duaAfterAdhanPopUp").style.zIndex = 109;
                setTimeout(() => {
                    _debug("Dua after Adhan","Timeout complete","Pop up hidden",logStyles.poster,false)
                    document.getElementById("duaAfterAdhan").classList.remove("duaAfterAdhanShow");
                    document.getElementById("duaAfterAdhan").classList.add("hidden");
                    document.getElementById("duaAfterAdhanPopUp").classList.remove("show");
                    document.getElementById("duaAfterAdhanPopUp").classList.add("hidden");
                    document.getElementById("duaAfterAdhanPopUp").style.zIndex = -1;
                    duaAfterAdhanCalc();
                }, hideTimeout)
            }, showTimeout)
            return;
        }
    }
}

// ~~~~~~~~~~~ Cell Phone PopUp ~~~~~~~~~~ //
function cellPhoneNoticeCalc() {
    //Date from PC, this will automatically have the correct TimeZone unlike pHp
    let dateTime = new Date();
    let timeMil = Date.now();

    //Send dateTime to midnight last & next
    let lastMidnight = dateTime.setHours(0, 0, 0, 0); // last midnight
    let nextMidnight = lastMidnight + 86400000; // next midnight

    //Arrays
    let cellTimes = [
        ["midnightL", lastMidnight],
        ["fajr", fajrJamaah],
        ["dhuhr", dhuhrJmaah],
        ["asr", asrJamaah],
        ["maghrib", maghribAthan],
        ["esha", eshaJamaah],
        ["midnightN", nextMidnight]
    ];
    let cellWithTimes = []; //Array with populate only with valid cell times

    //Check for any cell times that are not valid
    for (x = 0; x < cellTimes.length; x++) {
        if (cellTimes[x][1] != "~~~~" || cellTimes[x][1] != "#REF" || cellTimes[x][1] != null || cellTimes[x][1] != 0 || cellTimes[x][1] != "") {
            cellWithTimes.push([cellTimes[x][0], cellTimes[x][1]]);
        }
    }

    //Loop through cell with times to check for next cell time
    for (i = 0; i < cellWithTimes.length - 1; i++) {
        let cellTimeOne = lastMidnight + minutestoseconds(cellWithTimes[i][1]) * 1000;
        let cellTimeTwo = lastMidnight + minutestoseconds(cellWithTimes[i + 1][1]) * 1000;

        //pop up cell
        if (timeMil > cellTimeOne && timeMil < cellTimeTwo) {
            let showBefore = (cellTimeTwo) - 120000 - timeMil;
            let hideAfter = 180000;

            _debug("Cellphone Pop Up",cellWithTimes[i+1][0] + ":" + cellWithTimes[i+1][1],"Showing in: "+Math.round(showBefore / 1000 / 60) + "m",logStyles.function,false)
            
            
            setTimeout(() => {
                _debug("Cellphone Pop Up","Pop up activated - hiding in",Math.round(hideAfter / 1000 / 60) + "m",logStyles.function,false)
                document.getElementById("cellPhoneNotice").classList.remove("hidden");
                document.getElementById("cellPhoneNotice").classList.add("cellphoneShow");
                document.getElementById("cellPhonePopUp").classList.remove('hidden');
                document.getElementById("cellPhonePopUp").classList.add('show');
                document.getElementById("cellPhonePopUp").style.zIndex = 120
                setTimeout(() => {
                    _debug("Cellphone Pop Up","Hiding pop up","",logStyles.function,false)
                    document.getElementById("cellPhoneNotice").classList.add("hidden");
                    document.getElementById("cellPhoneNotice").classList.remove("cellphoneShow");
                    document.getElementById("cellPhonePopUp").classList.remove('show');
                    document.getElementById("cellPhonePopUp").classList.add('hidden');
                    document.getElementById("cellPhonePopUp").style.zIndex = -100
                    cellPhoneNoticeCalc();
                }, hideAfter)
            }, showBefore)
            return;
        }
    }
}

var lastTime = (new Date()).getTime();
function detectSleep(){
    var currentTime = (new Date()).getTime();
    if (currentTime > (lastTime + 20000*2)) {  // ignore small delays
        _debug("SLEEP","SLEEP DETECTED","RELOADING PAGE",logStyles.function,false)
        loadingScreen("Updating Times", true, false, 50)
        setTimeout(window.location.reload(true),3000);
    }
    lastTime = currentTime;
}



// !====================================================== //
// !=============== Slider Poster Functions ============== //
// !====================================================== //

// ~~~~~~~ Slider Show & Hide Population ~~~~~~ //
let sliderRunning    = false
let enabledA4Posters = []
function populateSlider(){
    enabledA4Posters = []
    //Ayah
    if(ayahTrueFalse === "Display") enabledA4Posters.push("ayahSlide")

    //Hadith
    if(hadithTrueFalse === "Display") enabledA4Posters.push("hadithSlide")

    //Sunnah
    if(sunnahTrueFalse === "Display") enabledA4Posters.push("sunnahSlide")

    //Community Braodcast Text
    if(commBroadTrueFalse === "Display") enabledA4Posters.push("commBroadcast")

    //Text Announcements
    let announcements = ["announceTrueFalse0", "announceTrueFalse1", announceTrueFalse2, announceTrueFalse3, announceTrueFalse4, announceTrueFalse5, announceTrueFalse6, announceTrueFalse7, announceTrueFalse8, announceTrueFalse9, announceTrueFalse10]
    for(a=2; a<announcements.length; a++){
        if(announcements[a] === "Display") enabledA4Posters.push(`announceSlide${a}`)
    }

    //A4 Posters
    let a4Posters  = ["posterTrueFalse0", posterTrueFalse1, posterTrueFalse2, posterTrueFalse3, posterTrueFalse4, posterTrueFalse5, posterTrueFalse6, posterTrueFalse7, posterTrueFalse8, posterTrueFalse9, posterTrueFalse10]
    let a4PosterId = ["posterImage0"    , posterImage1    , posterImage2    , posterImage3    , posterImage4    , posterImage5    , posterImage6    , posterImage7    , posterImage8    , posterImage9    , posterImage10    ]
    for(p=1; p<a4Posters.length; p++){
        if(a4Posters[p] === "Display") enabledA4Posters.push(`posterSlide${p}`)
        //document.getElementById(`posterImage${p}`).src = (a4PosterId[p].substr(0,4) == "http") ? a4PosterId[p] : 'https://drive.google.com/uc?id=' + a4PosterId[p];
        //document.getElementById(`posterImage${p}`).src = (a4PosterId[p].substr(0,4) == "http") ? a4PosterId[p] : 'https://lh3.google.com/u/0/d/' + a4PosterId[p];
        document.getElementById(`posterImage${p}`).src = (a4PosterId[p].substr(0,4) == "http") ? a4PosterId[p] : 'https://api.masjidboardlive.com/imageproxy?id=' + a4PosterId[p];

        $(`#posterImage${p}`).css({ "padding-top": posterTopPadding + 'px',  });
        $(`#posterImage${p}`).css({ "height"     : posterHeight + 'px',      });
        $(`#posterImage${p}`).css({ "width"      : posterWidth + 'px',       });
        $(`#posterImage${p}`).css({ "z-index"    : '40',                     });
    }

    //Community A4 Posters
    let a4PostersC  = ["posterTrueFalse0", commPosterTrueFalse1, commPosterTrueFalse2, commPosterTrueFalse3, commPosterTrueFalse4, commPosterTrueFalse5, commPosterTrueFalse6, commPosterTrueFalse7, commPosterTrueFalse8, commPosterTrueFalse9, commPosterTrueFalse10]
    let a4PosterIdC = ["posterImage0"    , commPosterImage1    , commPosterImage2    , commPosterImage3    , commPosterImage4    , commPosterImage5    , commPosterImage6    , commPosterImage7    , commPosterImage8    , commPosterImage9    , commPosterImage10    ]
    for(pc=1; pc<a4PostersC.length; pc++){
        if(a4PostersC[pc] === "Display") enabledA4Posters.push(`commPoster${pc}`)
        document.getElementById(`commPosterImage${pc}`).src = (a4PosterIdC[pc].substr(0,4) == "http") ? a4PosterIdC[pc] : 'https://api.masjidboardlive.com/imageproxy?id=' + a4PosterIdC[pc];
        $(`#commPosterImage${pc}`).css({ "padding-top": posterTopPadding + 'px',  });
        $(`#commPosterImage${pc}`).css({ "height"     : posterHeight + 'px',      });
        $(`#commPosterImage${pc}`).css({ "width"      : posterWidth + 'px',       });
        $(`#commPosterImage${pc}`).css({ "z-index"    : '40',                     });
    }

    //Sick list
    if (sickTrueFalse === "Display"){
        enabledA4Posters.push("sickSlide")
        let sicklist = sick1 + "<br>" + sick2 + "<br>" + sick3 + "<br>" + sick4 + "<br>" + sick5 + "<br>" + sick6 + "<br>" + sick7 + "<br>" + sick8 + "<br>" + sick9 + "<br>" + sick10;
        document.getElementById("sickC").innerHTML = sicklist;
    }

    //Eidgah 
    if (eidgahTrueFalse === "Display"){
        enabledA4Posters.push("eidSlide")
        document.getElementById("eidSlide").innerHTML = `
            <p class="sliderMainHeadings" id="eidSalaahHeading">${translate("eidSalaahHeading")}</p>
            <span id="eidSlideC">
                <p class="fSliderInfoHeadings" id="eidSalaahHeadingDate">${translate("eidSalaahHeadingDate")}</p>
                <p class="sliderTextWr" id="eidDate">${eidDate}</p>
                <p class="fSliderInfoHeadings" id="eidSalaahHeadingVenue">${translate("eidSalaahHeadingVenue")}</p>
                <div id="eidVenueC"><p class="sliderTextWr" id="eidVenue">${eidVenue}</p></div>
                <p class="fSliderInfoHeadings" id="eidSalaahHeadingAddress">${translate("eidSalaahHeadingAddress")}</p>
                <div id="eidAddressC"><p class="sliderTextWr" id="eidAddress">${eidAddress}</p></div>
                <p class="fSliderInfoHeadings" id="eidSalaahHeadingLecture">${translate("eidSalaahHeadingLecture")}</p>
                <p class="sliderTextWr" id="eidLecture">${eidLecture}</p>
                <p class="fSliderInfoHeadings" id="eidSalaahHeadingSalah">${translate("eidSalaahHeadingSalah")}</p>
                <p class="sliderTextWr" id="eidSalaah">${eidSalaah}</p>
            </span>`;
    }

    //Ladies Taleem
    let tf1 = (mobile == false);
    let tf2 = (taleemTrueFalse1 === "Display" && taleemTrueFalse2 === "Display" && mobile == false);
    if (taleemTrueFalse1 === "Display" && mobile == false){
        enabledA4Posters.push("taleemSlide")
        document.getElementById("taleemSlide").innerHTML = `
            <p class="sliderMainHeadings" id="ladiesTaleemHeading">${translate("ladiesTaleemHeading")}</p>
            <span id="taleemSlideC">    
                <p id="taleemTime1" class="fSliderInfoHeadings">${(tf1) ? taleemTime1 : ""}</p>
                <p id="taleemDate1">${(tf1) ? taleemDate1 : ""}</p>
                <div id="taleemResi1C">
                    <p id="taleemResi1">${(tf1) ? taleemResi1 : ""}</p>
                </div>
                <div id="taleemAddy1C">
                    <p id="taleemAddy1">${(tf1) ? taleemAddy1 : ""}</p>
                </div>
                <br>
                <p id="taleemTime2" class="fSliderInfoHeadings">${(tf2) ? taleemTime2 : ""}</p>
                <p id="taleemDate2">${(tf2) ? taleemDate2 : ""}</p>
                <div id="taleemResi2C">
                    <p id="taleemResi2">${(tf2) ? taleemResi2 : ""}</p>
                </div>
                <div id="taleemAddy2C">
                    <p id="taleemAddy2">${(tf2) ? taleemAddy2 : ""}</p>
                </div>
            </span>`;
    }

    //Dawah 3days & Gasht
    let tf = (gashtTrueFalse === "Display" && dawahTrueFalse === "Display")
    if(gashtTrueFalse === "Display"){
        enabledA4Posters.push("dawahSlide")
        document.getElementById("dawahSlide").innerHTML = `
            <span id="dawahSlideC">
                <p id="masjidTaleem">${masjidTaleem}</p>
                <div class="tabline"></div>
                <p class="gashtInOutHeadings" id="gashtInHeading">${translate("gashtInHeadingNew")}<span id="gashtDayIn" style="color: var(--gashtDay)">${gashtDayIn}</span></p>
                <div class="tabline"></div>
                <p class="gashtInOutHeadings" id="gashtOutHeading">${translate("gashtOutHeadingNew")}<span id="gashtDayOut" style="color: var(--gashtDay)">${gashtDayOut}</span></p>
                <div>
                    <p id="gashtOut1">${gashtOut1}</p>
                    <p id="gashtOut2">${gashtOut2}</p>
                </div>
                <div class="tabline"></div>
                <p id="threeDaysHeader">${(tf) ? translate("threeDaysHeader") : ""}</p>
                <div class="threeDaysBox">
                    <span id="threeDays1">${(tf) ? threeDays1 : ""}</span>
                    <span id="threeDays1Date">${(tf) ? threeDays1Date : ""}</span>
                </div>
                <div class="threeDaysBox">
                    <span id="threeDays2">${(tf) ? threeDays2 : ""}</span>
                    <span id="threeDays2Date">${(tf) ? threeDays2Date : ""}</span>
                </div>
            </span>`;
    }

    //Bank Slide
    if(bankTrueFalse === "Display"){
        enabledA4Posters.push("bankSlide")
        document.getElementById("bankSlide").innerHTML = `
            <span id="bankSlideC">
                <p class="sliderMainHeadings" id="bankHeader">${bankHeader}</p>
                <p class="fSliderInfoHeadings" id="bankSlideHeading">${translate("bankSlideHeading")}</p>
                <p class="sliderTextWr" id="bankName">${bankName}</p>
                <p class="fSliderInfoHeadings" id="bankSlideHeadingBC">${translate("bankSlideHeadingAN")}</p>
                <p class="sliderTextWr" id="bankCode">${accountName}</p>
                <p class="fSliderInfoHeadings" id="bankSlideHeadingAN">${(stringToBoolean(bankAUSBSB)) ? translate("bankSlideHeadingBC_AU") : translate("bankSlideHeadingBC")}</p>
                <p class="sliderTextWr" id="accountName">${bankCode}</p>
                <p class="fSliderInfoHeadings" id="bankSlideHeadingANum">${translate("bankSlideHeadingANum")}</p>
                <p class="sliderTextWr" id="accountNum">${accountNum}</p>
            </span>`;
    }

    //Nikaah Notice (Slider ONLY, not popup)
    if(nikahTrueFalse === "Display"){
        enabledA4Posters.push("nikahSlide")
        document.getElementById("nikahSlide").innerHTML = `
            <p class="nikahMainHeadings" id="nikaahHeading">${translate("nikaahHeading")}</p>
            <span id="nikaahSlideSpan">
                <p class="nikahInfoHeadings nikahHeadings" id="nikaahGroom">${(nikahNameOne =="") ? "" : translate("nikaahGroom")}</p>
                    <p class="nikahText" id="nikahNameOne">${nikahNameOne}</p>
                <p class="nikahInfoHeadings nikahHeadings" id="nikahGroomRelation">${translate(nikahGroomRelation)}</p>
                    <p class="nikahText" id="nikahRelationOne">${nikahRelationOne}</p>
                ${(nikahBride == "") ? "" : "<p class='nikahInfoHeadings nikahHeadings'>"+translate("nikaahBride")+"</p>"}
                    ${(nikahBride == "") ? "" : ("<p class='nikahText'>"+nikahBride)+"</p>"}
                ${(nikahBride=="") ? ("<p class='nikahInfoHeadings nikahHeadings'>"+translate("nikaahBride")+": "+"<span class='nikahInfoHeadings' id='nikahRelationTwo'>"+translate(nikahRelationTwo)+"</span></p>") : "<p class='nikahInfoHeadings nikahHeadings'>"+translate(nikahRelationTwo)+"</p>"}
                    <p class="nikahText" id="nikahNameTwo">${nikahNameTwo}</p>
                <p class="nikahInfoHeadings nikahHeadings" id="nikahDateHeading">${(nikahDate=="") ? "" : translate("nikahDateHeading")}</p>
                    <p class="nikahText" id="nikahDate">${nikahDate}</p>
                <p class="nikahInfoHeadings nikahHeadings" id="nikahTimeHeading">${(nikahTime=="") ? "" : translate("nikahTimeHeading")}</p>
                    <p class="nikahText" id="nikahTime">${nikahTime}</p>
                <!--<p style="font-size: 50px; font-weight:600; font-family: "Times New Roman", Times, serif">إنْ شَاءَ اللهُ</p>-->
            </span>`;
    }

    //Populate all slides
    //First Hide all slides
    if(document.querySelectorAll(".carousel-item")){
        document.querySelectorAll(".carousel-item").forEach((currentValue, index, arr) => {
            currentValue.classList.add("hidden");
            currentValue.classList.remove("carousel-item");
        })
    }
    if(document.querySelectorAll(".carousel  .active")){
        document.querySelectorAll(".carousel  .active").forEach((currentValue, index, arr) => {
            currentValue.classList.remove("active");
        })
    }

    //if array isnt null
    if(enabledA4Posters.length > 0){
        let sliderIndicators = ""
        for(i=0; i<enabledA4Posters.length; i++){
            document.getElementById(enabledA4Posters[i]).classList.remove("hidden");
            document.getElementById(enabledA4Posters[i]).classList.add("carousel-item");
            sliderIndicators = sliderIndicators + `<li data-target="#carousel" class="${enabledA4Posters[i]}" data-slide-to="${i}"></li>`
        }
        document.getElementById(enabledA4Posters[0]).classList.add("active");
        document.getElementById("sliderIndicators").innerHTML = sliderIndicators
        document.querySelector("#sliderIndicators > li").classList.add("active");

        //reset TextFix
        let textFitSlides = ['hadithSlide','ayahSlide','sunnahSlide','taleemSlide','nikahSlide','eidSlide','bankSlide','sickSlide','dawahSlide']
        textFitSlides.forEach((currentValue, index, arr)=>{
            namedFitTextGroups[currentValue][0].isTextFitted = false;
        })
        sliderMain()
    } else{
        console.log("All Posters Disabled")
    }
}

// ~~ Main Slider function + scroll bar ~~ //
function sliderMain() {
    if(sliderRunning) return
    sliderRunning = true;
    let progBar = document.getElementById('transition-timer-carousel-progress-bar');
    let left    = document.getElementById('left-carousel-control');
    let right   = document.getElementById('right-carousel-control');
    let slider  = document.getElementById('carousel');
    let count   = 1;

    _debug("Slider","Duration",slideDuration + "s",logStyles.function,false)

    //Start prog bar animation
    progBar.style.animation = `sliderProgBar ${(slideDuration == 0) ? slideDuration+10 : slideDuration}s linear 0s ${count} normal none running`;
    progBar.addEventListener("webkitAnimationEnd", () => {
        count++
        $('#carousel').carousel('next')
        progBar.style.animation = `sliderProgBar ${(slideDuration == 0) ? slideDuration+10 : slideDuration}s linear 0s ${count} normal none running`;
    });

    //Hover to pause animation
    slider.addEventListener("mouseover", () => {
        progBar.style.animationPlayState = "paused";
    })

    //Carry on animation after hover
    slider.addEventListener("mouseleave", () => {
        progBar.style.animationPlayState = "running";
    })

    //button for prev
    left.addEventListener("click", () => {
        if (debug) console.log("Previous Slide Manually Triggered")
        $('#carousel').carousel('prev')
    });
    //button for next
    right.addEventListener("click", () => {
        if (debug) console.log("Next Slide Manually Triggered")
        $('#carousel').carousel('next')
    });
    //trigger text fit on next slide, check if needed 1st
    $('#carousel').on('slid.bs.carousel', function(e) {
        let id = e.relatedTarget.id
        if (typeof namedFitTextGroups[id] != "undefined") {
            if (!namedFitTextGroups[id][0].isTextFitted) {
                mblTextFit(namedFitTextGroups[id])
                namedFitTextGroups[id][0].isTextFitted = true;
            }
        }
    });
}



// !====================================================== //
// !================== Pop Up Functions  ================= //
// !====================================================== //

// ~~~~~~~~~ Nikaah Pop Up ~~~~~~~~~ //
function nikahPopUpFunction() {
    if(!nikahPopUpTrueFalse) return
    let nikaahPopUpHTML = `
    <span id="nikaahPopUpSpan">
        <p class="nikahMainHeadings" id="nikaahHeadingPu">${translate("nikaahHeading")}</p>
        <p class="nikahInfoHeadings nikahHeadings" id="nikaahGroomPu">${(nikahNameOne=="") ? "" : translate("nikaahGroom")}</p>
            <p class="nikahText" id="nikahNameOnePu">${nikahNameOne}</p>
        <p class="nikahInfoHeadings nikahHeadings" id="nikahGroomRelationPu">${translate(nikahGroomRelation)}</p>
            <p class="nikahText" id="nikahRelationOnePu">${nikahRelationOne}</p>
        ${(nikahBride == "") ? "" : "<p class='nikahInfoHeadings nikahHeadings'>"+translate("nikaahBride")+"</p>"}
            ${(nikahBride == "") ? "" : ("<p class='nikahText'>"+nikahBride)+"</p>"}
        ${(nikahBride=="") ? ("<p class='nikahInfoHeadings nikahHeadings'>"+translate("nikaahBride")+": "+"<span class='nikahInfoHeadings' id='nikahRelationTwoPu'>"+translate(nikahRelationTwo)+"</span></p>") : "<p class='nikahInfoHeadings nikahHeadings'>"+translate(nikahRelationTwo)+"</p>"}
            <p class="nikahText" id="nikahNameTwoPu">${nikahNameTwo}</p>
        <p class="nikahInfoHeadings nikahHeadings" id="nikahDateHeadingPu">${(nikahDate=="") ? "" : translate("nikahDateHeading")}</p>
            <p class="nikahText" id="nikahDatePu">${nikahDate}</p>
        <p class="nikahInfoHeadings nikahHeadings" id="nikahTimeHeadingPu">${(nikahTime=="") ? "" : translate("nikahTimeHeading")}</p>
            <p class="nikahText" id="nikahTimePu">${nikahTime}</p>
        <!--<p style="font-size: 50px; font-weight:600; font-family: "Times New Roman", Times, serif">إنْ شَاءَ اللهُ</p>-->
    </span>
    `
    document.getElementById("nikahPopUp").innerHTML = nikaahPopUpHTML;

    var nikahDayBefore = (parseInt(nikahPopUp) - parseInt(timeZoneMilli)) - oneDayMilli;
    var nikahDayAfter = (parseInt(nikahPopUp) - parseInt(timeZoneMilli)) + oneDayMilli;

    if (nikahTrueFalse === "Display") {
        if (nikahDayBefore < nowDate && nowDate < nikahDayAfter) {
            _debug("Nikaah Pop Up","Showing today","",logStyles.function,false)
                //hidePoster("nikaah","Nikaah pop up")
            $('#nikahPopUp').removeClass('hidden');
            $('#nikahPopUp').addClass('nikahPopUp')
            $('#nikahPopUpBg').removeClass('hidden');
            $('#nikahPopUpBg').addClass('show');
            $('#nikahPopUpBg').css("z-index", "100")
            //mblTextFit(namedFitTextGroups["nikahPopUp"])
            setTimeout(() => { mblTextFit(namedFitTextGroups["nikahPopUp"]) }, 2000)
        } else {
            $('#nikahPopUp').removeClass('nikahPopUp')
            $('#nikahPopUp').addClass('hidden')
            $('#nikahPopUpBg').removeClass('show');
            $('#nikahPopUpBg').addClass('hidden');
        }
    } else {
        $('#nikahPopUp').addClass('hidden');
        $('#nikahPopUpBg').addClass('hidden');
    }
}

// ~~~~~~~~~~ Funeral Pop Up ~~~~~~~~~ //
function funeralBoard() {
    let funeralNotice = `
        <p class="funeralHeading flashfuneral" id="funeralNoticeH">${translate("funeralNoticeH")}</p>
        <p id="funeralArabic">إنّا للهِ و إنّا إلَيْهِ رَاجِعُوْن</p>
        <span id="funeralBoardSpan">
            <p class="fSliderInfoHeadings funeralHeadings" id="funeralNameH">${translate("funeralNameH")}</p>
                <div id="funeralNameC"><p class="fSliderText" id="funeralName">${funeralName}</p></div>
                <div id="funeralRelationC"><p class="fSliderText" id="funeralRelation">${funeralRelation}</p></div>
            <p class="fSliderInfoHeadings funeralHeadings" id="funeralAddyH">${translate("funeralAddyH")}</p>
                <div id="funeralAddyC"><p class="fSliderText" id="funeralAddy">${funeralAddy}</p></div>
                <div id="funeralPickupC"><p class="fSliderText" id="funeralPickup">${funeralPickup}</p></div>
            <p class="fSliderInfoHeadings funeralHeadings" id="funeralCemeteryH">${translate("funeralCemeteryH")}</p>
                <div id="funeralCemeteryC"><p class="fSliderText" id="funeralCemetery">${funeralCemetery}</p></div>
            <p class="fSliderInfoHeadings funeralHeadings" id="funeralSalaahVenueH">${translate("funeralSalaahVenueH")}</p>
                <div id="funeralSalaahVenueC"><p class="fSliderText" id="funeralSalaahVenue">${funeralSalaahVenue}</p></div>
                <div id="funeralSalaahTimeC"><p class="fSliderText" id="funeralSalaahTime">${funeralSalaahTime}</p></div>
            <br>
        </span>
        <p id="funeralDuaEnglish">${translate("funeralDuaEnglish")}</p>
    `
    
        if (funeralTrueFalse === "Display") {
            hidePoster("funeral","Funeral pop up", true)
            document.getElementById("funeralBoard").innerHTML = funeralNotice;
            document.getElementById("funeralBoard").classList.remove("hidden");
            document.getElementById("funeralBoard").classList.add("funeralShow");
            document.getElementById("funeralPopUp").classList.remove('hidden');
            document.getElementById("funeralPopUp").classList.add('show')
            document.getElementById("funeralPopUp").style.zIndex = 190
            mblTextFit(namedFitTextGroups["funeralBoard"])
        } else {
            hidePoster("funeral","Funeral pop up", false)
            document.getElementById("funeralBoard").classList.remove("funeralShow");
            document.getElementById("funeralBoard").classList.add("hidden");
            $("#funeralPopUp").removeClass('show');
            $("#funeralPopUp").addClass('hidden');
        }
}




// !====================================================== //
// !================= Full Screen Posters ================ //
// !====================================================== //

// ~~~ Declaration of poster variables ~~~ //
var jumuahDuroodState = 0;
var jumuahDuroodTotal = 10



//  Hiding big poster for pop ups or anywhere else  //
function hidePoster(posterid,message, truefalse) {
    //_debug("Hide Poster",message,truefalse,logStyles.warn,false)
    posterHide[posterid] = truefalse;

    //Check if any objects are hiding poster
    for (var key in posterHide) {
        if (posterHide.hasOwnProperty(key)) {
            if(posterHide[key] == true){
                $("#posterContainer").css("display", "none");
                $("#tableTime").switchClass("pu_tableTime", "tableTime", 0);
                //clocknamePositions();
                //showBigPoster("hideAll", 0, "hideAll")
                _debug("Hide Poster","Cannot show poster, hide == true ",posterHide[key],logStyles.warn,false)
                return
            }
        }
    }
}

// ~~~ Check if any posters are active ~~~ //
function checkPoster() {
    if (posters.iftaarPoster.state          == 1 ||
        posters.thursDuroodPoster.state     == 1 ||
        posters.jumuahSunnah.state          == 1 ||
        posters.ayaamPoster.state           == 1 ||
        posters.islamicMonthPoster.state    == 1 ||
        posters.muharramPoster1.state       == 1 ||
        posters.muharramPoster2.state       == 1 ||
        posters.muharramPoster3.state       == 1 ||
        posters.ashuraSpendPoster.state     == 1 ||
        posters.rajabPoster.state           == 1 ||
        posters.shabaanPoster1.state        == 1 ||
        posters.shabaanPoster2.state        == 1 ||
        posters.ramadanApproachPoster.state == 1 ||
        posters.ramadanPoster1.state        == 1 ||
        posters.ramadanPoster2.state        == 1 ||
        posters.ramadanPoster3.state        == 1 ||
        posters.laylatulQadrPoster.state    == 1 ||
        posters.eidDuaPoster.state          == 1 ||
        posters.hajjPoster1.state           == 1 ||
        posters.hajjPoster2.state           == 1 ||
        posters.hajjPoster3.state           == 1 ||
        posters.hajjPoster4.state           == 1 ||
        posters.hajjPoster5.state           == 1 ||
        posters.hajjPoster6.state           == 1 ||
        posters.hajjPoster7.state           == 1 ||
        posters.hajjPoster8.state           == 1 ||
        posters.hajjPoster9.state           == 1 ||
        posters.hajjPoster10.state          == 1 ||
        posters.tashreeqPoster.state        == 1 ||
        posters.bigPosterImage0.state       == 1 ||
        posters.bigPosterImage1.state       == 1 ||
        posters.bigPosterImage2.state       == 1 ||
        posters.bigPosterImage3.state       == 1 ||
        posters.bigPosterImage4.state       == 1 ||
        posters.bigPosterImage5.state       == 1 ||
        posters.bigPosterImage6.state       == 1 ||
        posters.bigPosterImage7.state       == 1 ||
        posters.bigPosterImage8.state       == 1 ||
        posters.bigPosterImage9.state       == 1 ||
        posters.cBigPosterImage0.state      == 1 ||
        posters.cBigPosterImage1.state      == 1 ||
        posters.cBigPosterImage2.state      == 1 ||
        posters.cBigPosterImage3.state      == 1 ||
        posters.cBigPosterImage4.state      == 1 ||
        posters.cBigPosterImage5.state      == 1 ||
        posters.cBigPosterImage6.state      == 1 ||
        posters.cBigPosterImage7.state      == 1 ||
        posters.cBigPosterImage8.state      == 1 ||
        posters.cBigPosterImage9.state      == 1 ||
        posters.moonVisibility.state        == 1 ||
        jumuahDuroodState                   == 1
    ) { return true } else return false
}

// ~~ Main poster function to show/hide ~~ //
async function showBigPoster(posterid, timeout, visibility) {
    if (debug && visibility != "hideAll") console.log(`%c | Name: %c${posterid} \n %c| ID: %c${posters[posterid].id} \n %c| Timeout: %c${timeout} \n %c| Visibility: %c${visibility}`,"",logStyles.info,"",logStyles.info,"",logStyles.info,"",logStyles.info)
    if(isBoardDisabled && posterid != "iftaarPoster") return;
    
    //Check if any objects are hiding poster
    for (var key in posterHide) {
        if (posterHide.hasOwnProperty(key)) {
            if(posterHide[key] == true){
                if (debug) console.log("Cannot show poster, hide == true " + posterHide[key])
                return
            }
        }
    }

    //check if posters exist
    /*if (posterid.slice(0, -1) == 'bigPosterImage'){
        if (await checkFile(posterid) == false || posters[posterid].id == "" || posters[posterid].id == null || posters[posterid].id == undefined) {
            _debug("Big Poster","Poster does not have an ID or does not exist in path",posterid,logStyles.error,false,false)
            return
        }
    } else {
        if (await checkFile(`https://premium.masjidboardlive.com/v2/images/posters/${lang}/${posterid}.jpg`) == false || posters[posterid].id == "" || posters[posterid].id == null || posters[posterid].id == undefined) {
            _debug("Perpetual Poster","Poster does not have an ID or does not exist in path",posterid,logStyles.error,false,false)
            return
        }
    }*/


    if (visibility == "visible" || visibility == "Display") {
        //if (debug) console.log("%c"+ posterid +" Poster Container Active", "color:green")
        if (posterid.slice(0, -1) == 'bigPosterImage' || posterid.slice(0, -1) == 'bigPosterImage1' || posterid.slice(0, -1) == 'cBigPosterImage') {
            //document.getElementById(posterid).src = (posters[posterid].id.substr(0,4) == "http") ? posters[posterid].id : "https://drive.google.com/uc?id="+posters[posterid].id;
            document.getElementById(posterid).src = (posters[posterid].id.substr(0,4) == "http") ? posters[posterid].id : "https://api.masjidboardlive.com/imageproxy?id="+posters[posterid].id;
            $('#' + posterid).css("display", "block");
        } else {
            if(posterid == "moonVisibility"){
                let newMoonBirth = new Date(moonBirth);
                let conjDay      = currentDate.getDate() - newMoonBirth.getDate()
                let url          = `moon_visibility_${islamicDateYear}_${(islamicDateNum > 5) ? islamicMonthNum+1 : islamicMonthNum}_${conjDay}`
                document.getElementById(posterid).src = `https://premium.masjidboardlive.com/v2/images/moon_map/${url}.png?${version}`;
            } else{
                document.getElementById(posterid).src = `https://premium.masjidboardlive.com/v2/images/posters/${lang}/${posterid}.jpg?${version}`;
            }
            $('#' + posterid).css("display", "block");
        }
        owwalHighlight();
        //delay for image to load
            posters[posterid].state = 1;
            isBigPopUpShowing = true;
        setTimeout(() => {
            $("#posterContainer").animate({ opacity: 1 });
            $("#tableTime").removeAttr("style")
            $("#tableTime").switchClass("tableTime", "pu_tableTime", 0)
            $("#posterContainer").css("display", "block");
            $('#' + posterid).css("opacity", 0);
            $('#' + posterid).css("display", "block");
            $('#' + posterid).animate({ opacity: 1 },500);

            //posters[posterid].state = 1;
            //isBigPopUpShowing = true;
            //Pop up date text fit
            //if (!namedFitTextGroups["popUpOwwalTimes"][0].isTextFitted && isBigPopUpShowing) {
                mblTextFit(namedFitTextGroups["popUpDate"])
                mblTextFit(namedFitTextGroups["popUpOwwalTimes"])
            //    namedFitTextGroups["popUpDate"      ][0].isTextFitted = true;
            //    namedFitTextGroups["popUpOwwalTimes"][0].isTextFitted = true;
            //}
        }, 700)

        //Hide after timeout
        setTimeout(() => {
            _debug("Big Poster","Timeout complete - Hiding container","",logStyles.function,false,false)
            //$('#' + posterid).css("display", "none");
            $('#' + posterid).animate({ opacity: 0});
            isBigPopUpShowing = false
            posters[posterid].state = 0;
            if (checkPoster() == false) {
                clocknamePositions();
                $("#tableTime").switchClass("pu_tableTime", "tableTime", 0);
                $("#posterContainer").animate({
                    opacity: 0,
                    complete: () => {
                        $("#posterContainer").css("display", "none");
                    }
                });
                $("#posterContainer").css("display", "none");
            }
        }, timeout);
        if (debug) console.log("%cHiding Poster in " + timeout / 1000 / 60 + " Minutes", "color: orange")
    }
    
    //hide if function posterHide is triggered
    else if (visibility == "hideAll") {
        _debug("Big Poster","Hide all triggered","",logStyles.warn,false,false)
        isBigPopUpShowing = false
        $("#posterContainer").animate({
            opacity: 0,
            complete: () => {
                $("#posterContainer").css("display", "none");
            }
        });
        clocknamePositions();
        $("#tableTime").switchClass("pu_tableTime", "tableTime", 0)
    }

    //Forcing a hide externally
    else {
        if (posters[posterid].state == 1) {
            $("#posterContainer").animate({
                opacity: 0,
                complete: () => {
                    $("#posterContainer").css("display", "none");
                }
            });
            isBigPopUpShowing = false
            clocknamePositions();
            $("#tableTime").switchClass("pu_tableTime", "tableTime", 0)
            //$('#' + posterid).css("display", "none");
            document.getElementById(posterid).src = "";
            posters[posterid].state = 0;
        }
    }
}

//  Iftar PopUp, this one needs to be on its own function  //
var iftaarIsShowing = false;
function iftarPopUpJ() {
    if(mblPosterHide.includes("iftaarPoster")) return;
    if(iftaarIsShowing) return;
    
    _debug("Iftaar poster","function executed","",logStyles.function,false,false)
    var iftarHide = minutestoseconds(maghribAthan) + 60;

    if ((islamicMonthNum == 1 || islamicMonthNum == 2 || islamicMonthNum == 3 || islamicMonthNum == 4 || islamicMonthNum == 5 || islamicMonthNum == 6 || islamicMonthNum == 7 || islamicMonthNum == 8 || islamicMonthNum == 9 || islamicMonthNum == 11) && (timeInSeconds > minutestoseconds(maghribAthan) && timeInSeconds <= iftarHide)) {
        showBigPoster("iftaarPoster", 60000, "visible");
        iftaarIsShowing = true;
    } else if (islamicMonthNum == 10 && islamicDateNum != 1 && timeInSeconds > minutestoseconds(maghribAthan) && timeInSeconds <= iftarHide) {
        showBigPoster("iftaarPoster", 60000, "visible");
        iftaarIsShowing = true;
    } else if (islamicMonthNum == 12 && islamicDateNum != 10 && islamicDateNum != 11 && islamicDateNum != 12 && islamicDateNum != 13 && timeInSeconds > minutestoseconds(maghribAthan) && timeInSeconds <= iftarHide) {
        showBigPoster("iftaarPoster", 60000, "visible");
        iftaarIsShowing = true;
    } else {
        showBigPoster("iftaarPoster", 0, "hide");
        iftaarIsShowing = false;
    }
}

// ~~~~ Ayaam-e-Tashreeq Poster PopUp ~~~~ //
/*function tashreeqPopUp() {
    //Guard if statement, if its not the right month day dont go through the rest of the function
    if (islamicMonthNum != 12 || [9, 10, 11, 12, 13].includes(islamicDateNum) || mblPosterHide.includes("tashreeqPoster")) return;

    //Date from PC, this will automatically have the correct TimeZone unlike pHp
    let dateTime = new Date();
    let timeMil = Date.now();

    //Send dateTime to midnight last & next
    let lastMidnight = dateTime.setHours(0, 0, 0, 0); // last midnight
    let nextMidnight = lastMidnight + 86400000; // next midnight

    //Arrays
    let tashreeqTimes = [
        ["midnightL", lastMidnight],
        ["fajr"     , fajrJamaah  ],
        ["dhuhr"    , dhuhrJamaah ],
        ["asr"      , asrJamaah   ],
        ["maghrib"  , maghribAthan],
        ["esha"     , eshaJamaah  ],
        ["midnightN", nextMidnight]
    ];
    let tashreeqWithTimes = []; //Array with populate only with valid cell times

    //Check for any tashreeq times that are not valid
    for (x = 0; x < tashreeqTimes.length; x++) {
        if (tashreeqTimes[x][1] != "~~~~" || tashreeqTimes[x][1] != "#REF" || tashreeqTimes[x][1] != null || tashreeqTimes[x][1] != 0 || tashreeqTimes[x][1] != "") {
            tashreeqWithTimes.push([tashreeqTimes[x][0], tashreeqTimes[x][1]]);
        }
    }

    //Loop through tashreeeq with times to check for next cell time
    for (i = 0; i < tashreeqWithTimes.length - 1; i++) {
        let tashreeqTimeOne = lastMidnight + (minutestoseconds(tashreeqWithTimes[i][1]) * 1000);
        let tashreeqTimeTwo = lastMidnight + (minutestoseconds(tashreeqWithTimes[i + 1][1]) * 1000);

        //pop up cell
        if (timeMil > tashreeqTimeOne && timeMil < tashreeqTimeTwo) {

            _debug("Tashreeq Poster","Next pop up: ",(tashreeqWithTimes[i + 1][0] + " - " + tashreeqWithTimes[i + 1][1]),logStyles.poster,false,false)

            //let show = (tashreeqWithTimes[i][0] != 'maghrib') ? (tashreeqTimeTwo) - 125000 - timeMil : (tashreeqTimeTwo) - 155000 - timeMil;
            let show = (tashreeqTimeTwo) - 120000 - timeMil;
            if(show < 0) show = 30000;
            console.log(show)
            //let hide = 180000;

            _debug("Tashreeq Poster","pop up for: ",(tashreeqWithTimes[i + 1][0] + " - Showing in: " + Math.round(show / 1000) + "s"),logStyles.poster,false,false)

            setTimeout(() => {
                //hidePoster("tashreeq","Tashreeq Poster Pop Up", true)
                _debug("Tashreeq Poster","pop up for: ","Pop Up Showing",logStyles.poster,false,false)
                showBigPoster("tashreeqPoster", 105000, "visible");
                setTimeout(()=>{tashreeqPopUp()},200000)
            }, show)
            return;
        }
    }

    var timeBefore = 125;
    var timeAfter = 280;

    var fajrtb = minutestoseconds(fajrJamaah) - timeBefore;
    var fajrta = minutestoseconds(fajrJamaah) + timeAfter;

    var dhuhrtb = minutestoseconds(dhuhrJamaah) - timeBefore;
    var dhuhrta = minutestoseconds(dhuhrJamaah) + timeAfter;

    var asrtb = minutestoseconds(asrJamaah) - timeBefore;
    var asrta = minutestoseconds(asrJamaah) + timeAfter;

    var magrbtb = minutestoseconds(maghribAthan) + 155;
    var magrbta = minutestoseconds(maghribAthan) + 560;

    var eshatb = minutestoseconds(eshaJamaah) - timeBefore;
    var eshata = minutestoseconds(eshaJamaah) + timeAfter;

    var timeNowSecs = timeInSeconds

    if (islamicMonthNum == 12) {
        if (islamicDateNum == 9 || islamicDateNum == 10 || islamicDateNum == 11 || islamicDateNum == 12 || islamicDateNum == 13) {
            if (timeNowSecs > fajrtb && fajrta > timeNowSecs) { showBigPoster("tashreeqPoster", 105000, "visible"); } else if (timeNowSecs > dhuhrtb && dhuhrta > timeNowSecs) { showBigPoster("tashreeqPoster", 105000, "visible"); } else if (timeNowSecs > asrtb && asrta > timeNowSecs) { showBigPoster("tashreeqPoster", 105000, "visible"); }
        }
        if (islamicDateNum == 10 || islamicDateNum == 11 || islamicDateNum == 12) {
            if (islamicMonthNum == 12 && timeNowSecs > magrbtb && magrbta > timeNowSecs) { showBigPoster("tashreeqPoster", 105000, "visible"); } else if (timeNowSecs > eshatb && eshata > timeNowSecs) { showBigPoster("tashreeqPoster", 105000, "visible"); }
        }
    } else { showBigPoster("tashreeqPoster", 0, "hidden"); }
}*/

// ~~~~~~ New poster pop up function ~~~~~ //
function fullDayPosters() {
    //Check if any objects are hiding poster
    console.groupCollapsed('%c----- Full Screen Poster Pop Up -----',logStyles.info)
    let activePosters     = [];
    let activeBigPosters  = [];
    let activeBPDurations = [];

    //Variables
    var timenow = timeInSeconds
    var totalBigPosterDuration = 0;

    //Big Posters 
    posters.bigPosterImage0.id = bigPosterImage0;
    posters.bigPosterImage1.id = bigPosterImage1;
    posters.bigPosterImage2.id = bigPosterImage2;
    posters.bigPosterImage3.id = bigPosterImage3;
    posters.bigPosterImage4.id = bigPosterImage4;
    posters.bigPosterImage5.id = bigPosterImage5;
    posters.bigPosterImage6.id = bigPosterImage6;
    posters.bigPosterImage7.id = bigPosterImage7;
    posters.bigPosterImage8.id = bigPosterImage8;
    posters.bigPosterImage9.id = bigPosterImage9;
    posters.cBigPosterImage0.id = cBigPosterImage0;
    posters.cBigPosterImage1.id = cBigPosterImage1;
    posters.cBigPosterImage2.id = cBigPosterImage2;
    posters.cBigPosterImage3.id = cBigPosterImage3;
    posters.cBigPosterImage4.id = cBigPosterImage4;
    posters.cBigPosterImage5.id = cBigPosterImage5;
    posters.cBigPosterImage6.id = cBigPosterImage6;
    posters.cBigPosterImage7.id = cBigPosterImage7;
    posters.cBigPosterImage8.id = cBigPosterImage8;
    posters.cBigPosterImage9.id = cBigPosterImage9;

    let tf = [bigPosterTrueFalse0, bigPosterTrueFalse1, bigPosterTrueFalse2, bigPosterTrueFalse3, bigPosterTrueFalse4, bigPosterTrueFalse5, bigPosterTrueFalse6, bigPosterTrueFalse7, bigPosterTrueFalse8, bigPosterTrueFalse9, cBigPosterTrueFalse0, cBigPosterTrueFalse1, cBigPosterTrueFalse2, cBigPosterTrueFalse3, cBigPosterTrueFalse4, cBigPosterTrueFalse5, cBigPosterTrueFalse6, cBigPosterTrueFalse7, cBigPosterTrueFalse8, cBigPosterTrueFalse9]
    let du = [bigPosterDuration0, bigPosterDuration1, bigPosterDuration2, bigPosterDuration3, bigPosterDuration4, bigPosterDuration5, bigPosterDuration6, bigPosterDuration7, bigPosterDuration8, bigPosterDuration9, cBigPosterDuration0, cBigPosterDuration1, cBigPosterDuration2, cBigPosterDuration3, cBigPosterDuration4, cBigPosterDuration5, cBigPosterDuration6, cBigPosterDuration7, cBigPosterDuration8, cBigPosterDuration9]

    for (i = 0; i < tf.length; i++) {
        if (tf[i] == "Display") {
            //_debug("Big Poster",'BigPoster' + i + ': ' + tf[i],"",logStyles.general,false,false)
            let currDuration = (isNaN(du[i])) ? (oldTotalPosterDuration/1000)/tf.length : (du[i] == 0) ? du[i] = 100 : Number(du[i]);
            totalBigPosterDuration = totalBigPosterDuration + currDuration
            activeBigPosters.push((i<=9) ? `bigPosterImage${i}` : `cBigPosterImage${i-10}`)
            activeBPDurations.push(currDuration)
        }
    }

    
    //AyaamBeed Poster Popup
    if (islamicDateNum == 12 && islamicMonthNum != 9 && islamicMonthNum != 12) {
        if (timenow > minutestoseconds(fajrStarts) && timenow < minutestoseconds(sunset)) {
            _debug("Big Poster","Ayyam Poster","1",logStyles.general,false,false)
            activePosters.push('ayaamPoster')
        }
    }
    if (islamicDateNum == 13 && islamicMonthNum != 9 && islamicMonthNum != 12) {
        if ((timenow > minutestoseconds(sunset) + 540) || (timenow < minutestoseconds(istiwa))) {
            _debug("Big Poster","Ayyam Poster","2",logStyles.general,false,false)
            activePosters.push('ayaamPoster')
        }
    }
    /*if ((islamicDateNum == 13 || islamicDateNum == 14) && islamicMonthNum != 9 && islamicMonthNum != 12) {
        _debug("Big Poster","Ayyam Poster","2",logStyles.general,false,false)
        activePosters.push('ayaamPoster')
    }
    if (islamicDateNum == 15 && islamicMonthNum != 9 && islamicMonthNum != 12) {
        if ((timenow > minutestoseconds(sunset) + 540) || (timenow < minutestoseconds(fajrStarts))) {
            _debug("Big Poster","Ayyam Poster","3",logStyles.general,false,false)
            activePosters.push('ayaamPoster')
        }
    }*/


    //Jumuah Posters (Durood & Sunnah)
    //Sunnah
    if (currentDay == 5 && timenow < minutestoseconds(jumuahTime3) + 0) {
        _debug("Big Poster","Jumuah Sunnah","",logStyles.general,false,false)
        activePosters.push('jumuahSunnah')
    }
    //Durood
    if (currentDay == 5 && timenow < minutestoseconds(sunset)) {
        _debug("Big Poster","Jumuah Durood","",logStyles.general,false,false)
            //push poster but calc random value 1st
    }

    //Thursday Durood Poster Popup
    if (currentDay == 4 && timenow > minutestoseconds(sunset) + 540) {
        _debug("Big Poster","Thursday Durood","",logStyles.general,false,false)
        activePosters.push('jumuahSunnah')
        activePosters.push('thursDuroodPoster')
    }

    //New Islamic Month Poster Popup (Dua for new month/year)
    if ((islamicDateNum == 1 && timeInSeconds > minutestoseconds(eshaJamaah) + 1800) || (islamicDateNum == 1 && timeInSeconds < minutestoseconds(maghribAthan))) {
        _debug("Big Poster","New Islamic Month","",logStyles.general,false,false)
        activePosters.push('islamicMonthPoster')
    }

    //Muharram Poster Popups
    if (islamicMonthNum == 1) {
        if (islamicDateNum == 2 || islamicDateNum == 5) {
            _debug("Big Poster","Muharram Poster","1",logStyles.general,false,false)
            activePosters.push('muharramPoster1')
        } else if (islamicDateNum == 8) {
            _debug("Big Poster","Muharram Poster","2",logStyles.general,false,false)
            activePosters.push('muharramPoster2')
        } else if (islamicDateNum == 9) {
            _debug("Big Poster","Muharram Poster","3",logStyles.general,false,false)
            activePosters.push('muharramPoster3')
        } else if (islamicDateNum == 10) {
            _debug("Big Poster","Muharram Poster","Ashura",logStyles.general,false,false)
            activePosters.push('ashuraSpendPoster')
        }
    }
        
    //Rajab Dua Poster Popup (Dua for rajab)
    if (islamicMonthNum == 7) {
        if ((islamicDateNum == 1 && timeInSeconds > minutestoseconds(eshaJamaah) + 1800) || (islamicDateNum == 1 && timeInSeconds < minutestoseconds(maghribAthan))) {
            _debug("Big Poster","Rajab Poster","",logStyles.general,false,false)
            activePosters.push('rajabPoster')
        } else if (islamicDateNum >= 2 && islamicDateNum <= 5) {
            _debug("Big Poster","Rajab Poster","",logStyles.general,false,false)
            activePosters.push('rajabPoster')
        }
    }

    //15 Shabaan Dua Poster Popup
    //if (islamicMonthNum == 8){
    //    if((islamicDateNum == 14 && timeInSeconds > timetoseconds(fajrStarts)) && (islamicDateNum == 15 && timeInSeconds < timetoseconds(fajrStarts))) {
    //        _debug("Big Poster","Shabaan Poster","",logStyles.general,false,false)
    //        activePosters.push('shabaanPoster1')
    //    }
    //}

    //Ramadhaan approach dua
    if (islamicMonthNum == 8){
        if(islamicDateNum == 18 || islamicDateNum == 25 || islamicDateNum == 27 || islamicDateNum == 29) {
            _debug("Big Poster","Ramadhaan Approach Poster","",logStyles.general,false,false)
            activePosters.push('ramadanApproachPoster')
        }
    }

    //Ramadan Duas Poster Popup
    if (islamicMonthNum == 9) {
        if(((islamicDateNum == 1 && timeInSeconds > minutestoseconds(eshaJamaah) + 1800) || (islamicDateNum == 1 && timeInSeconds < minutestoseconds(maghribAthan)))){
            _debug("Big Poster","Ramadan Poster","1",logStyles.general,false,false)
            activePosters.push('ramadanPoster1')
        }

        if (islamicDateNum >= 2 && islamicDateNum <= 10) {
            _debug("Big Poster","Ramadan Poster","1",logStyles.general,false,false)
            activePosters.push('ramadanPoster1')
        } else if (islamicDateNum > 10 && islamicDateNum <= 20) {
            _debug("Big Poster","Ramadan Poster","2",logStyles.general,false,false)
            activePosters.push('ramadanPoster2')
        } else if (islamicDateNum > 20 && islamicDateNum < 30) {
            _debug("Big Poster","Ramadan Poster","3",logStyles.general,false,false)
            activePosters.push('ramadanPoster3')
        }
        
        if (islamicDateNum => 20 && islamicDateNum <= 30){
            if(isEven(islamicDateNum) && timeInSeconds < minutestoseconds(sunset)){
                _debug("Big Poster","Laylatul Qadr Poster","",logStyles.general,false,false)
                activePosters.push('laylatulQadrPoster')
            } else if(isOdd(islamicDateNum) && timeInSeconds > minutestoseconds(sunset)){
                _debug("Big Poster","Laylatul Qadr Poster","",logStyles.general,false,false)
                activePosters.push('laylatulQadrPoster')
            } else if(isOdd(islamicDateNum) && timeInSeconds < minutestoseconds(fajrStarts)){
                _debug("Big Poster","Laylatul Qadr Poster","",logStyles.general,false,false)
                activePosters.push('laylatulQadrPoster')
            }
        }
    }

    //Eid Dua Poster Popup
    if ((islamicMonthNum == 10 && islamicDateNum == 1) || (islamicMonthNum == 12 && islamicDateNum == 10)) {
        _debug("Big Poster","Eid Poster","",logStyles.general,false,false)
        activePosters.push('eidDuaPoster')
    }

    //Dhul Hijjah Poster Popups
    if (islamicMonthNum == 12) {
        if (islamicDateNum < 10) {
            _debug("Big Poster","Hajj Poster",islamicDateNum,logStyles.general,false,false)
            activePosters.push("hajjPoster" + islamicDateNum)
        }
    }
    
    //Ayaam-e-Tashreeq
    /*if (islamicMonthNum != 12 || [9, 10, 11, 12, 13].includes(islamicDateNum) || mblPosterHide.includes("tashreeqPoster")){
        //Date from PC, this will automatically have the correct TimeZone.
        let dateTime = new Date();
        let timeMil = Date.now();
    
        //Send dateTime to midnight last & next
        let lastMidnight = dateTime.setHours(0, 0, 0, 0); // last midnight
        let nextMidnight = lastMidnight + 86400000; // next midnight
    
        //Arrays
        let tashreeqTimes = [
            ["midnightL", lastMidnight],
            ["fajr"     , fajrJamaah  ],
            ["dhuhr"    , dhuhrJamaah ],
            ["asr"      , asrJamaah   ],
            ["maghrib"  , maghribAthan],
            ["esha"     , eshaJamaah  ],
            ["midnightN", nextMidnight]
        ];
        let tashreeqWithTimes = []; //Array with populate only with valid cell times
    
        //Check for any tashreeq times that are not valid
        for (x = 0; x < tashreeqTimes.length; x++) {
            if (tashreeqTimes[x][1] != "~~~~" || tashreeqTimes[x][1] != "#REF" || tashreeqTimes[x][1] != null || tashreeqTimes[x][1] != 0 || tashreeqTimes[x][1] != "") {
                tashreeqWithTimes.push([tashreeqTimes[x][0], tashreeqTimes[x][1]]);
            }
        }
    
        //Loop through tashreeeq with times to check for next cell time
        for (i = 0; i < tashreeqWithTimes.length - 1; i++) {
            let tashreeqTimeOne = lastMidnight + (minutestoseconds(tashreeqWithTimes[i][1]) * 1000);
            let tashreeqTimeTwo = lastMidnight + (minutestoseconds(tashreeqWithTimes[i + 1][1]) * 1000);
            
            //pop up
            if(timeMil > tashreeqTimeOne && timeMil < tashreeqTimeTwo){
                if(timeMil > (tashreeqTimeTwo-900000) || timeMil > (tashreeqTimeTwo+900000)){
                    activePosters.push('tashreeqPoster')
                }
            }
        }
    }*/
    
    //Moon Visibility Pop Up
        let newMoonBirth = new Date(moonBirth);
            newMoonBirth.setFullYear(currentDate.getFullYear());
        let bestViewed   = new Date(bestVisibil);
            bestViewed.setFullYear(currentDate.getFullYear());
        let birthTime    = (newMoonBirth.getHours() * 3600) + (newMoonBirth.getMinutes() * 60) + newMoonBirth.getSeconds();
        let sunsetDate   = minutestoseconds(sunset)
        let asrTime      = minutestoseconds(asrJamaah)

        if(currentDate.getDate() == newMoonBirth.getDate() && currentDate.getDate() == bestViewed.getDate()){
            _debug("Big Poster","Moon Map","Curr Date = Birth Date && Curr Date = Best Viewed Date",logStyles.general,false,false)
            activePosters.push('moonVisibility')
        } else if(currentDate.getDate() == newMoonBirth.getDate()){
            if(birthTime > sunsetDate && timeInSeconds > asrTime){
                _debug("Big Poster","Moon Map","Birth Time > Sunset Time && Time Now > Asr Time",logStyles.general,false,false)
                activePosters.push('moonVisibility')
            } else if(birthTime < sunsetDate){
                _debug("Big Poster","Moon Map","Birth Time < Sunset Time",logStyles.general,false,false)
                activePosters.push('moonVisibility')
            }
        } else if((currentDate.getDate() == bestViewed.getDate()) && (currentDate.getMonth() == bestViewed.getMonth()) && (currentDate.getFullYear() == bestViewed.getFullYear())){
            _debug("Big Poster","Moon Map","Curr Date = Best Viewed Date",logStyles.general,false,false)
            activePosters.push('moonVisibility')
        }


    
    //Check for disabled posters
    //postersToHide = (mblPosterHide != "-" || mblPosterHide != null || configuration.length != 0) ? mblPosterHide.split(",") : []
    _debug("Big Poster","MBL Posters Disable",mblPosterHide,logStyles.function,false,false)
    for(ph=0; ph<mblPosterHide.length; ph++){
        let index = activePosters.indexOf(`${mblPosterHide[ph]}`);
        if (index > -1) { 
            activePosters.splice(index, 1); 
        }
    }


    // ====================================================== //
    // ================ Duration calculation ================ //
    // ====================================================== //
    let PPDuration = 0
    let BPqty = activeBigPosters.length
    let PPqty = activePosters.length
    let newTotalBigPosterDuration;
    let timeout;
    
    if(bigPosterSecondScreenActive){
        var elements = document.querySelectorAll('.fullScreenPosters');
        for(var i=0; i<elements.length; i++){
            elements[i].style.width  = "1910px";
            elements[i].style.height = "1070px";
        }
        PPDuration = activePosters.length*35;
        newTotalBigPosterDuration = activeBPDurations.reduce((a, b) => a + b, 0)
        timeout = Number(((newTotalBigPosterDuration+PPDuration)*1000));
    } else {
        //minimum for perpetual posters, without bigposters, is mblPosterTime
        if(BPqty == 0 && PPqty > 0){
            _debug("Big Poster","BigPoster Qty = 0 & Perpetual Posters = "+PPqty,"Perpetual Duration = mblPosterTime",logStyles.function,false,false)
            PPDuration = mblPosterTime;
        }
    
        //proportion bigposter if total is < 120 and we have no posters
        if(PPqty == 0 && BPqty > 0 && (totalBigPosterDuration < 120)){
            _debug("Big Poster","Perpetual = 0 & Bigposter = "+BPqty+" & totalBigPosterDuration < 120","Proportioning BIG POSTER",logStyles.function,false,false)
            for(i=0;i<activeBPDurations.length;i++){
                activeBPDurations[i] = (120/totalBigPosterDuration)*activeBPDurations[i];
            }
        }
    
        if(totalBigPosterDuration != 0 && totalBigPosterDuration < 120 && PPqty > 0){
            _debug("Big Poster","Bigposter < 120 && Perpetual > 0","Proportioning ALL posters",logStyles.function,false,false)
            BPPPDuration = (activePosters.length*35)+totalBigPosterDuration
            for(i=0; i<activeBPDurations.length; i++){
                activeBPDurations[i] = 120/BPPPDuration*activeBPDurations[i];
            }
            for(i=0; i<activePosters.length; i++){
                PPDuration = PPDuration + (mblPosterTime/BPPPDuration*35)
            }
        }
    
        if(totalBigPosterDuration > 120 && PPqty > 0){
            _debug("Big Poster","Big poster duration > 120","Perpetual posters set to 35s each",logStyles.function,false,false)
            PPDuration = activePosters.length*35;
        }
    
        newTotalBigPosterDuration = activeBPDurations.reduce((a, b) => a + b, 0)
        timeout = Number(((newTotalBigPosterDuration+PPDuration)*1000)+180000);
        
        /*
            - Perpetual without big posters 
                Use time from the sheet
            - Proportion bigposter if total is < 60 and we have no posters
                Use same function but change to 60
                (BP = totalBigPosterQty / 60) maybe?
            - 3rd if statement change 120 - 60 <=
                BP = totalBigPosterQty / 60
                PP = activePosters.length*30
            - 4th change 120 to 60 >
                PP = activePosters.length*30
        */

    }
        
    //Show posters in array
    console.group('Activated Posters')
    if(activePosters.length != 0 || activeBigPosters.length != 0){
        if(activeBigPosters.length != 0){
            _debug("Big Poster","Total Time",newTotalBigPosterDuration,logStyles.general,false,false)
            for(bp = 0; bp < activeBigPosters.length; bp++){
                let currentDuration = 0
                for(a=bp; a>=0; a--){
                    currentDuration = currentDuration + activeBPDurations[a]
                }
                _debug("Big Poster",activeBigPosters[bp],'Duration: ' + currentDuration,logStyles.general,false,false)
                showBigPoster(activeBigPosters[bp], currentDuration*1000, "visible");   
            }
        }
        if (activePosters.length != 0) {
            PPDuration = PPDuration/activePosters.length;
            console.log("PPDURATION: " + PPDuration)
            for (p = 0; p < activePosters.length; p++) {
                _debug("Big Poster",'PerpetualPoster: ' + activePosters[p],'Duration: ' + (PPDuration*(p+1)+totalBigPosterDuration),logStyles.general,false,false)
                showBigPoster(activePosters[p], (PPDuration*(p+1)+totalBigPosterDuration)*1000, "visible");
                console.log(activePosters[p])
                console.log((PPDuration*(p+1)+totalBigPosterDuration)*1000)
            }
        }
    }
    else {
        _debug("Big Poster","No Active Posters","Pop up not showing",logStyles.warn,false,false)
    }
    console.groupEnd('Activated Posters')

    if (debug) console.log("New big poster timeout = " + timeout)
    _debug("Big Poster","New big poster timeout = ",timeout,logStyles.function,false,false)
    
    console.groupEnd('%c----- Active Posters -----',logStyles.primary)
    
    clearInterval(fullDayPostersI);
    (timeout==0) ? timeout=180000 : timeout=timeout
    fullDayPostersI = setInterval(fullDayPosters, timeout);
    fullDayPostersI;

    if(debug) timeToShow(new Date(), timeout)
              //timeToHide(new Date(), (newTotalBigPosterDuration + PPDuration))
}

function timeToShow(timeStart,interval){
    let box = document.createElement("div");
    box.style.position        = "absolute"
    box.style.bottom          = 0
    box.style.left            = 0
    box.style.padding         = "5px"
    box.style.margin          = "auto"
    box.style.display         = "block"
    box.style.backgroundColor = "var(--background)";
    box.style.color           = "var(--tickerOdd)"
    box.innerHTML = "<p style='font-size:10px;'>0s</p>"
    document.body.appendChild(box);
    
    let timeLeftinterval = setInterval(()=>{
      let timeNow = new Date();
      let timeLeft = Math.abs(Math.round(((timeNow-timeStart)-interval)/1000))
      box.innerHTML = "<p style='font-size:10px;'>"+timeLeft + "s</p>"
      //console.log("Time remaining = " + timeLeft + "s")
      if(timeLeft <= 1 || timeLeft > (interval/1000)) {
        clearInterval(timeLeftinterval)
        box.remove()
      }
    },1000)
}
function timeToHide(startTime,totalTime){
  let timeLapsed = 0;
  new CircleProgress(`#posterStatus`, {
    max: 100,
    value: 100,
    //textFormat: "value",
    textFormat: ()=>{
      timeLapsed++
      let now     = new Date();
      let elapsed = now - startTime;
      let timer   = Math.max(totalTime - Math.floor(elapsed / 1000), 0);
      if(timer > 60){
        let date = new Date(null);
            date.setSeconds(timer); 
        return date.toISOString().substr(14, 5);
      } else{
         return timer + "s";
      }
    },
    animation: "linear",
    animationDuration: totalTime*1000
  })
}






// !====================================================== //
// !================== Global Functions ================== //
// !====================================================== //
// ~~~~~~~~~~~~~ Current Time ~~~~~~~~~~~~ //
function currentTime(t) {
    let currentDate = new Date();
    let timeC = document.getElementById("currentTime")
    h = checkTime(currentDate.getHours());
    m = checkTime(currentDate.getMinutes());
    s = checkTime(currentDate.getSeconds());

    let time = h + ":" + m + ":" + s;
    timeInSeconds = (currentDate.getHours() * 60 * 60) + (currentDate.getMinutes() * 60) + currentDate.getSeconds();

    let islamicTime = (timeInSeconds > minutestoseconds(sunset)) ? secondstotime(timeInSeconds - minutestoseconds(sunset)) : secondstotime(timeInSeconds + (86400 - minutestoseconds(sunset)))

    if (islamicLiveTime) {
        (salaah12HourFormat) ? timeC.innerHTML = convertArabicNumber(time24to12(islamicTime), "clock", ""): timeC.innerHTML = convertArabicNumber(islamicTime, "clock", "")
            //$("#tableTime").css({"top":"222px"})
    } else {
        (salaah12HourFormat) ? timeC.innerHTML = time24to12(time): timeC.innerHTML = time
            //$("#tableTime").css({"top":bigClockTopPad})
    }

    if (t = 's') {
        return (h * 60 * 60) + (m * 60) + currentDate.getSeconds()
    } else {
        return time;
    }
};

// ~~ Function to add leading 0 in time ~~ //
function checkTime(i) {
    if (i < 10) { i = "0" + i } return i;
}
// ~~ Function to convert 24hr to 12hr time ~~ //
function time24to12(time) {
    if(time == undefined) return time
    time = time.toString().match(/^([01]\d|2[0-3])(:)([0-5]\d)(:[0-5]\d)?$/) || [time];
    if (time.length > 1) {
        time = time.slice(1);
        time[0] = +time[0] % 12 || 12;
        //time[0] = checkTime(time[0]);
        //time[3] = checkTime(time[3]);
    } return time.join('');
}

//  Function to convert seconds to HH:MM:SS  //
function secondstotime(secs) {
    var t = new Date(1970, 0, 1);
    t.setSeconds(secs);
    var s = t.toTimeString().substr(0, 8);
    if (secs > 86399)
        s = Math.floor((t - Date.parse("1/1/70")) / 3600000) + s.substr(2);
    return s;
}

//  Function to convert HH:MM:SS to seconds  //
function timetoseconds(hhmmss) {
    if (typeof hhmmss === "string" && hhmmss.includes(":")) {
        var a = hhmmss.split(":");
        var timeinsecs = (+a[0]) * 60 * 60 + (+a[1]) * 60 + (+a[2]);
    } else {
        return hhmmss;
    }
}

// ~ Function to convert HH:MM to seconds  //
function minutestoseconds(mm) {
    if (typeof mm === "string" && mm.includes(":")) {
        var k = mm.split(":");
        return mininsecs = (+k[0]) * 60 * 60 + (+k[1]) * 60;
    } else {
        return 0;
    }
}

// ~ Function to check for error on times  //
function checkForError(time){
    if(time==undefined) return time;
    var bottomList  = ["00:00","#N/A", "#REF","#VALUE!"]
    var expStr      = bottomList.join("|");
    return time.replace(new RegExp('\(' + expStr + ')', 'gi'), '* * *').replace(/\s{2,}/g, '* * * *')
}

// ~~~ convert english to arabic number ~~ //
function convertArabicNumber(english, type, className) {
    var engString = english + '';
    var array = engString.split("");
    var arabic = []
    for (i = 0; i < array.length; i++) {
        if (array[i] == 1) arabic.push("١");
        else if (array[i] == 2) arabic.push("٢");
        else if (array[i] == 3) arabic.push("٣");
        else if (array[i] == 4) arabic.push("٤");
        else if (array[i] == 5) arabic.push("٥");
        else if (array[i] == 6) arabic.push("٦");
        else if (array[i] == 7) arabic.push("٧");
        else if (array[i] == 8) arabic.push("٨");
        else if (array[i] == 9) arabic.push("٩");
        else if (array[i] == 0) arabic.push("٠");
        else if (array[i] == ":" || array[i] == ".") {
            if (type == "ticker") arabic.push(`<span class="${className}" style='font-size: 25px;'>.</span>`);
            else if (type == "salaah") arabic.push(`<span class="${className}" style='font-size: 35px;'>.</span>`);
            else if (type == "clock") arabic.push(`<span class="${className}" style='font-size: 40px; text-shadow: 0px 0 0 #000;'>.</span>`);
            else arabic.push(":");
        }
    }
    return arabic.join('');
}

// ~~ Checking if file exists on server ~~ //
function checkFile(path) {
    return new Promise(resolve => {
        var bigPosterCheck = path.slice(18, -5)
        if(bigPosterCheck == "bigPosterImage"){
            resolve(true)
        }
        else{
            $.ajax({
                url: path,
                type: 'HEAD',
                error: () => {
                    resolve(false)
                },
                success: () => {
                    resolve(true)
                }
            })
        }
    })
}

// ~~ Converts srings to boolean values ~~ //
function stringToBoolean(string) {
    if(string == undefined) return false
    switch (string.toLowerCase().trim()) {
        case "true"     :
        case "Yes"      :
        case "yes"      :
        case "Hide"     :
        case "hide"     :
        case 1          :
        case "1"        :
            return true ;
        case "false"    :
        case "no"       :
        case "No"       :
        case "Display"  :
        case "display"  :
        case null       :
        case 0          :
        case "0"        :
            return false;
        default:
            return Boolean(string);
    }
}

function isEven(number) {
    return number % 2 === 0;
}

function isOdd(number) {
    return number % 2 !== 0;
}

function translateShortMonth(dateString) {
    let [day, shortMonth] = dateString.split(" ");
    let englishShortMonths = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];
    let monthIndex = englishShortMonths.indexOf(shortMonth);
    if (monthIndex === -1) {
        return dateString
    }
    let currLangMonthNames = JSON.parse(translations.monthShortname[lang])
    let translatedMonth = currLangMonthNames[monthIndex];

    return `${day} ${translatedMonth}`;
}

// ~~ Function for posting debug messages to the console ~~ //
// copy this for logs _debug(heading,message,variable,style,force,custom)
function _debug(heading,message,variable,style,force,custom){
    let time = `⏱ ${String(new Date().getHours()).padStart(2, '0')}:${String(new Date().getMinutes()).padStart(2, '0')}:${String(new Date().getSeconds()).padStart(2, '0')}`;

    if(heading == "startup"){
        if(debug){
            console.log("%c🚀 MasjidBoard Live Debug Mode 🚀", "font-size: 30px; font-weight: bold; color: #664d03; background-color: #fff3cd; padding: 15px; margin: 15px 0px; border-radius: 15px; border: thick solid #664d03")
            console.log("%cMBL Debug Colour Legend", "color: #FFF; background: #444; font-size: 1.2em; padding: 8px 16px; margin-top: 16px; border-radius: 12px 12px 0 0");
            console.log("%cTime Stamps%c - %cPoster Pop Ups%c - %cFunctions%c  - %cMBL API%c - %c General Info%c - %cErrors%c - %cWarnings", logStyles.time, "", logStyles.poster, "", logStyles.function, "", logStyles.api, "", logStyles.general, "", logStyles.error, "", logStyles.warn)
            console.log("%cSession Info", "color: #FFF; background: #444; font-size: 1.2em; padding: 8px 16px; margin-top: 16px; border-radius: 12px 12px 0 0");
            console.log(`%c| %cMasjid Name           |  %c${masjidName1 + " " + masjidName2}` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cSheet ID              |  %c${boardId}` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cSheet URL (NB)        |  %c${"https://docs.google.com/spreadsheets/d/"+boardId+"/"}` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cSheet Version         |  %c${mblVersion}` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cBoard Secret          |  %c${(mobile) ? urlParams.get('mid') : urlParams.get('id')}` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cLanuage               |  %c${translations.lang_name[lang] + " - " + lang + " - " + translations.lang_dir[lang] }` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cTheme                 |  %c${theme}` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cUpdate Interval       |  %c${sheetRefreshRate/1000}s` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cBig Poster Interval   |  %c${(parseInt(posterSetInterval, 10) + 1000)/1000/60}m` , "","font-weight: bold;",logStyles.info)
            console.log(`%c| %cBig Poster 2nd Screen |  %c${bigPosterSecondScreen}` , "","font-weight: bold;",((bigPosterSecondScreen) ? logStyles.green : logStyles.red))
            console.log(`%c| %cBig Posters Disabled  |  %c${posterDisabled}` , "","font-weight: bold;",((posterDisabled) ? logStyles.green : logStyles.red))
            console.log(`%c| %cMobile                |  %c${mobile}` , "","font-weight: bold;",((mobile) ? logStyles.green : logStyles.red))
            console.log(`%c| %cResolution            |  %c${window.screen.availWidth + " x " + window.screen.availHeight}` , "","font-weight: bold;",logStyles.info)
            console.log("%cLogs", "color: #FFF; background: #444; font-size: 1.2em; padding: 8px 16px; margin-top: 16px; border-radius: 12px 12px 0 0");
        } else{
            console.log("%c🕌 MasjidBoard Live 🕌", "font-size: 30px; font-weight: bold; color: #664d03; background-color: #fff3cd; padding: 15px; margin: 15px 0px; border-radius: 15px; border: thick solid #664d03")
        }
    }

    if (custom){
        console.log(message)
    } else if(custom == "group"){
        console.groupCollapsed(`%c${heading}`,`${style}`)
        console.log(`%c${time}%c - %c${heading}%c - ${message}: %c${variable}`,logStyles.time,"",`${style}`,"",`${((variable==true) ? logStyles.green : (variable==false) ? logStyles.red : logStyles.info)}`)
    } else if(debug || force) {
        console.log(`%c${time}%c - %c${heading}%c - ${message}: %c${variable}`,logStyles.time,"",`${style}`,"",`${((variable==true) ? logStyles.green : (variable==false) ? logStyles.red : logStyles.info)}`)
    } else {
        return
    }
}