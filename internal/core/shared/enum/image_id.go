package enum

type ImageId string

func (i ImageId) String() string {
	return string(i)
}

const (
	Almalinux864Bit                ImageId = "ALMALINUX_8_64BIT"
	Almalinux964Bit                ImageId = "ALMALINUX_9_64BIT"
	ArchLinux64Bit                 ImageId = "ARCH_LINUX_64BIT"
	Centos764Bit                   ImageId = "CENTOS_7_64BIT"
	Debian1064Bit                  ImageId = "DEBIAN_10_64BIT"
	Debian1164Bit                  ImageId = "DEBIAN_11_64BIT"
	Debian1264Bit                  ImageId = "DEBIAN_12_64BIT"
	Freebsd1364Bit                 ImageId = "FREEBSD_13_64BIT"
	Freebsd1464Bit                 ImageId = "FREEBSD_14_64BIT"
	RockyLinux864Bit               ImageId = "ROCKY_LINUX_8_64BIT"
	RockyLinux964Bit               ImageId = "ROCKY_LINUX_9_64BIT"
	Ubuntu200464Bit                ImageId = "UBUNTU_20_04_64BIT"
	Ubuntu220464Bit                ImageId = "UBUNTU_22_04_64BIT"
	Ubuntu240464Bit                ImageId = "UBUNTU_24_04_64BIT"
	WindowsServer2016Standard64Bit ImageId = "WINDOWS_SERVER_2016_STANDARD_64BIT"
	WindowsServer2019Standard64Bit ImageId = "WINDOWS_SERVER_2019_STANDARD_64BIT"
	WindowsServer2022Standard64Bit ImageId = "WINDOWS_SERVER_2022_STANDARD_64BIT"
)

func (i ImageId) Values() []string {
	var stringValues []string

	for _, imageId := range imageIds {
		stringValues = append(stringValues, string(imageId))
	}

	return stringValues
}

var imageIds = []ImageId{
	Almalinux864Bit,
	Almalinux964Bit,
	ArchLinux64Bit,
	Centos764Bit,
	Debian1264Bit,
	Debian1064Bit,
	Debian1164Bit,
	Debian1264Bit,
	Freebsd1364Bit,
	Freebsd1464Bit,
	RockyLinux864Bit,
	RockyLinux964Bit,
	Ubuntu200464Bit,
	Ubuntu220464Bit,
	Ubuntu240464Bit,
	WindowsServer2016Standard64Bit,
	WindowsServer2019Standard64Bit,
	WindowsServer2022Standard64Bit,
}

func NewImageId(value string) (ImageId, error) {
	return findEnumForString(value, imageIds, Almalinux864Bit)
}
