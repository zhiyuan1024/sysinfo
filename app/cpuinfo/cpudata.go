package cpuinfo

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type cpuData struct {
	now       time.Time
	user      int
	system    int
	nice      int
	idel      int
	ioWait    int
	softIRQ   int
	hardIRQ   int
	isSet     bool
	statsFile string
	oldData   *cpuData
}

type cpuCorsData []cpuData

func NewCPUData() *cpuData {
	data := new(cpuData)
	data.statsFile = "/proc/stat"
	return data
}

func (data *cpuData) Sub() *cpuData {
	if !data.isSet {
		return nil
	}
	var sub cpuData
	sub.now = data.now
	sub.user = data.user - data.oldData.user
	sub.system = data.system - data.oldData.system
	sub.nice = data.nice - data.oldData.nice
	sub.idel = data.idel - data.oldData.idel
	sub.ioWait = data.ioWait - data.oldData.ioWait
	sub.softIRQ = data.softIRQ - data.oldData.softIRQ
	sub.hardIRQ = data.hardIRQ - data.oldData.hardIRQ
	total := sub.user + sub.system + sub.nice + sub.ioWait + sub.idel + sub.hardIRQ + sub.softIRQ
	if total != 0 {
		sub.user = sub.user * 10000 / total
		sub.nice = sub.nice * 10000 / total
		sub.system = sub.system * 10000 / total
		sub.idel = sub.idel * 10000 / total
		sub.ioWait = sub.ioWait * 10000 / total
		sub.softIRQ = sub.softIRQ * 10000 / total
		sub.hardIRQ = sub.hardIRQ * 10000 / total
	}
	return &sub
}

func (data *cpuData) Dump() (string, error) {
	sub := data.Sub()
	if sub == nil {
		return "", nil
	}
	strTime := sub.now.Format("20060102150405")
	line := fmt.Sprintf("cpu: %s %d %d %d %d %d %d %d\n",
		strTime, sub.user, sub.system, sub.nice,
		sub.idel, sub.ioWait, sub.hardIRQ, sub.softIRQ)
	*(data.oldData) = *data
	return line, nil
}

func (data *cpuData) Collecting() error {
	f, err := os.Open(data.statsFile)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Split(line, " ")[0] == "cpu" {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read cpuinfo error, err = %v", err)
	}
	toks := strings.Split(line, " ")
	d := data
	data.isSet = true
	if data.oldData == nil {
		data.oldData = new(cpuData)
		d = data.oldData
		data.isSet = false
	}
	d.now = time.Now()
	if d.user, err = strconv.Atoi(toks[2]); err != nil {
		return fmt.Errorf("get user usage error, err = %v", err)
	}
	if d.nice, err = strconv.Atoi(toks[3]); err != nil {
		return fmt.Errorf("get nice usage error, err = %v", err)
	}
	if d.system, err = strconv.Atoi(toks[4]); err != nil {
		return fmt.Errorf("get system usage error, err = %v", err)
	}
	if d.idel, err = strconv.Atoi(toks[5]); err != nil {
		return fmt.Errorf("get idel usage error, err = %v", err)
	}
	if d.ioWait, err = strconv.Atoi(toks[6]); err != nil {
		return fmt.Errorf("get iowait usage error, err = %v", err)
	}
	if d.hardIRQ, err = strconv.Atoi(toks[7]); err != nil {
		return fmt.Errorf("get hardirq usage error, err = %v", err)
	}
	if d.softIRQ, err = strconv.Atoi(toks[8]); err != nil {
		return fmt.Errorf("get softirq usage error, err = %v", err)
	}
	return nil
}
