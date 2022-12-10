package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

func getTempoHandler(client mqtt.Client, msg mqtt.Message) {
	publishColors()
}

var (
	today        string
	tomorrow     string
	client       mqtt.Client
	tempoURL     string
	mqttURL      string
	mqttLogin    string
	mqttPassword string
)

type tempoAPIAnswer struct {
	CouleurJourJ  string `json:"couleurJourJ"`
	CouleurJourJ1 string `json:"couleurJourJ1"`
}

func publishColors() {
	token := client.Publish("mqtt-tempo/today", 2, false, today)
	token.Wait()
	log.Debugf("sent mqtt-tempo/today with %s", today)
	token = client.Publish("mqtt-tempo/tomorrow", 2, false, tomorrow)
	token.Wait()
	log.Debugf("sent mqtt-tempo/tomorrow with %s", tomorrow)
}

func publishWeKnowWhatTomorrowWillBe(t string) {
	token := client.Publish("mqtt-tempo/alarm", 2, false, t)
	token.Wait()
	log.Debugf("sent mqtt-tempo/alarm with %s", t)
}

func updateColors() {
	go func() {
		for {
			var dataFromEDF tempoAPIAnswer
			now := time.Now()
			date := now.Format("2006-01-02")
			url := fmt.Sprintf(tempoURL, date)
			log.Debugf("url: %s", url)
			client := &http.Client{}
			request, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Fatalf("udpateColors(): %v", err)
			}
			request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:70.0) Gecko/20100101 Firefox/70.0")
			resp, err := client.Do(request)
			if err != nil {
				log.Fatalf("udpateColors(): %v", err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("udpateColors(): %v", err)
			}
			json.Unmarshal(body, &dataFromEDF)
			if tomorrow == "ND" && dataFromEDF.CouleurJourJ != "ND" {
				publishWeKnowWhatTomorrowWillBe(dataFromEDF.CouleurJourJ)
			}
			if tomorrow != dataFromEDF.CouleurJourJ1 || today != dataFromEDF.CouleurJourJ {
				today = dataFromEDF.CouleurJourJ
				tomorrow = dataFromEDF.CouleurJourJ1
				publishColors()
			}
			log.Debugf("today is %s and tomorrow is %s", today, tomorrow)
			time.Sleep(10 * time.Minute)
		}
	}()

}

func main() {
	today = "ND"
	tomorrow = "ND"
	tempoURL = os.Getenv("TEMPO_URL")
	mqttURL = os.Getenv("TEMPO_MQTT_URL")
	mqttLogin = os.Getenv("TEMPO_MQTT_LOGIN")
	mqttPassword = os.Getenv("TEMPO_MQTT_PASSWORD")
	if os.Getenv("DEBUG") != "" {
		log.SetLevel(log.DebugLevel)
		log.Debugf("Debug mode")
	}
	co := mqtt.NewClientOptions().AddBroker(mqttURL)
	co.SetClientID("mqtt-tempo")
	co.SetPassword(mqttPassword)
	co.SetUsername(mqttLogin)
	co.SetAutoReconnect(true)
	client = mqtt.NewClient(co)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("main(): %v", token.Error())
	}

	if token := client.Subscribe("mqtt-tempo/get", 0, getTempoHandler); token.Wait() && token.Error() != nil {
		log.Fatalf("main(): %v", token.Error())
	}

	updateColors()

	// waiting for sigint or sigterm to stop
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go func() {
		<-sigs
		done <- true
	}()
	<-done
}
