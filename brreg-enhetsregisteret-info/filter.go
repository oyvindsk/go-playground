package main

import (
	"encoding/json"
	"fmt"
)


func main() {

	var jsonBlob = []byte(`[{"organisasjonsnummer":974761122,"navn":"STATENS LEGEMIDDELVERK","registreringsdatoEnhetsregisteret":"1995-08-09","organisasjonsform":"ORGL","hjemmeside":"www.legemiddelverket.no","registrertIFrivillighetsregisteret":"N","registrertIMvaregisteret":"N","registrertIForetaksregisteret":"N","registrertIStiftelsesregisteret":"N","antallAnsatte":250,"institusjonellSektorkode":{"kode":"6100","beskrivelse":"Statsforvaltningen"},"naeringskode1":{"kode":"84.120","beskrivelse":"Offentlig administrasjon tilknyttet helsestell, sosial virksomhet,undervisning, kirke, kultur og miljøvern"},"postadresse":{"adresse":"Postboks 63 Kalbakken","postnummer":"0901","poststed":"OSLO","kommunenummer":"0301","kommune":"OSLO","landkode":"NO","land":"Norge"},"forretningsadresse":{"adresse":"Sven Oftedals vei 6","postnummer":"0950","poststed":"OSLO","kommunenummer":"0301","kommune":"OSLO","landkode":"NO","land":"Norge"},"konkurs":"N","underAvvikling":"N","underTvangsavviklingEllerTvangsopplosning":"N","overordnetEnhet":983887406,"links":[{"rel":"self","href":"http://data.brreg.no/enhetsregisteret/enhet/974761122"},{"rel":"overordnetEnhet","href":"http://data.brreg.no/enhetsregisteret/enhet/983887406"}]},
        {"organisasjonsnummer":974767376,"navn":"ORKDAL KOMMUNE GRUNNSKOLE","stiftelsesdato":"1969-02-01","registreringsdatoEnhetsregisteret":"1997-02-20","organisasjonsform":"ORGL","registrertIFrivillighetsregisteret":"N","registrertIMvaregisteret":"N","registrertIForetaksregisteret":"N","registrertIStiftelsesregisteret":"N","antallAnsatte":250,"institusjonellSektorkode":{"kode":"6500","beskrivelse":"Kommuneforvaltningen"},"naeringskode1":{"kode":"84.120","beskrivelse":"Offentlig administrasjon tilknyttet helsestell, sosial virksomhet,undervisning, kirke, kultur og miljøvern"},"postadresse":{"adresse":"Postboks 83","postnummer":"7301","poststed":"ORKANGER","kommunenummer":"1638","kommune":"ORKDAL","landkode":"NO","land":"Norge"},"forretningsadresse":{"adresse":"Allfarveien 5","postnummer":"7300","poststed":"ORKANGER","kommunenummer":"1638","kommune":"ORKDAL","landkode":"NO","land":"Norge"},"konkurs":"N","underAvvikling":"N","underTvangsavviklingEllerTvangsopplosning":"N","overordnetEnhet":958731558,"links":[{"rel":"self","href":"http://data.brreg.no/enhetsregisteret/enhet/974767376"},{"rel":"overordnetEnhet","href":"http://data.brreg.no/enhetsregisteret/enhet/958731558"}]}]`)

    type Adr struct {
        Kommune string `json:"kommune"`
    }

    type Firma struct {
        Navn    string `json:"navn"`
        Ansatte int `json:"antallAnsatte"`
        Adrr    Adr `json:"forretningsadresse"`
    }


	var firmaer []Firma

	err := json.Unmarshal(jsonBlob, &firmaer)

	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("%+v", firmaer)


}
