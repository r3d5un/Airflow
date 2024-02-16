package brreg

type Unit struct {
	OrganisasjonsNummer                       string                   `json:"organisasjonsnummer"`
	Navn                                      string                   `json:"navn"`
	Organisasjonsform                         Organisasjonsform        `json:"organisasjonsform"`
	RegistreringsdatoEnhetsregister           string                   `json:"registreringsdatoEnhetsregister"`
	RegistrertIMvaRegisteret                  bool                     `json:"registrertIMvaregisteret"`
	Naeringskode1                             NaeringsKode1            `json:"naeringskode1"`
	Naeringskode2                             NaeringsKode2            `json:"naeringskode2"`
	AntallAnsatte                             int                      `json:"antallAnsatte"`
	ForretningsAdresse                        Address                  `json:"forretningsadresse"`
	StiftelsesDato                            string                   `json:"stiftelsesdato"`
	InstitusjonellSektorKode                  InstitusjonellSektorKode `json:"institusjonellSektorKode"`
	RegistrertIForetaksregisteret             bool                     `json:"registrertIForetaksregisteret"`
	RegistrertIStiftelsesregisteret           bool                     `json:"registrertIStiftelsesregisteret"`
	RegistrertIFrivillighetsregisteret        bool                     `json:"registrertIFrivillighetsregisteret"`
	SisteInnsendteAarsregnskap                string                   `json:"sisteInnsendteAarsregnskap"`
	Konkurs                                   bool                     `json:"konkurs"`
	UnderAvvikling                            bool                     `json:"underAvvikling"`
	UnderTvangsavviklingEllerTvangsopplosning bool                     `json:"underTvangsavviklingEllerTvangsopplosning"`
	Maalform                                  string                   `json:"maalform"`
	Links                                     []string                 `json:"links"`
}

type SubUnit struct {
	OrganisasjonsNummer             string            `json:"organisasjonsnummer"`
	Navn                            string            `json:"navn"`
	Organisasjonsform               Organisasjonsform `json:"organisasjonsform"`
	RegistreringsdatoEnhetsregister string            `json:"registreringsdatoEnhetsregister"`
	RegistrertIMvaRegisteret        bool              `json:"registrertIMvaregisteret"`
	Naeringskode1                   NaeringsKode1     `json:"naeringskode1"`
	AntallAnsatte                   int               `json:"antallAnsatte"`
	OverordnetEnhet                 string            `json:"overordnetEnhet"`
	OppstartsDato                   string            `json:"oppstartsDato"`
	Beliggenhetsadresse             Address           `json:"beliggenhetsadresse"`
	Links                           []string          `json:"links"`
}

type Organisasjonsform struct {
	Kode        string   `json:"kode"`
	Beskrivelse string   `json:"beskrivelse"`
	Links       []string `json:"links"`
}

type NaeringsKode1 struct {
	Kode        string `json:"kode"`
	Beskrivelse string `json:"beskrivelse"`
}

type NaeringsKode2 struct {
	Kode             string `json:"kode"`
	Beskrivelse      string `json:"beskrivelse"`
	HjelpeenhetsKode string `json:"hjelpeenhetskode"`
}

type Address struct {
	Land          string   `json:"land"`
	Landkode      string   `json:"landkode"`
	Postnummer    string   `json:"postnummer"`
	Poststed      string   `json:"poststed"`
	Adresse       []string `json:"adresse"`
	Kommune       string   `json:"kommune"`
	Kommunenummer string   `json:"kommunenummer"`
}

type InstitusjonellSektorKode struct {
	Kode        string `json:"kode"`
	Beskrivelse string `json:"beskrivelse"`
}
