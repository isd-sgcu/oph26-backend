package model

// FACULTY ENUM ================================

type Faculty string

const (
	EDU      Faculty = "edu"
	PSY      Faculty = "psy"
	DENT     Faculty = "dent"
	LAW      Faculty = "law"
	COMMARTS Faculty = "commarts"
	CBS      Faculty = "cbs"
	MD       Faculty = "md"
	PHARM    Faculty = "pharm"
	POLSCI   Faculty = "polsci"
	SCI      Faculty = "sci"
	SPSC     Faculty = "spsc"
	ENG      Faculty = "eng"
	FAA      Faculty = "faa"
	ECON     Faculty = "econ"
	ARCH     Faculty = "arch"
	AHS      Faculty = "ahs"
	VET      Faculty = "vet"
	ARTS     Faculty = "arts"
	SCII     Faculty = "scii"
	CUSAR    Faculty = "cusar"
)

var FacultySet = map[Faculty]struct{}{
	EDU:      {},
	PSY:      {},
	DENT:     {},
	LAW:      {},
	COMMARTS: {},
	CBS:      {},
	MD:       {},
	PHARM:    {},
	POLSCI:   {},
	SCI:      {},
	SPSC:     {},
	ENG:      {},
	FAA:      {},
	ECON:     {},
	ARCH:     {},
	AHS:      {},
	VET:      {},
	ARTS:     {},
	SCII:     {},
	CUSAR:    {},
}

func FacultyIsValid(input string) bool {
	_, ok := FacultySet[Faculty(input)]
	return ok
}

// Validate list of faculties
func FacultiesAreValid(input []string) bool {
	valid := true
	for _, i := range input {
		_, ok := FacultySet[Faculty(i)]
		if !ok {
			valid = false
		}
	}
	return valid
}

// PROVINCE ENUM ================================

type Province string

const (
	Krabi                 Province = "กระบี่"
	Bangkok               Province = "กรุงเทพมหานคร"
	Kanchanaburi          Province = "กาญจนบุรี"
	Kalasin               Province = "กาฬสินธุ์"
	KamphaengPhet         Province = "กำแพงเพชร"
	KhonKaen              Province = "ขอนแก่น"
	Chanthaburi           Province = "จันทบุรี"
	Chachoengsao          Province = "ฉะเชิงเทรา"
	Chonburi              Province = "ชลบุรี"
	ChaiNat               Province = "ชัยนาท"
	Chaiyaphum            Province = "ชัยภูมิ"
	Chumphon              Province = "ชุมพร"
	ChiangRai             Province = "เชียงราย"
	ChiangMai             Province = "เชียงใหม่"
	Trang                 Province = "ตรัง"
	Trat                  Province = "ตราด"
	Tak                   Province = "ตาก"
	NakhonNayok           Province = "นครนายก"
	NakhonPathom          Province = "นครปฐม"
	NakhonPhanom          Province = "นครพนม"
	NakhonRatchasima      Province = "นครราชสีมา"
	NakhonSiThammarat     Province = "นครศรีธรรมราช"
	NakhonSawan           Province = "นครสวรรค์"
	Nonthaburi            Province = "นนทบุรี"
	Narathiwat            Province = "นราธิวาส"
	Nan                   Province = "น่าน"
	BuengKan              Province = "บึงกาฬ"
	Buriram               Province = "บุรีรัมย์"
	PathumThani           Province = "ปทุมธานี"
	PrachuapKhiriKhan     Province = "ประจวบคีรีขันธ์"
	Prachinburi           Province = "ปราจีนบุรี"
	Pattani               Province = "ปัตตานี"
	PhraNakhonSiAyutthaya Province = "พระนครศรีอยุธยา"
	Phayao                Province = "พะเยา"
	PhangNga              Province = "พังงา"
	Phatthalung           Province = "พัทลุง"
	Phichit               Province = "พิจิตร"
	Phitsanulok           Province = "พิษณุโลก"
	Phetchaburi           Province = "เพชรบุรี"
	Phetchabun            Province = "เพชรบูรณ์"
	Phrae                 Province = "แพร่"
	Phuket                Province = "ภูเก็ต"
	MahaSarakham          Province = "มหาสารคาม"
	Mukdahan              Province = "มุกดาหาร"
	MaeHongSon            Province = "แม่ฮ่องสอน"
	Yasothon              Province = "ยโสธร"
	Yala                  Province = "ยะลา"
	RoiEt                 Province = "ร้อยเอ็ด"
	Ranong                Province = "ระนอง"
	Rayong                Province = "ระยอง"
	Ratchaburi            Province = "ราชบุรี"
	Lopburi               Province = "ลพบุรี"
	Lampang               Province = "ลำปาง"
	Lamphun               Province = "ลำพูน"
	Loei                  Province = "เลย"
	Sisaket               Province = "ศรีสะเกษ"
	SakonNakhon           Province = "สกลนคร"
	Songkhla              Province = "สงขลา"
	Satun                 Province = "สตูล"
	SamutPrakan           Province = "สมุทรปราการ"
	SamutSongkhram        Province = "สมุทรสงคราม"
	SamutSakhon           Province = "สมุทรสาคร"
	SaKaeo                Province = "สระแก้ว"
	Saraburi              Province = "สระบุรี"
	SingBuri              Province = "สิงห์บุรี"
	Sukhothai             Province = "สุโขทัย"
	SuphanBuri            Province = "สุพรรณบุรี"
	SuratThani            Province = "สุราษฎร์ธานี"
	Surin                 Province = "สุรินทร์"
	NongKhai              Province = "หนองคาย"
	NongBuaLamphu         Province = "หนองบัวลำภู"
	AngThong              Province = "อ่างทอง"
	AmnatCharoen          Province = "อำนาจเจริญ"
	UdonThani             Province = "อุดรธานี"
	Uttaradit             Province = "อุตรดิตถ์"
	UthaiThani            Province = "อุทัยธานี"
	UbonRatchathani       Province = "อุบลราชธานี"
)

