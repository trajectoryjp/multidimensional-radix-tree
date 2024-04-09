package lib

type Morton [16]byte

type BitValue byte

const (
	BitValueNull BitValue = 0b00
	BitBalueOn   BitValue = 0b11
	BitBalueOff  BitValue = 0b01
)

func ConvertIndexTo1DimMoton(zoom int, index int64) (morton Morton) {
	bit := int64(0x01)
	for z := 0; z < zoom; z++ {
		if index&bit == 0 {
			morton = morton.Set(z, false)
		} else {
			morton = morton.Set(z, true)
		}
	}
}
