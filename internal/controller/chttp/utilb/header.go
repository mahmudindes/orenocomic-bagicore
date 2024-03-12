package utilb

import (
	"slices"
	"sort"
	"strconv"
	"strings"
)

func HeaderAccept(h string) []string {
	type hav struct {
		t string
		q float64
	}

	tt, t0, t1 := make([]string, 0), make([]string, 0), make([]hav, 0)

	t0p := func() {
		sort.SliceStable(t0, func(i, j int) bool {
			bi, ai, _ := strings.Cut(t0[i], "/")
			bj, aj, _ := strings.Cut(t0[j], "/")
			if bi != "*" && bj == "*" {
				return true
			}
			if ai != "*" && aj == "*" {
				return true
			}
			return false
		})

		tt = append(tt, t0...)
		clear(t0)
	}

	for _, t := range strings.Split(h, ",") {
		ct, op, kp := strings.Cut(t, ";")

		if ct == "" || !strings.Contains(ct, "/") {
			continue
		}

		if kp && op != "" {
			tv, ok := strings.CutPrefix(op, "q=")
			if !ok || len(tv) > 3 {
				continue
			}

			if tv == "" {
				t0 = append(t0, strings.TrimSpace(ct))
				continue
			}

			ti, _ := strconv.ParseFloat(tv, 32)
			switch {
			case ti < 1.0:
				t1 = append(t1, hav{strings.TrimSpace(ct), ti})
			case ti == 1.0:
				t0 = append(t0, strings.TrimSpace(ct))
			}

			continue
		}

		if len(t0) > 0 {
			t0p()
		}

		tt = append(tt, strings.TrimSpace(ct))
	}

	if len(t0) > 0 {
		t0p()
	}

	if len(t1) > 0 {
		sort.SliceStable(t1, func(i, j int) bool {
			if t1[j].q <= t1[i].q {
				bi, ai, _ := strings.Cut(t1[i].t, "/")
				bj, aj, _ := strings.Cut(t1[j].t, "/")
				if bi != "*" && bj == "*" {
					return true
				}
				if ai != "*" && aj == "*" {
					return true
				}
			}
			return t1[i].q > t1[j].q
		})

		for _, t := range t1 {
			tt = append(tt, t.t)
		}
	}

	return tt
}

func HeaderAcceptFirst(v string, st ...string) string {
	for _, t := range HeaderAccept(v) {
		bf, af, ok := strings.Cut(t, "/")
		if !ok {
			continue
		}

		if af == "*" {
			if bf == "*" {
				return st[0]
			}

			for _, s := range st {
				if bf == strings.Split(s, "/")[0] {
					return s
				}
			}
		}

		if slices.Contains(st, t) {
			return t
		}
	}
	return ""
}
