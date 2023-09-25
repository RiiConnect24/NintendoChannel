package constants

// RatingGroup is the rating organization for a specific region.
type RatingGroup uint8

const (
	_ RatingGroup = iota
	// CERO is the RatingGroup for NTSC-J games.
	CERO

	// ESRB is the RatingGroup for NTSC-U games.
	ESRB

	// PEGI is the RatingGroup for PAL games.
	PEGI RatingGroup = 4
)

// RatingData contains the name and age for a rating
type RatingData struct {
	Name [11]uint16
	Age  uint8
}

var RatingsData = map[RatingGroup][]RatingData{
	CERO: {
		{Name: [11]uint16{'A'}, Age: 0},
		{Name: [11]uint16{'B'}, Age: 12},
		{Name: [11]uint16{'C'}, Age: 15},
		{Name: [11]uint16{'D'}, Age: 17},
		{Name: [11]uint16{'Z'}, Age: 18},
	},
	ESRB: {
		{Name: [11]uint16{'E', 'C'}, Age: 3},
		{Name: [11]uint16{'E'}, Age: 6},
		{Name: [11]uint16{'E', '1', '0'}, Age: 10},
		{Name: [11]uint16{'T'}, Age: 13},
		{Name: [11]uint16{'M'}, Age: 17},
	},
	PEGI: {
		{Name: [11]uint16{'3'}, Age: 3},
		{Name: [11]uint16{'7'}, Age: 7},
		{Name: [11]uint16{'1', '2'}, Age: 12},
		{Name: [11]uint16{'1', '6'}, Age: 16},
		{Name: [11]uint16{'1', '8'}, Age: 18},
	},
}

// Region is the Wii's region flags found in TMDs.
type Region int

const (
	Japan Region = iota
	PAL
	NTSC
)

type Language int

const (
	Japanese = iota
	English
	German
	French
	Spanish
	Italian
	Dutch
)

type RegionMeta struct {
	Region      Region
	Languages   []Language
	RatingGroup RatingGroup
}

var Regions = []RegionMeta{
	{
		Region:      Japan,
		Languages:   []Language{Japanese},
		RatingGroup: CERO,
	},
	{
		Region:      NTSC,
		Languages:   []Language{English, French, Spanish},
		RatingGroup: ESRB,
	},
	{
		Region:      PAL,
		Languages:   []Language{English, German, French, Spanish, Italian, Dutch},
		RatingGroup: PEGI,
	},
}

// ConsoleModels is the type of consoles the Nintendo Channel games has.
type ConsoleModels [3]byte

var (
	// RVL represents titles on the Wii.
	RVL ConsoleModels = [3]byte{'R', 'V', 'L'}

	// NTR represents titles on the DS and DS Lite.
	NTR ConsoleModels = [3]byte{'N', 'T', 'R'}

	// TWL represents titles on the DSi.
	TWL ConsoleModels = [3]byte{'T', 'W', 'L'}

	// CTR represents titles on the 3DS.
	CTR ConsoleModels = [3]byte{'C', 'T', 'R'}
)

// TitleGroupTypes represents a type of title a console's game is.
type TitleGroupTypes uint8

const (
	_ TitleGroupTypes = iota

	// Disc represents Wii disc games.
	Disc

	// WiiWare represents WiiWare games.
	WiiWare

	// WiiChannels represents Wii Channels such as Forecast and News.
	WiiChannels

	// DS represents DS games.
	DS

	// VirtualConsole represents Virtual Console games.
	VirtualConsole

	// DSi represents games that support DSi only features
	DSi

	// DSiWare represents DSiWare games
	DSiWare

	// ThreeDS represents 3DS games
	ThreeDS

	// ThreeDSDownloadSoftware represents 3DS Download Software
	ThreeDSDownloadSoftware

	// ThreeDSGameBoy represents GameBoy Virtual Console games on the 3DS.
	ThreeDSGameBoy

	// WiiU represents Wii U disc games.
	WiiU

	// Switch represents Switch games on cartridge.
	Switch
)

// TitleTypeData contains the metadata for a title type
type TitleTypeData struct {
	ConsoleModel ConsoleModels
	GroupID      TitleGroupTypes
	TypeID       uint8
	ConsoleName  string
}

