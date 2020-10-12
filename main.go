package main

import (
  "os"
  "io"
  "fmt"
  "bufio"
  "regexp"
  "errors"
  "strings"
  "strconv"
)

type Station struct {
  MAC string
  RXBytes int64
  RXPackets int64
  TXBytes int64
  TXPackets int64
  TXFailed int64
  Signal int64
  TXBitrate int64
  RXBitrate int64
  InactiveTime int64
  ConnectedTime int64
}

func readInput() ([]string, error)  {
  info, err := os.Stdin.Stat()
  if err != nil {
    return []string{}, err
  }

  if info.Mode() & os.ModeCharDevice != 0 {
    return []string{}, errors.New("No Input")
  }

  reader := bufio.NewReader(os.Stdin)
  var output []rune

  for {
    input, _, err := reader.ReadRune()
    if err != nil && err == io.EOF {
      break
    }
    output = append(output, input)
  }

  var lines []string
  var lineBuilder strings.Builder
  for _, tmpRune := range output {
    if tmpRune == '\n' {
      lines = append(lines, lineBuilder.String())
      lineBuilder.Reset()
      continue
    }
    lineBuilder.WriteRune(tmpRune)
  }

  return lines, nil
}

func getNumber(line string, numberRegex *regexp.Regexp) int64 {
  rawNumbers := numberRegex.FindString(line)
  number, err := strconv.ParseInt(rawNumbers, 10, 64)
  if err != nil {
    return 0
  }
  return number
}

func getBitrate(line string, floatRegex, bitRateRegex *regexp.Regexp) int64 {
  rawFloats := floatRegex.FindString(line)
  rawFloat, err := strconv.ParseFloat(rawFloats, 64)
  if err != nil {
    return 0
  }
  rawUnit := bitRateRegex.FindString(line)
  if len(rawUnit) > 5 {
    var multiplier = 1
    if rawUnit[0] == 'K' {
      multiplier = 1000
    } else if rawUnit[0] == 'M' {
      multiplier = 1000000
    } else if rawUnit[0] == 'G' {
      multiplier = 1000000000
    }

    rawFloat *= float64(multiplier)
  }
  return int64(rawFloat)
}

func getLineTime(line string, numberRegex *regexp.Regexp) int64 {
  rawDurationString := numberRegex.FindString(line)
  rawDuration, err := strconv.ParseInt(rawDurationString, 10, 64)
  if err != nil {
    return 0
  }

  multiplier := float64(1)
  if strings.Index(line, "ms") >= 0 {
    multiplier = 0.001
  }
  return int64(float64(rawDuration) * multiplier)
}

func convertToStations(lines []string) []Station {
  macAddressRegex := regexp.MustCompile(`([a-z,0-9]{2}:?){6}`)
  numberRegex := regexp.MustCompile(`[+,-]?[0-9]+`)
  floatRegex := regexp.MustCompile(`[+,-]?[0-9]+\.[0-9]+`)
  bitRateRegex := regexp.MustCompile(`[K,M,G]*Bit\/s`)

  result := make([]Station, 0)
  var tmpStation Station
  for i, line := range lines {
    // Get the current MAC address
    if stationStart := strings.Index(line, "Station"); stationStart >= 0 {
      if i > 0 {
        result = append(result, tmpStation)
      }
      address := macAddressRegex.FindString(line)
      tmpStation.MAC = address
      continue
    }
    // Get the current RX-Bytes
    if strings.Index(line, "rx bytes:") >= 0 {
      tmpStation.RXBytes = getNumber(line, numberRegex)
      continue
    }
    // Get the current RX-Packets
    if strings.Index(line, "rx packets:") >= 0 {
      tmpStation.RXPackets = getNumber(line, numberRegex)
      continue
    }
    // Get the current TX-Bytes
    if strings.Index(line, "tx bytes:") >= 0 {
      tmpStation.TXBytes = getNumber(line, numberRegex)
      continue
    }
    // Get the current TX-Packets
    if strings.Index(line, "tx packets:") >= 0 {
      tmpStation.TXPackets = getNumber(line, numberRegex)
      continue
    }
    // Get the current TX-Failed
    if strings.Index(line, "tx failed:") >= 0 {
      tmpStation.TXFailed = getNumber(line, numberRegex)
      continue
    }
    // Get the current Signal
    if strings.Index(line, "signal:") >= 0 {
      tmpStation.Signal = getNumber(line, numberRegex)
      continue
    }
    // Get the current TX-Bitrate
    if strings.Index(line, "tx bitrate:") >= 0 {
      tmpStation.TXBitrate = getBitrate(line, floatRegex, bitRateRegex)
      continue
    }
    // Get the current RX-Bitrate
    if strings.Index(line, "rx bitrate:") >= 0 {
      tmpStation.RXBitrate = getBitrate(line, floatRegex, bitRateRegex)
      continue
    }
    // Get the current Inactive-Time
    if strings.Index(line, "inactive time:") >= 0 {
      tmpStation.InactiveTime = getLineTime(line, numberRegex)
      continue
    }
    // Get the current Connected-Time
    if strings.Index(line, "connected time:") >= 0 {
      tmpStation.ConnectedTime = getLineTime(line, numberRegex)
      continue
    }
  }
  if len(tmpStation.MAC) > 0 {
    result = append(result, tmpStation)
  }

  return result
}

func stationsToMetrics(stations []Station) string {
  var metricsBuilder strings.Builder

  metricsBuilder.WriteString("# TYPE wifi_received_bytes_total counter\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_received_bytes_total{mac=\"%s\"} %d\n", station.MAC, station.RXBytes))
  }

  metricsBuilder.WriteString("# TYPE wifi_received_packets_total counter\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_received_packets_total{mac=\"%s\"} %d\n", station.MAC, station.RXPackets))
  }

  metricsBuilder.WriteString("# TYPE wifi_received_bitrate gauge\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_received_bitrate{mac=\"%s\"} %d\n", station.MAC, station.RXBitrate))
  }

  metricsBuilder.WriteString("# TYPE wifi_transmitted_bytes_total counter\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_transmitted_bytes_total{mac=\"%s\"} %d\n", station.MAC, station.TXBytes))
  }

  metricsBuilder.WriteString("# TYPE wifi_transmitted_packets_total counter\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_transmitted_packets_total{mac=\"%s\"} %d\n", station.MAC, station.TXPackets))
  }

  metricsBuilder.WriteString("# TYPE wifi_transmitted_failed_total counter\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_transmitted_failed_total{mac=\"%s\"} %d\n", station.MAC, station.TXFailed))
  }

  metricsBuilder.WriteString("# TYPE wifi_transmitted_bitrate gauge\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_transmitted_bitrate{mac=\"%s\"} %d\n", station.MAC, station.TXBitrate))
  }

  metricsBuilder.WriteString("# TYPE wifi_signal gauge\n")
  for _, station := range stations {
    metricsBuilder.WriteString(fmt.Sprintf("wifi_signal{mac=\"%s\"} %d\n", station.MAC, station.Signal))
  }

  return metricsBuilder.String()
}

func main() {
  lines, err := readInput()
  if err != nil {
    panic(err)
  }

  stations := convertToStations(lines)

  metricsResult := stationsToMetrics(stations)

  fmt.Printf("%s", metricsResult)
}
