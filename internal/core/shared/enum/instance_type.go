package enum

type InstanceType string

func (i InstanceType) String() string {
	return string(i)
}

func (i InstanceType) Values() []string {
	var stringValues []string

	for _, instanceType := range instanceTypes {
		stringValues = append(stringValues, string(instanceType))
	}

	return stringValues
}

const (
	InstanceTypeM3Large     InstanceType = "lsw.m3.large"
	InstanceTypeM3Xlarge    InstanceType = "lsw.m3.xlarge"
	InstanceTypeM32Xlarge   InstanceType = "lsw.m3.2xlarge"
	InstanceTypeM4Large     InstanceType = "lsw.m4.large"
	InstanceTypeM4Xlarge    InstanceType = "lsw.m4.xlarge"
	InstanceTypeM42Xlarge   InstanceType = "lsw.m4.2xlarge"
	InstanceTypeM44Xlarge   InstanceType = "lsw.m4.4xlarge"
	InstanceTypeM5Large     InstanceType = "lsw.m5.large"
	InstanceTypeM5Xlarge    InstanceType = "lsw.m5.xlarge"
	InstanceTypeM52Xlarge   InstanceType = "lsw.m5.2xlarge"
	InstanceTypeM54Xlarge   InstanceType = "lsw.m5.4xlarge"
	InstanceTypeM5ALarge    InstanceType = "lsw.m5a.large"
	InstanceTypeM5AXlarge   InstanceType = "lsw.m5a.xlarge"
	InstanceTypeM5A2Xlarge  InstanceType = "lsw.m5a.2xlarge"
	InstanceTypeM5A4Xlarge  InstanceType = "lsw.m5a.4xlarge"
	InstanceTypeM5A8Xlarge  InstanceType = "lsw.m5a.8xlarge"
	InstanceTypeM5A12Xlarge InstanceType = "lsw.m5a.12xlarge"
	InstanceTypeM6ALarge    InstanceType = "lsw.m6a.large"
	InstanceTypeM6AXlarge   InstanceType = "lsw.m6a.xlarge"
	InstanceTypeM6A2Xlarge  InstanceType = "lsw.m6a.2xlarge"
	InstanceTypeM6A4Xlarge  InstanceType = "lsw.m6a.4xlarge"
	InstanceTypeM6A8Xlarge  InstanceType = "lsw.m6a.8xlarge"
	InstanceTypeM6A12Xlarge InstanceType = "lsw.m6a.12xlarge"
	InstanceTypeM6A16Xlarge InstanceType = "lsw.m6a.16xlarge"
	InstanceTypeM6A24Xlarge InstanceType = "lsw.m6a.24xlarge"
	InstanceTypeC3Large     InstanceType = "lsw.c3.large"
	InstanceTypeC3Xlarge    InstanceType = "lsw.c3.xlarge"
	InstanceTypeC32Xlarge   InstanceType = "lsw.c3.2xlarge"
	InstanceTypeC34Xlarge   InstanceType = "lsw.c3.4xlarge"
	InstanceTypeC4Large     InstanceType = "lsw.c4.large"
	InstanceTypeC4Xlarge    InstanceType = "lsw.c4.xlarge"
	InstanceTypeC42Xlarge   InstanceType = "lsw.c4.2xlarge"
	InstanceTypeC44Xlarge   InstanceType = "lsw.c4.4xlarge"
	InstanceTypeC5Large     InstanceType = "lsw.c5.large"
	InstanceTypeC5Xlarge    InstanceType = "lsw.c5.xlarge"
	InstanceTypeC52Xlarge   InstanceType = "lsw.c5.2xlarge"
	InstanceTypeC54Xlarge   InstanceType = "lsw.c5.4xlarge"
	InstanceTypeC5ALarge    InstanceType = "lsw.c5a.large"
	InstanceTypeC5AXlarge   InstanceType = "lsw.c5a.xlarge"
	InstanceTypeC5A2Xlarge  InstanceType = "lsw.c5a.2xlarge"
	InstanceTypeC5A4Xlarge  InstanceType = "lsw.c5a.4xlarge"
	InstanceTypeC5A9Xlarge  InstanceType = "lsw.c5a.9xlarge"
	InstanceTypeC5A12Xlarge InstanceType = "lsw.c5a.12xlarge"
	InstanceTypeC6ALarge    InstanceType = "lsw.c6a.large"
	InstanceTypeC6AXlarge   InstanceType = "lsw.c6a.xlarge"
	InstanceTypeC6A2Xlarge  InstanceType = "lsw.c6a.2xlarge"
	InstanceTypeC6A4Xlarge  InstanceType = "lsw.c6a.4xlarge"
	InstanceTypeC6A8Xlarge  InstanceType = "lsw.c6a.8xlarge"
	InstanceTypeC6A12Xlarge InstanceType = "lsw.c6a.12xlarge"
	InstanceTypeC6A16Xlarge InstanceType = "lsw.c6a.16xlarge"
	InstanceTypeC6A24Xlarge InstanceType = "lsw.c6a.24xlarge"
	InstanceTypeR3Large     InstanceType = "lsw.r3.large"
	InstanceTypeR3Xlarge    InstanceType = "lsw.r3.xlarge"
	InstanceTypeR32Xlarge   InstanceType = "lsw.r3.2xlarge"
	InstanceTypeR4large     InstanceType = "lsw.r4.large"
	InstanceTypeR4Xlarge    InstanceType = "lsw.r4.xlarge"
	InstanceTypeR42Xlarge   InstanceType = "lsw.r4.2xlarge"
	InstanceTypeR5Large     InstanceType = "lsw.r5.large"
	InstanceTypeR5Xlarge    InstanceType = "lsw.r5.xlarge"
	InstanceTypeR52Xlarge   InstanceType = "lsw.r5.2xlarge"
	InstanceTypeR5ALarge    InstanceType = "lsw.r5a.large"
	InstanceTypeR5AXlarge   InstanceType = "lsw.r5a.xlarge"
	InstanceTypeR5A2Xlarge  InstanceType = "lsw.r5a.2xlarge"
	InstanceTypeR5A4Xlarge  InstanceType = "lsw.r5a.4xlarge"
	InstanceTypeR5A8Xlarge  InstanceType = "lsw.r5a.8xlarge"
	InstanceTypeR5A12Xlarge InstanceType = "lsw.r5a.12xlarge"
	InstanceTypeR6ALarge    InstanceType = "lsw.r6a.large"
	InstanceTypeR6AXlarge   InstanceType = "lsw.r6a.xlarge"
	InstanceTypeR6A2Xlarge  InstanceType = "lsw.r6a.2xlarge"
	InstanceTypeR6A4Xlarge  InstanceType = "lsw.r6a.4xlarge"
	InstanceTypeR6A8Xlarge  InstanceType = "lsw.r6a.8xlarge"
	InstanceTypeR6A12Xlarge InstanceType = "lsw.r6a.12xlarge"
	InstanceTypeR6A16Xlarge InstanceType = "lsw.r6a.16xlarge"
	InstanceTypeR6A24Xlarge InstanceType = "lsw.r6a.24xlarge"
)

