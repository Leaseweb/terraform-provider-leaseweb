package public_cloud

type Image struct {
	Id string
}

func NewImage(
	id string,
) Image {
	return Image{
		Id: id,
	}
}
