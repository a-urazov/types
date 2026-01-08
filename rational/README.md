# Пакет Rational

Пакет `rational` предоставляет реализацию рациональных чисел (дробей) с поддержкой различных математических операций без потери точности.

## Особенности

- Представление рациональных чисел в виде числителя и знаменателя
- Поддержка арифметических операций: сложение, вычитание, умножение, деление
- Автоматическое приведение к несократимому виду
- Поддержка больших чисел с помощью `math/big`
- Методы для сравнения, преобразования и анализа чисел

## Установка

```bash
# Пакет является частью модуля types
```

## Использование

### Создание рационального числа

```go
// Создание нового рационального числа 3/4
r1 := rational.New(3, 4)

// Создание рационального числа из целого числа
r2 := rational.NewFromInt(5)  // эквивалентно 5/1

// Создание из big.Int значений
num := big.NewInt(7)
den := big.NewInt(8)
r3 := rational.NewFromBigInt(num, den)  // 7/8
```

### Арифметические операции

```go
r1 := rational.New(1, 2)  // 1/2
r2 := rational.New(1, 3)  // 1/3

// Сложение
sum := r1.Add(r2)          // 1/2 + 1/3 = 5/6

// Вычитание
diff := r1.Subtract(r2)    // 1/2 - 1/3 = 1/6

// Умножение
prod := r1.Multiply(r2)    // 1/2 * 1/3 = 1/6

// Деление
quot := r1.Divide(r2)      // 1/2 ÷ 1/3 = 3/2
```

### Другие операции

```go
r := rational.New(3, 4)

// Приведение к несократимому виду (делается автоматически)
reduced := r.Reduce()

// Преобразование в другие типы
floatVal := r.ToFloat64()  // 0.75
intVal := r.ToInt64()      // 0 (целая часть)

// Сравнение
r1 := rational.New(1, 2)
r2 := rational.New(2, 4)   // тоже 1/2 после сокращения
isEqual := r1.Equals(r2)   // true

// Знак числа
sign := r.Sign()           // -1, 0 или 1
isPos := r.IsPositive()    // true если положительное
isNeg := r.IsNegative()    // true если отрицательное
isZero := r.IsZero()       // true если равно нулю

// Получение числителя и знаменателя
num := r.Numerator()       // big.Int
den := r.Denominator()     // big.Int

// Возведение в степень
squared := r.Power(2)      // (3/4)^2 = 9/16
inverse := r.Power(-1)     // обратное число = 4/3
```

## API

- `New(num, den int64) *Rational` - создает новое рациональное число
- `NewFromInt(value int64) *Rational` - создает рациональное число из целого
- `NewFromBigInt(num, den *big.Int) *Rational` - создает из big.Int значений
- `Add(other *Rational) *Rational` - сложение
- `Subtract(other *Rational) *Rational` - вычитание
- `Multiply(other *Rational) *Rational` - умножение
- `Divide(other *Rational) *Rational` - деление
- `Invert() *Rational` - мультипликативная инверсия (1/x)
- `Negate() *Rational` - аддитивная инверсия (-x)
- `Abs() *Rational` - абсолютное значение
- `Compare(other *Rational) int` - сравнение (-1, 0, 1)
- `Sign() int` - знак числа
- `IsZero(), IsPositive(), IsNegative()` - проверки свойств
- `ToFloat64(), ToInt64()` - преобразование к другим типам
- `Reduce() *Rational` - приведение к несократимому виду
- `Numerator(), Denominator()` - получение компонентов
- `Equals(other *Rational) bool` - проверка равенства
- `Clone() *Rational` - копирование
- `Power(exp int64) *Rational` - возведение в степень
- `String() string` - строковое представление
