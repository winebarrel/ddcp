package ddcp

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func ddOpes(opes map[string]string) (cmd_opes []string) {
	cmd_opes = []string{}

	for k, v := range opes {
		cmd_opes = append(cmd_opes, fmt.Sprintf("%s=%s", k, v))
	}

	return
}

func ddOpesList(source string, dest string, chunk_size int64, chunk_num int64, remainder int64) (opes_list [][]string) {
	if remainder > 0 {
		chunk_num++
	}

	opes_list = make([][]string, chunk_num)
	chunk_size_mb := chunk_size / (1024 * 1024)

	for i := int64(0); i < chunk_num; i++ {
		opes := map[string]string{
			"if":   source,
			"of":   dest,
			"conv": "notrunc",
			"bs":   "1m"}

		if i < chunk_num-1 {
			opes["count"] = strconv.FormatInt(chunk_size_mb, 10)
		}

		if i > 0 {
			offset := strconv.FormatInt(chunk_size_mb*i, 10)
			opes["skip"] = offset
			opes["seek"] = offset
		}

		opes_list[i] = ddOpes(opes)
	}

	return
}

func dd(opes []string, ch chan error) {
	out, err := exec.Command("dd", opes...).CombinedOutput()

	if err != nil {
		err = fmt.Errorf("'dd %s' is failed: %s", opes, out)
	}

	ch <- err
}

func runCmds(opes_list [][]string) (err error) {
	ch := make(chan error)

	for _, opes := range opes_list {
		go dd(opes, ch)
	}

	for _ = range opes_list {
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

	opes_list := ddOpesList(params.source, params.dest, params.chunk_size, chunk_num, remainder)

	return runCmds(opes_list)
}
