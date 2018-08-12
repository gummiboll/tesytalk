package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gummiboll/tesytalk/pkg/tesy"
)

func main() {
	var err error
	var timeSinceReset string
	addressFlag := flag.String("address", "", "Address of a tesy waterheater")
	verboseFlag := flag.Bool("verbose", false, "Print verbose information about the device")
	resetPowerUsageFlag := flag.Bool("resetpow", false, "Resets kWh calculator")
	setTempFlag := flag.Int("settemp", 0, "Set temperature, 1-70 are accepted")
	flag.Parse()

	if *addressFlag == "" {
		fmt.Println("Adress is required")
		os.Exit(1)
	}

	tesy := tesy.Tesy{Address: *addressFlag}

	if *resetPowerUsageFlag {
		fmt.Println("Resetting kWh calculator..")
		err = tesy.ResetPowerCalc()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(".. done!")
		os.Exit(0)
	}

	if *setTempFlag > 0 {
		if *setTempFlag > 70 {
			fmt.Println("Only values between 1-70 are accepted")
			os.Exit(1)
		} else {
			fmt.Println("Setting temperature..")
			err = tesy.SetTemperature(*setTempFlag)
			if err != nil {
				fmt.Println(fmt.Sprintf("Failed to set temperature to %d: %s", *setTempFlag, err))
				os.Exit(1)
			}
			fmt.Println(".. done!")
			os.Exit(0)
		}
	}

	if *verboseFlag {
		err = tesy.GetDevStatus()
		if err != nil {
			fmt.Println(fmt.Sprintf("Failed to communicate with your tesy device: %s", err))
			os.Exit(1)
		}
	}

	tesy.Status, err = tesy.GetStatus()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to retreive status for your tesy device: %s", err))
		os.Exit(1)
	}

	tesy.PowerUsage, err = tesy.GetCalcRes()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to retreive status for your tesy device: %s", err))
		os.Exit(1)
	}

	hoursSinceReset, err := tesy.TimeSinceReset()
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to calculate time since reset %s", err))
		os.Exit(1)
	}

	switch {
	case hoursSinceReset < 1:
		timeSinceReset = "Less then a hour"
	case hoursSinceReset > 24:
		timeSinceReset = fmt.Sprintf("%.1f days (%.1f hours)", hoursSinceReset/24, hoursSinceReset)
	default:
		timeSinceReset = fmt.Sprintf("%.1f hours", hoursSinceReset)
	}

	fmt.Println(fmt.Sprintf("Heater state:        %s", tesy.Status.HeaterState))
	fmt.Println(fmt.Sprintf("Current temperature: %.1f\u2103/%.1f\u2103", tesy.Status.Gradus, tesy.Status.RefGradus))
	fmt.Println(fmt.Sprintf("Power usage:         %.2f kWh since last reset (%s ago)", tesy.PowerUsage.KWH(), timeSinceReset))

	if *verboseFlag {
		fmt.Println(fmt.Sprintf("Heater capacity:     %d litres", tesy.Status.Volume))
		fmt.Println(fmt.Sprintf("DevID:               %s", tesy.DevID))
		fmt.Println(fmt.Sprintf("MacAddr:             %s", tesy.MacAddr))
		fmt.Println(fmt.Sprintf("Watts:               %d", tesy.Status.Watts))
	}
}
