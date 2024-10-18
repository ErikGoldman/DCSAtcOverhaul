package deepgramRecognizer

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	mp3 "github.com/braheezy/shine-mp3/pkg/mp3"

	interfaces "github.com/deepgram/deepgram-go-sdk/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/pkg/client/listen"

	"github.com/rs/zerolog/log"

	"github.com/dharmab/skyeye/pkg/recognizer"
)

type AtcDeepgramRecognizer struct {
	client *client.RESTClient
	apiKey string
}

func NewAtcDeepgramRecognizer(apiKey string) recognizer.Recognizer {
	client.Init(client.InitLib{
		LogLevel: client.LogLevelTrace, // LogLevelStandard / LogLevelFull / LogLevelTrace
	})

	c := client.NewREST(apiKey, &interfaces.ClientOptions{
		Host: "https://api.deepgram.com",
	})

	return &AtcDeepgramRecognizer{
		client: c,
		apiKey: apiKey,
	}
}

func (r *AtcDeepgramRecognizer) Debug_ReadFromWavFile(filename string) (string, error) {
	wavData, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read WAV file: %w", err)
	}

	text, err := r.RecognizeBytes(context.Background(), wavData)
	if err != nil {
		return "", fmt.Errorf("failed to recognize audio: %w", err)
	}
	return text, nil
}

func writeRawFile(pcm []float32) error {
	// Write raw PCM data to file
	rawFile, err := os.Create("incoming_pcm.raw")
	if err != nil {
		return fmt.Errorf("failed to create raw PCM file: %w", err)
	}
	defer rawFile.Close()

	// Convert float32 PCM data to strings and write them one per line
	for _, sample := range pcm {
		_, err := fmt.Fprintf(rawFile, "%f\n", sample)
		if err != nil {
			return fmt.Errorf("failed to write sample to file: %w", err)
		}
	}

	return nil
}

func readRawFile(filename string) ([]float32, error) {
	// Open the raw PCM file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open raw PCM file: %w", err)
	}
	defer file.Close()

	var pcm []float32
	var value float32

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		_, err := fmt.Sscanf(scanner.Text(), "%f", &value)
		if err != nil {
			return nil, fmt.Errorf("failed to parse float value: %w", err)
		}
		pcm = append(pcm, value)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return pcm, nil
}

func TestWriteWav() (string, error) {
	// Read the raw PCM file
	rawPCM, err := readRawFile("incoming_pcm.raw")
	if err != nil {
		return "", fmt.Errorf("failed to read raw PCM file: %w", err)
	}

	log.Info().Msgf("Read %d samples from raw PCM file", len(rawPCM), rawPCM[18424])

	// Convert PCM data to WAV format
	wavData, err := pcmToWav(16000, 16, rawPCM)
	if err != nil {
		return "", fmt.Errorf("failed to convert PCM to WAV: %w", err)
	}

	// Open the output WAV file
	wavFile, err := os.Create("output.wav")
	if err != nil {
		return "", fmt.Errorf("failed to create output WAV file: %w", err)
	}
	defer wavFile.Close()

	// Write the WAV data to the file
	_, err = wavFile.Write(wavData)
	if err != nil {
		return "", fmt.Errorf("failed to write WAV data to file: %w", err)
	}

	// Write Opus file
	opusFile, err := os.Create("output.mp3")
	if err != nil {
		return "", fmt.Errorf("failed to create output Opus file: %w", err)
	}
	defer opusFile.Close()

	// Convert PCM data to Opus format
	opusData, err := pcmToMp3(44100, 1, rawPCM)
	if err != nil {
		return "", fmt.Errorf("failed to convert PCM to Opus: %w", err)
	}

	// Write the Opus data to the file
	_, err = opusFile.Write(opusData)
	if err != nil {
		return "", fmt.Errorf("failed to write Opus data to file: %w", err)
	}

	return "Audio saved to output.wav", nil
}

