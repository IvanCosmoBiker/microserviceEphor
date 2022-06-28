package fiscal

import (
    interfaceFiscal "ephorservices/src/server/utils/interface/fiscal"
    ofdFerma "ephorservices/src/server/services/fiscal/fermaOfd"
    ofdOrange "ephorservices/src/server/services/fiscal/orange" 
    ofdNanokass "ephorservices/src/server/services/fiscal/nanokass"
)
// instance of type fiscal
var ofd ofdFerma.NewFermaOfdStruct
var Orange ofdOrange.NewOrangeStruct
var Nanokass ofdNanokass.NewNanokassStruct
func GetFiscal(fiscal int) (interfaceFiscal.Fiscal) {
    switch fiscal {
        case interfaceFiscal.Fr_EphorOrangeData,
        interfaceFiscal.Fr_ServerOrangeData,
        interfaceFiscal.Fr_EphorServerOrangeData,
        interfaceFiscal.Fr_OrangeData:
        return Orange.NewFiscal()
        fallthrough
        case interfaceFiscal.Fr_NanoKassa,
        interfaceFiscal.Fr_ServerNanoKassa:
        return Nanokass.NewFiscal()
        fallthrough
        case interfaceFiscal.Fr_OFD:
        return ofd.NewFiscal()
    }
    return nil
}