var ProvinceSet = map[Province]struct{}{
	Krabi:                 {},
	Bangkok:               {},
	Kanchanaburi:          {},
	Kalasin:               {},
	KamphaengPhet:         {},
	KhonKaen:              {},
	Chanthaburi:           {},
	Chachoengsao:          {},
	Chonburi:              {},
	ChaiNat:               {},
	Chaiyaphum:            {},
	Chumphon:              {},
	ChiangRai:             {},
	ChiangMai:             {},
	Trang:                 {},
	Trat:                  {},
	Tak:                   {},
	NakhonNayok:           {},
	NakhonPathom:          {},
	NakhonPhanom:          {},
	NakhonRatchasima:      {},
	NakhonSiThammarat:     {},
	NakhonSawan:           {},
	Nonthaburi:            {},
	Narathiwat:            {},
	Nan:                   {},
	BuengKan:              {},
	Buriram:               {},
	PathumThani:           {},
	PrachuapKhiriKhan:     {},
	Prachinburi:           {},
	Pattani:               {},
	PhraNakhonSiAyutthaya: {},
	Phayao:                {},
	PhangNga:              {},
	Phatthalung:           {},
	Phichit:               {},
	Phitsanulok:           {},
	Phetchaburi:           {},
	Phetchabun:            {},
	Phrae:                 {},
	Phuket:                {},
	MahaSarakham:          {},
	Mukdahan:              {},
	MaeHongSon:            {},
	Yasothon:              {},
	Yala:                  {},
	RoiEt:                 {},
	Ranong:                {},
	Rayong:                {},
	Ratchaburi:            {},
	Lopburi:               {},
	Lampang:               {},
	Lamphun:               {},
	Loei:                  {},
	Sisaket:               {},
	SakonNakhon:           {},
	Songkhla:              {},
	Satun:                 {},
	SamutPrakan:           {},
	SamutSongkhram:        {},
	SamutSakhon:           {},
	SaKaeo:                {},
	Saraburi:              {},
	SingBuri:              {},
	Sukhothai:             {},
	SuphanBuri:            {},
	SuratThani:            {},
	Surin:                 {},
	NongKhai:              {},
	NongBuaLamphu:         {},
	AngThong:              {},
	AmnatCharoen:          {},
	UdonThani:             {},
	Uttaradit:             {},
	UthaiThani:            {},
	UbonRatchathani:       {},
}

func ProvinceIsValid(input string) bool {
	_, ok := ProvinceSet[Province(input)]
	return ok
}

// OBJECTIVE ENUM ================================

type Objective string

const (
	Explorechula   Objective = "เยี่ยมชมบรรยากาศของมหาวิทยาลัย"
	OtherObjective Objective = "อื่น ๆ"
)

var ObjectiveSet = map[Objective]struct{}{
	Explorechula:   {},
	OtherObjective: {},
}

func ObjectiveIsValid(input string) bool {
	_, ok := ObjectiveSet[Objective(input)]
	return ok
}

// Validate list of objectives
func ObjectivesAreValid(input []string) bool {
	valid := true
	for _, i := range input {
		_, ok := ObjectiveSet[Objective(i)]
		if !ok {
			valid = false
		}
	}
	return valid
}

// NEWS SOURCES ENUM ================================

type NewsSource string

const (
	OtherNewsSource NewsSource = "อื่น ๆ"
)

var NewsSourceSet = map[NewsSource]struct{}{
	OtherNewsSource: {},
}

func NewsSourceIsValid(input string) bool {
	_, ok := NewsSourceSet[NewsSource(input)]
	return ok
}

// Validate list of news sources
func NewsSourcesAreValid(input []string) bool {
	valid := true
	for _, i := range input {
		_, ok := NewsSourceSet[NewsSource(i)]
		if !ok {
			valid = false
		}
	}
	return valid
}

// STUDY LEVEL ENUM ================================

type StudyLevel string

const (
	Elementary        StudyLevel = "ประถมศึกษา"
	Matthayom1        StudyLevel = "มัธยมศึกษาปีที่ 1"
	Matthayom2        StudyLevel = "มัธยมศึกษาปีที่ 2"
	Matthayom3        StudyLevel = "มัธยมศึกษาปีที่ 3"
	Matthayom4        StudyLevel = "มัธยมศึกษาปีที่ 4"
	Matthayom5        StudyLevel = "มัธยมศึกษาปีที่ 5"
	Matthayom6        StudyLevel = "มัธยมศึกษาปีที่ 6"
	Vocational1       StudyLevel = "ปวช. 1"
	Vocational2       StudyLevel = "ปวช. 2"
	Vocational3       StudyLevel = "ปวช. 3"
	HigherVocational1 StudyLevel = "ปวส. 1"
	HigherVocational2 StudyLevel = "ปวส. 2"
	HigherEducation   StudyLevel = "อุดมศึกษา"
	OtherEducation    StudyLevel = "อื่น ๆ"
)

var StudyLevelSet = map[StudyLevel]struct{}{
	Elementary:        {},
	Matthayom1:        {},
	Matthayom2:        {},
	Matthayom3:        {},
	Matthayom4:        {},
	Matthayom5:        {},
	Matthayom6:        {},
	Vocational1:       {},
	Vocational2:       {},
	Vocational3:       {},
	HigherVocational1: {},
	HigherVocational2: {},
	HigherEducation:   {},
	OtherEducation:    {},
}

func StudyLevelIsValid(input string) bool {
	_, ok := StudyLevelSet[StudyLevel(input)]
	return ok
}
