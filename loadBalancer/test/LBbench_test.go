package test

import (
	"bytes"
	"loadBalancer/internal"
	"net/http"
	"os"
	"runtime/pprof"
	"sync"
	"testing"
	"time"
)

func BenchmarkLoadBalancer(b *testing.B) {
	cpuFile, err := os.Create("cpu.out")
	if err != nil {
		b.Fatalf("failed to create CPU profile: %v", err)
	}
	defer cpuFile.Close()

	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		b.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()

	client := &http.Client{}

	// Завантажте конфігурацію
	LBconf, err := internal.LoadLBConfig("../cmd/LBConfig.yaml")
	if err != nil {
		b.Fatalf("failed to load config: %v", err)
	}
	sh, err := internal.LoadSHConfig("../cmd/SGConfig.yaml")
	if err != nil {
		b.Fatalf("failed to load config: %v", err)
	}

	var SHconf []*internal.SHConfig
	SHconf = append(SHconf, sh...)

	// Ініціалізуйте інфраструктуру
	infra, err := internal.NewInfra(LBconf, SHconf)
	if err != nil {
		b.Fatal(err)
	}
	if infra == nil {
		b.Fatal("infra nil")
	}

	muxer := infra.ServeMUX()

	// Запустіть сервер у фоновій горутині
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if err := infra.LBserver.Run(muxer); err != nil {
			b.Fatal(err)
		}
	}()

	// Дайте серверу трохи часу на запуск
	time.Sleep(100 * time.Millisecond)

	b.ResetTimer()

	for i := 0; i < 100; i++ {
		req, err := http.NewRequest("GET", "http://localhost:8080/api", bytes.NewBufferString("lalala"))
		if err != nil {
			b.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			b.Fatal(err)
		}
		resp.Body.Close() // Закрийте тіло відповіді
	}
}
