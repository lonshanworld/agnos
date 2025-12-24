package main

import (
	"fmt"
	"log"
	"time"

	"agnos_candidate_assignment/config"
	"agnos_candidate_assignment/database"
	"agnos_candidate_assignment/models"
	"agnos_candidate_assignment/repositories"
	"agnos_candidate_assignment/services"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	cfg := config.Load()

	db, err := database.NewPostgresConnectionNoMigrate(cfg)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	// Ensure schema exists for seeding.
	if err := db.AutoMigrate(&models.Hospital{}, &models.Staff{}, &models.Patient{}); err != nil {
		log.Fatalf("auto migrate (seed) failed: %v", err)
	}
	fmt.Println("schema ensured: hospitals, staff, patients")

	hospRepo := repositories.NewHospitalRepository(db)
	staffRepo := repositories.NewStaffRepository(db)

	authSvc := services.NewAuthService(staffRepo, hospRepo, cfg)

	// deterministic arrays for visibility
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
	}

	// fully hardcoded patient records (explicit values for every column except ID)
	patientData := []struct {
		Hospital   string
		First      string
		Middle     string
		Last       string
		DOB        time.Time
		PatientHN  string
		NationalID string
		PassportID string
		Phone      string
		Email      string
		Gender     models.Gender
	}{
		{"Central Hospital", "Narin", "Chai", "Sukcharoen", time.Date(1990, 3, 12, 0, 0, 0, 0, time.UTC), "HN00001", "NAT0000001", "PP000001", "0800000001", "narin.s1@central.example", models.Male},
		{"Central Hospital", "Somchai", "Boon", "Pradchaphet", time.Date(1985, 7, 2, 0, 0, 0, 0, time.UTC), "HN00002", "NAT0000002", "PP000002", "0810000002", "somchai.p2@central.example", models.Male},
		{"Central Hospital", "Somsri", "Mai", "Wattana", time.Date(1978, 11, 21, 0, 0, 0, 0, time.UTC), "HN00003", "NAT0000003", "PP000003", "0820000003", "somsri.w3@central.example", models.Female},
		{"Central Hospital", "Kamon", "Preecha", "Chaiyawan", time.Date(1992, 1, 9, 0, 0, 0, 0, time.UTC), "HN00004", "NAT0000004", "PP000004", "0830000004", "kamon.c4@central.example", models.Male},
		{"Central Hospital", "Pimchanok", "Nong", "Saengdao", time.Date(1988, 4, 30, 0, 0, 0, 0, time.UTC), "HN00005", "NAT0000005", "PP000005", "0840000005", "pim.s5@central.example", models.Female},
		{"Central Hospital", "Anan", "Sak", "Ratanapak", time.Date(1975, 10, 5, 0, 0, 0, 0, time.UTC), "HN00006", "NAT0000006", "PP000006", "0850000006", "anan.r6@central.example", models.Male},
		{"Central Hospital", "Ladda", "Kae", "Preecha", time.Date(1995, 6, 17, 0, 0, 0, 0, time.UTC), "HN00007", "NAT0000007", "PP000007", "0860000007", "ladda.p7@central.example", models.Female},
		{"Central Hospital", "Wichai", "Mon", "Klong", time.Date(1982, 8, 24, 0, 0, 0, 0, time.UTC), "HN00008", "NAT0000008", "PP000008", "0870000008", "wichai.k8@central.example", models.Male},
		{"Central Hospital", "Nicha", "Ying", "Kongthong", time.Date(1991, 12, 3, 0, 0, 0, 0, time.UTC), "HN00009", "NAT0000009", "PP000009", "0880000009", "nicha.k9@central.example", models.Female},
		{"Central Hospital", "Parinya", "Lek", "Taworn", time.Date(1984, 5, 14, 0, 0, 0, 0, time.UTC), "HN00010", "NAT0000010", "PP000010", "0890000010", "parinya.t10@central.example", models.Male},
		{"Central Hospital", "Malee", "Duang", "Jaroen", time.Date(1979, 9, 6, 0, 0, 0, 0, time.UTC), "HN00011", "NAT0000011", "PP000011", "0900000011", "malee.j11@central.example", models.Female},
		{"Central Hospital", "Sakda", "Jit", "Phong", time.Date(1987, 2, 28, 0, 0, 0, 0, time.UTC), "HN00012", "NAT0000012", "PP000012", "0910000012", "sakda.p12@central.example", models.Male},
		{"Central Hospital", "Kanya", "Sai", "Phrom", time.Date(1993, 3, 19, 0, 0, 0, 0, time.UTC), "HN00013", "NAT0000013", "PP000013", "0920000013", "kanya.p13@central.example", models.Female},
		{"Central Hospital", "Pongsak", "Noi", "Nakorn", time.Date(1980, 7, 11, 0, 0, 0, 0, time.UTC), "HN00014", "NAT0000014", "PP000014", "0930000014", "pongsak.n14@central.example", models.Male},
		{"Central Hospital", "Chompoo", "Wan", "Roi", time.Date(1994, 11, 29, 0, 0, 0, 0, time.UTC), "HN00015", "NAT0000015", "PP000015", "0940000015", "chompoo.r15@central.example", models.Female},
		{"Green Valley Hospital", "Thanapol", "Pit", "Srisuk", time.Date(1986, 6, 4, 0, 0, 0, 0, time.UTC), "HN00016", "NAT0000016", "PP000016", "0950000016", "thanapol.s16@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Krit", "Ton", "Chantra", time.Date(1990, 10, 20, 0, 0, 0, 0, time.UTC), "HN00017", "NAT0000017", "PP000017", "0960000017", "krit.c17@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Orn", "Mee", "Sawasdee", time.Date(1983, 1, 15, 0, 0, 0, 0, time.UTC), "HN00018", "NAT0000018", "PP000018", "0970000018", "orn.s18@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Tida", "Fon", "Kulap", time.Date(1996, 9, 8, 0, 0, 0, 0, time.UTC), "HN00019", "NAT0000019", "PP000019", "0980000019", "tida.k19@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Nattapong", "Bee", "Siri", time.Date(1981, 12, 25, 0, 0, 0, 0, time.UTC), "HN00020", "NAT0000020", "PP000020", "0990000020", "nattapong.s20@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Aroon", "Man", "Prayoon", time.Date(1977, 2, 3, 0, 0, 0, 0, time.UTC), "HN00021", "NAT0000021", "PP000021", "0800000021", "aroon.m21@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Benja", "Sak", "Lert", time.Date(1989, 4, 12, 0, 0, 0, 0, time.UTC), "HN00022", "NAT0000022", "PP000022", "0810000022", "benja.s22@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Chai", "Nop", "Suwan", time.Date(1992, 5, 30, 0, 0, 0, 0, time.UTC), "HN00023", "NAT0000023", "PP000023", "0820000023", "chai.n23@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Duang", "Korn", "Phan", time.Date(1976, 8, 19, 0, 0, 0, 0, time.UTC), "HN00024", "NAT0000024", "PP000024", "0830000024", "duang.k24@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Ekkachai", "Tao", "Ruen", time.Date(1984, 11, 11, 0, 0, 0, 0, time.UTC), "HN00025", "NAT0000025", "PP000025", "0840000025", "ekkachai.t25@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Fah", "Mai", "Loy", time.Date(1991, 3, 2, 0, 0, 0, 0, time.UTC), "HN00026", "NAT0000026", "PP000026", "0850000026", "fah.m26@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Gavin", "Leo", "Sila", time.Date(1993, 7, 7, 0, 0, 0, 0, time.UTC), "HN00027", "NAT0000027", "PP000027", "0860000027", "gavin.l27@greenvalley.example", models.Male},
		{"Sunrise Medical", "Hana", "Poon", "Kiet", time.Date(1988, 12, 12, 0, 0, 0, 0, time.UTC), "HN00028", "NAT0000028", "PP000028", "0870000028", "hana.p28@sunrise.example", models.Female},
		{"Sunrise Medical", "Ittipon", "Ram", "Noi", time.Date(1980, 6, 6, 0, 0, 0, 0, time.UTC), "HN00029", "NAT0000029", "PP000029", "0880000029", "ittipon.r29@sunrise.example", models.Male},
		{"Sunrise Medical", "Jira", "Som", "Boon", time.Date(1979, 9, 9, 0, 0, 0, 0, time.UTC), "HN00030", "NAT0000030", "PP000030", "0890000030", "jira.s30@sunrise.example", models.Male},
		{"Sunrise Medical", "Ketsara", "Ung", "Suk", time.Date(1994, 2, 2, 0, 0, 0, 0, time.UTC), "HN00031", "NAT0000031", "PP000031", "0900000031", "ketsara.u31@sunrise.example", models.Female},
		{"Sunrise Medical", "Lek", "Pong", "Art", time.Date(1992, 10, 10, 0, 0, 0, 0, time.UTC), "HN00032", "NAT0000032", "PP000032", "0910000032", "lek.p32@sunrise.example", models.Male},
		{"Sunrise Medical", "Malee", "Rin", "Chai", time.Date(1983, 5, 5, 0, 0, 0, 0, time.UTC), "HN00033", "NAT0000033", "PP000033", "0920000033", "malee.r33@sunrise.example", models.Female},
		{"Sunrise Medical", "Nate", "Som", "Korn", time.Date(1986, 8, 8, 0, 0, 0, 0, time.UTC), "HN00034", "NAT0000034", "PP000034", "0930000034", "nate.s34@sunrise.example", models.Male},
		{"Central Hospital", "Oon", "Kham", "Win", time.Date(1990, 4, 4, 0, 0, 0, 0, time.UTC), "HN00035", "NAT0000035", "PP000035", "0940000035", "oon.k35@central.example", models.Female},
		{"Central Hospital", "Pree", "Chan", "Som", time.Date(1978, 1, 1, 0, 0, 0, 0, time.UTC), "HN00036", "NAT0000036", "PP000036", "0950000036", "pree.c36@central.example", models.Female},
		{"Central Hospital", "Rin", "Phun", "Mai", time.Date(1982, 2, 2, 0, 0, 0, 0, time.UTC), "HN00037", "NAT0000037", "PP000037", "0960000037", "rin.p37@central.example", models.Male},
		{"Central Hospital", "Sutee", "Wan", "Rai", time.Date(1987, 3, 3, 0, 0, 0, 0, time.UTC), "HN00038", "NAT0000038", "PP000038", "0970000038", "sutee.w38@central.example", models.Male},
		{"Central Hospital", "Thida", "Nam", "Ploy", time.Date(1991, 4, 14, 0, 0, 0, 0, time.UTC), "HN00039", "NAT0000039", "PP000039", "0980000039", "thida.n39@central.example", models.Female},
		{"Central Hospital", "Udom", "Sak", "Tan", time.Date(1976, 5, 5, 0, 0, 0, 0, time.UTC), "HN00040", "NAT0000040", "PP000040", "0990000040", "udom.s40@central.example", models.Male},
		{"Green Valley Hospital", "Vipa", "Jai", "Nuan", time.Date(1989, 6, 6, 0, 0, 0, 0, time.UTC), "HN00041", "NAT0000041", "PP000041", "0800000041", "vipa.j41@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Worawut", "Lek", "Yen", time.Date(1984, 7, 7, 0, 0, 0, 0, time.UTC), "HN00042", "NAT0000042", "PP000042", "0810000042", "worawut.l42@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Xing", "Yu", "Lee", time.Date(1990, 8, 8, 0, 0, 0, 0, time.UTC), "HN00043", "NAT0000043", "PP000043", "0820000043", "xing.y43@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Yada", "Noi", "Kiet", time.Date(1979, 9, 9, 0, 0, 0, 0, time.UTC), "HN00044", "NAT0000044", "PP000044", "0830000044", "yada.n44@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Zin", "Pa", "Jit", time.Date(1992, 10, 10, 0, 0, 0, 0, time.UTC), "HN00045", "NAT0000045", "PP000045", "0840000045", "zin.p45@greenvalley.example", models.Male},
		{"Sunrise Medical", "Aom", "Pee", "Suk", time.Date(1980, 11, 11, 0, 0, 0, 0, time.UTC), "HN00046", "NAT0000046", "PP000046", "0850000046", "aom.p46@sunrise.example", models.Female},
		{"Sunrise Medical", "Boran", "Rai", "Thong", time.Date(1978, 12, 12, 0, 0, 0, 0, time.UTC), "HN00047", "NAT0000047", "PP000047", "0860000047", "boran.r47@sunrise.example", models.Male},
		{"Sunrise Medical", "Chada", "Pim", "Ruk", time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC), "HN00048", "NAT0000048", "PP000048", "0870000048", "chada.p48@sunrise.example", models.Female},
		{"Sunrise Medical", "Dara", "Noi", "Rin", time.Date(1983, 2, 2, 0, 0, 0, 0, time.UTC), "HN00049", "NAT0000049", "PP000049", "0880000049", "dara.n49@sunrise.example", models.Female},
		{"Sunrise Medical", "Eak", "Som", "Kut", time.Date(1986, 3, 3, 0, 0, 0, 0, time.UTC), "HN00050", "NAT0000050", "PP000050", "0890000050", "eak.s50@sunrise.example", models.Male},
		{"Central Hospital", "Fong", "Suk", "Mai", time.Date(1977, 4, 4, 0, 0, 0, 0, time.UTC), "HN00051", "NAT0000051", "PP000051", "0900000051", "fong.s51@central.example", models.Female},
		{"Central Hospital", "Ganya", "Tee", "Ploy", time.Date(1981, 5, 5, 0, 0, 0, 0, time.UTC), "HN00052", "NAT0000052", "PP000052", "0910000052", "ganya.t52@central.example", models.Female},
		{"Central Hospital", "Heng", "Lek", "Pon", time.Date(1982, 6, 6, 0, 0, 0, 0, time.UTC), "HN00053", "NAT0000053", "PP000053", "0920000053", "heng.l53@central.example", models.Male},
		{"Central Hospital", "Issara", "Mon", "Cha", time.Date(1989, 7, 7, 0, 0, 0, 0, time.UTC), "HN00054", "NAT0000054", "PP000054", "0930000054", "issara.m54@central.example", models.Male},
		{"Central Hospital", "Jintana", "Nok", "Pra", time.Date(1993, 8, 8, 0, 0, 0, 0, time.UTC), "HN00055", "NAT0000055", "PP000055", "0940000055", "jintana.n55@central.example", models.Female},
		{"Green Valley Hospital", "Kamonrat", "Lek", "Sun", time.Date(1975, 9, 9, 0, 0, 0, 0, time.UTC), "HN00056", "NAT0000056", "PP000056", "0950000056", "kamonrat.l56@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Lom", "Chai", "Bun", time.Date(1976, 10, 10, 0, 0, 0, 0, time.UTC), "HN00057", "NAT0000057", "PP000057", "0960000057", "lom.c57@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Mali", "Phai", "Jun", time.Date(1988, 11, 11, 0, 0, 0, 0, time.UTC), "HN00058", "NAT0000058", "PP000058", "0970000058", "mali.p58@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Narin", "Long", "Tao", time.Date(1990, 12, 12, 0, 0, 0, 0, time.UTC), "HN00059", "NAT0000059", "PP000059", "0980000059", "narin.l59@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Oat", "Ploy", "Chai", time.Date(1991, 1, 1, 0, 0, 0, 0, time.UTC), "HN00060", "NAT0000060", "PP000060", "0990000060", "oat.p60@greenvalley.example", models.Male},
		{"Sunrise Medical", "Ploy", "Nem", "Sut", time.Date(1979, 2, 2, 0, 0, 0, 0, time.UTC), "HN00061", "NAT0000061", "PP000061", "0800000061", "ploy.n61@sunrise.example", models.Female},
		{"Sunrise Medical", "Rapee", "Yui", "Som", time.Date(1984, 3, 3, 0, 0, 0, 0, time.UTC), "HN00062", "NAT0000062", "PP000062", "0810000062", "rapeey.u62@sunrise.example", models.Female},
		{"Sunrise Medical", "Sainan", "Tok", "Mee", time.Date(1986, 4, 4, 0, 0, 0, 0, time.UTC), "HN00063", "NAT0000063", "PP000063", "0820000063", "sainan.t63@sunrise.example", models.Male},
		{"Sunrise Medical", "Tarn", "Ploy", "Nok", time.Date(1992, 5, 5, 0, 0, 0, 0, time.UTC), "HN00064", "NAT0000064", "PP000064", "0830000064", "tarn.p64@sunrise.example", models.Female},
		{"Central Hospital", "Ubon", "Chai", "Nam", time.Date(1978, 6, 6, 0, 0, 0, 0, time.UTC), "HN00065", "NAT0000065", "PP000065", "0840000065", "ubon.c65@central.example", models.Male},
		{"Central Hospital", "Vichai", "Suk", "Pon", time.Date(1980, 7, 7, 0, 0, 0, 0, time.UTC), "HN00066", "NAT0000066", "PP000066", "0850000066", "vichai.s66@central.example", models.Male},
		{"Central Hospital", "Wipa", "Nai", "Lom", time.Date(1985, 8, 8, 0, 0, 0, 0, time.UTC), "HN00067", "NAT0000067", "PP000067", "0860000067", "wipa.n67@central.example", models.Female},
		{"Central Hospital", "Xena", "Ploy", "Rin", time.Date(1990, 9, 9, 0, 0, 0, 0, time.UTC), "HN00068", "NAT0000068", "PP000068", "0870000068", "xena.p68@central.example", models.Female},
		{"Central Hospital", "Ying", "Mai", "Sun", time.Date(1979, 10, 10, 0, 0, 0, 0, time.UTC), "HN00069", "NAT0000069", "PP000069", "0880000069", "ying.m69@central.example", models.Female},
		{"Central Hospital", "Zara", "Ploy", "Nok", time.Date(1982, 11, 11, 0, 0, 0, 0, time.UTC), "HN00070", "NAT0000070", "PP000070", "0890000070", "zara.p70@central.example", models.Female},
		{"Green Valley Hospital", "Aree", "Suk", "Mai", time.Date(1983, 12, 12, 0, 0, 0, 0, time.UTC), "HN00071", "NAT0000071", "PP000071", "0900000071", "aree.s71@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Boon", "Nai", "Tee", time.Date(1977, 1, 1, 0, 0, 0, 0, time.UTC), "HN00072", "NAT0000072", "PP000072", "0910000072", "boon.n72@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Chut", "Mee", "Pon", time.Date(1991, 2, 2, 0, 0, 0, 0, time.UTC), "HN00073", "NAT0000073", "PP000073", "0920000073", "chut.m73@greenvalley.example", models.Male},
		{"Green Valley Hospital", "Duang", "Ploy", "Sai", time.Date(1988, 3, 3, 0, 0, 0, 0, time.UTC), "HN00074", "NAT0000074", "PP000074", "0930000074", "duang.p74@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Ekk", "Nok", "Wan", time.Date(1976, 4, 4, 0, 0, 0, 0, time.UTC), "HN00075", "NAT0000075", "PP000075", "0940000075", "ekk.n75@greenvalley.example", models.Male},
		{"Sunrise Medical", "Fai", "Mee", "Thong", time.Date(1993, 5, 5, 0, 0, 0, 0, time.UTC), "HN00076", "NAT0000076", "PP000076", "0950000076", "fai.m76@sunrise.example", models.Female},
		{"Sunrise Medical", "Ganda", "Ploy", "Lek", time.Date(1981, 6, 6, 0, 0, 0, 0, time.UTC), "HN00077", "NAT0000077", "PP000077", "0960000077", "ganda.p77@sunrise.example", models.Female},
		{"Sunrise Medical", "Hoon", "Sai", "Tan", time.Date(1979, 7, 7, 0, 0, 0, 0, time.UTC), "HN00078", "NAT0000078", "PP000078", "0970000078", "hoon.s78@sunrise.example", models.Male},
		{"Sunrise Medical", "Irene", "Mai", "Ploy", time.Date(1990, 8, 8, 0, 0, 0, 0, time.UTC), "HN00079", "NAT0000079", "PP000079", "0980000079", "irene.m79@sunrise.example", models.Female},
		{"Sunrise Medical", "Jade", "Noi", "Rai", time.Date(1984, 9, 9, 0, 0, 0, 0, time.UTC), "HN00080", "NAT0000080", "PP000080", "0990000080", "jade.n80@sunrise.example", models.Female},
		{"Central Hospital", "Kae", "Pun", "Cher", time.Date(1985, 10, 10, 0, 0, 0, 0, time.UTC), "HN00081", "NAT0000081", "PP000081", "0800000081", "kae.p81@central.example", models.Female},
		{"Central Hospital", "Lon", "Art", "Mee", time.Date(1978, 11, 11, 0, 0, 0, 0, time.UTC), "HN00082", "NAT0000082", "PP000082", "0810000082", "lon.a82@central.example", models.Male},
		{"Central Hospital", "May", "Pim", "Noi", time.Date(1986, 12, 12, 0, 0, 0, 0, time.UTC), "HN00083", "NAT0000083", "PP000083", "0820000083", "may.p83@central.example", models.Female},
		{"Central Hospital", "Nok", "Sum", "Ban", time.Date(1992, 1, 1, 0, 0, 0, 0, time.UTC), "HN00084", "NAT0000084", "PP000084", "0830000084", "nok.s84@central.example", models.Female},
		{"Green Valley Hospital", "Ora", "Lek", "Kam", time.Date(1975, 2, 2, 0, 0, 0, 0, time.UTC), "HN00085", "NAT0000085", "PP000085", "0840000085", "ora.l85@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Pim", "Rai", "Tan", time.Date(1983, 3, 3, 0, 0, 0, 0, time.UTC), "HN00086", "NAT0000086", "PP000086", "0850000086", "pim.r86@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Quiz", "Ton", "Lek", time.Date(1989, 4, 4, 0, 0, 0, 0, time.UTC), "HN00087", "NAT0000087", "PP000087", "0860000087", "quiz.t87@greenvalley.example", models.Male},
		{"Sunrise Medical", "Rin", "Mai", "Pon", time.Date(1991, 5, 5, 0, 0, 0, 0, time.UTC), "HN00088", "NAT0000088", "PP000088", "0870000088", "rin.m88@sunrise.example", models.Female},
		{"Sunrise Medical", "Sia", "Ploy", "Nok", time.Date(1977, 6, 6, 0, 0, 0, 0, time.UTC), "HN00089", "NAT0000089", "PP000089", "0880000089", "sia.p89@sunrise.example", models.Female},
		{"Sunrise Medical", "Ton", "Ari", "Mai", time.Date(1982, 7, 7, 0, 0, 0, 0, time.UTC), "HN00090", "NAT0000090", "PP000090", "0890000090", "ton.a90@sunrise.example", models.Male},
		{"Central Hospital", "Umi", "Pee", "Kao", time.Date(1984, 8, 8, 0, 0, 0, 0, time.UTC), "HN00091", "NAT0000091", "PP000091", "0900000091", "umi.p91@central.example", models.Female},
		{"Central Hospital", "Vee", "Tan", "Som", time.Date(1979, 9, 9, 0, 0, 0, 0, time.UTC), "HN00092", "NAT0000092", "PP000092", "0910000092", "vee.t92@central.example", models.Male},
		{"Green Valley Hospital", "Wai", "Phu", "Nan", time.Date(1990, 10, 10, 0, 0, 0, 0, time.UTC), "HN00093", "NAT0000093", "PP000093", "0920000093", "wai.p93@greenvalley.example", models.Female},
		{"Green Valley Hospital", "Xan", "Lo", "Pim", time.Date(1981, 11, 11, 0, 0, 0, 0, time.UTC), "HN00094", "NAT0000094", "PP000094", "0930000094", "xan.l94@greenvalley.example", models.Male},
		{"Sunrise Medical", "Yen", "Mok", "Sut", time.Date(1976, 12, 12, 0, 0, 0, 0, time.UTC), "HN00095", "NAT0000095", "PP000095", "0940000095", "yen.m95@sunrise.example", models.Female},
		{"Sunrise Medical", "Zee", "Rai", "Pon", time.Date(1987, 1, 1, 0, 0, 0, 0, time.UTC), "HN00096", "NAT0000096", "PP000096", "0950000096", "zee.r96@sunrise.example", models.Male},
		{"Central Hospital", "Ari", "Nok", "Pun", time.Date(1988, 2, 2, 0, 0, 0, 0, time.UTC), "HN00097", "NAT0000097", "PP000097", "0960000097", "ari.n97@central.example", models.Male},
		{"Central Hospital", "Bee", "Suk", "Fan", time.Date(1975, 3, 3, 0, 0, 0, 0, time.UTC), "HN00098", "NAT0000098", "PP000098", "0970000098", "bee.s98@central.example", models.Female},
		{"Green Valley Hospital", "Cia", "Lom", "Pun", time.Date(1992, 4, 4, 0, 0, 0, 0, time.UTC), "HN00099", "NAT0000099", "PP000099", "0980000099", "cia.l99@greenvalley.example", models.Female},
		{"Sunrise Medical", "Dai", "Lek", "Fon", time.Date(1993, 5, 5, 0, 0, 0, 0, time.UTC), "HN00100", "NAT0000100", "PP000100", "0990000100", "dai.l100@sunrise.example", models.Male},
	}

	// preload hospital name -> id map to avoid repeated FindByName calls
	hospMap := make(map[string]uint, len(hospitals))
	for _, name := range hospitals {
		h, err := hospRepo.FindByName(name)
		if err == nil {
			fmt.Printf("hospital exists: %s\n", name)
			hospMap[name] = h.ID
			continue
		}
		h = &models.Hospital{Name: name}
		if err = hospRepo.Create(h); err != nil {
			log.Fatalf("failed to create hospital %s: %v", name, err)
		}
		fmt.Printf("created hospital: %s\n", name)
		hospMap[name] = h.ID
	}

	// create staff from deterministic list (use each staff's hardcoded password)
	staffCount := 0
	for _, s := range staffList {
		if _, err := authSvc.Register(s.Hospital, s.Username, s.Password); err != nil {
			fmt.Printf("skip staff %s@%s: %v\n", s.Username, s.Hospital, err)
			continue
		}
		fmt.Printf("created staff %s@%s\n", s.Username, s.Hospital)
		staffCount++
	}

	// create patients: build slice and insert in batches to reduce DB round-trips
	patients := make([]models.Patient, 0, len(patientData))
	for i, pd := range patientData {
		hid, ok := hospMap[pd.Hospital]
		if !ok {
			fmt.Printf("hospital not found for patient %s: skipping\n", pd.NationalID)
			continue
		}
		if i < 20 {
			fmt.Printf("sample patient %d: HOSP=%s HN=%s NAT=%s PAS=%s EMAIL=%s\n", i+1, pd.Hospital, pd.PatientHN, pd.NationalID, pd.PassportID, pd.Email)
		}
		patients = append(patients, models.Patient{
			HospitalID:   hid,
			FirstNameTH:  pd.First,
			MiddleNameTH: pd.Middle,
			LastNameTH:   pd.Last,
			DateOfBirth:  pd.DOB,
			PatientHN:    pd.PatientHN,
			NationalID:   pd.NationalID,
			PassportID:   pd.PassportID,
			PhoneNumber:  pd.Phone,
			Email:        pd.Email,
			Gender:       pd.Gender,
		})
	}

	// perform batch insert (no existence checks; run only once)
	if len(patients) > 0 {
		// use CreateInBatches to control batch size
		if err := db.CreateInBatches(patients, 50).Error; err != nil {
			log.Fatalf("failed to batch insert patients: %v", err)
		}
	}

	fmt.Printf("seeding complete: hospitals=%d, staff=%d, patients=%d\n", len(hospitals), staffCount, len(patients))
}
