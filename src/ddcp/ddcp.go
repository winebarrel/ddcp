package ddcp

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
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
			"bs":   "1047552"}

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

func Ddcp(params *DdcpParams) (err error) {
	src, src_err := os.Stat(params.source)

	if src_err != nil {
		err = fmt.Errorf("source file does not exist: %s", params.source)
		return
	}

	src_size := src.Size()

	_, dst_err := os.Stat(params.dest)

	if dst_err == nil {
		err = fmt.Errorf("dest file already exists: %s", params.dest)
		return
	}

	if src_size == 0 {
		out, cp_err := exec.Command("cp", params.source, params.dest).CombinedOutput()

		if cp_err != nil {
			err = fmt.Errorf("'cp %s %s' is failed: %s", params.source, params.dest, out)
			return
		}
	} else {
		remainder := src_size % params.chunk_size
		chunk_num := src_size / params.chunk_size

		opes_list := ddOpesList(params.source, params.dest, params.chunk_size, chunk_num, remainder)
		err = runCmds(opes_list)

		if err != nil {
			return
		}
	}

	if params.preserve {
		err = os.Chmod(params.dest, src.Mode())

		if err != nil {
			return
		}

		uid := src.Sys().(*syscall.Stat_t).Uid
		gid := src.Sys().(*syscall.Stat_t).Gid
		err = os.Chown(params.dest, int(uid), int(gid))

		if err != nil {
			return
		}
	}

	return
}
