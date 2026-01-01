package main

import (
	"fmt"
	"mangomon/config"
	"os"
)

func main() {
	content := `
# Some header
monitorrule=eDP-1,0.55,1,tile,0,1.00,0,0,1920,1080,60
# Comment
monitorrule=HDMI-A-1,0.50,1,tile,0,2.00,1920,0,3840,2160,144
`
	tmpFile := "test_config_v2.conf"
	os.WriteFile(tmpFile, []byte(content), 0644)
	defer os.Remove(tmpFile)

	parser, err := config.NewParser(tmpFile)
	if err != nil {
		panic(err)
	}

	rules, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d rules\n", len(rules))

	r1 := rules["eDP-1"]
	if r1.MFact != 0.55 {
		fmt.Printf("FAIL: eDP-1 MFact mismatch, got %f\n", r1.MFact)
	}
	if r1.NMaster != 1 {
		fmt.Printf("FAIL: eDP-1 NMaster mismatch, got %d\n", r1.NMaster)
	}

	r2 := rules["HDMI-A-1"]
	if r2.Scale != 2.0 {
		fmt.Printf("FAIL: HDMI-A-1 Scale mismatch, got %f\n", r2.Scale)
	}

	r2.X = 2000
	r2.Y = 100
	rules["HDMI-A-1"] = r2

	var saveRules []config.MonitorRule
	for _, v := range rules {
		saveRules = append(saveRules, v)
	}

	parser.Save(saveRules)

	bytes, _ := os.ReadFile(tmpFile)
	fmt.Println("--- Saved Content ---")
	fmt.Println(string(bytes))
}
