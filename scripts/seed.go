package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"agnos_candidate_assignment/config"
	"agnos_candidate_assignment/database"
	"agnos_candidate_assignment/models"

	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func strPtr(s string) *string { return &s }

func cleanStr(s string) string {
	if strings.IndexByte(s, 0) != -1 {
		log.Printf("sanitizing NUL bytes in string: %q", s)
	}
	return strings.ReplaceAll(s, "\x00", "")
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func genderFromInt(i int) models.Gender {
	if i == 0 {
		return models.Male
	}
	return models.Female
}

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := database.NewPostgresConnectionNoMigrate(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	if err := db.AutoMigrate(&models.Hospital{}, &models.Staff{}, &models.Patient{}); err != nil {
		log.Fatalf("failed to auto-migrate schema: %v", err)
	}

	fmt.Println("db migrated: hospitals, staffs, patients")

	hospitals := []string{"Central Hospital", "Green Valley Hospital", "Sunrise Medical"}

	staffList := []struct {
		Username string
		Hospital string
		Password string
	}{
		{"alice", "Central Hospital", "Alice!23"},
		{"bob", "Central Hospital", "Bob!23"},
		{"carol", "Central Hospital", "Carol!23"},
		{"david", "Green Valley Hospital", "David!23"},
		{"eva", "Green Valley Hospital", "Eva!23"},
		{"frank", "Green Valley Hospital", "Frank!23"},
		{"grace", "Sunrise Medical", "Grace!23"},
		{"henry", "Sunrise Medical", "Henry!23"},
		{"irene", "Sunrise Medical", "Irene!23"},
		{"jack", "Sunrise Medical", "Jack!23"},
		{"kate", "Central Hospital", "Kate!23"},
		{"luke", "Green Valley Hospital", "Luke!23"},
	}

	hospMap := make(map[string]uint, len(hospitals))
	for _, rawName := range hospitals {
		name := cleanStr(rawName)
		var h models.Hospital
		if err := db.Where("name = ?", name).First(&h).Error; err == nil {
			hospMap[name] = h.ID
			continue
		}
		h = models.Hospital{Name: name}
		if err := db.Create(&h).Error; err != nil {
			log.Fatalf("failed to create hospital %s: %v", name, err)
		}
		hospMap[name] = h.ID
	}

	staffCreated := 0
	for _, s := range staffList {
		uname := cleanStr(s.Username)
		hospName := cleanStr(s.Hospital)
		hid := hospMap[hospName]
		hash, err := bcrypt.GenerateFromPassword([]byte(s.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("skip staff %s: bcrypt error: %v\n", s.Username, err)
			continue
		}
		st := models.Staff{
			UserName:     uname,
			PasswordHash: string(hash),
			HospitalID:   hid,
		}
		if err := db.Create(&st).Error; err != nil {
			fmt.Printf("skip staff %s: %v\n", s.Username, err)
			continue
		}
		staffCreated++
	}

	patientsData := []struct {
		Hospital  string
		FirstTH   string
		MiddleTH  string
		LastTH    string
		FirstEN   string
		MiddleEN  string
		LastEN    string
		HN        string
		Nat       string
		PP        string
		Phone     string
		Email     string
		Year      int
		Month     int
		Day       int
		GenderInt int
	}{
		{"Central Hospital", "First1", "M1", "Last1", "F1", "ME1", "L1", "HN00001", "NAT00000001", "PP000001", "080000001", "patient1@example.com", 1971, 2, 2, 0},
		{"Green Valley Hospital", "First2", "M2", "Last2", "F2", "ME2", "L2", "HN00002", "NAT00000002", "PP000002", "080000002", "patient2@example.com", 1972, 3, 3, 1},
		{"Sunrise Medical", "First3", "M3", "Last3", "F3", "ME3", "L3", "HN00003", "NAT00000003", "PP000003", "080000003", "patient3@example.com", 1973, 4, 4, 0},
		{"Central Hospital", "First4", "M4", "Last4", "F4", "ME4", "L4", "HN00004", "NAT00000004", "PP000004", "080000004", "patient4@example.com", 1974, 5, 5, 1},
		{"Green Valley Hospital", "First5", "M5", "Last5", "F5", "ME5", "L5", "HN00005", "NAT00000005", "PP000005", "080000005", "patient5@example.com", 1975, 6, 6, 0},
		{"Sunrise Medical", "First6", "M6", "Last6", "F6", "ME6", "L6", "HN00006", "NAT00000006", "PP000006", "080000006", "patient6@example.com", 1976, 7, 7, 1},
		{"Central Hospital", "First7", "M7", "Last7", "F7", "ME7", "L7", "HN00007", "NAT00000007", "PP000007", "080000007", "patient7@example.com", 1977, 8, 8, 0},
		{"Green Valley Hospital", "First8", "M8", "Last8", "F8", "ME8", "L8", "HN00008", "NAT00000008", "PP000008", "080000008", "patient8@example.com", 1978, 9, 9, 1},
		{"Sunrise Medical", "First9", "M9", "Last9", "F9", "ME9", "L9", "HN00009", "NAT00000009", "PP000009", "080000009", "patient9@example.com", 1979, 10, 10, 0},
		{"Central Hospital", "First10", "M10", "Last10", "F10", "ME10", "L10", "HN00010", "NAT00000010", "PP000010", "080000010", "patient10@example.com", 1980, 11, 11, 1},
		{"Green Valley Hospital", "First11", "M11", "Last11", "F11", "ME11", "L11", "HN00011", "NAT00000011", "PP000011", "080000011", "patient11@example.com", 1981, 12, 12, 0},
		{"Sunrise Medical", "First12", "M12", "Last12", "F12", "ME12", "L12", "HN00012", "NAT00000012", "PP000012", "080000012", "patient12@example.com", 1982, 1, 13, 1},
		{"Central Hospital", "First13", "M13", "Last13", "F13", "ME13", "L13", "HN00013", "NAT00000013", "PP000013", "080000013", "patient13@example.com", 1983, 2, 14, 0},
		{"Green Valley Hospital", "First14", "M14", "Last14", "F14", "ME14", "L14", "HN00014", "NAT00000014", "PP000014", "080000014", "patient14@example.com", 1984, 3, 15, 1},
		{"Sunrise Medical", "First15", "M15", "Last15", "F15", "ME15", "L15", "HN00015", "NAT00000015", "PP000015", "080000015", "patient15@example.com", 1985, 4, 16, 0},
		{"Central Hospital", "First16", "M16", "Last16", "F16", "ME16", "L16", "HN00016", "NAT00000016", "PP000016", "080000016", "patient16@example.com", 1986, 5, 17, 1},
		{"Green Valley Hospital", "First17", "M17", "Last17", "F17", "ME17", "L17", "HN00017", "NAT00000017", "PP000017", "080000017", "patient17@example.com", 1987, 6, 18, 0},
		{"Sunrise Medical", "First18", "M18", "Last18", "F18", "ME18", "L18", "HN00018", "NAT00000018", "PP000018", "080000018", "patient18@example.com", 1988, 7, 19, 1},
		{"Central Hospital", "First19", "M19", "Last19", "F19", "ME19", "L19", "HN00019", "NAT00000019", "PP000019", "080000019", "patient19@example.com", 1989, 8, 20, 0},
		{"Green Valley Hospital", "First20", "M20", "Last20", "F20", "ME20", "L20", "HN00020", "NAT00000020", "PP000020", "080000020", "patient20@example.com", 1990, 9, 21, 1},
		{"Sunrise Medical", "First21", "M21", "Last21", "F21", "ME21", "L21", "HN00021", "NAT00000021", "PP000021", "080000021", "patient21@example.com", 1991, 10, 22, 0},
		{"Central Hospital", "First22", "M22", "Last22", "F22", "ME22", "L22", "HN00022", "NAT00000022", "PP000022", "080000022", "patient22@example.com", 1992, 11, 23, 1},
		{"Green Valley Hospital", "First23", "M23", "Last23", "F23", "ME23", "L23", "HN00023", "NAT00000023", "PP000023", "080000023", "patient23@example.com", 1993, 12, 24, 0},
		{"Sunrise Medical", "First24", "M24", "Last24", "F24", "ME24", "L24", "HN00024", "NAT00000024", "PP000024", "080000024", "patient24@example.com", 1994, 1, 25, 1},
		{"Central Hospital", "First25", "M25", "Last25", "F25", "ME25", "L25", "HN00025", "NAT00000025", "PP000025", "080000025", "patient25@example.com", 1995, 2, 26, 0},
		{"Green Valley Hospital", "First26", "M26", "Last26", "F26", "ME26", "L26", "HN00026", "NAT00000026", "PP000026", "080000026", "patient26@example.com", 1996, 3, 27, 1},
		{"Sunrise Medical", "First27", "M27", "Last27", "F27", "ME27", "L27", "HN00027", "NAT00000027", "PP000027", "080000027", "patient27@example.com", 1997, 4, 28, 0},
		{"Central Hospital", "First28", "M28", "Last28", "F28", "ME28", "L28", "HN00028", "NAT00000028", "PP000028", "080000028", "patient28@example.com", 1998, 5, 29, 1},
		{"Green Valley Hospital", "First29", "M29", "Last29", "F29", "ME29", "L29", "HN00029", "NAT00000029", "PP000029", "080000029", "patient29@example.com", 1999, 6, 30, 0},
		{"Sunrise Medical", "First30", "M30", "Last30", "F30", "ME30", "L30", "HN00030", "NAT00000030", "PP000030", "080000030", "patient30@example.com", 2000, 7, 1, 1},
		{"Central Hospital", "First31", "M31", "Last31", "F31", "ME31", "L31", "HN00031", "NAT00000031", "PP000031", "080000031", "patient31@example.com", 2001, 8, 2, 0},
		{"Green Valley Hospital", "First32", "M32", "Last32", "F32", "ME32", "L32", "HN00032", "NAT00000032", "PP000032", "080000032", "patient32@example.com", 2002, 9, 3, 1},
		{"Sunrise Medical", "First33", "M33", "Last33", "F33", "ME33", "L33", "HN00033", "NAT00000033", "PP000033", "080000033", "patient33@example.com", 2003, 10, 4, 0},
		{"Central Hospital", "First34", "M34", "Last34", "F34", "ME34", "L34", "HN00034", "NAT00000034", "PP000034", "080000034", "patient34@example.com", 2004, 11, 5, 1},
		{"Green Valley Hospital", "First35", "M35", "Last35", "F35", "ME35", "L35", "HN00035", "NAT00000035", "PP000035", "080000035", "patient35@example.com", 2005, 12, 6, 0},
		{"Sunrise Medical", "First36", "M36", "Last36", "F36", "ME36", "L36", "HN00036", "NAT00000036", "PP000036", "080000036", "patient36@example.com", 2006, 1, 7, 1},
		{"Central Hospital", "First37", "M37", "Last37", "F37", "ME37", "L37", "HN00037", "NAT00000037", "PP000037", "080000037", "patient37@example.com", 2007, 2, 8, 0},
		{"Green Valley Hospital", "First38", "M38", "Last38", "F38", "ME38", "L38", "HN00038", "NAT00000038", "PP000038", "080000038", "patient38@example.com", 2008, 3, 9, 1},
		{"Sunrise Medical", "First39", "M39", "Last39", "F39", "ME39", "L39", "HN00039", "NAT00000039", "PP000039", "080000039", "patient39@example.com", 2009, 4, 10, 0},
		{"Central Hospital", "First40", "M40", "Last40", "F40", "ME40", "L40", "HN00040", "NAT00000040", "PP000040", "080000040", "patient40@example.com", 2010, 5, 11, 1},
		{"Green Valley Hospital", "First41", "M41", "Last41", "F41", "ME41", "L41", "HN00041", "NAT00000041", "PP000041", "080000041", "patient41@example.com", 2011, 6, 12, 0},
		{"Sunrise Medical", "First42", "M42", "Last42", "F42", "ME42", "L42", "HN00042", "NAT00000042", "PP000042", "080000042", "patient42@example.com", 2012, 7, 13, 1},
		{"Central Hospital", "First43", "M43", "Last43", "F43", "ME43", "L43", "HN00043", "NAT00000043", "PP000043", "080000043", "patient43@example.com", 2013, 8, 14, 0},
		{"Green Valley Hospital", "First44", "M44", "Last44", "F44", "ME44", "L44", "HN00044", "NAT00000044", "PP000044", "080000044", "patient44@example.com", 2014, 9, 15, 1},
		{"Sunrise Medical", "First45", "M45", "Last45", "F45", "ME45", "L45", "HN00045", "NAT00000045", "PP000045", "080000045", "patient45@example.com", 2015, 10, 16, 0},
		{"Central Hospital", "First46", "M46", "Last46", "F46", "ME46", "L46", "HN00046", "NAT00000046", "PP000046", "080000046", "patient46@example.com", 2016, 11, 17, 1},
		{"Green Valley Hospital", "First47", "M47", "Last47", "F47", "ME47", "L47", "HN00047", "NAT00000047", "PP000047", "080000047", "patient47@example.com", 2017, 12, 18, 0},
		{"Sunrise Medical", "First48", "M48", "Last48", "F48", "ME48", "L48", "HN00048", "NAT00000048", "PP000048", "080000048", "patient48@example.com", 2018, 1, 19, 1},
		{"Central Hospital", "First49", "M49", "Last49", "F49", "ME49", "L49", "HN00049", "NAT00000049", "PP000049", "080000049", "patient49@example.com", 2019, 2, 20, 0},
		{"Green Valley Hospital", "First50", "M50", "Last50", "F50", "ME50", "L50", "HN00050", "NAT00000050", "PP000050", "080000050", "patient50@example.com", 2020, 3, 21, 1},
		{"Sunrise Medical", "First51", "M51", "Last51", "F51", "ME51", "L51", "HN00051", "NAT00000051", "PP000051", "080000051", "patient51@example.com", 1971, 4, 22, 0},
		{"Central Hospital", "First52", "M52", "Last52", "F52", "ME52", "L52", "HN00052", "NAT00000052", "PP000052", "080000052", "patient52@example.com", 1972, 5, 23, 1},
		{"Green Valley Hospital", "First53", "M53", "Last53", "F53", "ME53", "L53", "HN00053", "NAT00000053", "PP000053", "080000053", "patient53@example.com", 1973, 6, 24, 0},
		{"Sunrise Medical", "First54", "M54", "Last54", "F54", "ME54", "L54", "HN00054", "NAT00000054", "PP000054", "080000054", "patient54@example.com", 1974, 7, 25, 1},
		{"Central Hospital", "First55", "M55", "Last55", "F55", "ME55", "L55", "HN00055", "NAT00000055", "PP000055", "080000055", "patient55@example.com", 1975, 8, 26, 0},
		{"Green Valley Hospital", "First56", "M56", "Last56", "F56", "ME56", "L56", "HN00056", "NAT00000056", "PP000056", "080000056", "patient56@example.com", 1976, 9, 27, 1},
		{"Sunrise Medical", "First57", "M57", "Last57", "F57", "ME57", "L57", "HN00057", "NAT00000057", "PP000057", "080000057", "patient57@example.com", 1977, 10, 28, 0},
		{"Central Hospital", "First58", "M58", "Last58", "F58", "ME58", "L58", "HN00058", "NAT00000058", "PP000058", "080000058", "patient58@example.com", 1978, 11, 29, 1},
		{"Green Valley Hospital", "First59", "M59", "Last59", "F59", "ME59", "L59", "HN00059", "NAT00000059", "PP000059", "080000059", "patient59@example.com", 1979, 12, 30, 0},
		{"Sunrise Medical", "First60", "M60", "Last60", "F60", "ME60", "L60", "HN00060", "NAT00000060", "PP000060", "080000060", "patient60@example.com", 1980, 1, 1, 1},
		{"Central Hospital", "First61", "M61", "Last61", "F61", "ME61", "L61", "HN00061", "NAT00000061", "PP000061", "080000061", "patient61@example.com", 1981, 2, 2, 0},
		{"Green Valley Hospital", "First62", "M62", "Last62", "F62", "ME62", "L62", "HN00062", "NAT00000062", "PP000062", "080000062", "patient62@example.com", 1982, 3, 3, 1},
		{"Sunrise Medical", "First63", "M63", "Last63", "F63", "ME63", "L63", "HN00063", "NAT00000063", "PP000063", "080000063", "patient63@example.com", 1983, 4, 4, 0},
		{"Central Hospital", "First64", "M64", "Last64", "F64", "ME64", "L64", "HN00064", "NAT00000064", "PP000064", "080000064", "patient64@example.com", 1984, 5, 5, 1},
		{"Green Valley Hospital", "First65", "M65", "Last65", "F65", "ME65", "L65", "HN00065", "NAT00000065", "PP000065", "080000065", "patient65@example.com", 1985, 6, 6, 0},
		{"Sunrise Medical", "First66", "M66", "Last66", "F66", "ME66", "L66", "HN00066", "NAT00000066", "PP000066", "080000066", "patient66@example.com", 1986, 7, 7, 1},
		{"Central Hospital", "First67", "M67", "Last67", "F67", "ME67", "L67", "HN00067", "NAT00000067", "PP000067", "080000067", "patient67@example.com", 1987, 8, 8, 0},
		{"Green Valley Hospital", "First68", "M68", "Last68", "F68", "ME68", "L68", "HN00068", "NAT00000068", "PP000068", "080000068", "patient68@example.com", 1988, 9, 9, 1},
		{"Sunrise Medical", "First69", "M69", "Last69", "F69", "ME69", "L69", "HN00069", "NAT00000069", "PP000069", "080000069", "patient69@example.com", 1989, 10, 10, 0},
		{"Central Hospital", "First70", "M70", "Last70", "F70", "ME70", "L70", "HN00070", "NAT00000070", "PP000070", "080000070", "patient70@example.com", 1990, 11, 11, 1},
		{"Green Valley Hospital", "First71", "M71", "Last71", "F71", "ME71", "L71", "HN00071", "NAT00000071", "PP000071", "080000071", "patient71@example.com", 1991, 12, 12, 0},
		{"Sunrise Medical", "First72", "M72", "Last72", "F72", "ME72", "L72", "HN00072", "NAT00000072", "PP000072", "080000072", "patient72@example.com", 1992, 1, 13, 1},
		{"Central Hospital", "First73", "M73", "Last73", "F73", "ME73", "L73", "HN00073", "NAT00000073", "PP000073", "080000073", "patient73@example.com", 1993, 2, 14, 0},
		{"Green Valley Hospital", "First74", "M74", "Last74", "F74", "ME74", "L74", "HN00074", "NAT00000074", "PP000074", "080000074", "patient74@example.com", 1994, 3, 15, 1},
		{"Sunrise Medical", "First75", "M75", "Last75", "F75", "ME75", "L75", "HN00075", "NAT00000075", "PP000075", "080000075", "patient75@example.com", 1995, 4, 16, 0},
		{"Central Hospital", "First76", "M76", "Last76", "F76", "ME76", "L76", "HN00076", "NAT00000076", "PP000076", "080000076", "patient76@example.com", 1996, 5, 17, 1},
		{"Green Valley Hospital", "First77", "M77", "Last77", "F77", "ME77", "L77", "HN00077", "NAT00000077", "PP000077", "080000077", "patient77@example.com", 1997, 6, 18, 0},
		{"Sunrise Medical", "First78", "M78", "Last78", "F78", "ME78", "L78", "HN00078", "NAT00000078", "PP000078", "080000078", "patient78@example.com", 1998, 7, 19, 1},
		{"Central Hospital", "First79", "M79", "Last79", "F79", "ME79", "L79", "HN00079", "NAT00000079", "PP000079", "080000079", "patient79@example.com", 1999, 8, 20, 0},
		{"Green Valley Hospital", "First80", "M80", "Last80", "F80", "ME80", "L80", "HN00080", "NAT00000080", "PP000080", "080000080", "patient80@example.com", 2000, 9, 21, 1},
		{"Sunrise Medical", "First81", "M81", "Last81", "F81", "ME81", "L81", "HN00081", "NAT00000081", "PP000081", "080000081", "patient81@example.com", 2001, 10, 22, 0},
		{"Central Hospital", "First82", "M82", "Last82", "F82", "ME82", "L82", "HN00082", "NAT00000082", "PP000082", "080000082", "patient82@example.com", 2002, 11, 23, 1},
		{"Green Valley Hospital", "First83", "M83", "Last83", "F83", "ME83", "L83", "HN00083", "NAT00000083", "PP000083", "080000083", "patient83@example.com", 2003, 12, 24, 0},
		{"Sunrise Medical", "First84", "M84", "Last84", "F84", "ME84", "L84", "HN00084", "NAT00000084", "PP000084", "080000084", "patient84@example.com", 2004, 1, 25, 1},
		{"Central Hospital", "First85", "M85", "Last85", "F85", "ME85", "L85", "HN00085", "NAT00000085", "PP000085", "080000085", "patient85@example.com", 2005, 2, 26, 0},
		{"Green Valley Hospital", "First86", "M86", "Last86", "F86", "ME86", "L86", "HN00086", "NAT00000086", "PP000086", "080000086", "patient86@example.com", 2006, 3, 27, 1},
		{"Sunrise Medical", "First87", "M87", "Last87", "F87", "ME87", "L87", "HN00087", "NAT00000087", "PP000087", "080000087", "patient87@example.com", 2007, 4, 28, 0},
		{"Central Hospital", "First88", "M88", "Last88", "F88", "ME88", "L88", "HN00088", "NAT00000088", "PP000088", "080000088", "patient88@example.com", 2008, 5, 29, 1},
		{"Green Valley Hospital", "First89", "M89", "Last89", "F89", "ME89", "L89", "HN00089", "NAT00000089", "PP000089", "080000089", "patient89@example.com", 2009, 6, 30, 0},
		{"Sunrise Medical", "First90", "M90", "Last90", "F90", "ME90", "L90", "HN00090", "NAT00000090", "PP000090", "080000090", "patient90@example.com", 2010, 7, 1, 1},
		{"Central Hospital", "First91", "M91", "Last91", "F91", "ME91", "L91", "HN00091", "NAT00000091", "PP000091", "080000091", "patient91@example.com", 2011, 8, 2, 0},
		{"Green Valley Hospital", "First92", "M92", "Last92", "F92", "ME92", "L92", "HN00092", "NAT00000092", "PP000092", "080000092", "patient92@example.com", 2012, 9, 3, 1},
		{"Sunrise Medical", "First93", "M93", "Last93", "F93", "ME93", "L93", "HN00093", "NAT00000093", "PP000093", "080000093", "patient93@example.com", 2013, 10, 4, 0},
		{"Central Hospital", "First94", "M94", "Last94", "F94", "ME94", "L94", "HN00094", "NAT00000094", "PP000094", "080000094", "patient94@example.com", 2014, 11, 5, 1},
		{"Green Valley Hospital", "First95", "M95", "Last95", "F95", "ME95", "L95", "HN00095", "NAT00000095", "PP000095", "080000095", "patient95@example.com", 2015, 12, 6, 0},
		{"Sunrise Medical", "First96", "M96", "Last96", "F96", "ME96", "L96", "HN00096", "NAT00000096", "PP000096", "080000096", "patient96@example.com", 2016, 1, 7, 1},
		{"Central Hospital", "First97", "M97", "Last97", "F97", "ME97", "L97", "HN00097", "NAT00000097", "PP000097", "080000097", "patient97@example.com", 2017, 2, 8, 0},
		{"Green Valley Hospital", "First98", "M98", "Last98", "F98", "ME98", "L98", "HN00098", "NAT00000098", "PP000098", "080000098", "patient98@example.com", 2018, 3, 9, 1},
		{"Sunrise Medical", "First99", "M99", "Last99", "F99", "ME99", "L99", "HN00099", "NAT00000099", "PP000099", "080000099", "patient99@example.com", 2019, 4, 10, 0},
		{"Central Hospital", "First100", "M100", "Last100", "F100", "ME100", "L100", "HN00100", "NAT00000100", "PP000100", "080000100", "patient100@example.com", 2020, 5, 11, 1},
	}

	patients := make([]models.Patient, 0, len(patientsData))
	for _, p := range patientsData {
		hid := hospMap[p.Hospital]
		dob := time.Date(p.Year, time.Month(p.Month), p.Day, 0, 0, 0, 0, time.UTC)
		patients = append(patients, models.Patient{
			HospitalID:   hid,
			FirstNameTH:  strPtr(cleanStr(p.FirstTH)),
			MiddleNameTH: strPtr(cleanStr(p.MiddleTH)),
			LastNameTH:   strPtr(cleanStr(p.LastTH)),
			FirstNameEN:  strPtr(cleanStr(p.FirstEN)),
			MiddleNameEN: strPtr(cleanStr(p.MiddleEN)),
			LastNameEN:   strPtr(cleanStr(p.LastEN)),
			DateOfBirth:  dob,
			PatientHN:    cleanStr(p.HN),
			NationalID:   strPtr(cleanStr(p.Nat)),
			PassportID:   strPtr(cleanStr(p.PP)),
			PhoneNumber:  strPtr(cleanStr(p.Phone)),
			Email:        strPtr(cleanStr(p.Email)),
			Gender:       genderFromInt(p.GenderInt),
		})
	}

	if len(patients) > 0 {
		_ = db.Exec("DEALLOCATE ALL").Error
		sess := db.Session(&gorm.Session{PrepareStmt: false})

		for i := range patients {
			p := &patients[i]
			p.PatientHN = cleanStr(p.PatientHN)
			if p.NationalID != nil {
				v := cleanStr(*p.NationalID)
				p.NationalID = &v
			}
			if p.PassportID != nil {
				v := cleanStr(*p.PassportID)
				p.PassportID = &v
			}
			if p.PhoneNumber != nil {
				v := cleanStr(*p.PhoneNumber)
				p.PhoneNumber = &v
			}
			if p.Email != nil {
				v := cleanStr(*p.Email)
				p.Email = &v
			}
			if p.FirstNameTH != nil {
				v := cleanStr(*p.FirstNameTH)
				p.FirstNameTH = &v
			}
			if p.MiddleNameTH != nil {
				v := cleanStr(*p.MiddleNameTH)
				p.MiddleNameTH = &v
			}
			if p.LastNameTH != nil {
				v := cleanStr(*p.LastNameTH)
				p.LastNameTH = &v
			}
			if p.FirstNameEN != nil {
				v := cleanStr(*p.FirstNameEN)
				p.FirstNameEN = &v
			}
			if p.MiddleNameEN != nil {
				v := cleanStr(*p.MiddleNameEN)
				p.MiddleNameEN = &v
			}
			if p.LastNameEN != nil {
				v := cleanStr(*p.LastNameEN)
				p.LastNameEN = &v
			}
		}

		if err := sess.CreateInBatches(patients, 100).Error; err != nil {
			log.Fatalf("batch insert failed: %v", err)
		}
	}

	fmt.Printf("seeding complete: hospitals=%d, staff=%d, patients=%d\n", len(hospitals), staffCreated, len(patients))
}
