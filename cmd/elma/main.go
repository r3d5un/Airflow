package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	logger.Info("starting")
	fmt.Println("Hello, World!")

	logger.Info("downloading elma entries")
	resp, err := http.Get("https://hotell.difi.no/api/json/difi/elma/participants?")
	if err != nil {
		logger.Error("Error downloading elma entries", err)
		os.Exit(1)
	}

	logger.Info("unmarshalling elma entries")
	var elmaData ElmaResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&elmaData); err != nil {
		logger.Error("Error decoding elma entries", err)
		os.Exit(1)
	}

	fmt.Println(elmaData.Pages)

	// for _, entry := range elmaData.Entries {
	// 	fmt.Println(entry)
	// }
}

type ElmaResponseBody struct {
	Page    int     `json:"page"`
	Pages   int     `json:"pages"`
	Posts   int     `json:"posts"`
	Entries []Entry `json:"entries"`
}

type Entry struct {
	EHFPOSTAWARDG3P0630                   string `json:"EHF_POSTAWARD_G3_P06_3_0"`
	ArkivmeldingTekniskeTjenester         string `json:"Arkivmelding_TekniskeTjenester"`
	EHFPOSTAWARDG3P0230EO                 string `json:"EHF_POSTAWARD_G3_P02_3_0_EO"`
	PAYMENT09RESPONSE                     string `json:"PAYMENT_09_RESPONSE"`
	INNBYGGERPOSTFLYTTETPOST10            string `json:"INNBYGGERPOST_FLYTTETPOST_1_0"`
	PEPPOLBIS30ORDERINGEO                 string `json:"PEPPOLBIS_3_0_ORDERING_EO"`
	EHFPOSTAWARDG3P0130EO                 string `json:"EHF_POSTAWARD_G3_P01_3_0_EO"`
	EHFPOSTAWARDG3P0930EO                 string `json:"EHF_POSTAWARD_G3_P09_3_0_EO"`
	PAYMENT09                             string `json:"PAYMENT_09"`
	AvtaltAvtalt                          string `json:"Avtalt_avtalt"`
	EHFPOSTAWARDG3P0730                   string `json:"EHF_POSTAWARD_G3_P07_3_0"`
	PAYMENT10                             string `json:"PAYMENT_10"`
	INNBYGGERPOSTDIGITAL10                string `json:"INNBYGGERPOST_DIGITAL_1_0"`
	PEPPOLBISADVANCEDORDERING30CA         string `json:"PEPPOLBIS_ADVANCED_ORDERING_3_0_CA"`
	EInnsynResponse                       string `json:"eInnsyn_Response"`
	INNBYGGERPOSTUTSKRIFT10               string `json:"INNBYGGERPOST_UTSKRIFT_1_0"`
	Identifier                            string `json:"identifier"`
	PAYMENT04RESPONSE                     string `json:"PAYMENT_04_RESPONSE"`
	PEPPOLBIS30BILLING01UBL               string `json:"PEPPOLBIS_3_0_BILLING_01_UBL"`
	EHFADVORDER30PILOTRESPONSE            string `json:"EHF_ADVORDER_3_0_PILOT_RESPONSE"`
	EHFPOSTAWARDG3P0430                   string `json:"EHF_POSTAWARD_G3_P04_3_0"`
	PEPPOLBIS30INVOICERESPONSEEO          string `json:"PEPPOLBIS_3_0_INVOICE_RESPONSE_EO"`
	ArkivmeldingAdministrasjon            string `json:"Arkivmelding_Administrasjon"`
	ArkivmeldingHelseSosialOgOmsorg       string `json:"Arkivmelding_HelseSosialOgOmsorg"`
	AvtaltResponse                        string `json:"Avtalt_Response"`
	PAYMENT10RESPONSE                     string `json:"PAYMENT_10_RESPONSE"`
	EHFPOSTAWARDG3P0530                   string `json:"EHF_POSTAWARD_G3_P05_3_0"`
	PEPPOLBIS30CATALOGUEEO                string `json:"PEPPOLBIS_3_0_CATALOGUE_EO"`
	PEPPOLBIS30ORDERONLYEO                string `json:"PEPPOLBIS_3_0_ORDER_ONLY_EO"`
	Name                                  string `json:"name"`
	EHFPOSTAWARDG3P0930CA                 string `json:"EHF_POSTAWARD_G3_P09_3_0_CA"`
	PAYMENT03RESPONSE                     string `json:"PAYMENT_03_RESPONSE"`
	ArkivmeldingTaushetsbelagt            string `json:"Arkivmelding_Taushetsbelagt"`
	DigitalpostInfo                       string `json:"Digitalpost_Info"`
	ArkivmeldingKulturIdrettOgFritid      string `json:"Arkivmelding_KulturIdrettOgFritid"`
	EHFPOSTAWARDG3P0130CA                 string `json:"EHF_POSTAWARD_G3_P01_3_0_CA"`
	ArkivmeldingResponse                  string `json:"Arkivmelding_Response"`
	PEPPOLBIS30CATALOGUECA                string `json:"PEPPOLBIS_3_0_CATALOGUE_CA"`
	EInnsynJournalpost                    string `json:"eInnsyn_Journalpost"`
	PEPPOLBISADVANCEDORDERING30EO         string `json:"PEPPOLBIS_ADVANCED_ORDERING_3_0_EO"`
	ArkivmeldingNaturOgMiljoe             string `json:"Arkivmelding_NaturOgMiljoe"`
	PEPPOLBIS30DESPATCHADVICECA           string `json:"PEPPOLBIS_3_0_DESPATCH_ADVICE_CA"`
	ArkivmeldingNaeringsutvikling         string `json:"Arkivmelding_Naeringsutvikling"`
	PAYMENT02RESPONSE                     string `json:"PAYMENT_02_RESPONSE"`
	PEPPOLBIS30ORDERINGCA                 string `json:"PEPPOLBIS_3_0_ORDERING_CA"`
	EHFPOSTAWARDG3P0330                   string `json:"EHF_POSTAWARD_G3_P03_3_0"`
	PEPPOLBIS30ORDERAGREEMENTCA           string `json:"PEPPOLBIS_3_0_ORDER_AGREEMENT_CA"`
	DigitalpostVedtak                     string `json:"Digitalpost_Vedtak"`
	PEPPOLBIS30BILLING01CII               string `json:"PEPPOLBIS_3_0_BILLING_01_CII"`
	ArkivmeldingByggOgGeodata             string `json:"Arkivmelding_ByggOgGeodata"`
	EGOVGE10REQUEST                       string `json:"EGOV_GE_1_0_REQUEST"`
	ArkivmeldingOppvekstOgUtdanning       string `json:"Arkivmelding_OppvekstOgUtdanning"`
	EHFPOSTAWARDG3P0830                   string `json:"EHF_POSTAWARD_G3_P08_3_0"`
	EInnsynInnsynskrav                    string `json:"eInnsyn_Innsynskrav"`
	Icd                                   string `json:"Icd"`
	PAYMENT01RESPONSE                     string `json:"PAYMENT_01_RESPONSE"`
	PAYMENT04                             string `json:"PAYMENT_04"`
	EHFADVORDER30PILOTREQUEST             string `json:"EHF_ADVORDER_3_0_PILOT_REQUEST"`
	PAYMENT03                             string `json:"PAYMENT_03"`
	PAYMENT02                             string `json:"PAYMENT_02"`
	EGOVGE10RESPONSE                      string `json:"EGOV_GE_1_0_RESPONSE"`
	PAYMENT01                             string `json:"PAYMENT_01"`
	Regdate                               string `json:"regdate"`
	ArkivmeldingTrafikkReiserOgSamferdsel string `json:"Arkivmelding_TrafikkReiserOgSamferdsel"`
	ArkivmeldingSkatterOgAvgifter         string `json:"Arkivmelding_SkatterOgAvgifter"`
	EHFPOSTAWARDG3P0230CA                 string `json:"EHF_POSTAWARD_G3_P02_3_0_CA"`
}

type ElmaEntries []Entry
