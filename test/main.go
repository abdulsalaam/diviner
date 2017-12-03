package main

import (
	"diviner/fabsdk"
	"log"
)

func main() {
	sdk, err := fabsdk.NewSDK("fabric_config.yaml")
	if err != nil {
		log.Fatalln(err)
	}

	session, err := sdk.NewPreEnrolledUserSession("Diviner", "User1")
	if err != nil {
		log.Fatalln(err)
	}

	sc, err := sdk.NewSystemClient(session)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = fabsdk.GetChannel(sc, "divinerchannel", []string{"Diviner"})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("end")

}
