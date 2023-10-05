# GOSQL

database/sql (official package) extension

### Type 정의

1. type mapping
    1. GType: go 언어 타입
    2. DType: db sql 타입
    3. SType: gosql 명세에 사용하는 타입
    4. UType: SType, GType, DType 3개 모두 가지고 있는 데이터 형

   |   | SType   | GType           | DType (mysql) |
                  |---|---------|-----------------|---------------|
   | 1 | string  | string          | varchar       |
   | 2 | int     | int64           | bigint        |
   | 3 | float   | float64         | double        |
   | 4 | bytes   | []byte          | blob          |
   | 5 | uuid    | gyp.UUID        | char(36)      |
   | 6 | time    | time.Time       | datetime      |
   |   7 | decimal | decimal.Decimal | decimal       |


2. custom type
    * 사용자 지정 타입의 경우 다음 두 인터페이스를 구현해야 한다.
        1. driver.Valuer
        2. sql.Scanner
    * gosql 명세에 추가해야 한다.
        1. uuid 는 기본타입이 아니므로 명세에 추가해준다.