var instanceTypes = []InstanceType{
	"lsw.m3.large",
	"lsw.m3.xlarge",
	"lsw.m3.2xlarge",
	"lsw.m4.large",
	"lsw.m4.xlarge",
	"lsw.m4.2xlarge",
	"lsw.m4.4xlarge",
	"lsw.m5.large",
	"lsw.m5.xlarge",
	"lsw.m5.2xlarge",
	"lsw.m5.4xlarge",
	"lsw.m5a.large",
	"lsw.m5a.xlarge",
	"lsw.m5a.2xlarge",
	"lsw.m5a.4xlarge",
	"lsw.m5a.8xlarge",
	"lsw.m5a.12xlarge",
	"lsw.m6a.large",
	"lsw.m6a.xlarge",
	"lsw.m6a.2xlarge",
	"lsw.m6a.4xlarge",
	"lsw.m6a.8xlarge",
	"lsw.m6a.12xlarge",
	"lsw.m6a.16xlarge",
	"lsw.m6a.24xlarge",
	"lsw.c3.large",
	"lsw.c3.xlarge",
	"lsw.c3.2xlarge",
	"lsw.c3.4xlarge",
	"lsw.c4.large",
	"lsw.c4.xlarge",
	"lsw.c4.2xlarge",
	"lsw.c4.4xlarge",
	"lsw.c5.large",
	"lsw.c5.xlarge",
	"lsw.c5.2xlarge",
	"lsw.c5.4xlarge",
	"lsw.c5a.large",
	"lsw.c5a.xlarge",
	"lsw.c5a.2xlarge",
	"lsw.c5a.4xlarge",
	"lsw.c5a.9xlarge",
	"lsw.c5a.12xlarge",
	"lsw.c6a.large",
	"lsw.c6a.xlarge",
	"lsw.c6a.2xlarge",
	"lsw.c6a.4xlarge",
	"lsw.c6a.8xlarge",
	"lsw.c6a.12xlarge",
	"lsw.c6a.16xlarge",
	"lsw.c6a.24xlarge",
	"lsw.r3.large",
	"lsw.r3.xlarge",
	"lsw.r3.2xlarge",
	"lsw.r4.large",
	"lsw.r4.xlarge",
	"lsw.r4.2xlarge",
	"lsw.r5.large",
	"lsw.r5.xlarge",
	"lsw.r5.2xlarge",
	"lsw.r5a.large",
	"lsw.r5a.xlarge",
	"lsw.r5a.2xlarge",
	"lsw.r5a.4xlarge",
	"lsw.r5a.8xlarge",
	"lsw.r5a.12xlarge",
	"lsw.r6a.large",
	"lsw.r6a.xlarge",
	"lsw.r6a.2xlarge",
	"lsw.r6a.4xlarge",
	"lsw.r6a.8xlarge",
	"lsw.r6a.12xlarge",
	"lsw.r6a.16xlarge",
	"lsw.r6a.24xlarge",
}

func NewInstanceType(value string) (InstanceType, error) {
	return findEnumForString(value, instanceTypes, InstanceTypeC3Large)
}