var TitleTypesData = []TitleTypeData{
	{TypeID: 0, ConsoleModel: RVL, ConsoleName: "Wii", GroupID: Disc},
	{TypeID: 1, ConsoleModel: RVL, ConsoleName: "WiiWare", GroupID: WiiWare},
	{TypeID: 2, ConsoleModel: RVL, ConsoleName: "Wii Channels", GroupID: WiiChannels},
	{TypeID: 3, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console NES", GroupID: VirtualConsole},
	{TypeID: 4, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console Super NES", GroupID: VirtualConsole},
	{TypeID: 5, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console Nintendo 64", GroupID: VirtualConsole},
	{TypeID: 6, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console TurboGrafx16", GroupID: VirtualConsole},
	{TypeID: 7, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console Sega Genesis", GroupID: VirtualConsole},
	{TypeID: 8, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console Neo Geo", GroupID: VirtualConsole},
	{TypeID: 9, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console Master System", GroupID: VirtualConsole},
	{TypeID: 10, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console Commodore 64", GroupID: VirtualConsole},
	{TypeID: 11, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console Arcade", GroupID: VirtualConsole},
	{TypeID: 12, ConsoleModel: RVL, ConsoleName: "Wii Virtual Console MSX", GroupID: VirtualConsole},
	{TypeID: 13, ConsoleModel: NTR, ConsoleName: "Nintendo DS", GroupID: DS},
	{TypeID: 14, ConsoleModel: TWL, ConsoleName: "Nintendo DS", GroupID: DS},
	{TypeID: 15, ConsoleModel: TWL, ConsoleName: "Nintendo DSi", GroupID: DSi},
	{TypeID: 16, ConsoleModel: TWL, ConsoleName: "Nintendo DSiWare", GroupID: DSiWare},
	{TypeID: 17, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS", GroupID: ThreeDS},
	{TypeID: 18, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS Download Software", GroupID: ThreeDSDownloadSoftware},
	{TypeID: 19, ConsoleModel: CTR, ConsoleName: "New Nintendo 3DS", GroupID: ThreeDS},
	{TypeID: 20, ConsoleModel: CTR, ConsoleName: "New Nintendo 3DS Download Software", GroupID: ThreeDSDownloadSoftware},
	{TypeID: 21, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS Virtual Console NES", GroupID: ThreeDSGameBoy},
	{TypeID: 22, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS Virtual Console Game Boy", GroupID: ThreeDSGameBoy},
	{TypeID: 23, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS Virtual Console Game Boy Color", GroupID: ThreeDSGameBoy},
	{TypeID: 24, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS Virtual Console Game Boy Advance", GroupID: ThreeDSGameBoy},
	{TypeID: 25, ConsoleModel: CTR, ConsoleName: "Nintendo 3DS Virtual Console Game Gear", GroupID: ThreeDSGameBoy},
	{TypeID: 26, ConsoleModel: RVL, ConsoleName: "Wii U", GroupID: WiiU},
	{TypeID: 27, ConsoleModel: RVL, ConsoleName: "Wii U Download Software", GroupID: WiiU},
	{TypeID: 28, ConsoleModel: RVL, ConsoleName: "Wii U Virtual Console NES", GroupID: WiiU},
	{TypeID: 29, ConsoleModel: RVL, ConsoleName: "Wii U Virtual Console Super NES", GroupID: WiiU},
	{TypeID: 30, ConsoleModel: RVL, ConsoleName: "Wii U Virtual Console Nintendo 64", GroupID: WiiU},
	{TypeID: 31, ConsoleModel: RVL, ConsoleName: "Wii U Virtual Console Game Boy Advance", GroupID: WiiU},
	{TypeID: 32, ConsoleModel: RVL, ConsoleName: "Wii U Virtual Console Nintendo DS", GroupID: WiiU},
	{TypeID: 33, ConsoleModel: RVL, ConsoleName: "Wii U Virtual Console TurboGrafx16", GroupID: WiiU},
	{TypeID: 34, ConsoleModel: RVL, ConsoleName: "Wii U Virtual Console MSX", GroupID: WiiU},
	{TypeID: 35, ConsoleModel: RVL, ConsoleName: "Wii U Applications", GroupID: WiiU},
	{TypeID: 36, ConsoleModel: NTR, ConsoleName: "Nintendo Switch", GroupID: Switch},
	/*{TypeID: 37, ConsoleModel: NTR, ConsoleName: "Nintendo Switch Download Software", GroupID: Switch},*/
}

// TitleType is the classified type of title according to GameTDB
type TitleType uint8

const (
	_ TitleType = iota
	Wii
	WiiChannel
	NES
	SNES
	Nintendo64
	TurboGrafx16
	Genesis
	NeoGeo
	NintendoDS           TitleType = 13
	_WiiWare             TitleType = 1
	MasterSystem         TitleType = 9
	Commodore64          TitleType = 10
	VirtualConsoleArcade TitleType = 11
	NintendoDSi          TitleType = 15
	NintendoDSiWare      TitleType = 16
	NintendoThreeDS      TitleType = 17
	ThreeDSDownload      TitleType = 18
	NewThreeDS				TitleType = 19
	NewThreeDSDownload		TitleType = 20
	ThreeDSNES			TitleType = 21
	ThreeDSGameBoyColor	TitleType = 23
	ThreeDSGameBoyAdvance	TitleType = 24
	ThreeDSGameGead			TitleType = 25
	WiiUDisc				TitleType = 26
	WiiUDownload		  TitleType = 27
	WiiUNES					TitleType = 28
	WiiUSNES				TitleType = 29
	WiiUN64					TitleType = 30
	WiiUGBA					TitleType = 31
	WiiUDS					TitleType = 32
	WiiUTurboGrafx16		TitleType = 33
	WiiUMSX					TitleType = 34
	WiiUApplications		TitleType = 35
	SwitchPhysical			TitleType = 36
	/*SwitchDownload			TitleType = 37*/
)

var TitleTypeMap = map[string]TitleType{
	"Wii":       Wii,
	"Channel":   WiiChannel,
	"WiiWare":   _WiiWare,
	"VC-NES":    NES,
	"VC-SNES":   SNES,
	"VC-N64":    Nintendo64,
	"VC-SMS":    MasterSystem,
	"VC-MD":     Genesis,
	"VC-PCE":    TurboGrafx16,
	"VC-NEOGEO": NeoGeo,
	"VC-Arcade": VirtualConsoleArcade,
	"VC-C64":    Commodore64,
	"DS":        NintendoDS,
	"DSi":       NintendoDSi,
	"DSiWare":   NintendoDSiWare,
	"3DS":       NintendoThreeDS,
	"3DSWare":   ThreeDSDownload,
	"VC-GB":     ThreeDSDownload,
	"VC-GBC":    ThreeDSDownload,
	"VC-GBA":    ThreeDSDownload,
	"VC-GG":     ThreeDSDownload,
	"WiiU":			WiiUDisc,
	"eShop":		WiiUDownload,
	"Switch":		SwitchPhysical,
}

type Medal uint8

const (
	None Medal = iota
	Bronze
	Silver
	Gold
	Platinum
)

func GetPopularVideoQueryString(language Language) string {
	switch language {
	case Japanese:
		return `SELECT id, name_japanese, length, video_type FROM videos ORDER BY RAND() DESC`
	case English:
		return `SELECT id, name_english, length, video_type FROM videos ORDER BY RAND() DESC`
	case German:
		return `SELECT id, name_german, length, video_type FROM videos ORDER BY RAND() DESC`
	case French:
		return `SELECT id, name_french, length, video_type FROM videos ORDER BY RAND() DESC`
	case Spanish:
		return `SELECT id, name_spanish, length, video_type FROM videos ORDER BY RAND() DESC`
	case Italian:
		return `SELECT id, name_italian, length, video_type FROM videos ORDER BY RAND() DESC`
	case Dutch:
		return `SELECT id, name_dutch, length, video_type FROM videos ORDER BY RAND() DESC`
	default:
		// Will never reach here
		return ""
	}
}


func GetVideoQueryString(language Language) string {
	switch language {
	case Japanese:
		return `SELECT id, name_japanese, length, video_type FROM videos ORDER BY date_added DESC`
	case English:
		return `SELECT id, name_english, length, video_type FROM videos ORDER BY date_added DESC`
	case German:
		return `SELECT id, name_german, length, video_type FROM videos ORDER BY date_added DESC`
	case French:
		return `SELECT id, name_french, length, video_type FROM videos ORDER BY date_added DESC`
	case Spanish:
		return `SELECT id, name_spanish, length, video_type FROM videos ORDER BY date_added DESC`
	case Italian:
		return `SELECT id, name_italian, length, video_type FROM videos ORDER BY date_added DESC`
	case Dutch:
		return `SELECT id, name_dutch, length, video_type FROM videos ORDER BY date_added DESC`
	default:
		// Will never reach here
		return ""
	}
}
