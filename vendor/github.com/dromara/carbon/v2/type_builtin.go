package carbon

type (
	Timestamp      = TimestampType[timestampType]
	TimestampMilli = TimestampType[timestampMilliType]
	TimestampMicro = TimestampType[timestampMicroType]
	TimestampNano  = TimestampType[timestampNanoType]

	DateTime      = LayoutType[DateTimeType]
	DateTimeMicro = LayoutType[DateTimeMicroType]
	DateTimeMilli = LayoutType[DateTimeMilliType]
	DateTimeNano  = LayoutType[DateTimeNanoType]

	Date      = LayoutType[DateType]
	DateMilli = LayoutType[DateMilliType]
	DateMicro = LayoutType[DateMicroType]
	DateNano  = LayoutType[DateNanoType]

	Time      = LayoutType[TimeType]
	TimeMilli = LayoutType[TimeMilliType]
	TimeMicro = LayoutType[TimeMicroType]
	TimeNano  = LayoutType[TimeNanoType]
)

func NewTimestamp(c *Carbon) *Timestamp {
	return NewTimestampType[timestampType](c)
}
func NewTimestampMilli(c *Carbon) *TimestampMilli {
	return NewTimestampType[timestampMilliType](c)
}
func NewTimestampMicro(c *Carbon) *TimestampMicro {
	return NewTimestampType[timestampMicroType](c)
}
func NewTimestampNano(c *Carbon) *TimestampNano {
	return NewTimestampType[timestampNanoType](c)
}

func NewDateTime(c *Carbon) *DateTime {
	return NewLayoutType[DateTimeType](c)
}
func NewDateTimeMilli(c *Carbon) *DateTimeMilli {
	return NewLayoutType[DateTimeMilliType](c)
}
func NewDateTimeMicro(c *Carbon) *DateTimeMicro {
	return NewLayoutType[DateTimeMicroType](c)
}
func NewDateTimeNano(c *Carbon) *DateTimeNano {
	return NewLayoutType[DateTimeNanoType](c)
}

func NewDate(c *Carbon) *Date {
	return NewLayoutType[DateType](c)
}
func NewDateMilli(c *Carbon) *DateMilli {
	return NewLayoutType[DateMilliType](c)
}
func NewDateMicro(c *Carbon) *DateMicro {
	return NewLayoutType[DateMicroType](c)
}
func NewDateNano(c *Carbon) *DateNano {
	return NewLayoutType[DateNanoType](c)
}

func NewTime(c *Carbon) *Time {
	return NewLayoutType[TimeType](c)
}
func NewTimeMilli(c *Carbon) *TimeMilli {
	return NewLayoutType[TimeMilliType](c)
}
func NewTimeMicro(c *Carbon) *TimeMicro {
	return NewLayoutType[TimeMicroType](c)
}
func NewTimeNano(c *Carbon) *TimeNano {
	return NewLayoutType[TimeNanoType](c)
}

type timestampType int64

func (t timestampType) Precision() string {
	return PrecisionSecond
}

type timestampMilliType int64

func (t timestampMilliType) Precision() string {
	return PrecisionMillisecond
}
func (t timestampMilliType) DataType() string {
	return "timestamp(6)"
}

type timestampMicroType int64

func (t timestampMicroType) Precision() string {
	return PrecisionMicrosecond
}
func (t timestampMicroType) DataType() string {
	return "timestamp(6)"
}

type timestampNanoType int64

func (t timestampNanoType) Precision() string {
	return PrecisionNanosecond
}
func (t timestampNanoType) DataType() string {
	return "timestamp(6)"
}

type DateTimeType string

func (t DateTimeType) Layout() string {
	return DateTimeLayout
}
func (t DateTimeType) DataType() string {
	return "datetime"
}

type DateTimeMilliType string

func (t DateTimeMilliType) Layout() string {
	return DateTimeMilliLayout
}
func (t DateTimeMilliType) DataType() string {
	return "datetime(6)"
}

type DateTimeMicroType string

func (t DateTimeMicroType) Layout() string {
	return DateTimeMicroLayout
}
func (t DateTimeMicroType) DataType() string {
	return "datetime(6)"
}

type DateTimeNanoType string

func (t DateTimeNanoType) Layout() string {
	return DateTimeNanoLayout
}
func (t DateTimeNanoType) DataType() string {
	return "datetime(6)"
}

type DateType string

func (t DateType) Layout() string {
	return DateLayout
}
func (t DateType) DataType() string {
	return "date"
}

type DateMilliType string

func (t DateMilliType) Layout() string {
	return DateMilliLayout
}
func (t DateMilliType) DataType() string {
	return "datetime(6)"
}

type DateMicroType string

func (t DateMicroType) Layout() string {
	return DateMicroLayout
}
func (t DateMicroType) DataType() string {
	return "datetime(6)"
}

type DateNanoType string

func (t DateNanoType) Layout() string {
	return DateNanoLayout
}
func (t DateNanoType) DataType() string {
	return "datetime(6)"
}

type TimeType string

func (t TimeType) Layout() string {
	return TimeLayout
}
func (t TimeType) DataType() string {
	return "time"
}

type TimeMilliType string

func (t TimeMilliType) Layout() string {
	return TimeMilliLayout
}
func (t TimeMilliType) DataType() string {
	return "datetime(6)"
}

type TimeMicroType string

func (t TimeMicroType) Layout() string {
	return TimeMicroLayout
}
func (t TimeMicroType) DataType() string {
	return "datetime(6)"
}

type TimeNanoType string

func (t TimeNanoType) Layout() string {
	return TimeNanoLayout
}
func (t TimeNanoType) DataType() string {
	return "datetime(6)"
}
