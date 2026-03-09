# Go-Module (Internal Shared Library)

이 프로젝트는 고(Go) 기반 백엔드 컨테이너 서비스에서 공통으로 사용되는 라이브러리입니다.

## 📁 패키지 구성 및 역할
- **`/aws/s3`** : AWS S3와 관련된 파일들이 위치합니다.
- **`/config`** : 환경변수 설정과 관련된 파일들이 위치합니다.
- **`/db`** : 데이터베이스와 관련된 파일들이 위치합니다.
    - **`/clause`**: SQL Join, OrderBy 등 동적 쿼리 구문 생성을 담당합니다.
    - **`/cmd`**: Create, Update, Delete 쿼리 구문을 담당합니다.
    - **`/schema`**: 모델 구조체(Struct)의 태그를 분석하여 Select 필드를 자동 추출합니다.
- **`/idgen`** : 식별자 생성 로직과 관련된 유틸리티 파일들이 위치합니다.
- **`/validation`** : 유효성 검증에 사용되는 파일들이 위치합니다.

## 🚀 시작하기

### 1. 설치
```go
go get github.com/machinerd/go-module@vX.Y.Z
```


Requires Go >= 1.24 (toolchain supported)

### 2. import 적용

#### configuration 예시

* 초기 세팅

```go
// config.go
import (
  clause "github.com/machinerd/go-module/db/clause"
)

var cfg *module_config.Config

func ConfigLoad() {
  var err error
  cfg, err = module_config.New("./config/config.yml")
  if err != nil {
    panic(err)
  }

}

func GetConfig() *module_config.Config {
  return cfg
}
```

* config 설정 가져오기

```go
  k := config.GetConfig()

  bucket := k.String("s3_bucket")
```