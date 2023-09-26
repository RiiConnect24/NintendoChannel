package constants

import _ "embed"

var (
	//go:embed nc.bin
	DLList []byte

	//go:embed ratings/esrb/EC.jpg
	ECImage []byte

	//go:embed ratings/esrb/E.jpg
	EImage []byte

	//go:embed ratings/esrb/E10.jpg
	E10Image []byte

	//go:embed ratings/esrb/T.jpg
	TImage []byte

	//go:embed ratings/esrb/M.jpg
	MImage []byte

	//go:embed ratings/pegi/3.jpg
	PEGI3 []byte

	//go:embed ratings/pegi/7.jpg
	PEGI7 []byte

	//go:embed ratings/pegi/12.jpg
	PEGI12 []byte

	//go:embed ratings/pegi/16.jpg
	PEGI16 []byte

	//go:embed ratings/pegi/18.jpg
	PEGI18 []byte

	//go:embed ratings/cero/A.jpg
	CEROA []byte

	//go:embed ratings/cero/B.jpg
	CEROB []byte

	//go:embed ratings/cero/C.jpg
	CEROC []byte

	//go:embed ratings/cero/D.jpg
	CEROD []byte

	//go:embed ratings/cero/Z.jpg
	CEROZ []byte

	//go:embed ratings/esrb/EC-small.jpg
	ECImageSmall []byte

	//go:embed ratings/esrb/E-small.jpg
	EImageSmall []byte

	//go:embed ratings/esrb/E10-small.jpg
	E10ImageSmall []byte

	//go:embed ratings/esrb/T-small.jpg
	TImageSmall []byte

	//go:embed ratings/esrb/M-small.jpg
	MImageSmall []byte

	//go:embed ratings/pegi/3-small.jpg
	PEGI3Small []byte

	//go:embed ratings/pegi/7-small.jpg
	PEGI7Small []byte

	//go:embed ratings/pegi/12-small.jpg
	PEGI12Small []byte

	//go:embed ratings/pegi/16-small.jpg
	PEGI16Small []byte

	//go:embed ratings/pegi/18-small.jpg
	PEGI18Small []byte

	//go:embed ratings/cero/A-small.jpg
	CEROASmall []byte

	//go:embed ratings/cero/B-small.jpg
	CEROBSmall []byte

	//go:embed ratings/cero/C-small.jpg
	CEROCSmall []byte

	//go:embed ratings/cero/D-small.jpg
	CERODSmall []byte

	//go:embed ratings/cero/Z-small.jpg
	CEROZSmall []byte

	Images = map[RatingGroup][][]byte{
		CERO: {CEROA, CEROB, CEROC, CEROD, CEROZ},
		ESRB: {ECImage, EImage, E10Image, TImage, MImage},
		PEGI: {PEGI3, PEGI7, PEGI12, PEGI16, PEGI18},
	}

	ImagesSmall = map[RatingGroup][][]byte{
		CERO: {CEROASmall, CEROBSmall, CEROCSmall, CERODSmall, CEROZ},
		ESRB: {ECImageSmall, EImageSmall, E10ImageSmall, TImageSmall, MImage},
		PEGI: {PEGI3Small, PEGI7Small, PEGI12Small, PEGI16Small, PEGI18Small},
	}
)
