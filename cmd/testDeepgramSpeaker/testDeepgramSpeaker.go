package main

import (
	"encoding/json"
	"fmt"
	"os"

	deepgramspeaker "github.com/ErikGoldman/DCSAtcOverhaul/pkg/deepgramSpeaker"
)

func main() {
	// Read and parse the config.json file
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		panic("Failed to read config.json")
	}

	var configData struct {
		Deepgram struct {
			APIKey string `json:"api_key"`
		} `json:"deepgram"`
	}

	err = json.Unmarshal(configFile, &configData)
	if err != nil {
		panic("Failed to parse config.json")
	}

	synth := deepgramspeaker.NewSpeechSynthesizer(configData.Deepgram.APIKey)

	bytesChannel := make(chan []byte)
	err = synth.GenerateSpeech("aura-stella-en", "alpha one-one tower. loud and clear!", bytesChannel)
	if err != nil {
		panic(fmt.Sprintf("error generating speech: %s", err))
	}

	file, err := os.Create("output-deepgramtest.raw")
	if err != nil {
		panic(fmt.Sprintf("error creating file: %s", err))
	}
	defer file.Close()

	for {
		select {
		case data := <-bytesChannel:
			if data == nil {
				fmt.Printf("Got nil, exiting")
				return
			}

			fmt.Printf("Writing to file")
			_, err = file.Write(data)
			if err != nil {
				panic(fmt.Sprintf("error writing to file: %s", err))
			}
		}
	}
}
