package device

import (
    interfaceDevice "ephorservices/src/server/utils/interface/device"
    coolerDevice "ephorservices/src/server/services/device/cooler"
    automatDevice "ephorservices/src/server/services/device/automat"
)
// instance of type device
var cooler coolerDevice.NewCoolerStruct
var automat automatDevice.NewAutomatStruct

func GetDevice(device int) (interfaceDevice.Device) {
    switch device {
        case interfaceDevice.TypeCoffee,
        interfaceDevice.TypeSnack,
        interfaceDevice.TypeHoreca,
        interfaceDevice.TypeSodaWater,
        interfaceDevice.TypeMechanical,
        interfaceDevice.TypeComb:
        return automat.NewDevice()
        fallthrough
        case interfaceDevice.TypeCooler:
            return cooler.NewDevice()
    }
    return nil
}
