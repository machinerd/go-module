package audit

// 외부 테스트 패키지(audit_test)에서 비공개 함수를 쓰기 위한 연결. 테스트 빌드에만 포함된다.
var (
	DiffSnapshots = diffSnapshots
	BuildPatch    = buildPatch
	VerifyPatch   = verifyPatch
)
