package givenergy

import (
	"encoding/json"
	"github.com/dazfuller/azure-functions-golang-demo/internal/givenergy"
	"log"
	"net/http"
)

type GivEnergyManager struct {
	ApiKey string
}

func (ge *GivEnergyManager) GivEnergyHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("GivEnergy handler called")
	api := givenergy.NewAccountApi(ge.ApiKey)

	devices := api.GetCommunicationDevices()
	serialNumber := devices.Data[0].Inverter.Serial

	log.Printf("Collecting data for %s", serialNumber)
	latestData := api.GetLatestInverterData(serialNumber)
	summarisedData := latestData.Summarise()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summarisedData)
}
