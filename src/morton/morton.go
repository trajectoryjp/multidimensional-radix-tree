package morton

// index 0b1110-1101-1011
// byteごと反転　 0111-1011-1101
// morton [0]Off-On-On-On [1]On-Off-On-On [2]On-On-Off-On
//	(bit) [0]01111111     [1]11011111     [2]11110111

type Morton [16]byte

type BitValue byte

const (
	BitValueNull BitValue = 0b00
	BitBalueOn   BitValue = 0b11
	BitBalueOff  BitValue = 0b01
)

func CreateMorton(indexs Indexs, zoomSetLevel ZoomSetLevel, zoomSetLevelTable ZoomSetTable) (morton Morton) {

	for zs := ZoomSetLevel(0); zs < zoomSetLevel; zs++ {
		for dim := 0; dim < len(indexs); dim++ {
			digit := zoomSetLevelTable.GetZoom(zs, dim) // 有効桁数
			//m := ConvertIndexTo1DimMoton(digit, indexs[dim])
			//morton.Push(m)

			//畳み込み
		}
	}
	return morton
}

// 例
// index = 0b1011  ズームレベル4（ズームレベル1=1, 2=0, 3=1, 4=1)
// mortonはLSBがズームレベル1
// -> 0b00 11 11 01 11
func ConvertIndexTo1DimMoton(zoom ZoomLevel, index int64) (morton Morton) {
	bit := int64(0x01)
	for z := 0; z < int(zoom); z++ {
		if index&bit == 0 {
			morton = morton.Set(z, false)
		} else {
			morton = morton.Set(z, true)
		}
	}
	return morton
}

// zoomの位置にflagを設定する
func (mr Morton) Set(zoom int, flag bool) Morton {
	byteNum := zoom / 4
	bitNum := byte(zoom%4) * 2
	revBitNum := reverseByte(bitNum)
	if flag {
		mr[byteNum] = mr[byteNum] | revBitNum

	} else {
		mr[byteNum] = mr[byteNum] & ^revBitNum
	}
	mr[byteNum] = mr[byteNum] | 0b01
	return mr
}

func reverseByte(in byte) (out byte) {
	for k := 0; k < 8; k++ {
		b := in & 0x01
		out = out<<1 | b
		in = in >> 1
	}
	return out
}

/*
func (mr *Morton) Push(p Morton) {
	for bt := range p {
		for k := 0; k < 4; k++ {
			v := bt & 0b11
			mr.Append(v)
		}
	}
}

func (mr *Morton) Append(twoBit int) {

}
*/

func (mr *Morton) Digit() int {
	digit := 0
	for v := range *mr {
		checkMask := 0b01
		for b := 0; b < 8; b += 2 {
			if v&checkMask == 0 {
				return digit
			}
			digit++
			checkMask = checkMask << 2
		}
	}
	return 0
}

/*
// 子を追加する方向へのシフト（1桁上）
func (mr *Morton) UpShift() {
	for k := 0; k < len(mr)-1; k++ {
		mr[k] = mr[k]<<2 | (mr[k-1] >> 6 & 0b11)
	}
}

// 子を削除する方向へのシフト（1桁下）
func (mr *Morton) DownShift() {
	for k := 0; k < len(mr)-1; k++ {
		mr[k] = mr[k]<<2 | (mr[k-1] >> 6 & 0b11)
	}
}
*/

func (mr *Morton) ConvertToSuffix(shiftDigit int) (branchNumber int) {
	for d := 0; d < shiftDigit; d++ {
		v := (mr[0] | 0b11) >> 1
		branchNumber = branchNumber<<1 | int(v)
		mr.DownShift()
	}
	return branchNumber
}
