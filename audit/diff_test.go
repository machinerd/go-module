package audit_test

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"testing"

	"github.com/machinerd/go-module/audit"
	"github.com/wI2L/jsondiff"
)

// 실제로 틀린 이력이 저장된 경우 (komc-api-machine 기업 22367, patch v24)
// 배열에서 SOLUTIONS 를 지우면서 DISTRIBUTION, SUBCONTRACTING 순서가 바뀌었다.
func TestDiffSnapshotsDeleteAndReorder(t *testing.T) {
	old := map[string]any{"companyTypes": []any{"FACILITIES", "SOLUTIONS", "SUBCONTRACTING", "DISTRIBUTION"}}
	new := map[string]any{"companyTypes": []any{"FACILITIES", "DISTRIBUTION", "SUBCONTRACTING"}}

	// jsondiff 버그 재현 (Factorize + LCS): 적용 결과가 수정 후 상태와 다르다
	buggy, _ := jsondiff.Compare(old, new, jsondiff.Factorize(), jsondiff.LCS())
	if err := audit.VerifyPatch(old, new, buggy); err == nil {
		t.Log("jsondiff 의 Factorize+LCS 버그가 고쳐진 것으로 보임 (참고용)")
	} else {
		b, _ := json.Marshal(buggy)
		t.Logf("Factorize+LCS patch %s → %v", b, err)
	}

	patch, err := audit.DiffSnapshots(old, new)
	if err != nil {
		t.Fatal(err)
	}
	if err := audit.VerifyPatch(old, new, patch); err != nil {
		b, _ := json.Marshal(patch)
		t.Fatalf("DiffSnapshots patch %s: %v", b, err)
	}
}

func TestDiffSnapshots(t *testing.T) {
	series := func(id, priority int, name string) map[string]any {
		return map[string]any{"id": id, "priority": priority, "i18n": map[string]any{
			fmt.Sprint(id*10 + 1): map[string]any{"id": id*10 + 1, "locale": "ko", "name": name},
		}}
	}
	keyed := func(rows ...map[string]any) map[string]any {
		out := map[string]any{}
		for _, r := range rows {
			out[fmt.Sprint(r["id"])] = r
		}
		return map[string]any{"series": out}
	}

	tests := []struct {
		name     string
		old, new map[string]any
		want     string
	}{
		{
			name: "가운데 행 삭제 + priority 재번호",
			old:  keyed(series(25730, 1, "A"), series(25734, 2, "B"), series(25736, 3, "C")),
			new:  keyed(series(25730, 0, "A"), series(25736, 1, "C")),
			want: `[{"value":0,"op":"replace","path":"/series/25730/priority"},` +
				`{"op":"remove","path":"/series/25734"},` +
				`{"value":1,"op":"replace","path":"/series/25736/priority"}]`,
		},
		{
			name: "하위 relation 값 변경",
			old:  keyed(series(25734, 0, "A")),
			new:  keyed(series(25734, 0, "B")),
			want: `[{"value":"B","op":"replace","path":"/series/25734/i18n/257341/name"}]`,
		},
		{
			name: "배열 컬럼 가운데 값 삭제는 remove 하나",
			old:  map[string]any{"viewLanguages": []any{"en", "ko", "ja"}},
			new:  map[string]any{"viewLanguages": []any{"en", "ja"}},
			want: `[{"op":"remove","path":"/viewLanguages/1"}]`,
		},
		{
			name: "컬럼 값 변경은 replace",
			old:  map[string]any{"tel": "02-111", "fax": nil},
			new:  map[string]any{"tel": "02-222", "fax": nil},
			want: `[{"value":"02-222","op":"replace","path":"/tel"}]`,
		},
		{
			name: "변경 없음은 빈 patch",
			old:  keyed(series(25734, 0, "A")),
			new:  keyed(series(25734, 0, "A")),
			want: `null`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch, err := audit.BuildPatch(tt.old, tt.new)
			if err != nil {
				t.Fatalf("BuildPatch: %v", err)
			}
			got, _ := json.Marshal(patch)
			if string(got) != tt.want {
				t.Errorf("patch =\n %s\nwant\n %s", got, tt.want)
			}
		})
	}
}

// 배열의 삭제·추가·순서 변경·중복을 무작위로 섞어서, 만든 patch 가 항상 수정 후 상태를 재현하는지 확인한다.
func TestBuildPatchRandomArrays(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	values := []string{"A", "B", "C", "D", "E", "F"}
	randomArray := func() []any {
		n := rng.Intn(7)
		out := make([]any, n)
		for i := range out {
			out[i] = values[rng.Intn(len(values))]
		}
		return out
	}

	const runs = 20000
	lcsFailures := 0
	for i := 0; i < runs; i++ {
		old := map[string]any{"types": randomArray(), "tel": nil}
		new := map[string]any{"types": randomArray(), "tel": nil}

		if p, _ := audit.DiffSnapshots(old, new); audit.VerifyPatch(old, new, p) != nil {
			lcsFailures++
		}
		patch, err := audit.BuildPatch(old, new)
		if err != nil {
			t.Fatalf("old=%v new=%v: %v", old, new, err)
		}
		if err := audit.VerifyPatch(old, new, patch); err != nil {
			b, _ := json.Marshal(patch)
			t.Fatalf("old=%v new=%v patch=%s: %v", old, new, b, err)
		}
	}
	t.Logf("%d 회 중 LCS patch 검증 실패(대체 diff 사용) %d 회", runs, lcsFailures)
}
