package ddcp

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func ddCmdOpes(opes map[string]string) string {
	buf := bytes.NewBufferString("dd")

	for k, v := range opes {
		buf.WriteString(" ")
		buf.WriteString(k)
		buf.WriteString("=")
		buf.WriteString(v)
	}

	return buf.String()
}

func ddOpes(source string, dest string, bs string, count string) map[string]string {
	return map[string]string{
		"if":    source,
		"of":    dest,
		"conv":  "notrunc",
		"bs":    bs,
		"count": count}
}

func ddCmds(source string, dest string, chunk_size int64, chunk_num int64) (cmds []string) {
	cmds = make([]string, chunk_num)
	chunk_size_mb := chunk_size / (1024 * 1024)

	for i := int64(0); i < chunk_num; i++ {
		opes := ddOpes(source, dest, "1m", strconv.FormatInt(chunk_size_mb, 10))

		if i > 0 {
			offset := strconv.FormatInt(chunk_size_mb*i, 10)
			opes["skip"] = offset
			opes["seek"] = offset
		}

		cmds[i] = ddCmdOpes(opes)
	}

	return
}

func runCmd(cmd string, ch chan error) {
	out, err := exec.Command("sh", "-c", cmd).CombinedOutput()

	if err != nil {
		err = fmt.Errorf("'%s' is failed: %s", cmd, out)
	}

	ch <- err
}

func runCmds(cmds []string) (err error) {
	ch := make(chan error)

	for _, cmd := range cmds {
		go runCmd(cmd, ch)
	}

	for _ = range cmds {
		err = <-ch

		if err != nil {
			return
		}
	}

	return
}

func Ddcp(params *DdcpParams) error {
	src, src_err := os.Stat(params.source)

	if src_err != nil {
		return fmt.Errorf("source file does not exist: %s", params.source)
	}

	src_size := src.Size()

	if src_size == 0 {
		return fmt.Errorf("source file is empty: %s", params.source)
	}

	_, dst_err := os.Stat(params.dest)

	if dst_err == nil {
		return fmt.Errorf("dest file already exists: %s", params.dest)
	}

	remainder := src_size % params.chunk_size
	chunk_num := src_size / params.chunk_size

	cmds := ddCmds(params.source, params.dest, params.chunk_size, chunk_num)

	if remainder > 0 {
		opes := ddOpes(params.source, params.dest, strconv.FormatInt(remainder, 10), "1")
		offset := strconv.FormatInt(params.chunk_size*chunk_num, 10)
		opes["skip"] = offset
		opes["seek"] = offset
		cmds = append(cmds, ddCmdOpes(opes))
	}

	return runCmds(cmds)
}