func (r *AtcDeepgramRecognizer) RecognizeBytes(ctx context.Context, fileData []byte) (string, error) {
	log.Info().Msgf("Sending %d bytes to Deepgram recognition", len(fileData))

	// Prepare the request URL
	url := "https://api.deepgram.com/v1/listen?smart_format=false&language=en&model=nova-2"

	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(fileData))
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Token "+r.apiKey)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for successful status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code %d: %s", resp.StatusCode, string(body))
	}

	// Parse the JSON response
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Extract the transcription
	results, ok := result["results"].(map[string]interface{})

	log.Info().Msgf("Deepgram completed")

	if !ok {
		return "", fmt.Errorf("unexpected response format: missing 'results' field")
	}
	channels, ok := results["channels"].([]interface{})
	if !ok || len(channels) == 0 {
		return "", fmt.Errorf("unexpected response format: missing or empty 'channels' field")
	}
	alternatives, ok := channels[0].(map[string]interface{})["alternatives"].([]interface{})
	if !ok || len(alternatives) == 0 {
		return "", fmt.Errorf("unexpected response format: missing or empty 'alternatives' field")
	}
	transcript, ok := alternatives[0].(map[string]interface{})["transcript"].(string)
	if !ok {
		return "", fmt.Errorf("unexpected response format: missing 'transcript' field")
	}

	return transcript, nil
}

func (r *AtcDeepgramRecognizer) Recognize(ctx context.Context, pcm []float32, enableTranscriptionLogging bool) (string, error) {
	//writeRawFile(pcm)

	// Convert PCM data to WAV format
	log.Info().Msgf("Converting PCM to WAV")
	wavData, err := pcmToWav(16000, 16, pcm)
	if err != nil {
		return "", err
	}

	return r.RecognizeBytes(ctx, wavData)
}

func pcmToWav(sampleRate int, bitsPerSample int, pcm []float32) ([]byte, error) {
	// Convert float32 PCM to int16 samples
	samples := make([]int16, len(pcm))
	for i, sample := range pcm {
		samples[i] = int16(sample * 32767) // Scale to int16 range
	}

	// Create a buffer to hold the WAV file
	buf := new(bytes.Buffer)

	// Write WAV header
	binary.Write(buf, binary.LittleEndian, []byte("RIFF"))
	binary.Write(buf, binary.LittleEndian, int32(36+len(samples)*2)) // File size
	binary.Write(buf, binary.LittleEndian, []byte("WAVE"))
	binary.Write(buf, binary.LittleEndian, []byte("fmt "))
	binary.Write(buf, binary.LittleEndian, int32(16)) // Size of fmt chunk
	binary.Write(buf, binary.LittleEndian, int16(1))  // Audio format (PCM)
	binary.Write(buf, binary.LittleEndian, int16(1))  // Number of channels
	binary.Write(buf, binary.LittleEndian, int32(sampleRate))
	binary.Write(buf, binary.LittleEndian, int32(sampleRate*2)) // Byte rate
	binary.Write(buf, binary.LittleEndian, int16(2))            // Block align
	binary.Write(buf, binary.LittleEndian, int16(16))           // Bits per sample
	binary.Write(buf, binary.LittleEndian, []byte("data"))
	binary.Write(buf, binary.LittleEndian, int32(len(samples)*2)) // Size of data chunk

	// Write audio data
	for _, sample := range samples {
		binary.Write(buf, binary.LittleEndian, sample)
	}

	return buf.Bytes(), nil
}

func pcmToMp3(sampleRate int, channels int, pcm []float32) ([]byte, error) {
	enc := mp3.NewEncoder(sampleRate, channels)

	var buf bytes.Buffer

	// Convert float32 PCM to int16 samples
	samples := make([]int16, len(pcm))
	for i, sample := range pcm {
		samples[i] = int16(sample * 32767) // Scale to int16 range
	}
	// Encode the PCM data to MP3
	if err := enc.Write(&buf, samples); err != nil {
		return nil, fmt.Errorf("failed to encode PCM to MP3: %w", err)
	}

	return buf.Bytes(), nil
}

func extractPCMFromWAV(wavData []byte) ([]byte, error) {
	// Check if the file is long enough to be a valid WAV file
	if len(wavData) < 44 {
		return nil, fmt.Errorf("WAV data is too short")
	}

	// Verify the RIFF header
	if string(wavData[0:4]) != "RIFF" {
		return nil, fmt.Errorf("invalid WAV file: RIFF header not found")
	}

	// Verify the WAVE format
	if string(wavData[8:12]) != "WAVE" {
		return nil, fmt.Errorf("invalid WAV file: WAVE format not found")
	}

	// Find the data chunk
	dataStart := 0
	for i := 12; i < len(wavData)-8; i++ {
		if string(wavData[i:i+4]) == "data" {
			dataStart = i + 8
			break
		}
	}

	if dataStart == 0 {
		return nil, fmt.Errorf("invalid WAV file: data chunk not found")
	}

	// Extract the PCM data
	pcmData := wavData[dataStart:]

	return pcmData, nil
}
