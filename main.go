package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

// Video represents the JSON structure defining a video download
type Video struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run main.go <arquivo.json>")
		fmt.Println("Exemplo: go run main.go videos.json")
		return
	}

	jsonFileName := os.Args[1]

	// Lê o conteúdo do arquivo JSON
	fileContent, err := os.ReadFile(jsonFileName)
	if err != nil {
		fmt.Printf("Erro ao ler o arquivo %s: %v\n", jsonFileName, err)
		return
	}

	// Decodifica o JSON para o array de objetos Video
	var videos []Video
	if err := json.Unmarshal(fileContent, &videos); err != nil {
		fmt.Printf("Erro ao decodificar JSON: %v\n", err)
		return
	}

	// Verifica se a pasta "downloads" existe; se não, cria
	downloadDir := "downloads"
	if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
		fmt.Printf("Erro ao criar o diretório de downloads: %v\n", err)
		return
	}

	// Usamos um WaitGroup para sincronizar os downloads que executarão em paralelo
	var wg sync.WaitGroup

	for _, video := range videos {
		wg.Add(1)
		// Lança o download como uma goroutine para que sejam simultâneos
		go downloadVideo(video, downloadDir, &wg)
	}

	// Aguarda a finalização de todas as goroutines
	wg.Wait()
	fmt.Println("Todos os downloads foram concluídos!")
}

// downloadVideo invoca o ffmpeg para baixar o vídeo de fato a partir do m3u8
func downloadVideo(v Video, dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Garantir que o nome do arquivo tenha uma extensão .mp4 se não tiver
	fileName := v.Name
	if filepath.Ext(fileName) == "" {
		fileName += ".mp4"
	} else if filepath.Ext(fileName) == ".m3u8" {
	    // Se o usuário colocou .m3u8 no JSON de propósito, trocamos para .mp4
	    fileName = fileName[:len(fileName)-len(".m3u8")] + ".mp4"
	}

	filePath := filepath.Join(dir, fileName)
	fmt.Printf("Iniciando o download do vídeo: %s...\n", fileName)

	// O comando do ffmpeg para copiar o stream de vídeo e áudio do m3u8.
	// Usaremos -c copy para baixar e realizar o remux para mp4
	cmd := exec.Command("ffmpeg", "-i", v.URL, "-c", "copy", "-y", "-bsf:a", "aac_adtstoasc", filePath)

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Erro ao baixar o vídeo %s. Verifique se o ffmpeg está instalado ou se a URL é válida. Detalhes: %v\n", fileName, err)
		return
	}

	fmt.Printf("Download concluído com sucesso: %s\n", filePath)
}
