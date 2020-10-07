package log

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

func BuildContextWithMask(next http.Handler, mask func(fieldName, s string) string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if fieldConfig.Map != nil && len(*fieldConfig.Map) > 0 {
			var v interface{}
			err := json.NewDecoder(r.Body).Decode(&v)
			if err != nil {
				next.ServeHTTP(w, r)
			} else {
				m, ok := v.(map[string]interface{})
				if !ok {
					next.ServeHTTP(w, r)
				} else {
					var ctx context.Context
					ctx = r.Context()
					for k, e := range *fieldConfig.Map {
						if strings.Index(e, ".") >= 0 {
							v3 := ValueOf(v, e)
							if v3 != nil {
								s3, ok3 := v3.(string)
								if ok3 {
									if len(s3) > 0 {
										if mask != nil && fieldConfig.Masks != nil && len(*fieldConfig.Masks) > 0 {
											if Include(*fieldConfig.Masks, k) {
												ctx = context.WithValue(ctx, k, mask(k, s3))
											} else {
												ctx = context.WithValue(ctx, k, s3)
											}
										} else {
											ctx = context.WithValue(ctx, k, s3)
										}
									}
								} else {
									ctx = context.WithValue(ctx, k, v3)
								}
							}
						} else {
							x, ok2 := m[e]
							if ok2 && x != nil {
								s3, ok3 := x.(string)
								if ok3 {
									if len(s3) > 0 {
										if mask != nil && fieldConfig.Masks != nil && len(*fieldConfig.Masks) > 0 {
											if Include(*fieldConfig.Masks, k) {
												ctx = context.WithValue(ctx, k, mask(k, s3))
											} else {
												ctx = context.WithValue(ctx, k, s3)
											}
										} else {
											ctx = context.WithValue(ctx, k, s3)
										}
									}
								} else {
									ctx = context.WithValue(ctx, k, x)
								}
							}
						}
					}
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		} else {
			next.ServeHTTP(w, r)
		}
	}
	return http.HandlerFunc(fn)
}
func BuildContext(next http.Handler) http.Handler {
	return BuildContextWithMask(next, nil)
}
func Include(vs []string, v string) bool {
	for _, s := range vs {
		if v == s {
			return true
		}
	}
	return false
}
func ValueOf(m interface{}, path string) interface{} {
	arr := strings.Split(path, ".")
	i := 0
	var c interface{}
	c = m
	l1 := len(arr) - 1
	for i < len(arr) {
		key := arr[i]
		m2, ok := c.(map[string]interface{})
		if ok {
			c = m2[key]
		}
		if !ok || i >= l1 {
			return c
		}
		i++
	}
	return c
}

func Mask(s string, start int, end int, mask string) string {
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start+end >= len(s) {
		return strings.Repeat(mask, len(s))
	}
	return s[:start] + strings.Repeat(mask, len(s)-start-end) + s[len(s)-end:]
}

func MaskMargin(s string, start int, end int, mask string) string {
	if start >= end {
		return ""
	}
	if start < 0 {
		start = 0
	}
	if end < 0 {
		end = 0
	}
	if start >= len(s) {
		return strings.Repeat(mask, len(s))
	}
	if end >= len(s) {
		return strings.Repeat(mask, start) + s[start:]
	}
	return strings.Repeat(mask, start) + s[start:end] + strings.Repeat(mask, len(s)-end)
}